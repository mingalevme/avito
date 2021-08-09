package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/mingalevme/gologger"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

type HTMLDocumentGetter interface {
	Get(url string) (*goquery.Document, error)
}

type NetHTMLDocumentGetter struct {
	Logger log.Logger
}

func (g NetHTMLDocumentGetter) Get(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(errors.Wrap(err, "error while creating request object"))
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")
	req.Header.Add("Referer", "https://www.avito.ru/")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "ru,en;q=0.9,en-US;q=0.8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error while requesting remote content: %s", url))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			g.Logger.WithError(err).Errorf("error while closing HTTP Response body")
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Invalid response status: %d", resp.StatusCode))
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error while creating HTML document: %s", url))
	}
	return doc, nil
}

type StaticHTMLDocumentGetter struct {
	Body string
}

func (g StaticHTMLDocumentGetter) Get(url string) (*goquery.Document, error) {
	r := strings.NewReader(g.Body)
	return goquery.NewDocumentFromReader(r)
}
