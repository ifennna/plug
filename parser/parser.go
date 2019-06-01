package parser

import (
	"fmt"
	"plug/ast"
	"plug/lexer"
	"plug/token"
	"strconv"
)

// The arrangement of the following constants indicates their order of precedence
const (
	_ int = iota // give the following constants incrementing values from 0
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFuncs map[token.Type]prefixParseFunc
	infixParseFuncs  map[token.Type]infixParseFunc
}

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	parser.nextToken() // set currentToken
	parser.nextToken() // set peekToken

	parser.prefixParseFuncs = make(map[token.Type]prefixParseFunc)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)

	return parser
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) currentTokenIs(token token.Type) bool {
	return parser.currentToken.Type == token
}
func (parser *Parser) peekTokenIs(token token.Type) bool {
	return parser.peekToken.Type == token
}
func (parser *Parser) expectPeek(token token.Type) bool {
	if parser.peekTokenIs(token) {
		parser.nextToken()
		return true
	} else {
		parser.throwPeekError(token)
		return false
	}
}

func (parser *Parser) throwPeekError(token token.Type) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", token, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) registerPrefix(tokenType token.Type, function prefixParseFunc) {
	parser.prefixParseFuncs[tokenType] = function
}
func (parser *Parser) registerInfix(tokenType token.Type, function infixParseFunc) {
	parser.infixParseFuncs[tokenType] = function
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: parse expression(s)
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser.currentToken}

	parser.nextToken()

	// TODO: parse expression(s)
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.currentToken}

	statement.Expression = parser.parseExpression(LOWEST)

	// TODO: parse expression(s)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFuncs[parser.currentToken.Type]
	if prefix == nil {
		return nil
	}
	leftExpression := prefix()
	return leftExpression
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, message)
		return nil
	}

	literal.Value = value
	return literal
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !parser.currentTokenIs(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) Errors() []string {
	return parser.errors
}
