package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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
	SeverityNumber     int
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
	if o.connection == nil {
		return fmt.Errorf("[%s] connection isn't not established : %v", APP_NAME, err)
	}
	err = o.SendLog(message)
	if err == nil {
		return nil
	}
	return fmt.Errorf("[%s] retry limit exceeded, could not save log message", APP_NAME)
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

func (o *EvenscribeConnection) SendLog(message Log) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("[%s] failed to parse log message as json: %v", APP_NAME, err)
	}
	data = pad(data)
	_, err = o.connection.Write(data)
	if err != nil {
		return fmt.Errorf("[%s] failed to parse log message as json: %v", APP_NAME, err)
	}
	return nil
}

func main() {
	logEntry := Log{
		Timestamp:          time.Now().Unix(),
		TraceId:            "trace-id-123",
		SpanId:             "span-id-456",
		TraceFlags:         1,
		SeverityText:       "ERROR",
		SeverityNumber:     (rand.Intn(5) + 1),
		ServiceName:        "example-service",
		Body:               "This is a log message",
		ResourceAttributes: map[string]string{"env": "production", "version": "1.0.0"},
		LogAttributes:      map[string]string{"user_id": "12345", "operation": "create"},
	}
	ladder := [...]int{1_000, 10_000, 100_000, 1_000_000}
	olympus := New(ConnectionOptions{wait: false, retry: 0})
	olympus.Connect()
	for _, v := range ladder {
		start := time.Now()
		err_count := 0
		for i := 0; i < v; i++ {
			err := olympus.Log(logEntry)
			if err != nil {
				err_count++
			}
		}
		elapsed := time.Since(start)
		fmt.Printf("%d took %s with %d\n", v, elapsed, err_count)
	}
}
