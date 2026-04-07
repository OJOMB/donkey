package object

import "fmt"

// Boolean represents a true or false value in the Donkey programming language. It is used to store and manipulate logical values.
type Boolean struct {
	Value bool
}

// Type returns the type of the Boolean object, which is "BOOLEAN".
func (b *Boolean) Type() ObjectType { return TypeBoolean }

// Inspect returns a string representation of the Boolean object, which is "true" or "false" depending on the value of the Boolean.
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
