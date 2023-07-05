package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/atuleu/otelog"
	"github.com/atuleu/otelog/pkg/hooks"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var serviceName = flag.String("service", "logrus-export", "define the service name to send")
var serviceInstance = flag.String("instance", "instance1", "define the service instance to send")
var serviceVersion = flag.String("version", "v0.1.0", "define the service instance to send")
var period = flag.Duration("period", 2000*time.Millisecond, "period to generate log to")
var endpoint = flag.String("endpoint", "localhost:4317", "define the endpoint to send log to")

func main() {
	if err := execute(); err != nil {
		log.Fatalf("%s", err)
	}
}

func execute() error {
	flag.Parse()
	if err := setUpLogger(); err != nil {
		return err
	}

	setUpLogrusHook()

	for {
		logSomethingRandom()
		time.Sleep(*period)
	}
}

func setUpLogger() error {
	resource := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceName(*serviceName),
		semconv.ServiceVersion(*serviceVersion),
		semconv.ServiceInstanceID(*serviceInstance),
	)
	exporter, err := otelog.NewLogExporter(
		otelog.WithEndpoint(*endpoint),
		otelog.WithInsecure(),
		otelog.WithBatchLogProcessor(otelog.WithBatchTimeout(5*time.Second)),
		otelog.WithResource(resource),
	)
	if err != nil {
		return err
	}
	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(*endpoint),
		),
	)
	if err != nil {
		return err
	}

	otel.SetTracerProvider(sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	))

	otelog.SetLogExporter(exporter)
	return nil
}

func setUpLogrusHook() {
	hook := hooks.NewLogrusHook(hooks.FromLogrusLevel(logrus.InfoLevel))
	logrus.AddHook(hook)
}

func logSomethingRandom() {
	ctx := context.Background()
	var err error
	if rand.Int()%2 == 0 {
		err = errors.New("something bad happened")
	}

	if rand.Int()%2 == 0 {
		var span trace.Span
		ctx, span = otel.Tracer("logrus-export").Start(ctx, "someWorkWillBeDone")
		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, "process error")
				span.RecordError(err)
			}
			span.End()
		}()
	}

	// adding the context will automatically link the log entry to the
	// span.
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"value":   rand.Float64() * 42.0,
		"another": rand.Intn(42),
	})

	if err != nil {
		logger.WithError(err).Error("operation executed")
	} else {
		logger.Info("operation executed")
	}
}
