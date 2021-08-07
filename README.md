# Avito

Avito new items notifier

```shell
go build .
APP_PERSISTENCE_DRIVER=file \
APP_PERSISTENCE_FILE_FILENAME=/tmp/avito.json \
APP_LOG_STDOUT_LEVEL=error \
APP_NOTIFIER_DRIVER=telegram \
APP_NOTIFIER_TELEGRAM_TOKEN="MY_TELEGRAM_BOT_TOKEN" \
APP_NOTIFIER_TELEGRAM_CHAT_ID="MY_TELEGRAM_CHAT_ID" \
  ./avito check \
    "https://www.avito.ru/rossiya/bytovaya_elektronika?q=iphone+11" \
    "https://www.avito.ru/rossiya/telefony?q=iphone+12" \
    "https"
```

Available notifier (production) drivers:
- Telegram
- Slack (TODO)
- See internal/env/env.go:Notifier() for details

Available log (production) channels:
- Sentry
- Rollbar
- Slack
- Telegram (TODO)
- See internal/env/env.go:Logger() for details
