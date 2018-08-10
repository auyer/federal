package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/auyer/federal/comp"
	"github.com/auyer/federal/parse"
)

var fedExt = ".fed"

func cleanup(filename string) {
	os.Remove(filename + ".go")
}

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func make_args(options ...string) string {
	var args string
	for i, opt := range options {
		if len(opt) > 0 {
			args += opt
			if i < len(options)-1 {
				args += " "
			}
		}
	}
	return args
}

func printVersion() {
	fmt.Fprintln(os.Stderr, "federal v 0.1")
}

func main() {
	flag.Usage = func() {
		printVersion()
		fmt.Fprintln(os.Stderr, "\nUsage of:", os.Args[0])
		fmt.Fprintln(os.Stderr, os.Args[0], "[flags] <filename>")
		flag.PrintDefaults()
	}
	var (
		cc   = flag.String("cc", "go", "GO Compiler to use")
		cfl  = flag.String("cflags", "build", "Go Compiler Flags")
		keep = flag.Bool("keep", false, "Keep Intermediary files")
		ir   = flag.Bool("ir", false, "Stop at IR")
		ver  = flag.Bool("version", false, "Print version number and exit")
	)
	flag.Parse()

	if *ver {
		printVersion()
		os.Exit(1)
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)

	if filepath.Ext(filename) != fedExt {
		fatal("Source files should have the '.fed' extension")
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fatal(err)
	}

	filename = filename[:len(filename)-len(fedExt)]
	fmt.Println("> Parsing File")
	f := parse.ParseFile(filename, string(src))
	if f == nil {
		os.Exit(1)
	}

	fmt.Println("> Compiling:", filename)

	comp.CompileFile(filename, f)

	if *ir {
		fmt.Println("> Stopping at Intermediary Language")
		os.Exit(1)
	}
	/* compile to object code */
	var out []byte
	args := make_args(*cfl, filename+".go")
	out, err = exec.Command(*cc, strings.Split(args, " ")...).CombinedOutput()
	if err != nil {
		cleanup(filename)
		fatal(string(out), err)
	}
	if !*keep {
		cleanup(filename)
	}
}
