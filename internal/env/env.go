package env

import (
	"fmt"
	"github.com/mingalevme/avito/internal/notifier"
	"github.com/mingalevme/avito/internal/parser"
	"github.com/mingalevme/avito/internal/repository"
	"github.com/mingalevme/avito/internal/service"
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
	repository         repository.Repository
	notifier           notifier.Notifier
	parser             *parser.Parser
	htmlDocumentGetter parser.HTMLDocumentGetter
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
		Prefix: GetEnv("MINGALEVME_AVITO_GOLOGGER_ENV_NAMESAPCE", e.namespace),
		Env:    e.env,
	}.Create()
	if err != nil {
		panic(err)
	} else {
		e.logger = logger
	}
	return e.logger
}

func (e *Env) Parser() *parser.Parser {
	if e.parser != nil {
		return e.parser
	}
	e.parser = &parser.Parser{
		HTMLDocumentGetter: e.HTMLDocumentGetter(),
		Logger:             e.Logger(),
	}
	return e.parser
}

func (e *Env) HTMLDocumentGetter() parser.HTMLDocumentGetter {
	if e.htmlDocumentGetter != nil {
		return e.htmlDocumentGetter
	}
	e.htmlDocumentGetter = parser.NetHTMLDocumentGetter{
		Logger: e.Logger(),
	}
	return e.htmlDocumentGetter
}

func (e *Env) Repository() repository.Repository {
	if e.repository != nil {
		return e.repository
	}
	driver := e.getEnv("PERSISTENCE_DRIVER", "file")
	e.Logger().Debugf("Persistence: using driver: " + driver)
	if driver == "file" {
		e.repository = e.newFileRepository()
		return e.repository
	} else if driver == "in-memory" {
		e.repository = repository.NewInMemoryRepository()
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

func (e *Env) newFileRepository() *repository.FileRepository {
	f := e.getFileRepositoryFilename()
	if r, err := repository.NewFileRepository(f, e.Logger()); err != nil {
		panic(err)
	} else {
		return r
	}
}

func (e *Env) Notifier() notifier.Notifier {
	if e.notifier != nil {
		return e.notifier
	}
	e.notifier = e.newNotifierChannel(e.getEnv("NOTIFIER_CHANNEL", "stack"))
	return e.notifier
}

func (e *Env) newNotifierChannel(channel string) notifier.Notifier {
	switch channel {
	case "stack":
		return e.newStackNotifier()
	case "telegram":
		token := e.requireEnv("NOTIFIER_TELEGRAM_TOKEN")
		chatID := e.requireEnv("NOTIFIER_TELEGRAM_CHAT_ID")
		return notifier.NewTelegramNotifier(token, chatID, e.Logger())
	case "slack":
		webhookURL := e.requireEnv("NOTIFIER_SLACK_WEBHOOK_URL")
		return notifier.NewSlackNotifier(http.DefaultClient, webhookURL, e.Logger())
	case "stdout":
		return notifier.NewStdoutNotifier()
	case "logger":
		if level, err := gologger.ParseLevel(e.getEnv("NOTIFIER_LOGGER_LEVEL", "info")); err != nil {
			panic(err)
		} else {
			return notifier.NewLoggerNotifier(e.Logger(), level)
		}
	case "null":
		return notifier.NewNullNotifier()
	default:
		panic(errors.Errorf("unsupported notifier channel: %s", channel))
	}
}

func (e *Env) newStackNotifier() *notifier.StackNotifier {
	stack := notifier.NewStackNotifier(e.Logger())
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
