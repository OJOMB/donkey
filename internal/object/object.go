package object

// ObjectType is a string that represents the type of an object in the Donkey programming language.
type ObjectType string

const (
	// Object types

	// TypeInteger represents the type of an integer object
	TypeInteger = "INTEGER"
	// TypeBoolean represents the type of a boolean object
	TypeBoolean = "BOOLEAN"
	// TypeNowt represents the type of a nowt object (the equivalent of null or nil)
	TypeNowt = "NOWT"
)

// Object is the interface that all objects in the Donkey programming language must implement.
// It defines two methods: Type, which returns the type of the object, and Inspect, which returns a string representation of the object.
type Object interface {
	Type() ObjectType
	Inspect() string
}
