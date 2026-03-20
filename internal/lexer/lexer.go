package lexer

import "github.com/OJOMB/monkey/internal/lexer/tokens"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char literal under examination (ASCII single byte chars - Unicode not supported)
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() tokens.Token {
	l.skipWhitespace()

	var tok tokens.Token
	switch l.ch {
	case '=':
		tok = tokens.New(tokens.TokenTypeAssign, "=")
	case ';':
		tok = tokens.New(tokens.TokenTypeSemicolon, ";")
	case '(':
		tok = tokens.New(tokens.TokenTypeLParen, "(")
	case ')':
		tok = tokens.New(tokens.TokenTypeRParen, ")")
	case ',':
		tok = tokens.New(tokens.TokenTypeComma, ",")
	case '+':
		tok = tokens.New(tokens.TokenTypePlus, "+")
	case '-':
		tok = tokens.New(tokens.TokenTypeMinus, "-")
	case '<':
		tok = tokens.New(tokens.TokenTypeLT, "<")
	case '>':
		tok = tokens.New(tokens.TokenTypeGT, ">")
	case '*':
		tok = tokens.New(tokens.TokenTypeAsterisk, "*")
	case '/':
		tok = tokens.New(tokens.TokenTypeForwardSlash, "/")
	case '!':
		tok = tokens.New(tokens.TokenTypeBang, "!")
	case '{':
		tok = tokens.New(tokens.TokenTypeLBrace, "{")
	case '}':
		tok = tokens.New(tokens.TokenTypeRBrace, "}")
	case 0:
		tok = tokens.New(tokens.TokenTypeEOF, "")
	default:
		if l.isLetter(l.ch) {
			// if the current char is a letter then we want to read the whole identifier and return it as a token
			// we check the first byte is a letter rather than an alphanumeric since idents cannot start with a digit
			tok.Lexeme = l.readIdentifier()
			tok.Type = tokens.LookupIdent(tok.Lexeme)
			return tok
		} else if l.isDigit(l.ch) {
			// if the current char is a digit then we want to read the whole number and return it as a token
			tok.Type = tokens.TokenTypeInt
			tok.Lexeme = l.readNumber()
			return tok
		} else {
			tok = tokens.New(tokens.TokenTypeIllegal, string(l.ch))
		}
	}

	// advance our position in the input string so that the next call to NextToken will give us the next token
	l.readChar()

	return tok
}

// readChar is a helper method to give us the next char and advance our position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// if the read position has gone past the final input position we have finished lexing
		// we set the current char to ASCII NUL and return
		l.ch = 0
		return
	}

	l.ch = l.input[l.readPosition]
	l.position = l.readPosition

	// increment read position in anticipation of next call
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for l.isAlphanumeric(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for l.isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}
func (l *Lexer) isWhitespace() bool {
	white := l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r'
	return white
}

func (l *Lexer) skipWhitespace() {
	for l.isWhitespace() {
		l.readChar()
	}
}

func (l *Lexer) isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) isAlphanumeric(ch byte) bool {
	return l.isLetter(ch) || l.isDigit(ch)
}

func (l *Lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
