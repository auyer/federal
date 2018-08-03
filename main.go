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

var calcExt = ".fed"

func cleanup(filename string) {
	os.Remove(filename + ".c")
	os.Remove(filename + ".o")
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
	fmt.Fprintln(os.Stderr, "v 0.1")
}

func main() {
	flag.Usage = func() {
		printVersion()
		fmt.Fprintln(os.Stderr, "\nUsage of:", os.Args[0])
		fmt.Fprintln(os.Stderr, os.Args[0], "[flags] <filename>")
		flag.PrintDefaults()
	}
	var (
		cc   = flag.String("cc", "gcc", "C compiler to use")
		cfl  = flag.String("cflags", "-c -std=gnu99", "C compiler flags")
		cout = flag.String("cout", "--output=", "C compiler output flag")
		ld   = flag.String("ld", "gcc", "linker")
		ldf  = flag.String("ldflags", "", "linker flags")
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

	if filepath.Ext(filename) != calcExt {
		fatal("Source files should have the '.fed' extension")
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fatal(err)
	}

	filename = filename[:len(filename)-len(calcExt)]
	fmt.Println("Parsing File")
	f := parse.ParseFile(filename, string(src))
	if f == nil {
		os.Exit(1)
	}

	fmt.Println("Compiling:", filename)

	comp.CompileFile(filename, f)

	/* compile to object code */
	var out []byte
	args := make_args(*cfl, *cout+filename+".o", filename+".c")
	out, err = exec.Command(*cc, strings.Split(args, " ")...).CombinedOutput()
	if err != nil {
		cleanup(filename)
		fatal(string(out), err)
	}

	/* link to executable */
	args = make_args(*ldf, *cout+filename, filename+".o")
	out, err = exec.Command(*ld, strings.Split(args, " ")...).CombinedOutput()
	if err != nil {
		cleanup(filename)
		fatal(string(out), err)
	}
	cleanup(filename)
}
