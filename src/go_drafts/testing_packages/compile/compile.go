package compile

import (
	"fmt"
	"os"
	"os/exec"
)

type Init struct {
	Compiler string
	Args []string
	Attr os.ProcAttr
}

func (c *Init) Compile() (p *os.Process, err error) {	
	if c.Compiler, err = exec.LookPath(c.Compiler); err == nil {
		p, err := os.StartProcess(c.Compiler, c.Args, &c.Attr)
		if err == nil {
			return p, nil
		}
	}
	fmt.Println("OK", err)
	return nil, err
}