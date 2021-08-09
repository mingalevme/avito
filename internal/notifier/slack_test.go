package notifier

import (
	"github.com/mingalevme/avito/internal/model"
	"github.com/mingalevme/gologger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	webhookURL, ok := os.LookupEnv("MINGALEVME_AVITO_TEST_SLACK_WEBHOOK_URL")
	if !ok {
		t.Skip("Environment variable is not set")
	}
	n := NewSlackNotifier(http.DefaultClient, webhookURL, gologger.NewNullLogger())
	item := model.Item{
		ID:       1,
		URL:      "https://www.avito.ru/",
		Title:    "iPhone X",
		Price:    999999,
		ImageURL: "https://github.com/mingalevme/avito/blob/master/avito.png?raw=true",
	}
	err := n.Notify(item)
	assert.NoError(t, err)
}
