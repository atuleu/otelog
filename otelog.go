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
