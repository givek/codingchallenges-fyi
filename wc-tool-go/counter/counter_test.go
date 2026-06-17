package counter_test

import (
	"bytes"
	"fmt"
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
			name:     "empty string",
			input:    []byte(""),
			expected: 0,
		},
		{
			name:     "only ASCII characters",
			input:    []byte(fmt.Sprintf("Hello, \n Hope you are doing well!")),
			expected: 33,
		},
		{
			name:     "multi-byte UTF-8 characters (emojis and kanji)",
			input:    []byte("Hello 👋 世界"), // "👋" is 4 bytes, "世" and "界" are 3 bytes each
			expected: 17,
		},
		{
			name:     "large input (10 KB)",
			input:    bytes.Repeat([]byte("A"), 10240),
			expected: 10240,
		},
		{
			name:     "control characters and mixed newlines",
			input:    []byte("\x00\n\r\t"),
			expected: 4,
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
