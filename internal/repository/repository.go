package repository

import (
	"errors"
	"github.com/mingalevme/avito/internal/model"
)

var ErrNotFound = errors.New("item not found")

type Repository interface {
	GetAll() ([]model.Item, error)
	Get(id int) (model.Item, error)
	Has(id int) bool
	Add(item model.Item) error
	Sync() error
	Clean() error
}
