package main

import (
	"fmt"
	"os"
	"./testing_packages/compile_and_run"
)

func main() {
	compiler := "gcc"
	args := make([]string, 1)
	args[0] = "/Users/itar/Documents/GitHub/Judex/src/go_drafts/qwerty.c"
	
	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	c := compile_and_run.Init {
		Solution:		1,
		Format:			".c",
		Path:			"/Users/itar/Documents/GitHub/Judex/src/go_drafts/",
		Compiler:		compiler,
		//CompilerArgs:	args,
		CompilerAttr:	attr,
	}
	p, err := c.Compile()
	fmt.Println("Hello World!", p, err)
}
