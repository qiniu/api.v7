package qerr

import (
	"fmt"
)

// Error 接口以code, message, orig error的形式封装了底层的错误
// 调用Error() 和String()会返回具体的报错信息
type IError interface {
	// 继承golang error 接口
	error

	// 返回错误的分类信息
	Code() string

	// 返回错误的详细信息
	Message() string

	// 返回原始的报错信息
	OrigErr() error
}

// QiniuError 是Error接口的一个具体实现
type qiniuError struct {
	code     string
	message  string
	origErrs []error
}

func newQiniuError(code, message string, origErrs []error) *qiniuError {
	return &qiniuError{
		code:     code,
		message:  message,
		origErrs: origErrs,
	}
}

// Error 返回错误的信息
func (qe *qiniuError) Error() string {
	size := len(qe.origErrs)
	if size > 0 {
		return FormatError(qe.code, qe.message, "", errorList(qe.errs))
	}

	return FormatError(qe.code, qe.message, "", nil)
}

// 和Error() 功能一样
func (qe *qiniuError) String() string {
	return qe.Error()
}

// Code 返回错误信息分类码
func (qe *qiniuError) Code() string {
	return qe.code
}

// Message 返回具体的错误信息
func (qe *qiniuError) Message() string {
	return qe.message
}

// OrigErr 返回原始错误
func (qe *qiniuError) OrigErrs() []error {
	return qe.origErrs
}

// 返回error
func (qe *qiniuError) OrigErr() error {
	switch len(qe.origErrs) {
	case 0:
		return nil
	case 1:
		return qe.origErrs[0]
	default:
		if err, ok := qe.origErrs[0].(IError); ok {
			return newQiniuError(err.Code, err.Message(), qe.origErrs[1:])
		}
		return newQiniuError("BatchedErrors",
			"multiple errors occurred", qe.origErrs)
	}
}

// New 返回符合Error接口的一个实现
func New(code, message string, origErr error) Error {
	var errs []error
	if origErr != nil {
		errs = append(errs, origErr)
	}
	return newQiniuError(code, message, errs)
}

// RequestFailure 接口从错误信息Error中获取七牛请求中的特定信息， 比如reqid等信息
// 有可能错误信息中没有reqid信息， 比如request请求还没有达到我们的服务就报错了
type IRequestFailure interface {
	IError

	// http返回的状态码
	StatusCode() int

	// 服务接口返回的reqid信息， 又可能没有
	RequestID() string
}

type requestFailure struct {
	// 继承IError接口
	IError
	// http请求返回的状态码
	statusCode int

	// 服务端返回的reqid
	reqId string
}

// NewRequestFailure 返回一个具体的实现了IRequestFailure接口的实现
func NewRequestFailure(err IError, statusCode int, reqID string) IRequestFailure {
	return &requestFailure{
		IError:     IError,
		statusCode: statusCode,
		reqId:      reqID,
	}
}

func (rf *requestFailure) Error() string {
	extra := fmt.Sprintf("status code: %d, request id: %s",
		rf.statusCode, rf.requestID)
	return FormatError(r.Code(), r.Message(), extra, r.OrigErr())
}

// 符合stringer 接口
func (rf *requestFailure) String() string {
	return rf.Error()
}

// 返回http状态码
func (rf *requestFailure) StatusCode() string {
	return rf.statusCode
}

// 返回reqid
func (rf *requestFailure) RequestID() string {
	return rf.reqId
}

// 返回原始错误
func (rf *requestFailure) OrigErrs() []error {
	if b, ok := rf.IError.(qiniuError); ok {
		return b.OrigErrs()
	}
	return []error{rf.OrigErr()}
}

// 格式化输出错误信息
func FormatError(code, message, extra string, origErr error) string {
	msg := fmt.Sprintf("%s: %s", code, message)
	if extra != "" {
		msg = fmt.Sprintf("%s\n\t%s", msg, extra)
	}
	if origErr != nil {
		msg = fmt.Sprintf("%s\ncaused by: %s", msg, origErr.Error())
	}
	return msg
}

type errorList []error

func (e errorList) Error() string {
	msg := ""
	if size := len(e); size > 0 {
		for i := 0; i < size; i++ {
			msg += fmt.Sprintf("%s", e[i].Error())
		}
	}
	return msg
}
