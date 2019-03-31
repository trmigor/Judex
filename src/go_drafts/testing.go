package main

import (
	"fmt"
	"os"
	"./testing_packages/compile"
)

func main() {
	compiler := "gcc"
	args := make([]string, 4)
	args[0] = compiler
	args[1] = "/Users/itar/Documents/GitHub/Judex/src/go_drafts/qwerty.c"
	args[2] = "-o"
	args[3] = "/Users/itar/Documents/GitHub/Judex/src/go_drafts/qwerty.elf"
	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	c := compile.Init {
		Compiler:	compiler,
		Args:		args,
		Attr:		attr,
	}
	p, err := c.Compile()
	fmt.Println("Hello World!", p, err)
}
