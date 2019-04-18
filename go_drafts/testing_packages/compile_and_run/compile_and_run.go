/*
	Package for testing the compiled programming languages
*/
package compile_and_run

import (
	"os"				// OS syscalls
	"os/exec"			// Executive file search
	"strconv"			// String convertation
	"syscall"			// OS syscalls
	"time"				// Timing
	"bufio"				// Buffer input/output
	"bytes"				// Bytes manipulations
	"errors"			// Errors manipulation
	"encoding/csv"		// Protocoling
)


/*
	Structures
*/

// Structure with limits of a run
type Limits struct {
	TL				time.Duration	// -- time limit
	RTL 			time.Duration	// -- real time limit
	ML				int64			// -- memory limit
}

// Initializing structure
type Init struct {
	Solution     	int				// -- number of tested solution
	Format       	string			// -- solution file format (with a dot)
	Path         	string			// -- path to the solution file
	Compiler     	string			// -- chosen compiler
	CompilerArgs 	[]string		// -- compiler command line arguments
	CompilerAttr 	os.ProcAttr		// -- compiler attributes
	TestsPath    	string			// -- path to tests data files
	TestsNumber  	int				// -- number of tests
	RunArgs     	[]string		// -- run command line arguments
	RunAttr			os.ProcAttr		// -- run attributes
	RequiredRet		int				// -- required return value
	RunLimits		Limits			// -- run program limits
}

// Reporting structure
type Result struct {
	TestNumber		int				// -- number of used test
	ReturnedValue	int				// -- value returned from the process
	Verdict			string			// -- result of limits and returned value checking
	Checker			string			// -- result of answer checking
}


/*
	Methods
*/

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
		return nil, err
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

			// Redirection to resulting binary file
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
	It uses tests from Init.TestsPath like input if RunAttr.Files == <nil>.
	Running for different is doing in different goroutines
*/
func (c *Init) Run() (err error) {
	// Path to binary file
	run_name := c.Path + "res_" + strconv.Itoa(c.Solution) + ".elf"

	// Channel and list for testing results
	result_chan := make(chan Result)
	result_list := make([]Result, c.TestsNumber)

	for i := 1; i <= c.TestsNumber; i++ {
		// Redirecting of input, output and error output
		RunAttr := c.RunAttr
		if c.RunAttr.Files == nil {
			// Open tests file as input
			input, err := os.OpenFile(c.TestsPath + strconv.Itoa(i) + ".dat", os.O_RDONLY, 0744)
			if err != nil {
				return err
			}

			// Open output file
			output, err := os.OpenFile(c.Path + strconv.Itoa(i) + ".ans", os.O_RDWR|os.O_CREATE, 0766)
			if err != nil {
				return err
			}

			// Open error output file
			error, err := os.OpenFile(c.Path + strconv.Itoa(i) + ".err", os.O_RDWR|os.O_CREATE, 0766)
			if err != nil {
				return err
			}

			// Set input, output and error output files for a run
			RunAttr.Files = []*os.File{input, output, error}
		}

		// Set run arguments
		args := make([]string, 1)

		args[0] = run_name

		args = append(args, c.RunArgs...)

		// Running
		go c.run(run_name, args, &RunAttr, i, result_chan)
	}

	// Recieving information from goroutines
	for i := 1; i <= c.TestsNumber; i++ {
		result := <- result_chan
		result_list[result.TestNumber - 1] = result
	}

	// Protocoling
	protocol, err := os.OpenFile(c.Path + strconv.Itoa(c.Solution) + ".prot", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0766)
	csv_writer := csv.NewWriter(bufio.NewWriter(protocol))
	for i := range result_list {
		csv_writer.Write([]string{
			strconv.Itoa(result_list[i].TestNumber),
			strconv.Itoa(result_list[i].ReturnedValue),
			result_list[i].Verdict,
			result_list[i].Checker,
		})
		if err = csv_writer.Error(); err != nil {
			return err
		}
	}
	csv_writer.Flush()
	err = csv_writer.Error()
	return err
}


