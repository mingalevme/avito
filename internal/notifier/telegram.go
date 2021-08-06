package notifier

import (
	"fmt"
	"github.com/mingalevme/avito/internal/model"
	log "github.com/mingalevme/gologger"
	"io"
	"net/http"
	"net/url"
)

type TelegramNotifier struct {
	APIKey string
	ChatID string
	Logger log.Logger
}

func (n *TelegramNotifier) Notify(item model.Item) error {
	message := fmt.Sprintf("%s\n%s\n%s\n%s", item.Title, item.GetFormattedPrice(), item.URL, item.ImageURL)
	u := url.URL{
		Scheme: "https",
		Host:   "api.telegram.or",
		Path:   fmt.Sprintf("bot%s/sendMessage", n.APIKey),
	}
	q := u.Query()
	q.Add("chat_id", n.ChatID)
	q.Add("text", message)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != 200 {
		return err
	}
	return nil
}
