package notifier

import (
	"fmt"
	"github.com/mingalevme/avito/internal/model"
	"github.com/mingalevme/gologger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TelegramNotifier struct {
	Token  string
	ChatID string
	Logger gologger.Logger
}

func NewTelegramNotifier(token string, chatID string, logger gologger.Logger) *TelegramNotifier {
	return &TelegramNotifier{
		Token:  token,
		ChatID: chatID,
		Logger: logger,
	}
}

func (n *TelegramNotifier) Notify(item model.Item) error {
	n.Logger.WithField("item-id", item.ID).Debugf("[TelegramNotifier] Notifying is starting ...")
	message := fmt.Sprintf("%s\n%s\n%s\n%s", item.Title, item.GetFormattedPrice(), item.URL, item.ImageURL)
	u := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("bot%s/sendMessage", n.Token),
	}
	q := u.Query()
	q.Add("chat_id", n.ChatID)
	q.Add("text", message)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		n.Logger.WithField("item-id", item.ID).WithError(err).Errorf("[TelegramNotifier] Error while notifying")
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			n.Logger.WithField("item-id", item.ID).
				WithError(err).
				Errorf("[TelegramNotifier] error while reading response body of unsuccessful response")
		}
		n.Logger.WithField("item-id", item.ID).
			WithField("response-status", resp.StatusCode).
			WithField("response-body", string(body)).
			Errorf("[TelegramNotifier] Error while notifying")
		return err
	}
	n.Logger.WithField("item-id", item.ID).Debugf("[TelegramNotifier] Notifying has been finished")
	return nil
}
