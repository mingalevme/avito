package env

import (
	"fmt"
	"github.com/mingalevme/avito/internal/notifier"
	"github.com/mingalevme/avito/internal/parser"
	"github.com/mingalevme/avito/internal/repository"
	"github.com/mingalevme/avito/internal/service"
	"github.com/mingalevme/gologger"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type Env struct {
	prefix             string
	env                map[string]string
	logger             gologger.Logger
	repository         repository.Repository
	notifier           notifier.Notifier
	parser             *parser.Parser
	htmlDocumentGetter parser.HTMLDocumentGetter
	checker            *service.Checker
}

func New(prefix string, env map[string]string) *Env {
	clone := make(map[string]string, len(env))
	for k, v := range env {
		clone[k] = v
	}
	return &Env{
		prefix: prefix,
		env:    clone,
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
		Prefix: e.prefix,
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
	driver := e.getEnv("PERSISTENCE_DRIVER", "in-memory")
	e.Logger().Debugf("Persistence: using driver: " + driver)
	if driver == "file" {
		filename := e.requireEnv("PERSISTENCE_FILE_FILENAME")
		if r, err := repository.NewFileRepository(filename, e.Logger()); err != nil {
			panic(err)
		} else {
			e.repository = r
			return e.repository
		}
	} else if driver == "in-memory" {
		e.repository = repository.NewInMemoryRepository()
		return e.repository
	}
	panic(errors.New("unknown persistence driver: " + driver))
}

func (e *Env) Notifier() notifier.Notifier {
	if e.notifier != nil {
		return e.notifier
	}
	driver := e.getEnv("NOTIFIER_DRIVER", "stdout")
	switch driver {
	case "telegram":
		token := e.requireEnv("NOTIFIER_TELEGRAM_TOKEN")
		chatID := e.requireEnv("NOTIFIER_TELEGRAM_CHAT_ID")
		e.notifier = notifier.NewTelegramNotifier(token, chatID, e.Logger())
	case "stdout":
		e.notifier = notifier.NewStdoutNotifier()
	case "logger":
		if level, err := gologger.ParseLevel(e.getEnv("NOTIFIER_LOGGER_LEVEL", "info")); err != nil {
			panic(err)
		} else {
			e.notifier = notifier.NewLoggerNotifier(e.Logger(), level)
		}
	case "null":
		e.notifier = notifier.NewNullNotifier()
	default:
		panic(errors.Errorf("unsupported notifier driver: %s", driver))
	}
	return e.notifier
}

func (e *Env) getEnv(key string, def string) string {
	key = e.prefix + key
	if val, ok := e.env[key]; ok {
		return val
	} else {
		return def
	}
}

func (e *Env) requireEnv(key string) string {
	key = e.prefix + key
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
