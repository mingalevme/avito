package parser

import (
	log "github.com/mingalevme/gologger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	logger := log.NewStdoutLogger(log.LevelDebug)
	url := "https://www.avito.ru/rossiya?q=iPhone"
	p := Parser{
		HTMLDocumentGetter: NetHTMLDocumentGetter{
			Logger: logger,
		},
		Logger: logger,
	}
	items, err := p.Parse(url)
	assert.NoError(t, err)
	assert.Greater(t, len(items), 0)
}
