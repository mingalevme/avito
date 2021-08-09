package repository

import (
	"encoding/json"
	"fmt"
	"github.com/mingalevme/avito/pkg/model"
	"github.com/mingalevme/gologger"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

type FileRepository struct {
	filename string
	storage  *InMemoryRepository
	lock     sync.RWMutex
	logger   gologger.Logger
}

func NewFileRepository(filename string, logger gologger.Logger) (*FileRepository, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	r := &FileRepository{
		filename: filename,
		logger:   logger,
		storage:  NewInMemoryRepository(),
	}
	if err = r.init(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *FileRepository) GetAll() ([]model.Item, error) {
	return r.storage.GetAll()
}

func (r *FileRepository) Get(id int) (model.Item, error) {
	return r.storage.Get(id)
}

func (r *FileRepository) Has(id int) bool {
	return r.storage.Has(id)
}

func (r *FileRepository) Add(item model.Item) error {
	return r.storage.Add(item)
}

func (r *FileRepository) Size() int {
	return r.storage.Size()
}

func (r *FileRepository) Sync() error {
	if r.Size() == 0 {
		r.lock.Lock()
		defer r.lock.Unlock()
		return os.Remove(r.filename)
	}
	data, err := json.Marshal(r.storage.GetData())
	if err != nil {
		return errors.Wrap(err, "error while marshalling data to json")
	}
	proxy := r.filename + "-" + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	f, err := os.Create(proxy)
	if err != nil {
		return errors.Wrap(err, "error while creating proxy json file")
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			r.logger.WithError(err).Errorf("error while closing proxy json file")
		}
	}(f)
	n, err := f.Write(data)
	if err != nil {
		return errors.Wrap(err, "error while writing data to proxy json file")
	}
	if n != len(data) {
		return errors.Wrap(err, fmt.Sprintf("invalid bytes count wrote to to proxy json file: %d/%d", n, len(data)))
	}
	r.lock.Lock()
	err = os.Rename(proxy, r.filename)
	r.lock.Unlock()
	if err != nil {
		_ = os.Remove(proxy)
		return errors.Wrap(err, fmt.Sprintf("error while renaming proxy json file: %s", proxy))
	}
	return nil
}

func (r *FileRepository) Clean() error {
	return r.storage.Clean()
}

func (r *FileRepository) init() error {
	data, err := r.read()
	if err != nil {
		return errors.Wrap(err, "file repository: error while reading data from file while initializing")
	}
	for _, item := range data {
		if err := r.storage.Add(item); err != nil {
			return errors.Wrap(err, "file repository: error adding item while initializing")
		}
	}
	return nil
}

func (r *FileRepository) read() (map[int]model.Item, error) {
	data, err := ioutil.ReadFile(r.filename)
	if err != nil {
		return nil, errors.Wrap(err, "error while reading file: "+r.filename)
	}
	if len(data) == 0 {
		return map[int]model.Item{}, nil
	}
	var items map[int]model.Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, errors.Wrap(err, "error while unmarshalling json from file: "+r.filename)
	}
	if items == nil {
		items = map[int]model.Item{}
	}
	return items, nil
}
