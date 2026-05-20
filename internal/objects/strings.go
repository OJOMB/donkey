package objects

import "hash/fnv"

// String represents a string value in the Donkey programming language. It is used to store and manipulate text.
type String struct {
	Value string
}

// Type returns the type of the String object, which is "STRING".
func (s *String) Type() Type { return TypeString }

// Inspect returns a string representation of the String object, which is the value of the string itself.
func (s *String) Inspect() string { return s.Value }

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
