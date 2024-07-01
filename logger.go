package evenscribe

import (
	"log/slog"
	"net"
	"sync"

	stdlog "log"

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

func NewLogger(options *Options) error {
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

// Info logs with Severity of INFO.
func Info(msg string) {
	write(msg, INFO)
}

// InfoS logs with Severity of INFO and addition log fields.
//
// Example Usage: log.InfoS("Hello world", "log-attributes" , map[string]string{"userID":"1xbaaaaa"}, "trace-id", 12333, "span-id", 11111)
//
// Supported Arguments:
//
// log-attributes: map[string]string, key value pairs for structured logs
//
// trace-id: string, trace id for the log
//
// span-id: string, span id for the log
func InfoS(msg string, args ...any) {
	writeS(msg, INFO, args...)
}

// Error logs with Severity of ERROR.
func Error(msg string) {
	write(msg, ERROR)
}

// ErrorS logs with Severity of ERROR and addition log fields.
//
// Example Usage: log.ErrorS("Hello world", "log-attributes" , map[string]string{"userID":"1xbaaaaa"}, "trace-id", 12333, "span-id", 11111)
//
// Supported Arguments:
//
// log-attributes: map[string]string, key value pairs for structured logs
//
// trace-id: string, trace id for the log
//
// span-id: string, span id for the log
func ErrorS(msg string, args ...any) {
	writeS(msg, ERROR, args...)
}

// Debug logs with Severity of DEBUG.
func Debug(msg string) {
	write(msg, DEBUG)
}

// DebugS logs with Severity of DEBUG and addition log fields.
//
// Example Usage: log.DebugS("Hello world", "log-attributes" , map[string]string{"userID":"1xbaaaaa"}, "trace-id", 12333, "span-id", 11111)
//
// Supported Arguments:
//
// log-attributes: map[string]string, key value pairs for structured logs
//
// trace-id: string, trace id for the log
//
// span-id: string, span id for the log
func DebugS(msg string, args ...any) {
	writeS(msg, DEBUG, args...)
}

// Warn logs with Severity of WARN.
func Warn(msg string) {
	write(msg, WARN)
}

// WarnS logs with Severity of WARN and addition log fields.
//
// Example Usage: log.WarnS("Hello world", "log-attributes" , map[string]string{"userID":"1xbaaaaa"}, "trace-id", 12333, "span-id", 11111)
//
// Supported Arguments:
//
// log-attributes: map[string]string, key value pairs for structured logs
//
// trace-id: string, trace id for the log
//
// span-id: string, span id for the log
func WarnS(msg string, args ...any) {
	writeS(msg, WARN, args...)
}

// Fatal logs with Severity of FATAL.
func Fatal(msg string) {
	write(msg, FATAL)
}

// FatalS logs with Severity of FATAL and addition log fields.
//
// Example Usage: log.FatalS("Hello world", "log-attributes" , map[string]string{"userID":"1xbaaaaa"}, "trace-id", 12333, "span-id", 11111)
//
// Supported Arguments:
//
// log-attributes: map[string]string, key value pairs for structured logs
//
// trace-id: string, trace id for the log
//
// span-id: string, span id for the log
func FatalS(msg string, args ...any) {
	writeS(msg, FATAL, args...)
}

func pad(input []byte) []byte {
	const targetSize = 2000
	inputSize := len(input)
	padded := make([]byte, targetSize)
	copy(padded, input)
	for i := inputSize; i < targetSize; i++ {
		padded[i] = ' '
	}
	return padded
}

func writeS(msg string, severity Severity, args ...any) {
	logger := GetLogger()
	logger.mx.Lock()
	defer logger.mx.Unlock()
	l := NewLogBuilder().
		WithBody(msg).
		WithSeverity(severity).
		WithResourceAttributes(logger.ResourceAttributes).
		WithServiceName(logger.ServiceName).
		WithArgs(args...).
		Build()
	go WriteLogToSocket(l)
	if logger.PrintToStdout {
		stdlog.Printf("%+v", l)
	}
}

func write(msg string, severity Severity) {
	logger := GetLogger()
	logger.mx.Lock()
	defer logger.mx.Unlock()

	log := NewLogBuilder().
		WithBody(msg).
		WithSeverity(severity).
		WithServiceName(logger.ServiceName).
		WithResourceAttributes(logger.ResourceAttributes).
		WithLogAttributes(map[string]string{}).
		Build()

	if logger.PrintToStdout {
		stdlog.Printf("%+v", log)
	}
	go WriteLogToSocket(log)

}

func WriteLogToSocket(log *Log) {
	logger := GetLogger()
	logger.mx.Lock()
	defer logger.mx.Unlock()
	stringifiedLog, err := json.Marshal(log)
	if err != nil {
		slog.Error("Failed to marshal log: %v", err)
	}
	_, err = logger.conn.Write(pad(stringifiedLog))
	if err != nil {
		slog.Error("Failed to write to socket: %v", err)
	}
}
