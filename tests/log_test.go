package logs_test

import (
	log "github.com/evenscribe/evenscribe-go"
	"testing"
)

func BenchmarkLoggging(b *testing.B) {
	b.Run("Logging Efficiency Benchmark", func(b *testing.B) {
		err := log.NewLogger(
			log.NewOptions().
				WithDefaults().
				WithResourceAttributes(map[string]string{"evnviroment": "testing"}).
				WithServiceName("logger_test").
				WithPrintToStdout(false),
		)

		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.InfoS("Hello world", "log-attributes", map[string]string{"hello": "world"}, "trace-id", 12333, "span-id", 11111)
			}

		})

	})
}
