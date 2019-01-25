package client

import (
	"strconv"
	"time"

	"github.com/qiniu/api.v7/helper"
	"github.com/qiniu/api.v7/request"
)

// DefaultRetryer implements basic retry logic using exponential backoff for
// most services. If you want to implement custom retry logic, implement the
// request.Retryer interface or create a structure type that composes this
// struct and override the specific methods. For example, to override only
// the MaxRetries method:
//
//		type retryer struct {
//      client.DefaultRetryer
//    }
//
//    // This implementation always has 100 max retries
//    func (d retryer) MaxRetries() int { return 100 }
type DefaultRetryer struct {
	NumMaxRetries int
}

// MaxRetries 返回最大的重试次数
func (d DefaultRetryer) MaxRetries() int {
	return d.NumMaxRetries
}

// RetryDelay 返回请求的时间间隔， 当请求重试的时候
func (d DefaultRetryer) RetryDelay(r *request.Request) time.Duration {
	// Set the upper limit of delay in retrying at ~five minutes
	minTime := 30
	if delay, ok := getRetryDelay(r); ok {
		return delay
	}
	retryCount := r.RetryCount
	if retryCount > 13 {
		retryCount = 13
	}

	delay := (1 << uint(retryCount)) * (helper.SeededRand.Intn(minTime) + minTime)
	return time.Duration(delay) * time.Millisecond
}

// ShouldRetry 返回这个请求是否可以重试
func (d DefaultRetryer) ShouldRetry(r *request.Request) bool {
	if r.Retryable != nil {
		return r.Retryable
	}

	if r.HTTPResponse.StatusCode >= 500 && r.HTTPResponse.StatusCode != 501 {
		return true
	}
	return r.IsErrorRetryable()
}

// 查看请求的response中有没有Retry-After header
// 如果有， 按照header指示， 下一个请求按照这个延时
func getRetryDelay(r *request.Request) (time.Duration, bool) {
	if !canUseRetryAfterHeader(r) {
		return 0, false
	}

	delayStr := r.HTTPResponse.Header.Get("Retry-After")
	if len(delayStr) == 0 {
		return 0, false
	}

	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		return 0, false
	}

	return time.Duration(delay) * time.Second, true
}

// Retry-After header指示客户端需要等待多长时间后再发送请求
// 可以应用与该header 的statusCode有两种
// 503 - service unavailable
// 429 - Two many requests
func canUseRetryAfterHeader(r *request.Request) bool {
	switch r.HTTPResponse.StatusCode {
	case 429:
	case 503:
	default:
		return false
	}

	return true
}
