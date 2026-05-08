package objects

import "fmt"

type BuiltinFunction struct {
	Fn             Function
	Name           string
	Implementation func(args ...Object) Object
}

func (b *BuiltinFunction) Type() Type { return TypeBuiltin }

func (b *BuiltinFunction) Inspect() string {
	return fmt.Sprintf("builtin function: %s", b.Name)
}
