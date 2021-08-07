package cmd

import (
	"github.com/mingalevme/avito/internal/env"
	"github.com/pkg/errors"
	"sync"
)

type CheckCmd struct {
	URLS []string `arg name:"url" help:"URLs"`
}

func (r *CheckCmd) Run(e *env.Env) error {
	checker := e.Checker()
	var wg sync.WaitGroup
	for _, url := range r.URLS {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := checker.Check(url); err != nil {
				e.Logger().
					WithError(err).
					WithField("url", url).
					Errorf("error while checking url")
			}
		}(url)
	}
	wg.Wait()
	e.Logger().Debugf("Syncing repository")
	if err := e.Repository().Sync(); err != nil {
		return errors.Wrap(err, "error while syncing repository")
	}
	return nil
}
