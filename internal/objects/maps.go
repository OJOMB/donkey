package objects

import "strings"

// HashKey represents a key in a hashmap. It has two fields: Type, which is the type of the key, and Value, which is the value of the key.
type HashKey struct {
	Type  Type
	Value uint64
}

// HashPair represents a key-value pair in a hashmap. It has two fields: Key, which is the key of the pair, and Value, which is the value of the pair.
type HashPair struct {
	Key   Object
	Value Object
}

// Map represents a hashmap object in the Donkey programming language.
// It has a single field, Elements, which is a map that maps keys to values.
// Both keys and values are of type Object, which is the interface that all objects in the Donkey programming language must implement.
type Map struct {
	Pairs map[HashKey]HashPair
}

func (m *Map) Type() Type { return TypeMap }

func (m *Map) Inspect() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("{")
	i := 0
	for _, pair := range m.Pairs {
		strBuilder.WriteString(`"`)
		strBuilder.WriteString(pair.Key.Inspect())
		strBuilder.WriteString(`": `)

		if pair.Value.Type() == TypeString {
			strBuilder.WriteString(`"`)
		}

		strBuilder.WriteString(pair.Value.Inspect())
		if pair.Value.Type() == TypeString {
			strBuilder.WriteString(`"`)
		}

		if i < len(m.Pairs)-1 {
			strBuilder.WriteString(", ")
		}

		i++
	}

	strBuilder.WriteString("}")

	return strBuilder.String()
}

// Len returns the number of key-value pairs in the map.
func (m *Map) Len() int {
	return len(m.Pairs)
}

// IsEmpty returns true if the map has no key-value pairs, and false otherwise.
func (m *Map) IsEmpty() bool {
	return m.Len() == 0
}

// Get retrieves the value associated with the given key in the map.
func (m *Map) Get(key ObjectKey) (Object, bool) {
	if key == nil {
		return nil, false
	}

	pair, ok := m.Pairs[key.HashKey()]
	if !ok {
		return nil, false
	}

	return pair.Value, true
}

// Set sets the value associated with the given key in the map.
func (m *Map) Set(key ObjectKey, value Object) {
	if key == nil {
		return
	}

	m.Pairs[key.HashKey()] = HashPair{Key: key, Value: value}
}

// Delete removes the key-value pair associated with the given key from the map.
func (m *Map) Delete(key ObjectKey) {
	if key == nil {
		return
	}

	delete(m.Pairs, key.HashKey())
}

// GetMultiDimensional retrieves a value from the map at the specified keys. It supports nested maps, allowing for multi-dimensional access.
func (m *Map) GetMultiDimensional(indices []Object) (Object, error) {
	if len(indices) < 1 {
		return nil, ErrNoIndices
	}

	key, ok := indices[0].(ObjectKey)
	if !ok {
		return nil, ErrKeyNotHashable
	}

	value, ok := m.Get(key)
	if !ok {
		return nil, ErrKeyNotFound
	}

	if len(indices) == 1 {
		return value, nil
	}

	switch v := value.(type) {
	case Indexable:
		return v.GetMultiDimensional(indices[1:])
	default:
		return nil, ErrElementNotIndexable
	}
}

// SetMultiDimensional sets a value in the map at the specified keys. It supports nested maps, allowing for multi-dimensional access.
func (m *Map) SetMultiDimensional(indices []Object, value Object) error {
	if len(indices) < 1 {
		return ErrNoIndices
	}

	key, ok := indices[0].(ObjectKey)
	if !ok {
		return ErrKeyNotHashable
	}

	if len(indices) == 1 {
		m.Set(key, value)
		return nil
	}

	existingValue, ok := m.Get(key)
	if !ok {
		return ErrKeyNotFound
	}

	switch v := existingValue.(type) {
	case Indexable:
		err := v.SetMultiDimensional(indices[1:], value)
		if err != nil {
			return err
		}
	default:
		return ErrElementNotIndexable
	}

	m.Set(key, existingValue)
	return nil
}
