package objects

// Type is a string that represents the type of an object in the Donkey programming language.
type Type string

const (
	// Object types

	// TypeInteger represents the type of an integer object
	TypeInteger Type = "INTEGER"
	// TypeBoolean represents the type of a boolean object
	TypeBoolean Type = "BOOLEAN"
	// TypeString represents the type of a string object
	TypeString Type = "STRING"
	// TypeFunction represents the type of a function object
	TypeFunction Type = "FUNCTION"
	// TypeNowt represents the type of a nowt object (the equivalent of null or nil)
	TypeNowt Type = "NOWT"
	// TypeReturnValue represents the type of a return value object
	TypeReturnValue Type = "RETURN_VALUE"
)

// Object is the interface that all objects in the Donkey programming language must implement.
// It defines two methods: Type, which returns the type of the object, and Inspect, which returns a string representation of the object.
type Object interface {
	// Type returns the type of the object as an object Type.
	Type() Type
	// Inspect returns a string representation of the object, which is used for debugging and error messages.
	Inspect() string
}
