package resp_test

import (
	"fmt"
	"testing"

	"github.com/givek/codingchallenges-fyi/redis-go/internal/resp"
)

func TestParse_Array(t *testing.T) {
	// input := []byte("*3\r\n+cade\r\n:42\r\n*3\r\n+hello\r\n+world\r\n:64\r\n")
	// input := []byte("*2\r\n:42\r\n*1\r\n+hello\r\n")
	// input := []byte("*3\r\n:42\r\n+hello\r\n:64\r\n")
	// input := []byte("*5\r\n" + // 4
	// 	":42\r\n" + // 5
	// 	"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" + // 22
	// 	"$-1\r\n" + // 5
	// 	"$18\r\nhello\r\nworld\nredis\r\n" + // 25
	// 	"$-1\r\n", // 5
	// )
	input := []byte(
		"*4\r\n" + // Outer array of 4 elements // 4
			// Element 1: Nested array containing an error and a simple string
			"*2\r\n" + // 4
			"-ERR custom server error\r\n" + // 26
			"+OK\r\n" + // 5
			// Element 2: Deeply nested array with a null array and an integer
			"*2\r\n" + // 4
			"*-1\r\n" + // Null array / 5
			":-99999\r\n" + // Negative integer // 9
			// Element 3: Bulk string with binary/escaped data (newlines, quotes, unicode)
			"$38\r\n" + // 5
			"line1: payload\nline2: \r\n\t{\"json\":true}\r\n" + // 40
			// Element 4: Array containing an empty bulk string and a null bulk
			"*2\r\n" +
			"$0\r\n\r\n" + // Empty bulk string
			"$-1\r\n", // Null bulk string
	)

	result, size, err := resp.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Identifier() != resp.ArrayPrefix {
		t.Errorf(
			"expected identifier '%v', got %q",
			resp.ArrayPrefix,
			result.Identifier(),
		)
	}

	simpleStr, ok := result.(*resp.Array)
	if !ok {
		t.Fatalf("expected *resp.SimpleString, got %T", result)
	}

	fmt.Println("SizePQY -", size)

	for _, v := range simpleStr.Value() {
		fmt.Printf("%v - %q\n", string(v.Identifier()), v)
	}

}

// func TestParse_SimpleString(t *testing.T) {
// 	input := []byte("+OK\r\n")
//
// 	result, size, err := resp.Parse(input)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
//
// 	if result.Identifier() != resp.SimpleStringPrefix {
// 		t.Errorf(
// 			"expected identifier '%v', got %q",
// 			resp.SimpleStringPrefix,
// 			result.Identifier(),
// 		)
// 	}
//
// 	simpleStr, ok := result.(*resp.SimpleString)
// 	if !ok {
// 		t.Fatalf("expected *resp.SimpleString, got %T", result)
// 	}
//
// 	fmt.Println("SizePQY -", size)
//
// 	if simpleStr.Value() != "OK" {
// 		t.Errorf("expected value 'OK', got %q", simpleStr.Value())
// 	}
// }

//
// func TestParse_Null(t *testing.T) {
// 	input := []byte("$-1\r\n")
//
// 	result, err := resp.Parse(input)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
//
// 	if result.Identifier() != resp.BulkStringPrefix {
// 		t.Errorf(
// 			"expected identifier '%v', got %q",
// 			resp.BulkStringPrefix,
// 			result.Identifier(),
// 		)
// 	}
//
// 	_, ok := result.(*resp.Null)
// 	if !ok {
// 		t.Fatalf("expected *resp.Null, got %T", result)
// 	}
// }
