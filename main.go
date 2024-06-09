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
	ResourceAttributes map[string]string
	LogAttributes      map[string]string
	TraceId            string
	SpanId             string
	SeverityText       string
	ServiceName        string
	Body               string
	Timestamp          int64
	TraceFlags         uint32
	SeverityNumber     int32
}

// EvenscribeConnection represents a client for the EvenscribeConnection server daemon.
type EvenscribeConnection struct {
	connection        net.Conn
	exitChan          chan struct{}
	connectionOptions ConnectionOptions
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
		return fmt.Errorf("failed to connect to server socket; make sure the evenscribe server is running. : %v", err)
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
	e := o.Connect()
	if e != nil {
		println("Connect error;")
	}

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

	if string(ans) == "OK" || o.connectionOptions.retry == 0 {
		return nil
	}
	int := o.connectionOptions.retry
	for int > 0 {
		res, _ := o.Retry(message)
		if string(res) == "OK" {
			return nil
		}
		int--
	}
	return fmt.Errorf("retry limit exceeded, could not send log message")
}

// Retry sends a log message to the Olympus server daemon and returns the response.
func (o *EvenscribeConnection) Retry(message Log) (res []byte, err error) {
	e := o.Connect()
	if e != nil {
		println("Connect error;")
	}
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
