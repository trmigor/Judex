package main

import (
	"fmt"
	"os"
	"time"
	"./testing_packages/compileandrun"
)

func main() {
	compiler := "gcc"
	args := make([]string, 1)
	args[0] = "/Users/itar/Documents/GitHub/Judex/src/go_drafts/qwerty.c"
	
	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	c := compileandrun.Init {
		Solution:		1,
		Format:			".c",
		Path:			"/Users/itar/go/src/github.com/trmigor/Judex/go_drafts/",
		Compiler:		compiler,
		//CompilerArgs:	args,
		CompilerAttr:	attr,
		TestsPath:		"/Users/itar/go/src/github.com/trmigor/Judex/go_drafts/testing_packages/",
		TestsNumber:	1,
		RunLimits:		compileandrun.Limits {
							TL:		1*1000*1000*1000,
							RTL:	1*1000*1000*1000,
						},
	}
	p, err := c.Compile()
	fmt.Println("Hello World1", p, err)
	time.Sleep(100*1000*1000)
	err = c.Run()
	fmt.Println("Hello World2", err)
}
