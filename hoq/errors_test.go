package hoq

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrapErrWithCode(t *testing.T) {
	assert := require.New(t)
	err := WrapErrWithCode(400, errors.New("test"))
	assert.Equal("test", err.Error())
	assert.Equal(400, err.Code())

	err = NewErrWithCode(400, "test")
	assert.Equal("test", err.Error())
	assert.Equal(400, err.Code())
}
