package parser

import (
	"fmt"
	"strconv"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/tokens"
	"github.com/OJOMB/donkey/pkg/logs"
)

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
	tokens.TypePlus:         precedenceAdditive,
	tokens.TypeMinus:        precedenceAdditive,
	tokens.TypeForwardSlash: precedenceMultiplicative,
	tokens.TypeAsterisk:     precedenceMultiplicative,
	tokens.TypeLParen:       precedenceCall,
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

	parseFuncsPrefix map[tokens.Type]parseFuncPrefix
	parseFuncsInfix  map[tokens.Type]parseFuncInfix

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
		parseFuncsPrefix: make(map[tokens.Type]parseFuncPrefix),
		parseFuncsInfix:  make(map[tokens.Type]parseFuncInfix),
		Errors:           make([]string, 0),
		logger:           logger,
	}

	// register prefix parse functions for different token types
	p.RegisterPrefix(tokens.TypeIdent, p.parseExpressionIdentifier)
	p.RegisterPrefix(tokens.TypeInt, p.parseExpressionLiteralInteger)
	p.RegisterPrefix(tokens.TypeString, p.parseExpressionLiteralString)
	p.RegisterPrefix(tokens.TypeTrue, p.parseExpressionLiteralBoolean)
	p.RegisterPrefix(tokens.TypeFalse, p.parseExpressionLiteralBoolean)
	p.RegisterPrefix(tokens.TypeBang, p.parseExpressionPrefix)
	p.RegisterPrefix(tokens.TypeMinus, p.parseExpressionPrefix)
	p.RegisterPrefix(tokens.TypeLParen, p.parseExpressionGrouped)
	p.RegisterPrefix(tokens.TypeIf, p.parseExpressionIf)
	p.RegisterPrefix(tokens.TypeFunction, p.parseExpressionLiteralFunction)

	// register infix parse functions for different token types
	p.RegisterInfix(tokens.TypePlus, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeMinus, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeForwardSlash, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeAsterisk, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeEq, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeNotEq, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeLT, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeGT, p.parseExpressionInfix)

	// Read two tokens, so currToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p, nil
}

func (p *Parser) RegisterPrefix(tokType tokens.Type, fn parseFuncPrefix) {
	p.parseFuncsPrefix[tokType] = fn
}

