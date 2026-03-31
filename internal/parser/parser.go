package parser

import (
	"fmt"
	"strconv"

	"github.com/OJOMB/monkey/internal/ast"
	"github.com/OJOMB/monkey/internal/lexer"
	"github.com/OJOMB/monkey/internal/tokens"
	"github.com/OJOMB/monkey/pkg/logs"
)

const (
	// Precedence levels for parsing expressions
	_ int = iota
	// precedenceLowest is the lowest precedence level, used for expressions that don't have any operators.
	precedenceLowest
	// precedenceEquals is the precedence level for the equality operator (==).
	precedenceEquals
	// precedenceLessGreater is the precedence level for the less than and greater than operators (< and >).
	precedenceLessGreater
	// precedenceAdditive is the precedence level for addition and subtraction operators (+ and -).
	precedenceAdditive
	// precedenceMultiplicative is the precedence level for multiplication and division operators (* and /).
	precedenceMultiplicative
	// precedencePrefix is the precedence level for prefix operators, such as -X or !X.
	precedencePrefix
	// precedenceCall is the precedence level for function call expressions, such as myFunction(X).
	precedenceCall
)

var precedences = map[tokens.TokenType]int{
	tokens.TokenTypeEq:           precedenceEquals,
	tokens.TokenTypeNotEq:        precedenceEquals,
	tokens.TokenTypeLT:           precedenceLessGreater,
	tokens.TokenTypeGT:           precedenceLessGreater,
	tokens.TokenTypePlus:         precedenceAdditive,
	tokens.TokenTypeMinus:        precedenceAdditive,
	tokens.TokenTypeForwardSlash: precedenceMultiplicative,
	tokens.TokenTypeAsterisk:     precedenceMultiplicative,
	tokens.TokenTypeLParen:       precedenceCall,
}

type (
	parseFuncPrefix func() ast.Expression
	parseFuncInfix  func(expr ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	currToken tokens.Token
	peekToken tokens.Token

	Errors []string

	parseFuncsPrefix map[tokens.TokenType]parseFuncPrefix
	parseFuncsInfix  map[tokens.TokenType]parseFuncInfix

	logger logs.Logger
}

// New creates a new Parser instance with the given lexer and logger.
// It initializes the parser's state and registers the necessary parse functions for different token types.
// if logger is nil, the parser will use a null logger by default.
func New(l *lexer.Lexer, logger logs.Logger) (*Parser, error) {
	if logger == nil {
		// null logger to avoid nil pointer dereference
		logger = logs.NewNullSlogger()
	}

	if !l.IsInitialised() {
		return nil, ErrLexerUnitialized
	}

	p := &Parser{
		l:                l,
		parseFuncsPrefix: make(map[tokens.TokenType]parseFuncPrefix),
		parseFuncsInfix:  make(map[tokens.TokenType]parseFuncInfix),
		Errors:           make([]string, 0),
		logger:           logger,
	}

	// register prefix parse functions for different token types
	p.RegisterPrefix(tokens.TokenTypeIdent, p.parseExpressionIdentifier)
	p.RegisterPrefix(tokens.TokenTypeInt, p.parseExpressionLiteralInteger)
	p.RegisterPrefix(tokens.TokenTypeString, p.parseExpressionLiteralString)
	p.RegisterPrefix(tokens.TokenTypeTrue, p.parseExpressionLiteralBoolean)
	p.RegisterPrefix(tokens.TokenTypeFalse, p.parseExpressionLiteralBoolean)
	p.RegisterPrefix(tokens.TokenTypeBang, p.parseExpressionPrefix)
	p.RegisterPrefix(tokens.TokenTypeMinus, p.parseExpressionPrefix)
	p.RegisterPrefix(tokens.TokenTypeLParen, p.parseExpressionGrouped)

	// register infix parse functions for different token types
	p.RegisterInfix(tokens.TokenTypePlus, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TokenTypeMinus, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TokenTypeForwardSlash, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TokenTypeAsterisk, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TokenTypeEq, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TokenTypeNotEq, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TokenTypeLT, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TokenTypeGT, p.parseExpressionInfix)

	// Read two tokens, so currToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p, nil
}

func (p *Parser) RegisterPrefix(tokType tokens.TokenType, fn parseFuncPrefix) {
	p.parseFuncsPrefix[tokType] = fn
}

func (p *Parser) RegisterInfix(tokType tokens.TokenType, fn parseFuncInfix) {
	p.parseFuncsInfix[tokType] = fn
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prgrm := ast.NewProgram()
	for p.currToken.Type != tokens.TokenTypeEOF {
		if stmt := p.parseStatement(); stmt != nil {
			prgrm.Statements = append(prgrm.Statements, stmt)
		}

		p.nextToken()
	}

	return prgrm
}

func (p *Parser) parseStatement() ast.Statement {
	var stmt ast.Statement
	switch p.currToken.Type {
	case tokens.TokenTypeLet:
		return p.parseStatementLet()
	case tokens.TokenTypeReturn:
		return p.parseStatementReturn()
	// case tokens.TokenTypeIdent:
	// 	// TODO: can we add rebinding here? statements such as x = 42;
	// 	//
	// 	fallthrough
	default:
		// if the statement doesn't match any of the above types, we assume it's an expression statement and try to parse it as such
		// this would be something like foo(5 + 5); or 5 + 5; or "foobar";
		stmt = p.parseExpressionStatement()
	}

	return stmt
}

// parseStatementLet parses a let statement and returns an ast.LetStatement node.
func (p *Parser) parseStatementLet() *ast.StatementLet {
	// first token must be LET
	if p.currToken.Type != tokens.TokenTypeLet {
		return nil
	}

	var ls = &ast.StatementLet{Token: p.currToken}

	// next token must be ident
	if !p.expectPeek(tokens.TokenTypeIdent) {
		return nil
	}

	ls.Name = &ast.ExpressionIdentifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}

	// next token must be assign =
	if !p.expectPeek(tokens.TokenTypeAssign) {
		return nil
	}

	p.nextToken()

	// next we have the expression which will be the LetStatement.Value
	value := p.parseExpression(precedenceLowest)

	ls.Value = value

	return ls
}

