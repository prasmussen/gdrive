package drive

import (
    "io"
    "time"
    "sync"
    "golang.org/x/net/context"
)

const MaxIdleTimeout = time.Second * 120
const TimeoutTimerInterval = time.Second * 10

type timeoutReaderWrapper func(io.Reader) io.Reader

func getTimeoutReaderWrapperContext() (timeoutReaderWrapper, context.Context) {
    ctx, cancel := context.WithCancel(context.TODO())
    wrapper := func(r io.Reader) io.Reader {
         return getTimeoutReader(r, cancel)
    }
    return wrapper, ctx
}

func getTimeoutReaderContext(r io.Reader) (io.Reader, context.Context) {
    ctx, cancel := context.WithCancel(context.TODO())
    return getTimeoutReader(r, cancel), ctx
}

func getTimeoutReader(r io.Reader, cancel context.CancelFunc) io.Reader {
    return &TimeoutReader{
        reader: r,
        cancel: cancel,
        mutex: &sync.Mutex{},
    }
}

type TimeoutReader struct {
    reader io.Reader
    cancel context.CancelFunc
    lastActivity time.Time
    timer *time.Timer
    mutex *sync.Mutex
    done bool
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

    if time.Since(self.lastActivity) > MaxIdleTimeout {
        self.cancel()
        self.mutex.Unlock()
        return
    }

    self.mutex.Unlock()
    self.startTimer()
}
