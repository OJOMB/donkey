package objects

import "strings"

// Map represents a hashmap object in the Donkey programming language.
// It has a single field, Elements, which is a map that maps keys to values.
// Both keys and values are of type Object, which is the interface that all objects in the Donkey programming language must implement.
type Map struct {
	Elements map[Object]Object
}

func (l *Map) Type() Type { return TypeMaps }

func (l *Map) Inspect() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("{")
	i := 0
	for key, value := range l.Elements {
		strBuilder.WriteString(key.Inspect())
		strBuilder.WriteString(": ")
		strBuilder.WriteString(value.Inspect())
		if i < len(l.Elements)-1 {
			strBuilder.WriteString(", ")
		}

		i++
	}

	strBuilder.WriteString("}")

	return strBuilder.String()
}
