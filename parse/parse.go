package parse

import (
	"os"

	"github.com/auyer/federal/ast"
	"github.com/auyer/federal/scan"
	"github.com/auyer/federal/token"
)

func ParseFile(filename, src string) *ast.Source {
	var p parser
	p.init(filename, src)
	f := p.parseSource()
	if p.errors.Count() > 0 {
		p.errors.Print()
		return nil
	}
	return f
}

type parser struct {
	file    *token.Source
	errors  scan.ErrorList
	scanner scan.Scanner

	pos token.Pos
	tok token.Token
	lit string
}

func (p *parser) addError(msg string) {
	p.errors.Add(p.file.Position(p.pos), msg)
	if p.errors.Count() >= 10 {
		p.errors.Print()
		os.Exit(1)
	}
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.addError("Expected '" + tok.String() + "' got '" + p.lit + "'")
	}
	p.next()
	return pos
}

func (p *parser) init(fname, src string) {
	p.file = token.NewSource(fname, src)
	p.scanner.Init(p.file, src)
	p.next()
}

func (p *parser) next() {
	p.lit, p.tok, p.pos = p.scanner.Scan()
}

func (p *parser) parseBasicLit() *ast.BasicLit {
	return &ast.BasicLit{LitPos: p.pos, Kind: p.tok, Lit: p.lit}
}

func (p *parser) parseBinaryLitExpr(open token.Pos) *ast.BinaryExpr {
	pos := p.pos
	op := p.tok
	p.next()

	var list []ast.Expr
	for p.tok != token.RPAREN && p.tok != token.EOF {
		list = append(list, p.parseGenExpr())
	}
	if len(list) < 2 {
		p.addError("binary expression must have at least two operands")
	}
	end := p.expect(token.RPAREN)
	return &ast.BinaryExpr{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		Op:    op,
		OpPos: pos,
		List:  list,
	}
}

func (p *parser) parseBinaryExpr(open token.Pos) *ast.BinaryExpr {
	pos := p.pos
	op := p.tok
	p.next()
	open = p.expect(token.LPAREN)

	var list []ast.Expr
	for p.tok != token.RPAREN && p.tok != token.EOF {
		list = append(list, p.parseGenExpr())
	}
	if len(list) < 2 {
		p.addError("binary expression must have at least two operands")
	}
	end := p.expect(token.RPAREN)
	return &ast.BinaryExpr{
		Expression: ast.Expression{
			Opening: open,
			Closing: end,
		},
		Op:    op,
		OpPos: pos,
		List:  list,
	}
}

func (p *parser) parseGenExpr() ast.Expr {
	var expr ast.Expr

	// TODO ADD PRINT AND MAIN HERE
	switch p.tok {
	case token.LPAREN:
		expr = p.parseExprParen()
	case token.DO:
		expr = p.parseExprDo()
	case token.INTEGER:
		expr = p.parseBasicLit()
		p.next()
	default:
		p.addError("Expected '" + token.LPAREN.String() + "' or '" +
			token.INTEGER.String() + "' got '" + p.lit + "'")
		p.next()
	}

	return expr
}

func (p *parser) parseExprParen() ast.Expr {
	return p.parseExpr(p.expect(token.LPAREN))
}

func (p *parser) parseExprDo() ast.Expr {
	p.expect(token.DO)
	expr := p.parseExpr(p.expect(token.LPAREN))
	_ = p.expect(token.RPAREN)
	return expr
}

func (p *parser) parseExpr(pos token.Pos) ast.Expr {
	var expr ast.Expr
	switch p.tok {
	case token.ADD, token.SUB, token.MUL, token.DIV, token.REM:
		expr = p.parseBinaryExpr(pos)
	case token.ADDLIT, token.SUBLIT, token.MULLIT, token.DIVLIT, token.REMLIT:
		expr = p.parseBinaryLitExpr(pos)
	default:
		p.addError("Expected binary operator but got '" + p.lit + "'")
	}
	return expr
}

func (p *parser) parseSource() *ast.Source {
	var expr ast.Expr
	expr = p.parseGenExpr()
	if p.tok != token.EOF {
		p.addError("Expected EOF, got '" + p.lit + "'")
	}
	return &ast.Source{Root: expr}
}
