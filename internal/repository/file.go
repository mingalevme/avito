package repository

import (
	"encoding/json"
	"fmt"
	"github.com/mingalevme/avito/internal/model"
	gologger "github.com/mingalevme/gologger"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

type FileRepository struct {
	filename string
	data     map[int]model.Item
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
	}
	if err = r.init(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *FileRepository) GetAll() ([]model.Item, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	items := make([]model.Item, len(r.data))
	i := 0
	for _, item := range r.data {
		items[i] = item
		i++
	}
	return items, nil
}

func (r *FileRepository) Get(id int) (model.Item, error) {
	item, ok := r.data[id]
	if !ok {
		return model.Item{}, ErrNotFound
	}
	return item, nil
}

func (r *FileRepository) Has(id int) bool {
	_, ok := r.data[id]
	return ok
}

func (r *FileRepository) Add(item model.Item) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[item.ID] = item
	return nil
}

func (r *FileRepository) Sync() error {
	r.lock.Lock()
	data, err := json.Marshal(r.data)
	r.lock.Unlock()
	if err != nil {
		return errors.Wrap(err, "error while marshalling data to json")
	}
	proxy := r.filename + "-" + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	f, err := os.Create(proxy)
	if err != nil {
		return errors.Wrap(err, "error while creating proxy json file")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
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
	if err := os.Remove(r.filename); err != nil {
		panic(err)
	}
	return nil
}

func (r *FileRepository) init() error {
	if items, err := r.read(); err != nil {
		return err
	} else {
		r.data = items
		return nil
	}
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
