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
	// TokenTypeLT represents the less than operator token.
	TokenTypeLT TokenType = "<"
	// TokenTypeGT represents the less than operator token.
	TokenTypeGT TokenType = ">"
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

	///////////////
	// keywords //
	/////////////

	// TokenTypeFunction represents the 'fn' keyword token.
	TokenTypeFunction TokenType = "FUNCTION"
	// TokenTypeLet represents the 'let' keyword token.
	TokenTypeLet TokenType = "LET"
	// TokenTypeTrue represents the boolean value true
	TokenTypeTrue TokenType = "true"
	// TokenTypeFalse represents the boolean value true
	TokenTypeFalse TokenType = "false"
	// TokenTypeIf represents the control flow keyword if.
	TokenTypeIf TokenType = "if"
	// TokenTypeElif represents the control flow keyword elif.
	TokenTypeElif TokenType = "elif"
	// TokenTypeElse represents the control flow keyword else.
	TokenTypeElse TokenType = "else"
	// TokenTypeReturn represents the control flow keyword return.
	TokenTypeReturn TokenType = "return"
)

var keywords = map[string]TokenType{
	"fn":     TokenTypeFunction,
	"let":    TokenTypeLet,
	"if":     TokenTypeIf,
	"elif":   TokenTypeElif,
	"else":   TokenTypeElse,
	"return": TokenTypeReturn,
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
