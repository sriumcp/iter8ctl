package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletePath(t *testing.T) {
	p1 := CompletePath("", "a")
	p2 := CompletePath("../", "utils/a")
	p3 := CompletePath("", "b")
	assert.Equal(t, p1, p2)
	assert.NotEqual(t, p2, p3)
}