func (p *Parser) RegisterInfix(tokType tokens.Type, fn parseFuncInfix) {
	p.parseFuncsInfix[tokType] = fn
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prgrm := ast.NewProgram()
	for p.currToken.Type != tokens.TypeEOF {
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
	case tokens.TypeBinder:
		return p.parseStatementLet()
	case tokens.TypeReturn:
		return p.parseStatementReturn()
	// case tokens.TypeIdent:
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

// parseStatementLet parses a var statement and returns an ast.LetStatement node.
func (p *Parser) parseStatementLet() *ast.StatementBind {
	// first token must be var
	if p.currToken.Type != tokens.TypeBinder {
		return nil
	}

	var ls = &ast.StatementBind{Token: p.currToken}

	// next token must be ident
	if !p.expectPeek(tokens.TypeIdent) {
		return nil
	}

	ls.Name = &ast.ExpressionIdentifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}

	// next token must be assign =
	if !p.expectPeek(tokens.TypeAssign) {
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
	if p.currToken.Type != tokens.TypeReturn {
		return nil
	}

	rs := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	rs.ReturnValue = p.parseExpression(precedenceLowest)

	return rs
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	if p.currToken.Type == tokens.TypeSemicolon {
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
	// if the precedence of the peek token is higher than the current precedence, then we enter the loop
	// we continue to parse infix expressions until we reach a semicolon or a token with lower precedence than the current precedence
	for p.peekToken.Type != tokens.TypeSemicolon && precedence < p.peekPrecedence() {
		infixFunc, ok := p.parseFuncsInfix[p.peekToken.Type]
		if !ok {
			// TODO should this not be an error instead of just returning the left expression?
			// in this case we have a valid left expression, peekPrecedence has returned a value but we don't know how to parse the next token as an infix
			// in what case would this be valid?

			return leftExp
		}

		p.nextToken()

		// so if we have an infix function, we need to parse the infix expression and update the left expression to be the result of the infix expression
		// for example, if we have 5 + 5 * 5, then we would first parse the left expression as 5
		// then we would see the + token and parse the infix expression as (5 + 5)
		// and then we would see the * token and parse the infix expression as ((5 + 5) * 5)
		// if we then peek a token with lower precedence than the current precedence, we would stop parsing infix expressions and return the left expression, which would be ((5 + 5) * 5) in this example
		leftExp = infixFunc(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFuncError(tokType tokens.Type) {
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
func (p *Parser) peekError(tokType tokens.Type) {
	errMsg := fmt.Sprintf("expected next token type to be %s, got %s", tokType, p.peekToken.Type)
	p.Errors = append(p.Errors, errMsg)
}

// expectPeek checks if the peek token is of the expected type, and if so advances the parser's tokens.
// it returns true if the peek token is of the expected type, and false otherwise.
func (p *Parser) expectPeek(tokType tokens.Type) bool {
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
	return &ast.ExpressionLiteralBoolean{
		Token: p.currToken,
		Value: p.currToken.Type == tokens.TypeTrue,
	}
}

// parseExpressionLiteralString parses a string literal expression and returns an ast.ExpressionLiteralString node representing the parsed string literal.
// It expects the current token to be a string token, and it will set the Value field of the ExpressionLiteralString node to the lexeme of the current token.
func (p *Parser) parseExpressionLiteralString() ast.Expression {
	return &ast.ExpressionLiteralString{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}
}

// parseExpressionLiteralFunction parses a function literal expression and returns an ast.ExpressionLiteralFunction node representing the parsed function literal.
// It expects the current token to be the "fn" keyword, followed by a parameter list enclosed in parentheses, and a function body enclosed in braces.
func (p *Parser) parseExpressionLiteralFunction() ast.Expression {
	lit := &ast.ExpressionLiteralFunction{Token: p.currToken}

	if !p.expectPeek(tokens.TypeLParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(tokens.TypeLBrace) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

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
	if !p.expectPeek(tokens.TypeRParen) {
		return nil
	}

	return expr
}

func (p *Parser) parseExpressionIf() ast.Expression {
	expr := &ast.ExpressionIf{
		Branches: make([]ast.ConditionalBranch, 0),
	}

	for {
		branch := ast.ConditionalBranch{
			Token: p.currToken,
		}

		if !p.expectPeek(tokens.TypeLParen) {
			return nil
		}

		p.nextToken()
		branch.Condition = p.parseExpression(precedenceLowest)

		if !p.expectPeek(tokens.TypeRParen) {
			return nil
		}

		if !p.expectPeek(tokens.TypeLBrace) {
			return nil
		}

		branch.Consequence = p.parseBlockStatement()

		expr.Branches = append(expr.Branches, branch)

		if p.peekToken.Type != tokens.TypeElif {
			break
		}

		p.nextToken()
	}

	if p.peekToken.Type == tokens.TypeElse {
		p.nextToken()

		if !p.expectPeek(tokens.TypeLBrace) {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.StatementBlock {
	block := &ast.StatementBlock{Statements: make([]ast.Statement, 0)}

	p.nextToken()
	for p.currToken.Type != tokens.TypeRBrace && p.currToken.Type != tokens.TypeEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

// parseFunctionParameters parses the parameters of a function literal and returns a slice of ExpressionIdentifier nodes representing the parameters.
// It expects the current token to be the opening parenthesis of the parameter list, and it will advance the parser's tokens as it parses the parameters.
func (p *Parser) parseFunctionParameters() []*ast.ExpressionIdentifier {
	identifiers := make([]*ast.ExpressionIdentifier, 0)

	// in the case of no params
	if p.peekToken.Type == tokens.TypeRParen {
		p.nextToken()
		return identifiers
	}

	if p.peekToken.Type != tokens.TypeIdent {
		p.peekError(tokens.TypeIdent)
		return nil
	}

	// we're now at the first parameter, so we can start parsing
	p.nextToken()

	ident := &ast.ExpressionIdentifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}
	identifiers = append(identifiers, ident)

	// fn(param1, param2, param3) {	<statements> }
	for p.peekToken.Type == tokens.TypeComma {
		// advances to the comma
		p.nextToken()
		// advances to the next parameter after the comma
		p.nextToken()

		ident := &ast.ExpressionIdentifier{
			Token: p.currToken,
			Value: p.currToken.Lexeme,
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(tokens.TypeRParen) {
		return nil
	}

	return identifiers
}
