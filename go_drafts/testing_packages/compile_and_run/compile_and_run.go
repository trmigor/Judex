/*
	Package for testing the compiled programming languages
*/
package compile_and_run

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type Limits struct {
	TL int
	ML syscall.Rlimit
	
}

/*
	Initializing structure Init
	
	Init.Solution		-- number of testing solution
	Init.Format			-- solution file format (with a dot)
	Init.Path			-- path to the solution file
	Init.Compiler		-- chosen compiler
	Init.CompilerArgs	-- compiler command line arguments
	Init.CompilerAttr	-- compiler attributes
	Init.TestsPath		-- path to tests data files
	Init.RunArgs		-- solution command line arguments
	Init.RunAttr		-- run attributes
	
*/
type Init struct {
	Solution int
	Format string
	Path string
	Compiler string
	CompilerArgs []string
	CompilerAttr os.ProcAttr
	TestsPath string
	RunArgs []string
	RunAttr os.ProcAttr
}

/*
	Method that compiles source code into the binary file
	in ".elf" format. The source must be named like
	"sol_<solution_number><format>", where <format> matches
	the chosen programming language. The resuling binary
	file will be named like "res_<solution_number>.elf"
*/
func (c *Init) Compile() (p *os.Process, err error) {
	// Input format: "sol_<solution_number><format>"
	path := c.Path + "sol_" + strconv.Itoa(c.Solution) + c.Format
	
	// Output format: "res_<solution_number>.elf"
	res := c.Path + "res_" + strconv.Itoa(c.Solution) + ".elf"
	
	// Checking the existence of a source file
	source, err := os.Open(path)
	if err != nil {
		return p, err
	}
	err = source.Close()
	
	switch c.Compiler {
		
		/* 
			C
			<format>: ".c", ".i",
		*/
		case "gcc":
		
			/*
				If the compiler exists in the PATH environment variable,
				finds the relevant path to it
			*/
			if c.Compiler, err = exec.LookPath(c.Compiler); err == nil {
				args := make([]string, 3)
				
				args[0] = c.Compiler
				
				/*
					Redirection to resulting binary file
				*/
				args[1] = "-o"
				args[2] = res
				
				args = append(args, c.CompilerArgs...)
				args = append(args, path)
				
				// Calling of a compiler
				p, err = os.StartProcess(c.Compiler, args, &c.CompilerAttr)
			}
	}
	return p, err
}

/*
	Method that runs the binary file like "res_<solution_number>.elf".
	It uses tests from Init.TestsPath like input if RunAttr == <nil>.
	
*/
func (c *Init) Run() (p *os.Process, err error) {
	
	return p, err
}