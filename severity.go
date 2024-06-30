package evenscribe

// Available severity levels
//
// TRACE 1-4
// A fine-grained debugging event. Typically disabled in default configurations.
//
// DEBUG 5-8
// A debugging event.
//
// INFO 9-12
// An informational event. Indicates that an event happened.
//
// WARN 13-16
// A warning event. Not an error but is likely more important than an informational event.
//
// ERROR 17-20
// An error event. Something went wrong.
//
// FATAL 21-24
// A fatal error such as application or system crash.
type Severity int32

const (
	TRACE Severity = iota + 1
	TRACE2
	TRACE3
	TRACE4
	DEBUG
	DEBUG2
	DEBUG3
	DEBUG4
	INFO
	INFO2
	INFO3
	INFO4
	WARN
	WARN2
	WARN3
	WARN4
	ERROR
	ERROR2
	ERROR3
	ERROR4
	FATAL
	FATAL2
	FATAL3
	FATAL4
)

var severityTextMap = map[Severity]string{
	TRACE:  "TRACE",
	TRACE2: "TRACE",
	TRACE3: "TRACE",
	TRACE4: "TRACE",
	DEBUG:  "DEBUG",
	DEBUG2: "DEBUG",
	DEBUG3: "DEBUG",
	DEBUG4: "DEBUG",
	INFO:   "INFO",
	INFO2:  "INFO",
	INFO3:  "INFO",
	INFO4:  "INFO",
	WARN:   "WARN",
	WARN2:  "WARN",
	WARN3:  "WARN",
	WARN4:  "WARN",
	ERROR:  "ERROR",
	ERROR2: "ERROR",
	ERROR3: "ERROR",
	ERROR4: "ERROR",
	FATAL:  "FATAL",
	FATAL2: "FATAL",
	FATAL3: "FATAL",
	FATAL4: "FATAL",
}

func (s Severity) SeverityTextAndNumber() (string, int32) {
	return severityTextMap[s], int32(s)
}
