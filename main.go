package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

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

// Olympus represents a client for the Olympus server daemon.
type Olympus struct {
	conn     net.Conn
	exitChan chan struct{}
}

var wg sync.WaitGroup

// New creates a new instance of the Olympus client.
func New() *Olympus {
	return &Olympus{
		exitChan: make(chan struct{}),
	}
}

// Start establishes a connection to the Olympus server daemon and starts a goroutine
// to handle the connection and listen for the exit signal.
func (o *Olympus) Start() error {
	conn, err := net.Dial("unix", "/tmp/olympus_socket.sock")
	if err != nil {
		return err
	}

	o.conn = conn

	return nil
}

// Stop sends a signal to gracefully stop the Olympus client instance.
func (o *Olympus) Stop() {
	close(o.exitChan)
}

// Log sends a log message to the Olympus server daemon.
func (o *Olympus) Log(message Log) error {
	defer wg.Done()
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
	answer := make([]byte, 2)
	o.conn.Read(answer)

	return nil
}

func RunParallel(n int) {
	olympus := New()
	err := olympus.Start()
	if err != nil {
		log.Fatalf("Failed to start Olympus client: %v", err)
	}

	logEntry := Log{
		Timestamp:          strconv.FormatInt(time.Now().Unix(), 10),
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

	for i := 0; i < n; i++ {
		wg.Add(1)
		go olympus.Log(logEntry)
		wg.Wait()
	}

	olympus.Stop()
}

func main() {
	arr := []int{1}

	number_of_log := 1
	for _, v := range arr {
		start := time.Now().UnixMilli()
		for i := 0; i < v; i++ {
			RunParallel(number_of_log)
		}
		elapsed := time.Now().UnixMilli() - start
		fmt.Printf("It took %d ms for %d client to run %d querry\n", elapsed, v, number_of_log)
	}
}
