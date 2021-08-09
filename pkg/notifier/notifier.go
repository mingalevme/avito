package notifier

import (
	"github.com/mingalevme/avito/pkg/model"
)

type Notifier interface {
	Notify(item model.Item) error
}
