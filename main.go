package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// LogMessage represents a log message with a message string and a severity level.
type LogMessage struct {
	Message string
	Level   LogLevel
}

// Olympus represents a client for the Olympus server daemon.
type Olympus struct {
	conn     net.Conn
	exitChan chan struct{}
}

// New creates a new instance of the Olympus client.
func New() *Olympus {
	return &Olympus{
		exitChan: make(chan struct{}),
	}
}

// Start establishes a connection to the Olympus server daemon and starts a goroutine
// to handle the connection and listen for the exit signal.
func (o *Olympus) Start() error {
	conn, err := net.Dial("unix", "/tmp/unixsock")
	if err != nil {
		return err
	}

	o.conn = conn

	go func() {
		defer o.conn.Close()

		for {
			select {
			case <-o.exitChan:
				return
			default:
				// Handle other cases, if any
			}
		}
	}()

	return nil
}

// Stop sends a signal to gracefully stop the Olympus client instance.
func (o *Olympus) Stop() {
	close(o.exitChan)
}

// Log sends a log message to the Olympus server daemon.
func (o *Olympus) Log(message LogMessage) error {
	if o.conn == nil {
		return fmt.Errorf("connection to Olympus server is not established")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal log message: %v", err)
	}

	_, err = o.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send log message: %v", err)
	}

	return nil
}

func main() {
	olympus := New()

	err := olympus.Start()
	if err != nil {
		log.Fatalf("Failed to start Olympus client: %v", err)
	}

	message := LogMessage{
		Message: "Hello, Olympus!",
		Level:   INFO,
	}

	err = olympus.Log(message)
	if err != nil {
		log.Printf("Error logging message: %v", err)
	}

	olympus.Stop()
}
