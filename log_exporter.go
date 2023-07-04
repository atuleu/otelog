package otelog

import (
	logs "go.opentelemetry.io/proto/otlp/logs/v1"
)

// A LogExporter is used to export log for example, to an Open
// Telemetry collector.
type LogExporter interface {
	Export(log *logs.LogRecord)
}

// Creates a Log exporter that exports nothing.
func NoopLogExporter() LogExporter {
	return &noopLogExporter{}
}

type noopLogExporter struct{}

func (p *noopLogExporter) Export(log *logs.LogRecord) {}
