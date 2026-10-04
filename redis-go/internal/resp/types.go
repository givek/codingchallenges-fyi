package resp

import (
	"errors"
)

var ErrInvalidRespDataType = errors.New("invalid resp data type")

const (
	SimpleStringPrefix = '+'
	ErrorPrefix        = '-'
	IntPrefix          = ':'
	BulkStringPrefix   = '$'
	ArrayPrefix        = '*'
)

type DataType interface {
	Identifier() rune
}

type Error struct {
	message string
}

func NewError(message string) *Error {
	return &Error{
		message: message,
	}
}

type Int struct {
	value int64
}

func NewInt(value int64) *Int {
	return &Int{value: value}
}

func (n *Int) Identifier() rune {
	return IntPrefix
}

func (ss *Error) Identifier() rune {
	return ErrorPrefix
}

func (err *Error) Message() string {
	return err.message
}

type SimpleString struct {
	value string
}

func NewSimpleString(value string) *SimpleString {
	return &SimpleString{
		value: value,
	}
}

func (ss *SimpleString) Identifier() rune {
	return SimpleStringPrefix
}

func (ss *SimpleString) Value() string {
	return ss.value
}

type BulkString struct {
	value string
}

func NewBulkString(value string) *BulkString {
	return &BulkString{
		value: value,
	}
}

func (ss *BulkString) Identifier() rune {
	return BulkStringPrefix
}

func (ss *BulkString) Value() string {
	return ss.value
}

type Array struct {
	value []DataType
}

func NewArray(value []DataType) *Array {
	return &Array{
		value: value,
	}
}

func (ss *Array) Identifier() rune {
	return ArrayPrefix
}

func (ss *Array) Value() []DataType {
	return ss.value
}

type Null struct {
	id rune
}

func NewNull(id rune) *Null {
	return &Null{id: id}
}

func (n *Null) Identifier() rune {
	return n.id
}
