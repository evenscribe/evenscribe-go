package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const APP_NAME = "evenscribe"
const SOCKET_PATH = "/tmp/olympus_socket.sock"

type ConnectionOptions struct {
	wait  bool // should you wait for a response
	retry int8 // if 0 then no retry
}

// Log represents a log entry in the logs table.
type Log struct {
	Timestamp          int64
	TraceId            string
	SpanId             string
	TraceFlags         uint32
	SeverityText       string
	SeverityNumber     int32
	ServiceName        string
	Body               string
	ResourceAttributes map[string]string
	LogAttributes      map[string]string
}

// EvenscribeConnection represents a client for the EvenscribeConnection server daemon.
type EvenscribeConnection struct {
	connection        net.Conn
	connectionOptions ConnectionOptions
	exitChan          chan struct{}
}

// New creates a new instance of the Olympus client.
func New(connectionOptions ConnectionOptions) *EvenscribeConnection {
	return &EvenscribeConnection{
		connectionOptions: connectionOptions,
		exitChan:          make(chan struct{}),
	}
}

// Connect establishes a connection to the Olympus server daemon
// to handle the connection and listen for the exit signal.
func (o *EvenscribeConnection) Connect() error {
	conn, err := net.Dial("unix", SOCKET_PATH)
	if err != nil {
		return fmt.Errorf("[%s] failed to connect to server socket; make sure the evenscribe server is running. : %v", APP_NAME, err)
	}
	o.connection = conn
	return nil
}

// Stop sends a signal to gracefully stop the Olympus client instance.
func (o *EvenscribeConnection) Stop() {
	close(o.exitChan)
}

// Log sends a log message to the Olympus server daemon.
func (o *EvenscribeConnection) Log(message Log) (err error) {
	var attempt int8
	for attempt = 0; attempt <= o.connectionOptions.retry; attempt++ {
		err = o.Connect()
		if err != nil {
			return fmt.Errorf("[%s] connection couldn't not established : %v", APP_NAME, err)
		}
		if o.connection == nil {
			return fmt.Errorf("[%s] connection isn't not established : %v", APP_NAME, err)
		}
		err = o.SendLog(message)

		if err == nil {
			return nil
		}

		if attempt == 0 {
			fmt.Printf("[%s] Failed to save log message.", APP_NAME)
		} else {
			fmt.Printf("[%s] Retrying to save log failed (%d/%d).\n", APP_NAME, attempt, o.connectionOptions.retry)
		}
	}
	return fmt.Errorf("[%s] retry limit exceeded, could not save log message", APP_NAME)
}

func (o *EvenscribeConnection) SendLog(message Log) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("[%s] failed to parse log message as json: %v", APP_NAME, err)
	}

	_, err = o.connection.Write(data)
	if err != nil {
		return fmt.Errorf("[%s] failed to write to socket : %v", APP_NAME, err)
	}

	if !o.connectionOptions.wait {
		return nil
	}

	ans := make([]byte, 2)
	o.connection.Read(ans)

	if string(ans) == "OK" {
		return nil
	}

	return fmt.Errorf("[%s] failed to write log message: %v", APP_NAME, err)
}

func main() {
	olympus := New(ConnectionOptions{wait: true, retry: 3})
	logEntry := Log{
		Timestamp:          time.Now().Unix(),
		TraceId:            "trace-id-123",
		SpanId:             "span-id-456",
		TraceFlags:         1,
		SeverityText:       "ERROR",
		SeverityNumber:     3,
		ServiceName:        "example-service",
		Body:               "This is a log message",
		ResourceAttributes: map[string]string{"env": "production", "version": "1.0.0"},
		LogAttributes:      map[string]string{"user_id": "12345", "operation": "create"},
	}
	err := olympus.Log(logEntry)
	if err != nil {
		fmt.Println(err)
	}
}
