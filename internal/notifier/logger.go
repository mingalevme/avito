package notifier

import (
	"github.com/mingalevme/avito/internal/model"
	log "github.com/mingalevme/gologger"
)

type LoggerNotifier struct {
	Logger log.Logger
	Level  log.Level
}

func NewLoggerNotifier(logger log.Logger, level log.Level) *LoggerNotifier {
	return &LoggerNotifier{
		Logger: logger,
		Level:  level,
	}
}

func (n *LoggerNotifier) Notify(item model.Item) error {
	n.Logger.
		WithField("Title", item.Title).
		WithField("Url", item.URL).
		WithField("Price", item.GetFormattedPrice).
		WithField("Image", item.ImageURL).
		Log(n.Level, "New Item")
	return nil
}
