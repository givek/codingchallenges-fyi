package resp_test

import (
	"testing"

	"github.com/givek/codingchallenges-fyi/redis-go/internal/resp"
)

func TestParse_NestedAndComplexArray(t *testing.T) {
	input := []byte(
		"*4\r\n" + // Outer array of 4 elements

			// Element 1: Nested array containing an error and a simple string
			"*2\r\n" +
			"-ERR custom server error\r\n" +
			"+OK\r\n" +

			// Element 2: Deeply nested array with a null array and a negative integer
			"*2\r\n" +
			"*-1\r\n" + // Null array
			":-99999\r\n" +

			// Element 3: Bulk string with binary/escaped data (newlines, tabs, JSON)
			"$38\r\n" +
			"line1: payload\nline2: \r\n\t{\"json\":true}\r\n" +

			// Element 4: Array containing an empty bulk string and a null bulk string
			"*2\r\n" +
			"$0\r\n\r\n" + // Empty bulk string
			"$-1\r\n", // Null bulk string
	)

	result, size, err := resp.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert that the parser consumed the entire input buffer
	if size != int64(len(input)) {
		t.Errorf("expected consumed size %d, got %d", len(input), size)
	}

	if result.Identifier() != resp.ArrayPrefix {
		t.Errorf(
			"expected identifier '%v', got %q",
			resp.ArrayPrefix,
			result.Identifier(),
		)
	}

	arr, ok := result.(*resp.Array)
	if !ok {
		t.Fatalf("expected *resp.Array, got %T", result)
	}

	// Assert the outer array structure length
	expectedLen := 4
	if len(arr.Value()) != expectedLen {
		t.Errorf("expected array length %d, got %d", expectedLen, len(arr.Value()))
	}
}

func TestParse_SimpleString(t *testing.T) {
	input := []byte("+OK\r\n")

	result, size, err := resp.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if size != int64(len(input)) {
		t.Errorf("expected size %d, got %d", len(input), size)
	}

	if result.Identifier() != resp.SimpleStringPrefix {
		t.Errorf("expected identifier '%v', got %q", resp.SimpleStringPrefix, result.Identifier())
	}

	simpleStr, ok := result.(*resp.SimpleString)
	if !ok {
		t.Fatalf("expected *resp.SimpleString, got %T", result)
	}

	if simpleStr.Value() != "OK" {
		t.Errorf("expected value 'OK', got %q", simpleStr.Value())
	}
}

func TestParse_SimpleError(t *testing.T) {
	input := []byte("-ERR unknown command 'foobar'\r\n")

	result, size, err := resp.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if size != int64(len(input)) {
		t.Errorf("expected size %d, got %d", len(input), size)
	}

	if result.Identifier() != resp.ErrorPrefix {
		t.Errorf("expected identifier '%v', got %q", resp.ErrorPrefix, result.Identifier())
	}

	errVal, ok := result.(*resp.Error) // Adjust type name if your codebase uses resp.SimpleError
	if !ok {
		t.Fatalf("expected *resp.Error, got %T", result)
	}

	expectedMsg := "ERR unknown command 'foobar'"
	if errVal.Message() != expectedMsg {
		t.Errorf("expected value %q, got %q", expectedMsg, errVal.Message())
	}
}

func TestParse_Integer(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected int64
	}{
		{"Positive Integer", []byte(":1000\r\n"), 1000},
		{"Zero", []byte(":0\r\n"), 0},
		{"Negative Integer", []byte(":-42\r\n"), -42},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, size, err := resp.Parse(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if size != int64(len(tc.input)) {
				t.Errorf("expected size %d, got %d", len(tc.input), size)
			}

			if result.Identifier() != resp.IntPrefix {
				t.Errorf("expected identifier '%v', got %q", resp.IntPrefix, result.Identifier())
			}

			intVal, ok := result.(*resp.Int)
			if !ok {
				t.Fatalf("expected *resp.Integer, got %T", result)
			}

			if intVal.Value() != tc.expected {
				t.Errorf("expected value %d, got %d", tc.expected, intVal.Value())
			}
		})
	}
}

func TestParse_BulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"Normal Bulk String", []byte("$5\r\nhello\r\n"), "hello"},
		{"Empty Bulk String", []byte("$0\r\n\r\n"), ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, size, err := resp.Parse(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if size != int64(len(tc.input)) {
				t.Errorf("expected size %d, got %d", len(tc.input), size)
			}

			if result.Identifier() != resp.BulkStringPrefix {
				t.Errorf("expected identifier '%v', got %q", resp.BulkStringPrefix, result.Identifier())
			}

			bulkStr, ok := result.(*resp.BulkString)
			if !ok {
				t.Fatalf("expected *resp.BulkString, got %T", result)
			}

			if string(bulkStr.Value()) != tc.expected {
				t.Errorf("expected value %q, got %q", tc.expected, string(bulkStr.Value()))
			}
		})
	}
}

func TestParse_NullBulkString(t *testing.T) {
	input := []byte("$-1\r\n")

	result, size, err := resp.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if size != int64(len(input)) {
		t.Errorf("expected size %d, got %d", len(input), size)
	}

	if result.Identifier() != resp.BulkStringPrefix { // Or NullPrefix depending on design
		t.Errorf("expected identifier '%v', got %q", resp.BulkStringPrefix, result.Identifier())
	}

	_, ok := result.(*resp.Null)
	if !ok {
		t.Fatalf("expected *resp.Null, got %T", result)
	}
}

func TestParse_NullArray(t *testing.T) {
	input := []byte("*-1\r\n")

	result, size, err := resp.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if size != int64(len(input)) {
		t.Errorf("expected size %d, got %d", len(input), size)
	}

	_, ok := result.(*resp.Null) // Or specialized null array type depending on implementation
	if !ok {
		t.Fatalf("expected *resp.Null, got %T", result)
	}
}

func TestParse_ErrorsAndMalformed(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"Empty Input", []byte("")},
		{"Unknown Prefix", []byte("X5\r\nhello\r\n")},
		{"Missing CRLF", []byte("+OK")},
		{"Invalid Integer Format", []byte(":abc\r\n")},
		{"Bulk String Length Mismatch", []byte("$5\r\nhi\r\n")},
		{"Incomplete Array", []byte("*2\r\n+one\r\n")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := resp.Parse(tc.input)
			if err == nil {
				t.Errorf("expected error for malformed input %q, got nil", tc.input)
			}
		})
	}
}

func TestParse_TrailingData(t *testing.T) {
	// Ensures that Parse accurately returns the consumed size when extra trailing bytes exist
	input := []byte("+OK\r\n+EXTRA\r\n")

	result, size, err := resp.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedSize := int64(5) // "+OK\r\n" length
	if size != expectedSize {
		t.Errorf("expected size to be %d, got %d", expectedSize, size)
	}

	simpleStr, ok := result.(*resp.SimpleString)
	if !ok {
		t.Fatalf("expected *resp.SimpleString, got %T", result)
	}

	if simpleStr.Value() != "OK" {
		t.Errorf("expected 'OK', got %q", simpleStr.Value())
	}
}
