package evenscribe

import "time"

type Log struct {
	ResourceAttributes map[string]string `json:"resource_attributes"`
	LogAttributes      map[string]string `json:"log_attributes"`
	TraceId            string            `json:"trace_id"`
	SpanId             string            `json:"span_id"`
	SeverityText       string            `json:"severity_text"`
	ServiceName        string            `json:"service_name"`
	Body               string            `json:"body"`
	Timestamp          int64             `json:"timestamp"`
	TraceFlags         uint32            `json:"trace_flags"`
	SeverityNumber     int32             `json:"severity_number"`
}

func BuildLog(logger *Logger, severity Severity, message string) *Log {
	severityText, severityNumber := severity.SeverityTextAndNumber()
	return &Log{Body: message,
		SeverityText:       severityText,
		SeverityNumber:     severityNumber,
		ServiceName:        logger.ServiceName,
		Timestamp:          time.Now().UnixNano(),
		ResourceAttributes: logger.ResourceAttributes,
	}
}
