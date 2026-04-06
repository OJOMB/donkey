package parser

import "github.com/OJOMB/donkey/internal/tokens"

const (
	// Precedence levels for parsing expressions

	// precedenceLowest is the lowest precedence level, used for expressions that don't have any operators.
	precedenceLowest int = iota
	// precedenceEquals is the precedence level for the equality operator (==).
	precedenceEquals
	// precedenceLessGreater is the precedence level for the less than and greater than operators (< and >).
	precedenceLessGreater
	// precedenceAdditive is the precedence level for addition and subtraction operators (+ and -).
	precedenceAdditive
	// precedenceMultiplicative is the precedence level for multiplication and division operators (* and /).
	precedenceMultiplicative
	// precedenceExponentiation is the precedence level for the exponentiation operator (^).
	precedenceExponentiation
	// precedencePrefix is the precedence level for prefix operators, such as -X or !X.
	precedencePrefix
	// precedenceCall is the precedence level for function call expressions, such as myFunction(X).
	precedenceCall
)

var precedences = map[tokens.Type]int{
	tokens.TypeEq:           precedenceEquals,
	tokens.TypeNotEq:        precedenceEquals,
	tokens.TypeLT:           precedenceLessGreater,
	tokens.TypeGT:           precedenceLessGreater,
	tokens.TypeBang:         precedencePrefix,
	tokens.TypePlus:         precedenceAdditive,
	tokens.TypeMinus:        precedenceAdditive,
	tokens.TypeForwardSlash: precedenceMultiplicative,
	tokens.TypeAsterisk:     precedenceMultiplicative,
	tokens.TypePercent:      precedenceMultiplicative,
	tokens.TypeCaret:        precedenceExponentiation,
	tokens.TypeLParen:       precedenceCall,
}
