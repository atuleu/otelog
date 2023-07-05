// Package hooks provide hooks and integration for otelgo to various
// logging library.
//
// Currently only `github.com/sirupsen/logrus` is supported. See
// NewLogrusHook().
package hooks

import (
	"sort"

	"github.com/atuleu/otelog"
	"github.com/atuleu/otelog/internal/utils"
	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel/trace"
	common "go.opentelemetry.io/proto/otlp/common/v1"
	logs "go.opentelemetry.io/proto/otlp/logs/v1"
)

type logrusHook struct {
	exporter otelog.LogExporter
	levels   []logrus.Level
}

func (l *logrusHook) Levels() []logrus.Level {
	return l.levels
}

func (l *logrusHook) Fire(entry *logrus.Entry) error {
	l.exporter.Export(reportFromLogrus(entry))
	return nil
}

// NewLogrusHook creates a new logrus.Hook that will export all
// logrus.Entry to the global otelog registered LogExporter. By
// default no levels are enabled, and must be set with FromLogrusLevel
// or WithLogrusLevels.
//
// If a context.Context containing a valid otel.SpanContext is
// provided to the logrus.Entry, the exported LogRecord will be
// automatically linked with the span.
func NewLogrusHook(options ...LogrusOption) logrus.Hook {
	opts := newLogrusOptions(options...)

	return &logrusHook{
		levels:   ([]logrus.Level)(opts),
		exporter: otelog.GetLogExporter(),
	}
}

func spanIDToSlice(s trace.SpanID) []byte {
	return s[:]
}

func traceIDToSlice(t trace.TraceID) []byte {
	return t[:]
}

func reportFromLogrus(entry *logrus.Entry) *logs.LogRecord {
	record := &logs.LogRecord{
		TimeUnixNano:         uint64(entry.Time.UnixNano()),
		ObservedTimeUnixNano: uint64(entry.Time.UnixNano()),
		Body: &common.AnyValue{
			Value: &common.AnyValue_StringValue{
				StringValue: entry.Message,
			},
		},
		Attributes: mapLogrusFields(entry.Data),
	}
	record.SeverityNumber, record.SeverityText = mapLogrusSeverity(entry.Level)

	spanContext := trace.SpanContextFromContext(entry.Context)
	if spanContext.IsValid() == true {
		record.SpanId = spanIDToSlice(spanContext.SpanID())
		record.TraceId = traceIDToSlice(spanContext.TraceID())
		record.Flags = uint32(spanContext.TraceFlags() & 0xff)
	}

	return record
}

func mapLogrusFields(fields map[string]interface{}) []*common.KeyValue {
	if len(fields) == 0 {
		return nil
	}
	res := make([]*common.KeyValue, 0, len(fields))

	for k, v := range fields {
		res = append(res, &common.KeyValue{Key: k, Value: utils.ValueFromGo(v)})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Key < res[j].Key
	})

	return res

}

var logrusSeverities = map[logrus.Level]struct {
	value logs.SeverityNumber
	name  string
}{
	logrus.TraceLevel: {1, "TRACE"},
	logrus.DebugLevel: {5, "DEBUG"},
	logrus.InfoLevel:  {9, "INFO"},
	logrus.WarnLevel:  {13, "WARN"},
	logrus.ErrorLevel: {17, "ERROR"},
	logrus.FatalLevel: {21, "FATAL"},
	logrus.PanicLevel: {23, "FATAL3"},
}

func mapLogrusSeverity(level logrus.Level) (logs.SeverityNumber, string) {
	v, ok := logrusSeverities[level]
	if ok == false {
		return 0, "UNSPECIFIED"
	}
	return v.value, v.name
}
