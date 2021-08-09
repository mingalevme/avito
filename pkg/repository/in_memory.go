package repository

import (
	"github.com/mingalevme/avito/pkg/model"
	"sync"
)

type InMemoryRepository struct {
	data map[int]model.Item
	lock sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: map[int]model.Item{},
	}
}

func (r *InMemoryRepository) GetAll() ([]model.Item, error) {
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

func (r *InMemoryRepository) Get(id int) (model.Item, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	item, ok := r.data[id]
	if !ok {
		return model.Item{}, ErrNotFound
	}
	return item, nil
}

func (r *InMemoryRepository) Has(id int) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.data[id]
	return ok
}

func (r *InMemoryRepository) Add(item model.Item) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[item.ID] = item
	return nil
}

func (r *InMemoryRepository) Size() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.data)
}

func (r *InMemoryRepository) Sync() error {
	return nil
}

func (r *InMemoryRepository) Clean() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data = map[int]model.Item{}
	return nil
}

func (r *InMemoryRepository) GetData() map[int]model.Item {
	r.lock.RLock()
	defer r.lock.RUnlock()
	items := make(map[int]model.Item)
	for id, item := range r.data {
		items[id] = item
	}
	return items
}
