package Request

import (
	"time"
)

type Retryer interface {
	ShouldRetry(*Request) bool
	MaxRetries() int
	RetryDelay(*Request) time.Duration
}
