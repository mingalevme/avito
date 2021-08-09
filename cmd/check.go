package cmd

import (
	"fmt"
	"github.com/mingalevme/avito/pkg/env"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"sync"
)

type CheckCmd struct {
	URLS  []string `arg optional name:"url" help:"URLs"`
	Files []string `optional name:"file" short:"f" help:"Paths to files to read URLs from" type:"path"`
}

func (r *CheckCmd) Run(e *env.Env) error {
	checker := e.Checker()
	var wg sync.WaitGroup
	urls := append(r.URLS, r.readFiles(e)...)
	for _, url := range urls {
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

func (r *CheckCmd) readFiles(e *env.Env) []string {
	result := make([]string, 0)
	for _, filename := range r.Files {
		e.Logger().WithField("filename", filename).Debugf("File reading starting ...")
		var urls []string
		if filename == "-" {
			urls = r.readStdin()
		} else {
			urls = r.readFile(e, filename)
		}
		e.Logger().WithField("filename", filename).WithField("urls", urls).Debugf("File reading finished")
		result = append(result, urls...)
	}
	return result
}

func (r *CheckCmd) readFile(e *env.Env, filename string) []string {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(errors.Wrap(err, "error while reading file: "+filename))
	}
	e.Logger().WithField("filename", filename).WithField("content", string(input)).Debugf("File content")
	return r.splitFileInputToURLS(string(input))
}

func (r *CheckCmd) readStdin() []string {
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		panic(errors.Wrap(err, "error while reading from stdin"))
	}
	return r.splitFileInputToURLS(input)
}

func (r *CheckCmd) splitFileInputToURLS(input string) []string {
	urls := make([]string, 0)
	for _, url := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		urls = append(urls, url)
	}
	return urls
}
