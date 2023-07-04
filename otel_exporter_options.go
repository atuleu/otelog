package otelog

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type logExporterOptions struct {
	endpoint   string
	credential credentials.TransportCredentials

	resource  *resource.Resource
	scope     instrumentation.Scope
	processor LogProcessor
}

type LogExporterOption interface {
	apply(*logExporterOptions)
}

type resourceOption resource.Resource

func (o resourceOption) apply(opts *logExporterOptions) {
	opts.resource = (*resource.Resource)(o)
}

func WithResource(r *resource.Resource) LogExporterOption {
	return (resourceOption)(r)
}

type scopeOption instrumentation.Scope

func (o scopeOption) apply(opts *logExporterOptions) {
	opts.scope = (instrumentation.Scope)(o)
}

func WithScope(r instrumentation.Scope) LogExporterOption {
	return (scopeOption)(r)
}

type processorOption struct {
	p LogProcessor
}

func (o processorOption) apply(opts *logExporterOptions) {
	opts.processor = o.p
}

func WithSyncer() LogExporterOption {
	return processorOption{&syncProcessor{}}
}

func WithBatchLogProcessor(options ...BatchLogProcessorOption) LogExporterOption {
	return processorOption{newBatchProcessor(options...)}
}

type endpointOptions string

func (o endpointOptions) apply(opts *logExporterOptions) {
	opts.endpoint = (string)(o)
}

func WithEndpoint(endpoint string) LogExporterOption {
	return (endpointOptions)(endpoint)
}

type credentialOptions struct {
	c credentials.TransportCredentials
}

func (o credentialOptions) apply(opts *logExporterOptions) {
	opts.credential = o.c
}

func WithInsecure() LogExporterOption {
	return credentialOptions{insecure.NewCredentials()}
}

func WithTLSCredentials(c credentials.TransportCredentials) LogExporterOption {
	return credentialOptions{c}
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
