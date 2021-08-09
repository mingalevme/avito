package notifier

import (
	"github.com/mingalevme/avito/pkg/model"
)

type NullNotifier struct {

}

func NewNullNotifier() *NullNotifier {
	return &NullNotifier{}
}

func (n *NullNotifier) Notify(item model.Item) error {
	return nil
}
