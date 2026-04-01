package lexer

import "github.com/OJOMB/donkey/internal/tokens"

type Lexer struct {
	// input is the string input that the lexer will tokenize.
	input string
	// position is the current position in the input (points to current char)
	position int
	// readPosition is the current reading position in input (after current char)
	readPosition int
	// ch is the current char literal under examination (ASCII single byte chars - Unicode not supported)
	ch byte
}

// New creates a new Lexer instance with the given input string and returns a pointer to it. The lexer is initialised and ready to return tokens when NextToken is called.
func New(input string) *Lexer {
	l := &Lexer{input: input}

	// read the first char so that the lexer is ready to return the first token when NextToken is called
	l.readChar()

	return l
}

// IsInitialised returns true if the lexer has been initialised with an input string and is ready to return tokens, and false otherwise.
func (l *Lexer) IsInitialised() bool {
	return l != nil && (l.position != 0 || l.readPosition > 1 || l.ch != 0)
}

func (l *Lexer) NextToken() tokens.Token {
	l.skipWhitespace()

	var tok tokens.Token
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.New(tokens.TypeEq, "==")
			break
		}

		tok = tokens.New(tokens.TypeAssign, "=")
	case ';':
		tok = tokens.New(tokens.TypeSemicolon, ";")
	case '(':
		tok = tokens.New(tokens.TypeLParen, "(")
	case ')':
		tok = tokens.New(tokens.TypeRParen, ")")
	case ',':
		tok = tokens.New(tokens.TypeComma, ",")
	case '+':
		tok = tokens.New(tokens.TypePlus, "+")
	case '-':
		tok = tokens.New(tokens.TypeMinus, "-")
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.New(tokens.TypeLTEQ, "<=")
			break
		}

		tok = tokens.New(tokens.TypeLT, "<")
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.New(tokens.TypeGTEQ, ">=")
			break
		}

		tok = tokens.New(tokens.TypeGT, ">")
	case '*':
		tok = tokens.New(tokens.TypeAsterisk, "*")
	case '/':
		tok = tokens.New(tokens.TypeForwardSlash, "/")
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.New(tokens.TypeNotEq, "!=")
			break
		}

		tok = tokens.New(tokens.TypeBang, "!")
	case '{':
		tok = tokens.New(tokens.TypeLBrace, "{")
	case '}':
		tok = tokens.New(tokens.TypeRBrace, "}")
	case '"':
		// we have encountered the opening speech marks of a string literal, so we want to read the whole string literal and return it as a token
		tok.Type = tokens.TypeString
		tok.Lexeme = l.readString()
		return tok

	case 0:
		// if the current char is ASCII NUL then we have reached the end of the input and we return an EOF token
		tok = tokens.New(tokens.TypeEOF, "")
	default:
		if l.isLetter(l.ch) {
			// if the current char is a letter then we want to read the whole identifier and return it as a token
			// we check the first byte is a letter rather than an alphanumeric since idents cannot start with a digit
			tok.Lexeme = l.readIdentifier()
			tok.Type = tokens.LookupIdent(tok.Lexeme)
			return tok
		} else if l.isDigit(l.ch) {
			// if the current char is a digit then we want to read the whole number and return it as a token
			tok.Type = tokens.TypeInt
			tok.Lexeme = l.readNumber()
			return tok
		} else {
			tok = tokens.New(tokens.TypeIllegal, string(l.ch))
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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
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

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	strLiteral := l.input[position:l.position]

	// we have now read the closing speech marks of the string literal so we want to advance our position so that the next call to NextToken will give us the next token after the string literal
	l.readChar()

	return strLiteral
}
