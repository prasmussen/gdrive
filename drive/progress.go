package drive

import (
    "io"
    "fmt"
    "time"
)

type Progress struct {
    Writer io.Writer
    Reader io.Reader
    Size int64
    progress int64
    rate int64
    rateProgress int64
    rateUpdated time.Time
    updated time.Time
    done bool
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

    // Update rate every 3 seconds
    if self.rateUpdated.Add(time.Second * 3).Before(now) {
        self.rate = calcRate(newProgress - self.rateProgress, self.rateUpdated, now)
        self.rateUpdated = now
        self.rateProgress = newProgress
    }

    // Draw progress every second
    if self.updated.Add(time.Second).Before(now) || isLast {
        self.Draw(isLast)
    }

    // Update last draw time
    self.updated = now

    // Mark as done if error occurs
    self.done = isLast

    return n, err
}

func (self *Progress) Draw(isLast bool) {
    if self.done {
        return
    }

    // Clear line
    fmt.Fprintf(self.Writer, "\r%50s", "")

    // Print progress
    fmt.Fprintf(self.Writer, "\r%s/%s", formatSize(self.progress, false), formatSize(self.Size, false))

    // Print rate
    if self.rate > 0 {
        fmt.Fprintf(self.Writer, ", Rate: %s/s", formatSize(self.rate, false))
    }

    if isLast {
        fmt.Fprintf(self.Writer, "\n")
    }
}
