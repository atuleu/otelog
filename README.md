# otelog

Open Telemetry log support for golang. Basically it just provide a
missing LogExporter for `go.opentelemetry.io/otel`. The package also
provides a hook for `github.com/sirupsen/logrus`

This library is just a substitute while waiting on the official
implementation of the LogExporter in Open Telemetry. Indeed, in some
use case, like long lived batch job, you do not want to export your
log as span events, but rather real entries, that may be linked to the
span who started the span jobs.

This implementation is far to be error prone or optimized. For example
you cannot share your trace and log connection to the opentelemetry
connector. If this is a concern, wait or contribute to opentelemetry
to add this support.

Currently only tested with SigNoz and logrus.

## Example

See `examples/logrus-exporter` for an example of the API.
