package evenscribe

import (
	"time"
)

type LogBuilder struct {
	resourceAttributes map[string]string
	logAttributes      map[string]string
	traceId            string
	spanId             string
	severityText       string
	serviceName        string
	body               string
	traceFlags         uint32
	severityNumber     int32
}

func NewLogBuilder() *LogBuilder {
	return &LogBuilder{}
}

func (l *LogBuilder) WithResourceAttributes(resourceAttributes map[string]string) *LogBuilder {
	l.resourceAttributes = resourceAttributes
	return l
}

func (l *LogBuilder) WithLogAttributes(logAttributes map[string]string) *LogBuilder {
	l.logAttributes = logAttributes
	return l
}

func (l *LogBuilder) WithTraceId(traceId string) *LogBuilder {
	l.traceId = traceId
	return l
}

func (l *LogBuilder) WithSpanId(spanId string) *LogBuilder {
	l.spanId = spanId
	return l
}

func (l *LogBuilder) WithSeverity(severity Severity) *LogBuilder {
	l.severityText, l.severityNumber = severity.SeverityTextAndNumber()
	return l
}

func (l *LogBuilder) WithServiceName(serviceName string) *LogBuilder {
	l.serviceName = serviceName
	return l
}

func (l *LogBuilder) WithBody(body string) *LogBuilder {
	l.body = body
	return l
}

func (l *LogBuilder) WithTraceFlags(traceFlags uint32) *LogBuilder {
	l.traceFlags = traceFlags
	return l
}

func (l *LogBuilder) WithArgs(args ...any) *LogBuilder {
	if len(args)%2 != 0 {
		return l
	}

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		value := args[i+1]
		switch key {
		case "log-attributes":
			if attrMap, ok := value.(map[string]string); ok {
				l.WithLogAttributes(attrMap)
			}
		case "trace-id":
			if traceID, ok := value.(string); ok {
				l.WithTraceId(traceID)
			}
		case "span-id":
			if spanID, ok := value.(string); ok {
				l.WithSpanId(spanID)
			}
		case "trace-flags":
			if traceFlags, ok := value.(uint32); ok {
				l.WithTraceFlags(traceFlags)
			}
		default:
			continue
		}
	}
	return l
}

func (l *LogBuilder) Build() *Log {
	return &Log{
		ResourceAttributes: l.resourceAttributes,
		LogAttributes:      l.logAttributes,
		TraceId:            l.traceId,
		SpanId:             l.spanId,
		SeverityText:       l.severityText,
		ServiceName:        l.serviceName,
		Body:               l.body,
		Timestamp:          time.Now().Unix(),
		TraceFlags:         l.traceFlags,
		SeverityNumber:     l.severityNumber,
	}
}
