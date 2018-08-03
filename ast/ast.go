package ast

import (
	"github.com/auyer/federal/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Expr interface {
	Node
	exprNode()
}

type BasicLit struct {
	LitPos token.Pos
	Kind   token.Token
	Lit    string
}

type BinaryExpr struct {
	Expression
	Op    token.Token
	OpPos token.Pos
	List  []Expr
}

type Expression struct {
	Opening token.Pos
	Closing token.Pos
}

type Source struct {
	Root Expr
}

func (b *BasicLit) Pos() token.Pos   { return b.LitPos }
func (e *Expression) Pos() token.Pos { return e.Opening }
func (f *Source) Pos() token.Pos     { return f.Root.Pos() }

func (b *BasicLit) End() token.Pos   { return b.LitPos + token.Pos(len(b.Lit)) }
func (e *Expression) End() token.Pos { return e.Closing }
func (f *Source) End() token.Pos     { return f.Root.End() }

func (b *BasicLit) exprNode()   {}
func (e *Expression) exprNode() {}