/*
	Helping functions
*/

// Function, that runs a process and checks limits
func (c *Init) run(run_name string, RunArgs []string, RunAttr *os.ProcAttr,
		test_num int, result_chan chan Result) {
	// Setting timer
	start_time := time.Now()

	// Starting process
	pid, err := os.StartProcess(run_name, RunArgs, RunAttr)

	if err != nil {
		result_chan <- Result {
			TestNumber:		test_num,
			ReturnedValue:	-1,
			Verdict:		"Cannot start run: " + err.Error(),
		}
		return
	}

	// If it take too long, kill
	done := make(chan string)
	go kill_too_long(pid, c.RunLimits, done)

	// Scanning process info
	p_state, err := pid.Wait()

	if err != nil {
		result_chan <- Result {
			TestNumber:		test_num,
			ReturnedValue:	-1,
			Verdict:		"Cannot wait for running end: " + err.Error(),
		}
		return
	}

	// Stopping timer
	end_time := time.Now()
	
	_ = <- done

	// Standart result
	result := Result {
		TestNumber:		test_num,
		ReturnedValue:	p_state.ExitCode(),
		Verdict:		"OK",
	}

	// Process resources usage
	status := p_state.SysUsage().(*syscall.Rusage)

	// Checking RTL
	if real_time := end_time.Sub(start_time); real_time > c.RunLimits.RTL && c.RunLimits.RTL != 0 {
		result.ReturnedValue = -1
		result.Verdict = "Real time limit exceeded"
	}

	// Checking TL
	if time.Duration(status.Utime.Usec) > c.RunLimits.TL && c.RunLimits.TL != 0 {
		result.ReturnedValue = -1
		result.Verdict = "Time limit exceeded"
	}

	// Checking ML
	if status.Maxrss > c.RunLimits.ML && c.RunLimits.ML != 0 {
		result.ReturnedValue = -1
		result.Verdict = "Memory limit exceeded"
	}

	// Checking returned value
	if result.ReturnedValue >= 0 && result.ReturnedValue != c.RequiredRet {
		result.Verdict = "Run-time error (" + strconv.Itoa(result.ReturnedValue) + ")"
	}

	// Checking answers
	eq, err := checker(c.Path + strconv.Itoa(test_num) + ".ans", c.TestsPath + strconv.Itoa(test_num) + ".res")
	if eq && err == nil {
		result.Checker = "OK"
	}

	if !eq {
		result.Checker = "Wrong answer"
	}

	if err != nil {
		result.Checker = "Cannot check: " + err.Error()
	}

	// Result sending
	result_chan <- result
	return
}

// Function that kills process if it takes too long
func kill_too_long(pid *os.Process, RunLimits Limits, done chan string) {
	// Wait twice a real time limit
	time.Sleep(2*RunLimits.RTL)

	// If process is still active, kill it
	if _, err := os.FindProcess(pid.Pid); err == nil {
		pid.Kill()
		done <- "Process took too long"
	} else {
		done <- "OK"
	}
}

// Checks two files for equivalence
func checker(file_name_1, file_name_2 string) (res bool, err error) {
	file_1, err := os.OpenFile(file_name_1, os.O_RDONLY, 0766)
	if err != nil {
		return false, errors.New("checker: cannot open file: " + file_name_1)
	}
	defer file_1.Close()

	file_2, err := os.OpenFile(file_name_2, os.O_RDONLY, 0766)
	if err != nil {
		return false, errors.New("checker: cannot open file: " + file_name_2)
	}
	defer file_2.Close()

	first_scan := bufio.NewScanner(file_1)
	second_scan := bufio.NewScanner(file_2)

	for first_scan.Scan() {
		second_scan.Scan()
		if !bytes.Equal(first_scan.Bytes(), second_scan.Bytes()) {
			return false, nil
		}
	}

	return true, nil
}