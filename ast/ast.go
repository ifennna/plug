package ast

import "plug/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// Every program is the root node of an AST. Plug is composed of statements
func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token		// the token.LET token
	Name *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode() {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()  {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) statementNode() {}
func (returnStatement *ReturnStatement) TokenLiteral() string {return returnStatement.Token.Literal}