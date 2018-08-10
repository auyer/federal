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

// CompileFile creates the necessary files for compiling
func CompileFile(fname string, parsed *ast.Source) {

	var c compiler
	fp, err := os.Create(fname + ".go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fp.Close()
	c.fp = fp
	c.compSource(parsed)
}

// compNode checks what the current node is
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

// compBinaryExpr reads the binary expressions from the tree, and optmizes them
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

// compSource prints the necessary data into the IR file
func (c *compiler) compSource(f *ast.Source) {
	fmt.Fprintln(c.fp, "package main")
	fmt.Fprintln(c.fp, `import "fmt"`)
	fmt.Fprintln(c.fp, `import "strconv"`)
	fmt.Fprintln(c.fp, "func main() {")
	fmt.Fprintf(c.fp, "fmt.Println(strconv.Itoa(%d))\n", c.compNode(f.Root))
	fmt.Fprintln(c.fp, "}")
}
