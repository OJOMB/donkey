package ast

import "github.com/OJOMB/monkey/internal/tokens"

type LiteralInteger struct {
	Token tokens.Token
	Value int
}

func (li *LiteralInteger) expressionNode()     {}
func (li *LiteralInteger) TokenLexeme() string { return li.Token.Lexeme }

func (li *LiteralInteger) String() string {
	return li.Token.Lexeme
}

type LiteralString struct {
	Token tokens.Token
	Value string
}

func (ls *LiteralString) expressionNode()     {}
func (ls *LiteralString) TokenLexeme() string { return ls.Token.Lexeme }

func (ls *LiteralString) String() string {
	return ls.Token.Lexeme
}

type LiteralBoolean struct {
	Token tokens.Token
	Value bool
}

func (lb *LiteralBoolean) expressionNode()     {}
func (lb *LiteralBoolean) TokenLexeme() string { return lb.Token.Lexeme }

func (lb *LiteralBoolean) String() string {
	return lb.Token.Lexeme
}
