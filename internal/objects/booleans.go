package objects

import "fmt"

// Boolean represents a true or false value in the Donkey programming language. It is used to store and manipulate logical values.
type Boolean struct {
	Value bool
}

// Type returns the type of the Boolean object, which is "BOOLEAN".
func (b *Boolean) Type() Type { return TypeBoolean }

// Inspect returns a string representation of the Boolean object, which is "true" or "false" depending on the value of the Boolean.
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}
