package counter_test

import (
	"bytes"
	"testing"

	"github.com/givek/codingchallenges-fyi/wc-tool-go/counter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSizeInBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected int64
	}{
		{
			name:     "hello", // TODO
			input:    []byte("Hello"),
			expected: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sizeInBytes, err := counter.SizeInBytes(bytes.NewBuffer(test.input))

			require.NoError(t, err, "failed to parse input for test case:")

			assert.Equal(t, test.expected, sizeInBytes, "expect: %v, got: %v", test.expected, sizeInBytes)
		})
	}
}
