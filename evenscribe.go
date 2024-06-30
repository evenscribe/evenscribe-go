package evenscribe

import (
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/goccy/go-json"
)

var loggerInstance *Logger

type Logger struct {
	conn net.Conn
	mx   sync.Mutex
	*Options
}

func GetLogger() *Logger {
	return loggerInstance
}

func SetLogger(logger *Logger) {
	loggerInstance = logger
}

func New(options *Options) error {
	SetLogger(&Logger{Options: options})

	GetLogger().mx.Lock()
	defer GetLogger().mx.Unlock()

	socket, err := net.Dial("unix", GetLogger().SocketAddr)
	if err != nil {
		return err
	}
	GetLogger().conn = socket
	return nil
}

func NewWithDefaultOptions() error {
	SetLogger(&Logger{Options: NewOptions().WithDefaults()})

	GetLogger().mx.Lock()
	defer GetLogger().mx.Unlock()

	socket, err := net.Dial("unix", GetLogger().SocketAddr)
	if err != nil {
		return err
	}

	GetLogger().conn = socket
	return nil
}

func Info(msg string) {
	logger := GetLogger()

	logger.mx.Lock()
	defer logger.mx.Unlock()

	go WriteToSocket(msg, INFO)
	if logger.PrintToStdout {
		slog.Debug(msg)
	}
}

func Error(msg string) {
	logger := GetLogger()

	logger.mx.Lock()
	defer logger.mx.Unlock()

	go WriteToSocket(msg, ERROR)
	if logger.PrintToStdout {
		slog.Debug(msg)
	}
}

func Debug(msg string) {
	logger := GetLogger()

	logger.mx.Lock()
	defer logger.mx.Unlock()

	go WriteToSocket(msg, DEBUG)
	if logger.PrintToStdout {
		slog.Debug(msg)
	}
}

func Warn(msg string) {
	logger := GetLogger()

	logger.mx.Lock()
	defer logger.mx.Unlock()

	go WriteToSocket(msg, WARN)
	if logger.PrintToStdout {
		slog.Warn(msg)
	}
}

func Fatal(msg string) {
	logger := GetLogger()

	logger.mx.Lock()
	defer logger.mx.Unlock()

	go WriteToSocket(msg, FATAL)
	if logger.PrintToStdout {
		slog.Warn(msg)
	}
}

func WriteToSocket(msg string, severity Severity) {
	logger := GetLogger()

	logger.mx.Lock()
	defer logger.mx.Unlock()

	log := BuildLog(logger, severity, msg)
	stringifiedLog, err := json.Marshal(log)
	if err != nil {
		slog.Error("Failed to marshal log: %v", err)
	}
	_, err = GetLogger().conn.Write(stringifiedLog)
	if err != nil {
		slog.Error("Failed to write to socket: %v", err)
	}
	if severity == FATAL {
		os.Exit(1)
	}
}
