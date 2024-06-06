package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const SOCKET_PATH = "/tmp/olympus_socket.sock"

type ConnectionOptions struct {
	wait  bool // should you wait for a response
	retry int8 // if 0 then no retry
}

// Log represents a log entry in the logs table.
type Log struct {
	Timestamp          string
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

// EvenscribeConnection represents a client for the Evenscribe server daemon.
type EvenscribeConnection struct {
	connection        net.Conn
	connectionOptions ConnectionOptions
	exitChan          chan struct{}
}

// New creates a new instance of the Evenscribe client.
func New(connectionOptions ConnectionOptions) *EvenscribeConnection {
	return &EvenscribeConnection{
		connectionOptions: connectionOptions,
		exitChan:          make(chan struct{}),
	}
}

// Start establishes a connection to the Evenscribe server daemon
func (o *EvenscribeConnection) Start() error {
	conn, err := net.Dial("unix", SOCKET_PATH)
	if err != nil {
		return fmt.Errorf("failed to connect to server socket; make sure the evenscribe server is running. : %v", err)
	}
	o.connection = conn
	return nil
}

// Stop sends a signal to gracefully stop the Evenscribe client instance.
func (o *EvenscribeConnection) Stop() {
	close(o.exitChan)
}

// Log sends a log message to the Evenscribe server daemon.
func (o *EvenscribeConnection) Log(message Log) (err error) {
	if o.connection == nil {
		return fmt.Errorf("connection couldn't not established : %v", err)
	}
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to parse log message as json: %v", err)
	}
	_, err = o.connection.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write log message: %v", err)
	}
	if !o.connectionOptions.wait {
		return nil
	}
	ans := make([]byte, 2)
	o.connection.Read(ans)
	if string(ans) == "OK" {
		return nil
	}
	if string(ans) == "NO" {
		if o.connectionOptions.retry > 0 {
			var int = o.connectionOptions.retry
			for int > 0 {
				res, _ := o.Retry(message)
				if string(res) == "OK" {
					return nil
				}
				int--
			}
			return fmt.Errorf("retry limit exceeded, could not send log message")
		}
	}
	return nil
}

// Retry sends a log message to the Evenscribe server daemon and returns the response.
func (o *EvenscribeConnection) Retry(message Log) (res []byte, err error) {
	if o.connection == nil {
		return res, fmt.Errorf("connection is not established")
	}
	data, err := json.Marshal(message)
	if err != nil {
		return res, err
	}
	_, err = o.connection.Write(data)
	if err != nil {
		return res, err
	}
	ans := make([]byte, 2)
	o.connection.Read(ans)
	return ans, nil
}
