package notifier

import (
	"github.com/mingalevme/avito/pkg/model"
	log "github.com/mingalevme/gologger"
	"sync"
)

type StackNotifier struct {
	Notifiers []Notifier
	Logger    log.Logger
}

func NewStackNotifier(logger log.Logger) *StackNotifier {
	return &StackNotifier{
		Notifiers: []Notifier{},
		Logger:    logger,
	}
}

func (n *StackNotifier) Notify(item model.Item) error {
	var wg sync.WaitGroup
	for _, notifier := range n.Notifiers {
		wg.Add(1)
		go func(notifier Notifier) {
			defer wg.Done()
			err := notifier.Notify(item)
			if err != nil {
				n.Logger.WithError(err).Errorf("error while notifying")
			}
		}(notifier)
	}
	wg.Wait()
	return nil
}

func (n *StackNotifier) AddNotifier(notifier Notifier) {
	n.Notifiers = append(n.Notifiers, notifier)
}
