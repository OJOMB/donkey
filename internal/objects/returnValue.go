package objects

// ReturnValue is a struct that represents a return value in the Donkey programming language.
// It contains a single field, Value, which is an Object that represents the value being returned.
type ReturnValue struct {
	Value Object
}

// Type returns the type of the ReturnValue object, which is always TypeReturnValue.
func (rv *ReturnValue) Type() Type {
	return TypeReturnValue
}

// Inspect returns a string representation of the ReturnValue object, which includes the string representation of the value being returned.
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
