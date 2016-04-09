package drive

import (
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"time"
)

const MaxErrorRetries = 5

func isBackendOrRateLimitError(err error) bool {
	return isBackendError(err) || isRateLimitError(err)
}

func isBackendError(err error) bool {
	if err == nil {
		return false
	}

	ae, ok := err.(*googleapi.Error)
	return ok && ae.Code >= 500 && ae.Code <= 599
}

func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	ae, ok := err.(*googleapi.Error)
	return ok && ae.Code == 403
}

func isTimeoutError(err error) bool {
	return err == context.Canceled
}

func exponentialBackoffSleep(try int) {
	seconds := pow(2, try)
	time.Sleep(time.Duration(seconds) * time.Second)
}
