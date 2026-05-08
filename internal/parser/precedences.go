package parser

import "github.com/OJOMB/donkey/internal/tokens"

// I have simply copied the precedences used in the C programming language.
const (
	precedenceLowest         int = iota
	precedenceLogicalOr          // ||
	precedenceLogicalAnd         // &&
	precedenceBitwiseOr          // |
	precedenceBitwiseXor         // ^
	precedenceBitwiseAnd         // &
	precedenceEquals             // == !=
	precedenceLessGreater        // < > <= >=
	precedenceShift              // << >>
	precedenceAdditive           // + -
	precedenceMultiplicative     // * / %
	precedenceExponentiation     // **
	precedencePrefix             // -X !X ~X
	precedenceCall               // myFunc(X)
)

var precedences = map[tokens.Type]int{
	tokens.TypeLogicalOr:  precedenceLogicalOr,
	tokens.TypeLogicalAnd: precedenceLogicalAnd,

	tokens.TypeBitwiseOr:  precedenceBitwiseOr,
	tokens.TypeBitwiseXor: precedenceBitwiseXor,
	tokens.TypeBitwiseAnd: precedenceBitwiseAnd,

	tokens.TypeEQ:    precedenceEquals,
	tokens.TypeNotEQ: precedenceEquals,

	tokens.TypeLT:   precedenceLessGreater,
	tokens.TypeGT:   precedenceLessGreater,
	tokens.TypeLTEQ: precedenceLessGreater,
	tokens.TypeGTEQ: precedenceLessGreater,

	tokens.TypePlus:  precedenceAdditive,
	tokens.TypeMinus: precedenceAdditive,

	tokens.TypeForwardSlash: precedenceMultiplicative,
	tokens.TypeAsterisk:     precedenceMultiplicative,
	tokens.TypePercent:      precedenceMultiplicative,

	tokens.TypeExponent: precedenceExponentiation,

	tokens.TypeBang:   precedencePrefix,
	tokens.TypeLParen: precedenceCall,
}
