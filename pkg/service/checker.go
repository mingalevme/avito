package service

import (
	"fmt"
	"github.com/mingalevme/avito/pkg/model"
	"github.com/mingalevme/avito/pkg/notifier"
	"github.com/mingalevme/avito/pkg/parser"
	"github.com/mingalevme/avito/pkg/repository"
	log "github.com/mingalevme/gologger"
	"github.com/pkg/errors"
	"sync"
)

type Checker struct {
	Parser     *parser.Parser
	Repository repository.Repository
	Notifier   notifier.Notifier
	Logger     log.Logger
}

func NewChecker(p *parser.Parser, r repository.Repository, n notifier.Notifier, l log.Logger) *Checker {
	return &Checker{
		Parser:     p,
		Repository: r,
		Notifier:   n,
		Logger:     l,
	}
}

func (c Checker) Check(sourceUrl string) error {
	c.Logger.Infof("Starting parsing: %s ...", sourceUrl)
	items, err := c.Parser.Parse(sourceUrl)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error while parsing data from url: %s", sourceUrl))
	}
	if len(items) == 0 {
		return fmt.Errorf("empty result while parsing url: %s", sourceUrl)
	}
	var wg sync.WaitGroup
	for _, item := range items {
		if c.Repository.Has(item.ID) {
			c.Logger.Debugf("Item already exist: %s %s", item.Title, item.URL)
			continue
		}
		c.Logger.Debugf("Adding new item: %s %s", item.Title, item.URL)
		if err = c.Repository.Add(item); err != nil {
			c.Logger.WithError(err).Error("error while adding item to repository")
			continue
		}
		c.Logger.Debugf("New item has been added: %s %s", item.Title, item.URL)
		c.Logger.Debugf("Notifying new item: %s %s", item.Title, item.URL)
		wg.Add(1)
		go func(item model.Item) {
			defer wg.Done()
			if err := c.Notifier.Notify(item); err != nil {
				c.Logger.
					WithError(err).
					WithField("title", item.Title).
					WithField("url", item.URL).
					Errorf("error while notifying")
			}
		}(item)
	}
	wg.Wait()
	c.Logger.Debugf("Repository syncing has been finished successfully")
	c.Logger.Infof("Parsing has been finished: %s ...", sourceUrl)
	return nil
}
