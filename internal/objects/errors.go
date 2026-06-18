package objects

import "fmt"

// ErrorValue is a struct that represents an error value in the Donkey programming language.
// It contains a single field, Message, which is a string that describes the error.
type ErrorValue struct {
	Message string
}

// Type returns the type of the ErrorValue object, which is always TypeError.
func (ev *ErrorValue) Type() Type {
	return TypeErrorValue
}

// Inspect returns a string representation of the ErrorValue object, which includes the error message.
func (ev *ErrorValue) Inspect() string {
	return "ERROR: " + ev.Message
}

var (
	// ErrNoIndices is an error that indicates that no indices were provided for a list operation.
	ErrNoIndices = fmt.Errorf("no indices provided")
	// ErrIndexOutOfRange is an error that indicates that an index provided for a list operation is out of range.
	ErrIndexOutOfRange = fmt.Errorf("index out of range")
	// ErrKeyNotFound is an error that indicates that a key was not found in a map.
	ErrKeyNotFound = fmt.Errorf("key not found")
	// ErrKeyNotHashable is an error that indicates that a key provided for a map operation is not hashable.
	ErrKeyNotHashable = fmt.Errorf("key is not hashable")
	// ErrElementNotIndexable is an error that indicates that an element at a specified index is not indexable, but was expected to be one.
	ErrElementNotIndexable = fmt.Errorf("element at index is not indexable")
)
