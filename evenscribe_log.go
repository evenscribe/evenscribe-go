package evenscribe

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
