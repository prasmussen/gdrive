package drive

import (
	"golang.org/x/net/context"
	"io"
	"sync"
	"time"
)

const TimeoutTimerInterval = time.Second * 10

type timeoutReaderWrapper func(io.Reader) io.Reader

func getTimeoutReaderWrapperContext(timeout time.Duration) (timeoutReaderWrapper, context.Context) {
	ctx, cancel := context.WithCancel(context.TODO())
	wrapper := func(r io.Reader) io.Reader {
		// Return untouched reader if timeout is 0
		if timeout == 0 {
			return r
		}

		return getTimeoutReader(r, cancel, timeout)
	}
	return wrapper, ctx
}

func getTimeoutReaderContext(r io.Reader, timeout time.Duration) (io.Reader, context.Context) {
	ctx, cancel := context.WithCancel(context.TODO())

	// Return untouched reader if timeout is 0
	if timeout == 0 {
		return r, ctx
	}

	return getTimeoutReader(r, cancel, timeout), ctx
}

func getTimeoutReader(r io.Reader, cancel context.CancelFunc, timeout time.Duration) io.Reader {
	return &TimeoutReader{
		reader:         r,
		cancel:         cancel,
		mutex:          &sync.Mutex{},
		maxIdleTimeout: timeout,
	}
}

type TimeoutReader struct {
	reader         io.Reader
	cancel         context.CancelFunc
	lastActivity   time.Time
	timer          *time.Timer
	mutex          *sync.Mutex
	maxIdleTimeout time.Duration
	done           bool
}

func (self *TimeoutReader) Read(p []byte) (int, error) {
	if self.timer == nil {
		self.startTimer()
	}

	self.mutex.Lock()

	// Read
	n, err := self.reader.Read(p)

	self.lastActivity = time.Now()
	self.done = (err != nil)

	self.mutex.Unlock()

	if self.done {
		self.stopTimer()
	}

	return n, err
}

func (self *TimeoutReader) startTimer() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if !self.done {
		self.timer = time.AfterFunc(TimeoutTimerInterval, self.timeout)
	}
}

func (self *TimeoutReader) stopTimer() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.timer != nil {
		self.timer.Stop()
	}
}

func (self *TimeoutReader) timeout() {
	self.mutex.Lock()

	if self.done {
		self.mutex.Unlock()
		return
	}

	if time.Since(self.lastActivity) > self.maxIdleTimeout {
		self.cancel()
		self.mutex.Unlock()
		return
	}

	self.mutex.Unlock()
	self.startTimer()
}
