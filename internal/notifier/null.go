package notifier

import (
	"github.com/mingalevme/avito/internal/model"
)

type NullNotifier struct {

}

func NewNullNotifier() *NullNotifier {
	return &NullNotifier{}
}

func (n *NullNotifier) Notify(item model.Item) error {
	return nil
}
