package evenscribe

import (
	"log"
	"net"
)

var loggerInstance *Logger

type Logger struct {
	conn net.Conn
	*Options
}

func New(options *Options) error {
	loggerInstance = &Logger{Options: options}
	socket, err := net.Dial("unix", loggerInstance.SocketAddr)
	if err != nil {
		return err
	}
	loggerInstance.conn = socket
	return nil
}

func NewWithDefaultOptions() error {
	loggerInstance = &Logger{Options: NewOptions().WithDefaults()}
	socket, err := net.Dial("unix", loggerInstance.SocketAddr)
	if err != nil {
		return err
	}
	loggerInstance.conn = socket
	return nil
}

func Print(v ...any) {
	go func() {
		// Send the log message to the socket
	}()
	if loggerInstance.PrintToStdout {
		log.Print(v...)
	}
}

func Println(v ...any) {
	go func() {
		// Send the log message to the socket
	}()
	if loggerInstance.PrintToStdout {
		log.Println(v...)
	}
}
