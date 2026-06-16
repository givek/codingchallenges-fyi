package counter_test

import (
	"bytes"
	"testing"

	"github.com/givek/codingchallenges-fyi/wc-tool-go/counter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSizeInBytes(t *testing.T) {
	r := bytes.NewBuffer([]byte("Hello"))

	bc, err := counter.SizeInBytes(r)

	require.NoError(t, err, "failed to parse input for test case:")

	expected := int64(5)

	assert.Equal(t, bc, expected, "expect: %v, got: %v", expected, bc)
}
