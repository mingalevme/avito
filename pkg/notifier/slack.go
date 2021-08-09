package notifier

import (
	"context"
	"fmt"
	"github.com/mingalevme/avito/pkg/model"
	"github.com/mingalevme/gologger"
	"github.com/slack-go/slack"
	"net/http"
)

type SlackNotifier struct {
	HTTPClient *http.Client
	WebhookUrl string
	Logger     gologger.Logger
}

func NewSlackNotifier(httpClient *http.Client, webhookUrl string, logger gologger.Logger) *SlackNotifier {
	return &SlackNotifier{
		HTTPClient: httpClient,
		WebhookUrl: webhookUrl,
		Logger:     logger,
	}
}

func (n *SlackNotifier) Notify(item model.Item) error {
	n.Logger.WithField("item-id", item.ID).Debugf("[SlackNotifier] Notifying is starting ...")
	msg := &slack.WebhookMessage{
		Text:        fmt.Sprintf("%s\n%s\n%s", item.Title, item.GetFormattedPrice(), item.URL),
	}
	err := slack.PostWebhookCustomHTTPContext(context.Background(), n.WebhookUrl, n.HTTPClient, msg)
	if err != nil {
		n.Logger.WithField("item-id", item.ID).WithError(err).Errorf("[SlackNotifier] Error while notifying")
		return err
	}
	n.Logger.WithField("item-id", item.ID).Debugf("[SlackNotifier] Notifying has been finished")
	return nil
}
