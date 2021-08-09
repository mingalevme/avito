# Avito

Simple Avito new search items notifier

## Examples

### Docker

```shell
docker build --target avito -t avito .
docker run --rm \
  -v "$(pwd):/var/lib/avito" \
  -e "APP_PERSISTENCE_FILE_FILENAME=/var/lib/avito/avito.json" \
  -e "APP_LOG_STDOUT_LEVEL=error" \
  avito check \
    "https://www.avito.ru/rossiya/bytovaya_elektronika?q=iphone+11" \
    "https://www.avito.ru/rossiya/telefony?q=iphone+12"
```

### Source 

```shell
go build .
APP_PERSISTENCE_DRIVER=file \
APP_PERSISTENCE_FILE_FILENAME=/tmp/avito.json \
APP_LOG_STDOUT_LEVEL=error \
APP_NOTIFIER_CHANNEL=telegram \
APP_NOTIFIER_TELEGRAM_TOKEN="MY_TELEGRAM_BOT_TOKEN" \
APP_NOTIFIER_TELEGRAM_CHAT_ID="MY_TELEGRAM_CHAT_ID" \
  ./avito check \
    "https://www.avito.ru/rossiya/bytovaya_elektronika?q=iphone+11" \
    "https://www.avito.ru/rossiya/telefony?q=iphone+12"
```
| INFO |
| :--- |
| **APP_** prefix in env vars name can be changed via **MINGALEVME_AVITO_ENV_NAMESPACE** env var. Log env vars name prefix can be changed via **MINGALEVME_AVITO_GOLOGGER_ENV_NAMESPACE** |

| INFO |
| :--- |
| It is recommended to do the first run with `APP_NOTIFIER_CHANNEL=null` |

## Repository drivers
Default is **file** (**$HOME/avito.json**).

### File
PERSISTENCE_DRIVER=file
PERSISTENCE_FILE_FILENAME=/tmp/avito.json

### ImMemory
PERSISTENCE_DRIVER=in-memory

## Notifier drivers

Default is **stdout**.

### Stack
APP_NOTIFIER_CHANNEL=stack
APP_NOTIFIER_CHANNELS=stdout,telegram,slack # default "stdout"

### Telegram
APP_NOTIFIER_CHANNEL=telegram
APP_NOTIFIER_TELEGRAM_TOKEN="XXX:YYY"
APP_NOTIFIER_TELEGRAM_CHAT_ID="ZZZ"

### Slack
APP_NOTIFIER_CHANNEL=slack
APP_NOTIFIER_SLACK_WEBHOOK_URL="https://hooks.slack.com/services/XXX/YYY/ZZZ"

### Stdout
APP_NOTIFIER_CHANNEL=stdout

### Logger
APP_NOTIFIER_CHANNEL=logger
APP_NOTIFIER_LOGGER_LEVEL=info # default "info"

### Null / NoOp
APP_NOTIFIER_CHANNEL=null

## Log channels

Default is **stack**.

### Stack
APP_LOG_CHANNEL=stack
APP_LOG_CHANNELS=stdout,sentry,rollbar # default "stdout"

### Sentry
APP_LOG_CHANNEL=sentry # Or APP_SENTRY_DSN
APP_LOG_SENTRY_LEVEL="info" # default "warning"
APP_LOG_SENTRY_DSN="https://XXX@YYY.ingest.sentry.io/ZZZ"
APP_LOG_SENTRY_DEBUG="1" # default "" (an empty string -> i.e. false)
APP_LOG_SENTRY_ENV="debug" # default "production"

### Rollbar
APP_LOG_CHANNEL=rollbar
APP_LOG_ROLLBAR_LEVEL="error" # default "warning"
APP_LOG_ROLLBAR_TOKEN="XXX" # Or APP_ROLLBAR_TOKEN
LOG_ROLLBAR_ENV="debug" # default "production"

### Stdout
APP_LOG_CHANNEL=stdout

### Stderr
APP_LOG_CHANNEL=stderr

### Array (InMemory)
APP_LOG_CHANNEL=array

### null (NoOp)
APP_LOG_CHANNEL=null
