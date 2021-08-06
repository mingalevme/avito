package cmd

import (
	"github.com/mingalevme/avito/internal/env"
)

type CheckCmd struct {
	URLS []string `arg name:"url" help:"URLs"`
}

func (r *CheckCmd) Run(e *env.Env) error {
	checker := e.Checker()
	for _, url := range r.URLS {
		if err := checker.Check(url); err != nil {
			e.Logger().
				WithError(err).
				WithField("url", url).
				Errorf("error while checking url")
		}
	}
	return nil
}
