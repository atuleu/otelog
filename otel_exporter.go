package otelog

import (
	"context"

	"github.com/atuleu/otelog/internal/utils"
	collector "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	common "go.opentelemetry.io/proto/otlp/common/v1"
	logs "go.opentelemetry.io/proto/otlp/logs/v1"
	resource "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/grpc"
)

type otelExporter struct {
	logClient collector.LogsServiceClient

	processor LogProcessor
	scope     *common.InstrumentationScope
	resource  *resource.Resource
}

func (e *otelExporter) Export(record *logs.LogRecord) {
	e.processor.batch(record, e.sendBatch)
}

func (e *otelExporter) sendBatch(records []*logs.LogRecord) {
	e.logClient.Export(context.Background(),
		&collector.ExportLogsServiceRequest{
			ResourceLogs: []*logs.ResourceLogs{
				{
					Resource: e.resource,
					ScopeLogs: []*logs.ScopeLogs{
						{
							Scope:      e.scope,
							LogRecords: records,
						},
					},
				},
			},
		})
}

func buildScope(opts logExporterOptions) *common.InstrumentationScope {
	return &common.InstrumentationScope{
		Name:    opts.scope.Name,
		Version: opts.scope.Version,
	}
}

func buildResource(opts logExporterOptions) *resource.Resource {
	if opts.resource == nil {
		return nil
	}
	return &resource.Resource{
		Attributes: utils.KeyValues(opts.resource.Attributes()),
	}

}

// Creates a new LogExporter that will export LogRecord to the
// specified endpoint. The endpoint address must be specified with
// WithEndpoint(). Credential must be specified with either
// WithInsecure() or WithTLSCredential().
func NewLogExporter(options ...LogExporterOption) (LogExporter, error) {
	opts := newOtelLogExporterOptions(options...)

	if opts.conn == nil {
		var err error
		opts.conn, err = grpc.Dial(opts.endpoint,
			grpc.WithTransportCredentials(opts.credential),
		)
		if err != nil {
			return nil, err
		}
	}

	client := collector.NewLogsServiceClient(opts.conn)

	return &otelExporter{
		logClient: client,
		resource:  buildResource(opts),
		scope:     buildScope(opts),
		processor: opts.processor,
	}, nil

}
