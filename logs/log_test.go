package logs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger_log(t *testing.T) {
	log := NewHLogger()
	log.Debug("DEBUG")
	log.Warn("WARN")
	log.Info("HELLO")
	log.Error("ERROR")
	Error("1", 2, 3)
}

func TestNoNewLog(t *testing.T) {
	assert.NotPanics(t, func() {
		msgs := []interface{}{1, "2", "3", "4", 5.0}
		Debug(msgs...)
		Info(msgs...)
		Warn(msgs...)
		Error(msgs...)
	})
}

func TestError(t *testing.T) {
	assert.NotPanics(t, func() {
		e := errors.New("test")
		Error(e)
	})
}
