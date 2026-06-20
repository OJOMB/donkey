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

// GetMultiDimensional retrieves a character from the string at the specified index. If the index is out of range, it returns an error.
func (s *String) GetMultiDimensional(indices []Object) (Object, error) {
	if len(indices) != 1 {
		return nil, ErrNoIndices
	}

	indexObj := indices[0]
	intIndex, ok := indexObj.(*Integer)
	if !ok {
		return nil, ErrElementNotIndexable
	}

	idx := intIndex.Value
	if idx < 0 || idx >= len(s.Value) {
		return nil, ErrIndexOutOfRange
	}

	// cast string to rune slice to handle multi-byte characters correctly
	runes := []rune(s.Value)
	return &String{Value: string(runes[idx])}, nil
}

// SetMultiDimensional sets a character in the string at the specified index to the provided value. If the index is out of range, it returns an error.
func (s *String) SetMultiDimensional(indices []Object, value Object) error {
	if len(indices) != 1 {
		return ErrNoIndices
	}

	indexObj := indices[0]
	intIndex, ok := indexObj.(*Integer)
	if !ok {
		return ErrElementNotIndexable
	}

	idx := int(intIndex.Value)
	if idx < 0 || idx >= len(s.Value) {
		return ErrIndexOutOfRange
	}

	valueStr, ok := value.(*String)
	if !ok {
		return ErrElementNotIndexable
	}

	// cast string to rune slice to handle multi-byte characters correctly
	runes := []rune(s.Value)
	runes[idx] = []rune(valueStr.Value)[0]
	s.Value = string(runes)

	return nil
}
