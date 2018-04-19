package uikit

import (
	"net/http"
)

// R return struct
type R struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// NewResult make new Result
func NewResult() *R {
	return &R{
		Data: make(map[string]interface{}, 0),
	}
}

// Err make normal error result
func (r *R) Err(err error) *R {
	r.Code = http.StatusBadRequest

	if err != nil {
		r.Msg = err.Error()
	}
	return r
}

// ErrMsg  make normal error result
func (r *R) ErrMsg(msg string) *R {
	r.Code = http.StatusBadRequest
	r.Msg = msg
	return r
}

// ErrCodeMsg set user define error msg
func (r *R) ErrCodeMsg(httpStatusCode int, msg string) *R {
	r.Code = httpStatusCode
	r.Msg = msg
	return r
}

// Ok init result
func (r *R) Ok(msg ...string) *R {
	r.Code = http.StatusOK
	r.Msg = ""

	for _, value := range msg {
		r.Msg += value
	}
	return r
}

// Put use map[string]interface{} to store result data
func (r *R) Put(key string, data interface{}) *R {
	r.Data[key] = data
	return r
}

// SetPage add pagination info to result
func (r *R) SetPage(page *Pagination) *R {
	if page == nil {
		return r
	}
	r.Put("page", page)
	return r
}

// Reset reset
func (r *R) Reset() *R {
	r.Code = 0
	r.Msg = ""
	r.Data = nil
	return r
}
