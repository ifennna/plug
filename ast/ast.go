package ast

import (
	"bytes"
	"plug/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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
func (program *Program) String() string {
	var out bytes.Buffer

	for _, statement := range program.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
func (letStatement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(letStatement.TokenLiteral() + " ")
	out.WriteString(letStatement.Name.String())
	out.WriteString(" = ")
	if letStatement.Value != nil {
		out.WriteString(letStatement.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) String() string       { return identifier.Value }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) statementNode()       {}
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }
func (returnStatement *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(returnStatement.TokenLiteral() + " ")
	out.WriteString(" = ")
	if returnStatement.ReturnValue != nil {
		out.WriteString(returnStatement.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expressionStatement *ExpressionStatement) statementNode() {}
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}
func (expressionStatement *ExpressionStatement) String() string {
	if expressionStatement.Expression != nil {
		return expressionStatement.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (blockStatement *BlockStatement) statementNode() {}
func (blockStatement *BlockStatement) TokenLiteral() string {
	return blockStatement.Token.Literal
}
func (blockStatement *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range blockStatement.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (prefixExp *PrefixExpression) expressionNode()      {}
func (prefixExp *PrefixExpression) TokenLiteral() string { return prefixExp.Token.Literal }
func (prefixExp *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(prefixExp.Operator)
	out.WriteString(prefixExp.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (infixExp *InfixExpression) expressionNode()      {}
func (infixExp *InfixExpression) TokenLiteral() string { return infixExp.Token.Literal }
func (infixExp *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(infixExp.Left.String())
	out.WriteString(" " + infixExp.Operator + " ")
	out.WriteString(infixExp.Right.String())
	out.WriteString(")")

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (literal *IntegerLiteral) expressionNode()      {}
func (literal *IntegerLiteral) TokenLiteral() string { return literal.Token.Literal }
func (literal *IntegerLiteral) String() string       { return literal.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (bool *Boolean) expressionNode()      {}
func (bool *Boolean) TokenLiteral() string { return bool.Token.Literal }
func (bool *Boolean) String() string       { return bool.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExp *IfExpression) expressionNode()      {}
func (ifExp *IfExpression) TokenLiteral() string { return ifExp.Token.Literal }
func (ifExp *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ifExp.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifExp.Consequence.String())
	if ifExp.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ifExp.Alternative.String())
	}

	return out.String()
}
