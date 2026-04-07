package object

// Nowt represents the absence of a value, similar to null or nil in other languages. It is used to indicate that a variable or expression does not have a meaningful value.
type Nowt struct{}

func (n *Nowt) Type() ObjectType { return TypeNowt }
func (n *Nowt) Inspect() string  { return "nowt" }
