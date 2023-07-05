// Package otelog provides global function to register an Open
// Telemetry Collector endpoint to export log records to.
//
// This package is meant to be a temporarly replacement while waiting
// a more thoroughly implemented and tested package provided by
// https://pkg.go.dev/go.opentelemetry.io/ . Please use it in
// production only if you really needs it.
//
// # Basics
//
// You will need to SetLogExporter() with a NewLogExporter()
// indicating the Open Telemetry Collector Endpoint.
//
// Then you can use hooks in `github.com/atuleu/otelog/pkg/hooks` to
// integrate your logging library. Currently only
// `github.com/sirupsen/logrus` integration is provided.
package otelog

var globalExporter LogExporter = NoopLogExporter()

// SetLogExporter sets the global LogExporter to exporter. This method
// is not concurrently safe.
func SetLogExporter(exporter LogExporter) {
	globalExporter = exporter
}

// GetLogExporter gets the global LogExporter registered with
// SetLogExporter. If none was registered a NoopLogExporter will be
// returned. It is concurrently safe to call GetLogExporter from
// multiple go routines.
func GetLogExporter() LogExporter {
	return globalExporter
}
