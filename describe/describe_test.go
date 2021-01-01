package describe

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
	Cmd := Builder(&a, &b, &c)
	Cmd.Usage()
	assert.Greater(t, c.Len(), 0)
	assert.NoError(t, Cmd.Error())
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		{"-name", "helloworld"},
		// {"invalid"},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		Cmd := Builder(&a, &b, &c)
		Cmd.ParseArgs(args)
		assert.Error(t, Cmd.Error())
	}
}
