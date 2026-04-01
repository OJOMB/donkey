package tokens

const (
	// TokenTypeIllegal represents an illegal token.
	TokenTypeIllegal TokenType = "ILLEGAL"
	// TokenTypeEOF represents the end of file token.
	TokenTypeEOF TokenType = "EOF"

	// TokenTypeIdent represents an identifier token.
	TokenTypeIdent TokenType = "IDENT"
	// TokenTypeInt represents an integer token.
	TokenTypeInt TokenType = "INT"
	// TokenTypeString represents a string token.
	TokenTypeString TokenType = "STRING"

	///////////////////////////
	// Symbols and operators //
	///////////////////////////

	// TokenTypeAssign represents the assignment operator token.
	TokenTypeAssign TokenType = "="
	// TokenTypePlus represents the addition operator token.
	TokenTypePlus TokenType = "+"
	// TokenTypeMinus represents the subtraction operator token.
	TokenTypeMinus TokenType = "-"
	// TokenTypeBang represents the logical not operator token.
	TokenTypeBang TokenType = "!"
	// TokenTypeAsterisk represents the multiplication operator token.
	TokenTypeAsterisk TokenType = "*"
	// TokenTypeForwardSlash represents the division operator token.
	TokenTypeForwardSlash TokenType = "/"
	// TokenTypeComma represents the comma token.
	TokenTypeComma TokenType = ","
	// TokenTypeSemicolon represents the semicolon token.
	TokenTypeSemicolon TokenType = ";"

	// TokenTypeLParen represents the left parenthesis token.
	TokenTypeLParen TokenType = "("
	// TokenTypeRParen represents the right parenthesis token.
	TokenTypeRParen TokenType = ")"
	// TokenTypeLBrace represents the left brace token.
	TokenTypeLBrace TokenType = "{"
	// TokenTypeRBrace represents the right brace token.
	TokenTypeRBrace TokenType = "}"
	// TokenTypeEq represents the equality operator token.
	TokenTypeEq TokenType = "=="
	// TokenTypeNotEq represents the not equal operator token.
	TokenTypeNotEq TokenType = "!="
	// TokenTypeLT represents the less than operator token.
	TokenTypeLT TokenType = "<"
	// TokenTypeLTEQ represents the less than or equal to operator token.
	TokenTypeLTEQ TokenType = "<="
	// TokenTypeGTEQ represents the greater than or equal to operator token.
	TokenTypeGTEQ TokenType = ">="
	// TokenTypeGT represents the greater than operator token.
	TokenTypeGT TokenType = ">"
	// TokenTypeSpeechMarks represents the speech marks token used for opening string literals.
	TokenTypeSpeechMarks TokenType = `"`

	///////////////
	// keywords //
	/////////////

	// TokenTypeFunction represents the 'fn' keyword token.
	TokenTypeFunction TokenType = "FUNCTION"
	// TokenTypeBinder represents the 'var' keyword token.
	TokenTypeBinder TokenType = "VAR"
	// TokenTypeTrue represents the boolean value true
	TokenTypeTrue TokenType = "TRUE"
	// TokenTypeFalse represents the boolean value false
	TokenTypeFalse TokenType = "FALSE"
	// TokenTypeIf represents the control flow keyword if.
	TokenTypeIf TokenType = "IF"
	// TokenTypeElif represents the control flow keyword elif.
	TokenTypeElif TokenType = "ELIF"
	// TokenTypeElse represents the control flow keyword else.
	TokenTypeElse TokenType = "ELSE"
	// TokenTypeReturn represents the control flow keyword return.
	TokenTypeReturn TokenType = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     TokenTypeFunction,
	"var":    TokenTypeBinder,
	"if":     TokenTypeIf,
	"elif":   TokenTypeElif,
	"else":   TokenTypeElse,
	"return": TokenTypeReturn,
	"true":   TokenTypeTrue,
	"false":  TokenTypeFalse,
}

// LookupIdent checks if the given identifier is a keyword and returns the appropriate token type.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return TokenTypeIdent
}

// TokenType represents the type of a token.
type TokenType string

// Token represents a lexical input token.
type Token struct {
	Type   TokenType
	Lexeme string
}

func New(tokType TokenType, tokLit string) Token {
	return Token{
		Type:   tokType,
		Lexeme: tokLit,
	}
}
