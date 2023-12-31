package otelog

import (
	"sync"
	"sync/atomic"
	"time"

	logs "go.opentelemetry.io/proto/otlp/logs/v1"
)

type logBatchCallback func([]*logs.LogRecord)

// LogProcessor process incoming LogRecord and batches them if needed.
type LogProcessor interface {
	// batch adds the record to the current batch. If the current
	// batch is ready to be sent, callback with the current batch
	// content will be called.
	batch(record *logs.LogRecord, callback logBatchCallback)
}

type syncProcessor struct{}

func (b *syncProcessor) batch(record *logs.LogRecord, callback logBatchCallback) {
	callback([]*logs.LogRecord{record})
}

type batchProcessorOptions struct {
	MaxQueueSize int
	BatchTimeout time.Duration
}

// BatchLogProcessorOption is an Option for WithBatchLogProcessor()
type BatchLogProcessorOption interface {
	apply(options *batchProcessorOptions)
}

type batchLogTimeout time.Duration

func (t batchLogTimeout) apply(opts *batchProcessorOptions) {
	opts.BatchTimeout = time.Duration(t)
}

// WithBatchTimeout sets the timeout after which a batch will be
// exported regardless of the number of queued LogRecord.
func WithBatchTimeout(timeout time.Duration) BatchLogProcessorOption {
	return batchLogTimeout(timeout)
}

type batchQueueSize int

func (s batchQueueSize) apply(opts *batchProcessorOptions) {
	opts.MaxQueueSize = int(s)
}

// WithMaxQueueSize sets the max queue size before a batch will be
// emitted.
func WithMaxQueueSize(size int) BatchLogProcessorOption {
	return batchQueueSize(size)
}

func newBatchProcessorOptions(options ...BatchLogProcessorOption) batchProcessorOptions {
	res := batchProcessorOptions{
		MaxQueueSize: 512,
		BatchTimeout: 1000 * time.Millisecond,
	}
	for _, o := range options {
		o.apply(&res)
	}

	return res
}

type batchProcessor struct {
	mx   sync.RWMutex
	cond *sync.Cond

	timeout time.Duration
	buffer  []*logs.LogRecord
	size    atomic.Int32
}

func newBatchProcessor(options ...BatchLogProcessorOption) LogProcessor {
	opts := newBatchProcessorOptions(options...)

	res := &batchProcessor{
		timeout: opts.BatchTimeout,
		buffer:  make([]*logs.LogRecord, opts.MaxQueueSize),
	}

	res.cond = sync.NewCond(res.mx.RLocker())
	return res
}

func (b *batchProcessor) batch(record *logs.LogRecord, callback logBatchCallback) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	var newSize int
	for newSize = int(b.size.Add(1)); newSize > len(b.buffer); {
		b.cond.Wait()
	}

	b.buffer[newSize-1] = record

	if newSize == 1 {
		go func() {
			time.Sleep(b.timeout)
			b.process(callback)
		}()
	}

	if newSize == len(b.buffer) {
		go b.process(callback)
	}
}

func (b *batchProcessor) process(callback logBatchCallback) {
	b.mx.Lock()
	defer b.mx.Unlock()

	size := int(b.size.Load())
	if size == 0 {
		return
	}

	defer b.cond.Broadcast()

	batch := b.buffer[:size]
	b.buffer = make([]*logs.LogRecord, len(b.buffer))
	b.size.Store(0)

	go callback(batch)
}
