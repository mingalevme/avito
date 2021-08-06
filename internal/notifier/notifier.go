package notifier

import (
	"github.com/mingalevme/avito/internal/model"
)

type Notifier interface {
	Notify(item model.Item) error
}
