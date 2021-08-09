package env

import (
	"fmt"
	notifier2 "github.com/mingalevme/avito/pkg/notifier"
	parser2 "github.com/mingalevme/avito/pkg/parser"
	repository2 "github.com/mingalevme/avito/pkg/repository"
	"github.com/mingalevme/avito/pkg/service"
	"github.com/mingalevme/gologger"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
)

type Env struct {
	namespace          string
	env                map[string]string
	logger             gologger.Logger
	repository         repository2.Repository
	notifier           notifier2.Notifier
	parser             *parser2.Parser
	htmlDocumentGetter parser2.HTMLDocumentGetter
	checker            *service.Checker
}

func New(namespace string, env map[string]string) *Env {
	clone := make(map[string]string, len(env))
	for k, v := range env {
		clone[k] = v
	}
	return &Env{
		namespace: namespace,
		env:       clone,
	}
}

func (e *Env) Checker() *service.Checker {
	if e.checker != nil {
		return e.checker
	}
	e.checker = service.NewChecker(e.Parser(), e.Repository(), e.Notifier(), e.Logger())
	return e.checker
}

func (e *Env) Logger() gologger.Logger {
	if e.logger != nil {
		return e.logger
	}
	logger, err := gologger.Creator{
		Prefix: GetEnv("MINGALEVME_AVITO_GOLOGGER_ENV_NAMESPACE", e.namespace),
		Env:    e.env,
	}.Create()
	if err != nil {
		panic(err)
	} else {
		e.logger = logger
	}
	return e.logger
}

func (e *Env) Parser() *parser2.Parser {
	if e.parser != nil {
		return e.parser
	}
	e.parser = &parser2.Parser{
		HTMLDocumentGetter: e.HTMLDocumentGetter(),
		Logger:             e.Logger(),
	}
	return e.parser
}

func (e *Env) HTMLDocumentGetter() parser2.HTMLDocumentGetter {
	if e.htmlDocumentGetter != nil {
		return e.htmlDocumentGetter
	}
	e.htmlDocumentGetter = parser2.NetHTMLDocumentGetter{
		Logger: e.Logger(),
	}
	return e.htmlDocumentGetter
}

func (e *Env) Repository() repository2.Repository {
	if e.repository != nil {
		return e.repository
	}
	driver := e.getEnv("PERSISTENCE_DRIVER", "file")
	e.Logger().Debugf("Persistence: using driver: " + driver)
	if driver == "file" {
		e.repository = e.newFileRepository()
		return e.repository
	} else if driver == "in-memory" {
		e.repository = repository2.NewInMemoryRepository()
		return e.repository
	}
	panic(errors.New("unknown persistence driver: " + driver))
}

func (e *Env) getFileRepositoryFilename() string {
	filename := e.getEnv("PERSISTENCE_FILE_FILENAME", "")
	if filename != "" {
		return filename
	}
	if homeDir, err := os.UserHomeDir(); err != nil {
		panic(errors.Wrap(err, "error while resolving home dir"))
	} else {
		return fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), "avito.json")
	}
}

func (e *Env) newFileRepository() *repository2.FileRepository {
	f := e.getFileRepositoryFilename()
	e.Logger().Debugf("Persistence: using file: " + f)
	if r, err := repository2.NewFileRepository(f, e.Logger()); err != nil {
		panic(err)
	} else {
		return r
	}
}

func (e *Env) Notifier() notifier2.Notifier {
	if e.notifier != nil {
		return e.notifier
	}
	e.notifier = e.newNotifierChannel(e.getEnv("NOTIFIER_CHANNEL", "stack"))
	return e.notifier
}

func (e *Env) newNotifierChannel(channel string) notifier2.Notifier {
	switch channel {
	case "stack":
		return e.newStackNotifier()
	case "telegram":
		token := e.requireEnv("NOTIFIER_TELEGRAM_TOKEN")
		chatID := e.requireEnv("NOTIFIER_TELEGRAM_CHAT_ID")
		return notifier2.NewTelegramNotifier(token, chatID, e.Logger())
	case "slack":
		webhookURL := e.requireEnv("NOTIFIER_SLACK_WEBHOOK_URL")
		return notifier2.NewSlackNotifier(http.DefaultClient, webhookURL, e.Logger())
	case "stdout":
		return notifier2.NewStdoutNotifier()
	case "logger":
		if level, err := gologger.ParseLevel(e.getEnv("NOTIFIER_LOGGER_LEVEL", "info")); err != nil {
			panic(err)
		} else {
			return notifier2.NewLoggerNotifier(e.Logger(), level)
		}
	case "null":
		return notifier2.NewNullNotifier()
	default:
		panic(errors.Errorf("unsupported notifier channel: %s", channel))
	}
}

func (e *Env) newStackNotifier() *notifier2.StackNotifier {
	stack := notifier2.NewStackNotifier(e.Logger())
	channels := strings.Split(e.getEnv("NOTIFIER_STACK_CHANNELS", "stdout"), ",")
	for _, channel := range channels {
		channel = strings.TrimSpace(channel)
		if channel == "" {
			continue
		}
		if channel == "stack" {
			panic(errors.Errorf("stack channel recursion"))
		}
		stack.AddNotifier(e.newNotifierChannel(channel))
	}
	return stack
}

func (e *Env) getEnv(key string, def string) string {
	key = e.namespace + key
	if val, ok := e.env[key]; ok {
		return val
	} else {
		return def
	}
}

func (e *Env) requireEnv(key string) string {
	key = e.namespace + key
	if val, ok := e.env[key]; ok {
		return val
	} else {
		panic(fmt.Errorf("env-var %s does not set", key))
	}
}

func GetOSEnvMap() map[string]string {
	m := map[string]string{}
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		m[variable[0]] = variable[1]
	}
	return m
}

func GetEnv(key string, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
