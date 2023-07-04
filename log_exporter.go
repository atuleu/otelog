package otelog

import (
	logs "go.opentelemetry.io/proto/otlp/logs/v1"
)

type LogRecord logs.LogRecord

// A LogExporter is used to export log for example, to an Open
// Telemetry collector.
type LogExporter interface {
	Export(log *LogRecord)
}

// Creates a Log exporter that exports nothing.
func NoopLogExporter() LogExporter {
	return &noopLogExporter{}
}

type noopLogExporter struct{}

func (p *noopLogExporter) Export(log *LogRecord) {}
