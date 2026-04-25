package tokens

const (
	// TypeIllegal represents an illegal token.
	TypeIllegal Type = "ILLEGAL"
	// TypeEOF represents the end of file token.
	TypeEOF Type = "EOF"
	// TypeIdent represents an identifier token.
	TypeIdent Type = "IDENT"
	// TypeInt represents an integer token.
	TypeInt Type = "INT"
	// TypeString represents a string token.
	TypeString Type = "STRING"

	///////////////////////////
	// Symbols and operators //
	///////////////////////////

	// TypeAssign represents the assignment operator token.
	TypeAssign Type = "="
	// TypePlus represents the addition operator token.
	TypePlus Type = "+"
	// TypeMinus represents the subtraction operator token.
	TypeMinus Type = "-"
	// TypeBang represents the logical not operator token.
	TypeBang Type = "!"
	// TypeAsterisk represents the multiplication operator token.
	TypeAsterisk Type = "*"
	// TypeForwardSlash represents the division operator token.
	TypeForwardSlash Type = "/"
	// TypePercent represents the modulus operator token.
	TypePercent Type = "%"
	// TypeCaret represents the exponentiation operator token.
	TypeCaret Type = "^"
	// TypeComma represents the comma token.
	TypeComma Type = ","
	// TypeSemicolon represents the semicolon token.
	TypeSemicolon Type = ";"

	// TypeLParen represents the left parenthesis token.
	TypeLParen Type = "("
	// TypeRParen represents the right parenthesis token.
	TypeRParen Type = ")"
	// TypeLBrace represents the left brace token.
	TypeLBrace Type = "{"
	// TypeRBrace represents the right brace token.
	TypeRBrace Type = "}"
	// TypeEq represents the equality operator token.
	TypeEq Type = "=="
	// TypeNotEq represents the not equal operator token.
	TypeNotEq Type = "!="
	// TypeLT represents the less than operator token.
	TypeLT Type = "<"
	// TypeLTEQ represents the less than or equal to operator token.
	TypeLTEQ Type = "<="
	// TypeGTEQ represents the greater than or equal to operator token.
	TypeGTEQ Type = ">="
	// TypeGT represents the greater than operator token.
	TypeGT Type = ">"
	// TypeSpeechMarks represents the speech marks token used for opening string literals.
	TypeSpeechMarks Type = `"`
	// TypeBitwiseAnd represents the bitwise and operator token.
	TypeBitwiseAnd Type = "&"
	// TypeBitwiseOr represents the bitwise or operator token.
	TypeBitwiseOr Type = "|"
	// TypeLogicalAnd represents the logical and operator token.
	TypeLogicalAnd Type = "&&"
	// TypeLogicalOr represents the logical or operator token.
	TypeLogicalOr Type = "||"

	///////////////
	// keywords //
	/////////////

	// TypeFunction represents the 'fn' keyword token.
	TypeFunction Type = "FUNCTION"
	// TypeBind represents the 'var' keyword token.
	TypeBind Type = "VAR"
	// TypeTrue represents the boolean value true
	TypeTrue Type = "TRUE"
	// TypeFalse represents the boolean value false
	TypeFalse Type = "FALSE"
	// TypeIf represents the control flow keyword if.
	TypeIf Type = "IF"
	// TypeElif represents the control flow keyword elif.
	TypeElif Type = "ELIF"
	// TypeElse represents the control flow keyword else.
	TypeElse Type = "ELSE"
	// TypeReturn represents the control flow keyword return.
	TypeReturn Type = "RETURN"
)

var keywords = map[string]Type{
	"fn":     TypeFunction,
	"var":    TypeBind,
	"if":     TypeIf,
	"elif":   TypeElif,
	"else":   TypeElse,
	"return": TypeReturn,
	"true":   TypeTrue,
	"false":  TypeFalse,
}

// LookupIdent checks if the given identifier is a keyword and returns the appropriate token type.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return TypeIdent
}

// Type represents the type of a token.
type Type string

// Token represents a lexical input token.
type Token struct {
	Type   Type
	Lexeme string
}

// New creates a new Token with the given type and literal value.
func New(tokenType Type, tokenLiteral string) Token {
	return Token{
		Type:   tokenType,
		Lexeme: tokenLiteral,
	}
}
