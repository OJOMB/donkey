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

func (l *Map) Type() Type { return TypeMap }

func (l *Map) Inspect() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("{")
	i := 0
	for _, pair := range l.Pairs {
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

		if i < len(l.Pairs)-1 {
			strBuilder.WriteString(", ")
		}

		i++
	}

	strBuilder.WriteString("}")

	return strBuilder.String()
}