// parseStatementReturn parses a return statement and returns an ast.ReturnStatement node.
func (p *Parser) parseStatementReturn() *ast.ReturnStatement {
	if p.currToken.Type != tokens.TokenTypeReturn {
		return nil
	}

	rs := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	rs.ReturnValue = p.parseExpression(precedenceLowest)

	return rs
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	if p.currToken.Type == tokens.TokenTypeSemicolon {
		return nil
	}

	stmt := &ast.StatementExpression{Token: p.currToken}
	stmt.Expression = p.parseExpression(precedenceLowest)

	return stmt
}

// parseExpression parses an expression based on the precedence of the current token and returns an ast.Expression node representing the parsed expression.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// first we need to find the prefix parse function for the current token type
	prefixFunc := p.parseFuncsPrefix[p.currToken.Type]
	if prefixFunc == nil {
		p.noPrefixParseFuncError(p.currToken.Type)
		return nil
	}

	leftExp := prefixFunc()

	// next check if the expression ends here or if we have more to parse
	// if the peek token is a semicolon, then we have reached the end of the expression and we can return the left expression
	// if the precedence of the peek token is higher than the current precedence, then we need to parse the infix expression
	// we continue to parse infix expressions until we reach a semicolon or a token with lower precedence than the current precedence
	for p.peekToken.Type != tokens.TokenTypeSemicolon && precedence < p.peekPrecedence() {
		infixFunc := p.parseFuncsInfix[p.peekToken.Type]
		if infixFunc == nil {
			// TODO should this not be an error instead of just returning the left expression?
			// in this case we have a valid left expression but we don't know how to parse the next token as an infix
			// in what case would this be valid?

			return leftExp
		}

		p.nextToken()

		leftExp = infixFunc(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFuncError(tokType tokens.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokType)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) parseExpressionIdentifier() ast.Expression {
	return &ast.ExpressionIdentifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}
}

// peekError adds an error message to the parser's Errors slice indicating that the next token was not of the expected type.
func (p *Parser) peekError(tokType tokens.TokenType) {
	errMsg := fmt.Sprintf("expected next token type to be %s, got %s", tokType, p.peekToken.Type)
	p.Errors = append(p.Errors, errMsg)
}

// expectPeek checks if the peek token is of the expected type, and if so advances the parser's tokens.
// it returns true if the peek token is of the expected type, and false otherwise.
func (p *Parser) expectPeek(tokType tokens.TokenType) bool {
	if p.peekToken.Type != tokType {
		p.peekError(tokType)
		return false
	}

	p.nextToken()
	return true
}

func (p *Parser) parseExpressionLiteralInteger() ast.Expression {
	lit := &ast.ExpressionLiteralInteger{Token: p.currToken}

	value, err := strconv.Atoi(p.currToken.Lexeme)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Lexeme)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseExpressionLiteralBoolean() ast.Expression {
	lit := &ast.ExpressionLiteralBoolean{Token: p.currToken}

	value := p.currToken.Type == tokens.TokenTypeTrue

	lit.Value = value

	return lit
}

func (p *Parser) parseExpressionLiteralString() ast.Expression {
	lit := &ast.ExpressionLiteralString{Token: p.currToken}

	lit.Value = p.currToken.Lexeme

	return lit
}

func (p *Parser) parseExpressionPrefix() ast.Expression {
	expr := &ast.ExpressionPrefix{
		Token:    p.currToken,
		Operator: p.currToken.Lexeme,
	}

	p.nextToken()

	expr.Right = p.parseExpression(precedencePrefix)

	return expr
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return precedenceLowest
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return precedenceLowest
}

func (p *Parser) parseExpressionInfix(left ast.Expression) ast.Expression {
	expr := &ast.ExpressionInfix{
		Token:    p.currToken,
		Operator: p.currToken.Lexeme,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseExpressionGrouped() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(precedenceLowest)

	if !p.expectPeek(tokens.TokenTypeRParen) {
		return nil
	}

	return expr
}
