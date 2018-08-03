package comp

import (
	"fmt"
	"os"
	"strconv"

	"github.com/auyer/federal/ast"
	"github.com/auyer/federal/token"
)

type compiler struct {
	fp *os.File
}

func CompileFile(fname string, parsed *ast.Source) {

	var c compiler
	fp, err := os.Create(fname + ".c")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fp.Close()
	c.fp = fp
	c.compFile(parsed)
}

func (c *compiler) compNode(node ast.Node) int {
	switch n := node.(type) {
	case *ast.BasicLit:
		i, err := strconv.Atoi(n.Lit)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return i
	case *ast.BinaryExpr:
		return c.compBinaryExpr(n)
	default:
		return 0 /* can't be reached */
	}
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) int {
	var tmp int

	tmp = c.compNode(b.List[0])

	for _, node := range b.List[1:] {
		switch b.Op {
		case token.ADD:
			tmp += c.compNode(node)
		case token.SUB:
			tmp -= c.compNode(node)
		case token.MUL:
			tmp *= c.compNode(node)
		case token.DIV:
			tmp /= c.compNode(node)
		case token.REM:
			tmp %= c.compNode(node)
		case token.ADDLIT:
			tmp += c.compNode(node)
		case token.SUBLIT:
			tmp -= c.compNode(node)
		case token.MULLIT:
			tmp *= c.compNode(node)
		case token.DIVLIT:
			tmp /= c.compNode(node)
		case token.REMLIT:
			tmp %= c.compNode(node)
		}
	}

	return tmp
}

func (c *compiler) compFile(f *ast.Source) {
	fmt.Fprintln(c.fp, "#include <stdio.h>")
	fmt.Fprintln(c.fp, "int main(void) {")
	fmt.Fprintf(c.fp, "printf(\"%%d\\n\", %d);\n", c.compNode(f.Root))
	fmt.Fprintln(c.fp, "return 0;")
	fmt.Fprintln(c.fp, "}")
}
