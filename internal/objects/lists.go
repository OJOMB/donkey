package objects

import "strings"

// List represents a list object in the Donkey programming language. It contains a slice of objects that are the elements of the list.
type List struct {
	Elements []Object
}

func (l *List) Type() Type { return TypeList }

func (l *List) Inspect() string {
	var out strings.Builder
	out.WriteString("[")
	for i, elem := range l.Elements {
		if i > 0 {
			out.WriteString(", ")
		}

		out.WriteString(elem.Inspect())
	}

	out.WriteString("]")

	return out.String()
}
