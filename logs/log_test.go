package logs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogger_log(t *testing.T) {
	log := NewHLogger()
	log.Debug("DEBUG")
	log.Warn("WARN")
	log.Info("HELLO")
	log.Error("ERROR")
	Error("1", 2, 3)
}

var testLogger = NewHLogger()

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

func TestSetLevel(t *testing.T) {
	assert.NotPanics(t, func() {
		testLogger.Debug("TEST")
		testLogger.Info("TEST")
		testLogger.SetLevel(LevelInfo)
		testLogger.Debug("TEST")
		testLogger.Info("TEST")

		Debug("TEST")
		Info("TEST")
		SetLevel(LevelInfo)
		Debug("TEST")
		Info("TEST")
		//there should only six lines
	})
	time.Sleep(1 * time.Millisecond)
}

func TestNewHLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		x := NewHLogger()
		assert.NotNil(t, x)
	})
}

func TestGetAlignedLevel(t *testing.T) {
	length := len(testLogger.getAlignedLevel(LevelDebug))
	for i, _ := range testLogger.levelS {
		if len(testLogger.getAlignedLevel(LogLevel(i))) != length {
			t.Error("not all levels long the same")
		}
	}
}
