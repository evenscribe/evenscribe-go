package evenscribe

// Options to build the Evenscribe client
type Options struct {
	// Number of retries, Default: 5
	Retry int
	//Max execution time, Default: 1000ms
	MaxExecutionTime int
	//Address of the unix socket, Default: /tmp/evenscribe.sock
	SocketAddr string
	// Name of the service, *Required*
	//
	// Example: "acme-auth-service"
	ServiceName string
	// Resource Attributes, Default: nil
	//
	// Example:
	// []{
	// 	"environment": "production",
	// 	"region":      "us-west-1",
	// }
	ResourceAttributes map[string]string
	// Print to stdout,
	//
	// Default: true
	PrintToStdout bool
}

func NewOptions() *Options {
	return &Options{}
}

// Build Options with defaults
func (o *Options) WithDefaults() *Options {
	o.Retry = 5
	o.MaxExecutionTime = 1000
	o.SocketAddr = "/tmp/olympus_socket.sock"
	o.ServiceName = ""
	o.PrintToStdout = true
	return o
}

func (o *Options) WithRetry(retry int) *Options {
	o.Retry = retry
	return o
}

func (o *Options) WithMaxExecutionTime(maxExecutionTime int) *Options {
	o.MaxExecutionTime = maxExecutionTime
	return o
}

func (o *Options) WithSocketAddr(socketAddr string) *Options {
	o.SocketAddr = socketAddr
	return o
}

func (o *Options) WithServiceName(serviceName string) *Options {
	o.ServiceName = serviceName
	return o
}

func (o *Options) WithResourceAttributes(resourceAttributes map[string]string) *Options {
	o.ResourceAttributes = resourceAttributes
	return o
}

func (o *Options) WithPrintToStdout(printToStdout bool) *Options {
	o.PrintToStdout = printToStdout
	return o
}
