package drive

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

const MaxDrawInterval = time.Second * 1
const MaxRateInterval = time.Second * 3

func getProgressReader(r io.Reader, w io.Writer, size int64) io.Reader {
	// Don't wrap reader if output is discarded or size is too small
	if w == ioutil.Discard || (size > 0 && size < 1024*1024) {
		return r
	}

	return &Progress{
		Reader: r,
		Writer: w,
		Size:   size,
	}
}

type Progress struct {
	Writer       io.Writer
	Reader       io.Reader
	Size         int64
	progress     int64
	rate         int64
	rateProgress int64
	rateUpdated  time.Time
	updated      time.Time
	done         bool
}

func (self *Progress) Read(p []byte) (int, error) {
	// Read
	n, err := self.Reader.Read(p)

	now := time.Now()
	isLast := err != nil

	// Increment progress
	newProgress := self.progress + int64(n)
	self.progress = newProgress

	// Initialize rate state
	if self.rateUpdated.IsZero() {
		self.rateUpdated = now
		self.rateProgress = newProgress
	}

	// Update rate every x seconds
	if self.rateUpdated.Add(MaxRateInterval).Before(now) {
		self.rate = calcRate(newProgress-self.rateProgress, self.rateUpdated, now)
		self.rateUpdated = now
		self.rateProgress = newProgress
	}

	// Draw progress every x seconds
	if self.updated.Add(MaxDrawInterval).Before(now) || isLast {
		self.draw(isLast)
		self.updated = now
	}

	// Mark as done if error occurs
	self.done = isLast

	return n, err
}

func (self *Progress) draw(isLast bool) {
	if self.done {
		return
	}

	self.clear()

	// Print progress
	fmt.Fprintf(self.Writer, "%s", formatSize(self.progress, false))

	// Print total size
	if self.Size > 0 {
		fmt.Fprintf(self.Writer, "/%s", formatSize(self.Size, false))
	}

	// Print rate
	if self.rate > 0 {
		fmt.Fprintf(self.Writer, ", Rate: %s/s", formatSize(self.rate, false))
	}

	if isLast {
		self.clear()
	}
}

func (self *Progress) clear() {
	fmt.Fprintf(self.Writer, "\r%50s\r", "")
}
