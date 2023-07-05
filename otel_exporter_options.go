package otelog

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type logExporterOptions struct {
	conn       *grpc.ClientConn
	endpoint   string
	credential credentials.TransportCredentials

	resource  *resource.Resource
	scope     instrumentation.Scope
	processor LogProcessor
}

// LogExporterOption is an option to use with NewLogExporter().
type LogExporterOption interface {
	apply(*logExporterOptions)
}

type logExporterOptionFunc func(*logExporterOptions)

func (f logExporterOptionFunc) apply(opts *logExporterOptions) {
	f(opts)
}

// Sets the resource associated with the LogExporter
func WithResource(r *resource.Resource) LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.resource = r
	})
}

type scopeOption instrumentation.Scope

// Sets the scope associated with the LogExporter
func WithScope(s instrumentation.Scope) LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.scope = s
	})
}

// Sets a LogProcessor that sends every event immediatly. It
// should be avoided in production.
func WithSyncer() LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.processor = &syncProcessor{}
	})
}

// Sets a LogProcessor that batches logs before exporting them.
func WithBatchLogProcessor(options ...BatchLogProcessorOption) LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.processor = newBatchProcessor(options...)
	})
}

// Sets the Open Telemetry collector endpoint to export logs to.
func WithEndpoint(endpoint string) LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.endpoint = endpoint
	})
}

// Sets no credential for the OpenTelemetry endpoint.
func WithInsecure() LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.credential = insecure.NewCredentials()
	})
}

// Sets the gRPC TLS credential to use for the OpenTelemetry endpoint.
func WithTLSCredentials(c credentials.TransportCredentials) LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.credential = c
	})
}

func WithGRPCConn(conn *grpc.ClientConn) LogExporterOption {
	return logExporterOptionFunc(func(opts *logExporterOptions) {
		opts.conn = conn
	})
}

func newOtelLogExporterOptions(options ...LogExporterOption) logExporterOptions {
	opts := logExporterOptions{
		processor: newBatchProcessor(),
	}

	for _, o := range options {
		o.apply(&opts)
	}

	return opts
}
