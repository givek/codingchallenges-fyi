package resp

import (
	"bytes"
	"fmt"
	"strconv"
)

func isNull(b []byte) (bool, int64) {
	if len(b) != 5 {
		return false, 0
	}

	return ((b[0] == BulkStringPrefix || b[0] == ArrayPrefix) &&
		b[1] == '-' &&
		b[2] == '1' &&
		b[3] == '\r' &&
		b[4] == '\n'), int64(len(b))
}

func parseSimpleString(input []byte) (*SimpleString, int64, error) {
	strEndIdx := bytes.Index(input, []byte{'\r', '\n'})

	if strEndIdx == -1 {
		return nil, 0, fmt.Errorf("invalid simple string")
	}

	strValue := input[:strEndIdx]

	return NewSimpleString(string(strValue)), int64(strEndIdx) + 2, nil
}

func parseBulkString(input []byte) (*BulkString, int64, error) {
	strEndIdx := bytes.Index(input, []byte{'\r', '\n'})

	if strEndIdx == -1 {
		return nil, 0, fmt.Errorf("invalid bulk string")
	}

	sizeStr := input[:strEndIdx]
	size, err := strconv.ParseInt(string(sizeStr), 10, 64)
	if err != nil {
		return nil, 0, err
	}

	se := int64(strEndIdx) + 2

	xs := se + size
	// A bluk string must end with CRLF (\r\n)
	if xs+1 >= int64(len(input)) {
		return nil, 0, fmt.Errorf("invalid bulk string")
	}

	if input[xs] != '\r' && input[xs+1] != '\n' {
		return nil, 0, fmt.Errorf("invalid bulk string")
	}

	val := string(input)[se:xs]

	return NewBulkString(val), xs + 2, nil
}

func parseArray(input []byte) (*Array, int64, error) {
	strEndIdx := bytes.Index(input, []byte{'\r', '\n'})

	if strEndIdx == -1 {
		return nil, 0, fmt.Errorf("invalid array")
	}

	sizeStr := input[:strEndIdx]
	size, err := strconv.ParseInt(string(sizeStr), 10, 64)
	if err != nil {
		return nil, 0, err
	}

	se := int64(strEndIdx) + 2

	val := input[se:]

	res := make([]DataType, size)

	offset := int64(0)
	for i := range res {
		// This will keep parsing the 1st element, we need maybe return the amount of bytes we parsed

		vxi := val[offset:]

		// fmt.Printf("%q\n", string(val))
		// fmt.Printf("Start ---- %v - %v - %q\n", i, offset, string(vxi))
		// fmt.Printf("%v - %v - %q\n", i, offset, string(val[:offset]))
		// fmt.Println("")
		v, size, err := Parse(vxi)
		if err != nil {
			return nil, 0, err
		}

		// fmt.Println("End ---- SizeXYZ - ", size)

		// offset += size + 2
		offset += size

		// fmt.Printf("WithoutOffset: %v - %q\n", offset, string(val))
		// fmt.Printf("WithOffset: %v - %q\n", offset, string(val[offset:]))

		res[i] = v
	}

	return NewArray(res), se + offset, nil
}

func parseError(input []byte) (*Error, int64, error) {
	strEndIdx := bytes.Index(input, []byte{'\r', '\n'})

	if strEndIdx == -1 {
		return nil, 0, fmt.Errorf("invalid error")
	}

	strValue := input[:strEndIdx]

	return NewError(string(strValue)), int64(strEndIdx) + 2, nil
}

func parseInt(input []byte) (*Int, int64, error) {
	strEndIdx := bytes.Index(input, []byte{'\r', '\n'})

	if strEndIdx == -1 {
		return nil, 0, fmt.Errorf("invalid integer")
	}

	strValue := input[:strEndIdx]

	val, err := strconv.ParseInt(string(strValue), 10, 64)
	if err != nil {
		return nil, 0, err
	}

	return &Int{value: val}, int64(strEndIdx) + 2, nil
}

func Parse(resp []byte) (DataType, int64, error) {
	if len(resp) < 1 {
		return nil, 0, ErrInvalidRespDataType
	}

	dataTypeId := resp[0]

	fmt.Printf("DataType: %q - %q\n", string(dataTypeId), string(resp))

	if len(resp) >= 5 {
		if ok, s := isNull(resp[:5]); ok {
			return NewNull(rune(dataTypeId)), s, nil
		}
	}

	remaining := resp[1:]

	switch dataTypeId {

	case SimpleStringPrefix:
		v, s, e := parseSimpleString(remaining)
		return v, s + 1, e

	case ErrorPrefix:
		v, s, e := parseError(remaining)
		return v, s + 1, e

	case IntPrefix:
		v, s, e := parseInt(remaining)
		return v, s + 1, e

	case BulkStringPrefix:
		v, s, e := parseBulkString(remaining)
		return v, s + 1, e

	case ArrayPrefix:
		v, s, e := parseArray(remaining)
		return v, s + 1, e

	default:
		return nil, 0, ErrInvalidRespDataType
	}
}
