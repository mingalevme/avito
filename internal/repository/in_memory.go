package repository

import (
	"github.com/mingalevme/avito/internal/model"
)

type InMemoryRepository struct {
	Data map[int]model.Item
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		Data: map[int]model.Item{},
	}
}

func (r *InMemoryRepository) GetAll() ([]model.Item, error) {
	items := make([]model.Item, len(r.Data))
	i := 0
	for _, item := range r.Data {
		items[i] = item
		i++
	}
	return items, nil
}

func (r *InMemoryRepository) Get(id int) (model.Item, error) {
	item, ok := r.Data[id]
	if !ok {
		return model.Item{}, ErrNotFound
	}
	return item, nil
}

func (r *InMemoryRepository) Has(id int) bool {
	_, ok := r.Data[id]
	return ok
}

func (r *InMemoryRepository) Add(item model.Item) error {
	r.Data[item.ID] = item
	return nil
}

func (r *InMemoryRepository) Sync() error {
	return nil
}

func (r *InMemoryRepository) Clean() error {
	r.Data = map[int]model.Item{}
	return nil
}
