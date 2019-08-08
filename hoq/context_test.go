package hoq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContext(t *testing.T) {
	assert.NotPanics(t, func() {
		c := Context{}
		c.Done()
		c.Deadline()
		c.Err()
		c.Value("")
	})
}
