package objects

import (
	"strings"
)

// List represents a list object in the Donkey programming language. It contains a slice of objects that are the elements of the list.
type List struct {
	Elements []Object
}

// Type returns the type of the List object, which is always TypeList.
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

// Len returns the number of elements in the list.
func (l *List) Len() int {
	return len(l.Elements)
}

// IsEmpty returns true if the list has no elements, and false otherwise.
func (l *List) IsEmpty() bool {
	return l.Len() == 0
}

// Get retrieves an element from the list at the specified index. If the index is out of range, it returns an error.
func (l *List) Get(idx int) (Object, error) {
	if l.IsEmpty() || idx < 0 || idx >= l.Len() {
		return nil, ErrIndexOutOfRange
	}

	return l.Elements[idx], nil
}

// Set sets an element in the list at the specified index to the provided value. If the index is out of range, it returns an error.
func (l *List) Set(idx int, value Object) error {
	if l.IsEmpty() || idx < 0 || idx >= l.Len() {
		return ErrIndexOutOfRange
	}

	l.Elements[idx] = value

	return nil
}

// GetMultiDimensional retrieves an element from the list at the specified indices. It supports nested lists, allowing for multi-dimensional access.
func (l *List) GetMultiDimensional(indices []Object) (Object, error) {
	if len(indices) < 1 {
		return nil, ErrNoIndices
	}

	// Get the first index and check if it's an integer
	indexObj := indices[0]
	intIndex, ok := indexObj.(*Integer)
	if !ok {
		return nil, ErrElementNotIndexable
	}

	value, err := l.Get(int(intIndex.Value))
	if err != nil {
		return nil, err
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

// Set sets an element in the list at the specified indices to the provided value. It supports nested lists, allowing for multi-dimensional access.
func (l *List) SetMultiDimensional(indices []Object, value Object) error {
	if len(indices) < 1 {
		return ErrNoIndices
	}

	indexObj := indices[0]
	intIndex, ok := indexObj.(*Integer)
	if !ok {
		return ErrElementNotIndexable
	}

	idx := int(intIndex.Value)
	if l.IsEmpty() || idx < 0 || idx >= l.Len() {
		return ErrIndexOutOfRange
	}

	if len(indices) == 1 {
		return l.Set(idx, value)
	}

	nestedList, ok := l.Elements[idx].(*List)
	if !ok {
		return ErrElementNotIndexable
	}

	return nestedList.SetMultiDimensional(indices[1:], value)
}
