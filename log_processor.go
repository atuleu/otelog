package otelog

import (
	"sync"
	"sync/atomic"
	"time"
)

type logBatchCallback func([]*LogRecord)

type LogProcessor interface {
	// batch adds the record to the current batch. If the current
	// batch is ready to be sent, callback with the current batch
	// content will be called.
	batch(record *LogRecord, callback logBatchCallback)
}

type syncProcessor struct{}

func (b *syncProcessor) batch(record *LogRecord, callback logBatchCallback) {
	callback([]*LogRecord{record})
}

type batchProcessorOptions struct {
	MaxQueueSize int
	BatchTimeout time.Duration
}

type BatchLogProcessorOption interface {
	apply(options *batchProcessorOptions)
}

type batchLogTimeout time.Duration

func (t batchLogTimeout) apply(opts *batchProcessorOptions) {
	opts.BatchTimeout = time.Duration(t)
}

func WithBatchTimeout(timeout time.Duration) BatchLogProcessorOption {
	return batchLogTimeout(timeout)
}

type batchQueueSize int

func (s batchQueueSize) apply(opts *batchProcessorOptions) {
	opts.MaxQueueSize = int(s)
}

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
	buffer  []*LogRecord
	size    atomic.Int32
}

func newBatchProcessor(options ...BatchLogProcessorOption) LogProcessor {
	opts := newBatchProcessorOptions(options...)

	res := &batchProcessor{
		timeout: opts.BatchTimeout,
		buffer:  make([]*LogRecord, opts.MaxQueueSize),
	}

	res.cond = sync.NewCond(res.mx.RLocker())
	return res
}

func (b *batchProcessor) batch(record *LogRecord, callback logBatchCallback) {
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
	b.buffer = make([]*LogRecord, len(b.buffer))
	b.size.Store(0)

	go callback(batch)
}
