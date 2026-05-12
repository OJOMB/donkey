package parser

import (
	"fmt"
	"strconv"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/tokens"
	"github.com/OJOMB/donkey/pkg/logs"
)

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
// It initializes the parser's state and registers the necessary parse functions for prefix and infix expressions.
// if logger is nil, the parser will use a null logger by default.
func New(l *lexer.Lexer, logger logs.Logger) (*Parser, error) {
	if logger == nil {
		// null logger to avoid nil pointer dereference
		logger = logs.NewNullLogger()
	}

	if !l.IsInitialised() {
		return nil, ErrLexerUnitialized
	}

	p := &Parser{
		l:                l,
		parseFuncsPrefix: make(map[tokens.Type]parseFuncPrefix),
		parseFuncsInfix:  make(map[tokens.Type]parseFuncInfix),
		Errors:           make([]string, 0),
		logger:           logger.With("component", "parser"),
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
	p.RegisterPrefix(tokens.TypeContinue, p.parseExpressionKeyword)
	p.RegisterPrefix(tokens.TypeBreak, p.parseExpressionKeyword)
	p.RegisterPrefix(tokens.TypeLBracket, p.parseExpressionLiteralList)

	// register infix parse functions for different token types
	p.RegisterInfix(tokens.TypePlus, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeMinus, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeLogicalAnd, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeLogicalOr, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeForwardSlash, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeAsterisk, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeEQ, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeNotEQ, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeLT, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeGT, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeLParen, p.parseExpressionCall)
	p.RegisterInfix(tokens.TypePercent, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeExponent, p.parseExpressionInfix)

	p.RegisterInfix(tokens.TypeBitwiseAnd, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeBitwiseOr, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeBitwiseXor, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeBitwiseShiftLeft, p.parseExpressionInfix)
	p.RegisterInfix(tokens.TypeBitwiseShiftRight, p.parseExpressionInfix)

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
	case tokens.TypeBind:
		return p.parseStatementBind()
	case tokens.TypeReturn:
		return p.parseStatementReturn()
	case tokens.TypeLBrace:
		return p.parseBlockStatement()
	case tokens.TypeWhile:
		return p.parseStatementWhile()
	case tokens.TypeFor:
		return p.parseStatementFor()
	case tokens.TypeFunction:
		return p.parseStatementFunction()
	default:
		// <IDENT><INCREMENT> or <IDENT><DECREMENT>
		if p.isIncrementOrDecrement() {
			// here we have an identifier followed by an increment or decrement token, so we can assume this is an increment or decrement statement like foo++ or bar--;
			return p.parseStatementReBind()
		}

		if p.peekToken.Type == tokens.TypeAssign {
			// here we have an identifier followed by an assign token
			// so we can assume this is a rebind statement like foo = 5; or bar = "hello";
			return p.parseStatementReBind()
		}

		// if the statement doesn't match any of the above types, we assume it's an expression statement and try to parse it as such
		// this would be something like foo(5 + 5); or 5 + 5; or "foobar";
		stmt = p.parseExpressionStatement()
	}

	return stmt
}

func (p *Parser) isIncrementOrDecrement() bool {
	return (p.currToken.Type == tokens.TypeIdent && p.peekToken.Type == tokens.TypeIncrement) ||
		(p.currToken.Type == tokens.TypeIdent && p.peekToken.Type == tokens.TypeDecrement)
}

