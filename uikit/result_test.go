package uikit

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResult(t *testing.T) {
	var result = NewResult()
	assert.Equal(t, 0, result.Code)
	assert.Equal(t, "", result.Msg)
	assert.Equal(t, 0, len(result.Data), "should be inited")

}

func TestR_Ok(t *testing.T) {
	result := NewResult().Ok()
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "", result.Msg)

}
func TestR_Err(t *testing.T) {
	result := NewResult().Err(nil)
	assert.Equal(t, http.StatusBadRequest, result.Code)
	assert.Equal(t, "", result.Msg)
	result.Err(assert.AnError)
	assert.Equal(t, assert.AnError.Error(), result.Msg)
}

func TestR_ErrCodeMsg(t *testing.T) {
	result := NewResult().ErrCodeMsg(404, "not found")
	assert.Equal(t, 404, result.Code)
	assert.Equal(t, "not found", result.Msg)
}

func TestR_Put(t *testing.T) {
	result := NewResult().Put("test", "test")
	assert.Equal(t, "test", result.Data["test"])
}
