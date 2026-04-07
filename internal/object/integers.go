package object

import "fmt"

// Integer represents a whole number value in the Donkey programming language. It is used to store and manipulate numeric data.
type Integer struct {
	Value int64
}

// Type returns the type of the Integer object, which is "INTEGER".
func (i *Integer) Type() ObjectType { return TypeInteger }

// Inspect returns a string representation of the Integer object, which is the string form of the integer value.
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