// parseStatementBind parses a var statement and returns an ast.LetStatement node.
func (p *Parser) parseStatementBind() *ast.StatementBind {
	// first token must be var
	if p.currToken.Type != tokens.TypeBind {
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
func (p *Parser) parseStatementReturn() *ast.StatementReturn {
	if p.currToken.Type != tokens.TypeReturn {
		return nil
	}

	rs := &ast.StatementReturn{Token: p.currToken}

	p.nextToken()

	rs.Value = p.parseExpression(precedenceLowest)

	return rs
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	if p.currToken.Type == tokens.TypeSemicolon {
		return nil
	}

	stmt := &ast.StatementExpression{Token: p.currToken}

	// <IDENT>[<EXPRESSION>] is the only valid expression statement that starts with an identifier that is not a function call or a rebind
	// so we can check for that specific case before trying to parse a generic expression statement
	if p.currToken.Type == tokens.TypeIdent && p.peekToken.Type == tokens.TypeLBracket {
		// here we have an identifier followed by a left bracket, so we can assume this is an index expression being used as a statement, like myList[0];
		stmt.Expression = p.parseExpressionIndex()
		return stmt
	}

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

	// I really don't like this, but it saves us from having to write p.nextToken() every time we want to check for a specific token type
	// the name expectPeek does not imply the fn will advance the tokens, but it does.
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

func (p *Parser) parseExpressionCall(function ast.Expression) ast.Expression {
	expr := &ast.ExpressionCall{
		Token:    p.currToken,
		Function: function,
	}

	expr.Arguments = p.parseCallArguments()

	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := make([]ast.Expression, 0)

	// in the case of no arguments, we should have a right paren immediately after the left paren
	if p.peekToken.Type == tokens.TypeRParen {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(precedenceLowest))

	for p.peekToken.Type == tokens.TypeComma {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(precedenceLowest))
	}

	if !p.expectPeek(tokens.TypeRParen) {
		return nil
	}

	return args
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

func (p *Parser) parseStatementReBind() *ast.StatementRebind {
	if p.peekToken.Type != tokens.TypeAssign {
		if p.isIncrementOrDecrement() {
			return p.desugarIncrementDecrement()
		}

		return nil
	}

	reBind := &ast.StatementRebind{Token: p.currToken}

	reBind.Name = &ast.ExpressionIdentifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}

	p.nextToken() // advance to the assign token
	p.nextToken() // advance to the expression being assigned

	reBind.Value = p.parseExpression(precedenceLowest)
	if reBind.Value == nil {
		return nil
	}

	// when we rebind in the step of a for loop we don't have a semicolon after the expression, but in all other cases we do, so we need to check for that
	if p.peekToken.Type == tokens.TypeSemicolon {
		if !p.expectPeek(tokens.TypeSemicolon) {
			p.logger.Debug("expected ; after rebind expression, got %s instead", p.peekToken.Type)
			return nil
		}
	}

	return reBind
}

func (p *Parser) desugarIncrementDecrement() *ast.StatementRebind {
	// here we want to take a statement like foo++ and desugar it to foo = foo + 1
	// or take a statement like bar-- and desugar it to bar = bar - 1
	ident := &ast.ExpressionIdentifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}

	var operator string
	var operatorType tokens.Type
	switch p.peekToken.Type {
	case tokens.TypeIncrement:
		operator = "+"
		operatorType = tokens.TypePlus
	case tokens.TypeDecrement:
		operator = "-"
		operatorType = tokens.TypeMinus
	default:
		return nil
	}

	infixExpr := &ast.ExpressionInfix{
		Token:    tokens.New(operatorType, operator),
		Operator: operator,
		Left:     ident,
		Right: &ast.ExpressionLiteralInteger{
			Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
			Value: 1,
		},
	}

	reBind := &ast.StatementRebind{
		Token: p.currToken,
		Name:  ident,
		Value: infixExpr,
	}

	// we need to advance the tokens past the increment or decrement token
	p.nextToken() // advance to the increment or decrement token

	if !p.expectPeek(tokens.TypeSemicolon) {
		p.logger.Debug("expected ; after increment/decrement expression, got %s instead", p.peekToken.Type)
		return nil
	}

	return reBind
}

func (p *Parser) parseStatementWhile() ast.Statement {
	stmt := &ast.StatementWhile{
		Token: p.currToken,
	}

	// conditions are optional in while loops
	// so we only parse the condition if we see a left parenthesis after the while token
	if p.peekToken.Type == tokens.TypeLParen {
		p.nextToken()
		p.nextToken()
		stmt.Condition = p.parseExpression(precedenceLowest)

		if !p.expectPeek(tokens.TypeRParen) {
			p.logger.Debug("expected ) after while condition, got %s instead", p.peekToken.Type)
			return nil
		}
	}

	if !p.expectPeek(tokens.TypeLBrace) {
		p.logger.Debug("expected { after while condition, got %s instead", p.peekToken.Type)
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseStatementFor() ast.Statement {
	stmt := &ast.StatementFor{
		Token: p.currToken,
	}

	// for loops have the syntax for (initializer; condition; step) { body }
	// so we need to parse the initializer, condition, and step expressions before parsing the body

	if !p.expectPeek(tokens.TypeLParen) {
		p.logger.Debug("expected ( after for, got %s instead", p.peekToken.Type)
		return nil
	}

	// the initializer must be a var statement
	if !p.expectPeek(tokens.TypeBind) {
		p.logger.Debug("expected var statement as initializer in for loop, got %s instead", p.peekToken.Type)
		return nil
	}

	stmt.Initializer = p.parseStatementBind()

	// next we need to parse the condition expression, which is not optional in for loops
	if !p.expectPeek(tokens.TypeSemicolon) {
		p.logger.Debug("expected ; after initializer in for loop, got %s instead", p.peekToken.Type)
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(precedenceLowest)

	if !p.expectPeek(tokens.TypeSemicolon) {
		p.logger.Debug("expected ; after condition in for loop, got %s instead", p.peekToken.Type)
		return nil
	}

	// next we need to parse the step expression, which is not optional in for loops
	p.nextToken()
	stmt.Step = p.parseStatement()

	if !p.expectPeek(tokens.TypeRParen) {
		p.logger.Debug("expected step expression in for loop, got ) instead")
		return nil
	}

	if !p.expectPeek(tokens.TypeLBrace) {
		p.logger.Debug("expected { after for loop step expression, got %s instead", p.peekToken.Type)
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseExpressionKeyword() ast.Expression {
	return &ast.ExpressionKeyword{
		Token:   p.currToken,
		Keyword: p.currToken.Lexeme,
	}
}

// parseStatementFunction parses a function declaration statement and returns an ast.StatementFunctionBind node representing the parsed function declaration.
// THIS IS A NAMED FUNCTION, NOT AN ANONYMOUS FUNCTION LITERAL. The syntax for a function declaration is fn <name>(<parameters>) { <body> }.
func (p *Parser) parseStatementFunction() ast.Statement {
	fn := &ast.StatementFunctionBind{
		Token: p.currToken,
		Value: &ast.ExpressionLiteralFunction{
			Token: p.currToken,
		},
	}

	if p.peekToken.Type != tokens.TypeIdent {
		// we can try to parse this as an anonymous function literal
		expr := p.parseExpressionLiteralFunction()
		if expr == nil {
			p.logger.Debug("expected function name identifier after fn, got %s instead", p.peekToken.Type)
			return nil
		}

		return &ast.StatementExpression{
			Token:      tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
			Expression: expr,
		}
	}

	p.nextToken() // advance to the function name identifier

	fn.Name = &ast.ExpressionIdentifier{
		Token: p.currToken,
		Value: p.currToken.Lexeme,
	}

	if !p.expectPeek(tokens.TypeLParen) {
		p.logger.Debug("expected ( after function name in function declaration, got %s instead", p.peekToken.Type)
		return nil
	}

	fn.Value.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(tokens.TypeLBrace) {
		p.logger.Debug("expected { after function parameters in function declaration, got %s instead", p.peekToken.Type)
		return nil
	}

	fn.Value.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseExpressionLiteralList() ast.Expression {
	lit := &ast.ExpressionLiteralList{Token: p.currToken}

	// in the case of an empty list, we should have a right bracket immediately after the left bracket
	if p.peekToken.Type == tokens.TypeRBracket {
		p.nextToken()
		return lit
	}

	p.nextToken()
	lit.Elements = p.parseExpressionList(tokens.TypeRBracket)

	return lit
}

func (p *Parser) parseExpressionList(end tokens.Type) []ast.Expression {
	exprs := make([]ast.Expression, 0)
	exprs = append(exprs, p.parseExpression(precedenceLowest))
	for p.peekToken.Type == tokens.TypeComma {
		p.nextToken()
		p.nextToken()
		exprs = append(exprs, p.parseExpression(precedenceLowest))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return exprs
}

func (p *Parser) parseExpressionIndex() ast.Expression {
	expr := &ast.ExpressionIndex{
		Token: p.peekToken,
		Left: &ast.ExpressionIdentifier{
			Token: p.currToken,
			Value: p.currToken.Lexeme,
		},
	}

	p.nextToken() // advance to l bracket token
	p.nextToken() // advance to the expression inside the brackets
	expr.Index = p.parseExpression(precedenceLowest)
	if expr.Index == nil {
		p.logger.Debug("failed to parse index expression inside brackets")
		return nil
	}

	if !p.expectPeek(tokens.TypeRBracket) {
		p.logger.Debug("expected ] after index expression, got %s instead", p.peekToken.Type)
		return nil
	}

	return expr
}
