package notifier

import (
	"fmt"
	"github.com/mingalevme/avito/pkg/model"
)

type StdoutNotifier struct {
}

func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{}
}

func (n *StdoutNotifier) Notify(item model.Item) error {
	fmt.Printf("[Stdout Notifier] New item: %s %s %s %s\n", item.Title, item.URL, item.GetFormattedPrice(), item.ImageURL)
	return nil
}
