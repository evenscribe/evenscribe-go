package logs_test

import (
	"testing"
	"time"

	log "github.com/evenscribe/evenscribe-go"
)

func TestBuildLog(t *testing.T) {
	err := log.NewLogger(
		log.NewOptions().
			WithDefaults().
			WithResourceAttributes(map[string]string{"evnviroment": "testing"}).
			WithServiceName("logger_test").
			WithPrintToStdout(false),
	)
	if err != nil {
		t.Fatal(err)
	}

	for range 1000 {
		// log.InfoS("Hello world", "log-attributes", map[string]string{"hello": "world"}, "trace-id", 12333, "span-id", 11111)
		log.Info("Hello world")
	}
	time.Sleep(6 * time.Minute)
}
