// Package compileandrun : package for testing the compiled programming languages
package compileandrun

import (
	"fmt"
	"os"				// OS syscalls
	"os/exec"			// Executive file search
	"strconv"			// String convertation
	"syscall"			// OS syscalls
	"time"				// Timing
	"bufio"				// Buffer input/output
	"bytes"				// Bytes manipulations
	"errors"			// Errors manipulation
	"encoding/csv"		// Protocoling
	"runtime"			// Goroutines control
)


/*
	Structures
*/

// Limits : Structure with limits of a run
type Limits struct {
	TL				time.Duration	// -- time limit
	RTL 			time.Duration	// -- real time limit
	ML				int64			// -- memory limit
}

// Init : Initializing structure
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

// Result : Reporting structure
type Result struct {
	TestNumber		int				// -- number of used test
	ReturnedValue	int				// -- value returned from the process
	Time			float64			// -- used time
	Memory			int64			// -- used memory
	Verdict			string			// -- result of limits and returned value checking
	Checker			string			// -- result of answer checking
}


/*
	Methods
*/


// Compile : Method that compiles source code into the binary file
// in ".elf" format. The source must be named like
// "sol_<solution_number><format>", where <format> matches
// the chosen programming language. The resuling binary
// file will be named like "res_<solution_number>.elf"
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

// Run : Method that runs the binary file like "res_<solution_number>.elf".
// It uses tests from Init.TestsPath like input if RunAttr.Files == <nil>.
// Running for different tests is doing in different goroutines
func (c *Init) Run() (err error) {
	// Path to binary file
	runName := c.Path + "res_" + strconv.Itoa(c.Solution) + ".elf"

	// Channel and list for testing results
	resultChan := make(chan Result)
	resultList := make([]Result, c.TestsNumber)

	// Number of CPUs
	numCPU := runtime.NumCPU()
	for i := 1; i <= c.TestsNumber; {
		for j := 0; j < numCPU - 1 && i <= c.TestsNumber; j++ {
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

			args[0] = runName

			args = append(args, c.RunArgs...)

			// Running
			go c.run(runName, args, &RunAttr, i, resultChan)
			
			// Recieving information from goroutines
			result := <- resultChan

			// Checking answers
			eq, err := checker(c.Path + strconv.Itoa(c.TestsNumber) + ".ans",
								c.TestsPath + strconv.Itoa(c.TestsNumber) + ".res")
			if eq && err == nil {
				result.Checker = "OK"
			}

			if !eq {
				result.Checker = "Wrong answer"
			}

			if err != nil {
				result.Checker = "Cannot check: " + err.Error()
			}

			resultList[result.TestNumber - 1] = result
			i++
		}
	}

	// Protocoling
	protocol, err := os.OpenFile(c.Path + strconv.Itoa(c.Solution) + ".prot", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0766)
	csvWriter := csv.NewWriter(bufio.NewWriter(protocol))
	for i := range resultList {
		csvWriter.Write([]string{
			strconv.Itoa(resultList[i].TestNumber),
			strconv.Itoa(resultList[i].ReturnedValue),
			fmt.Sprint(resultList[i].Time),
			fmt.Sprint(resultList[i].Memory),
			resultList[i].Verdict,
			resultList[i].Checker,
		})
		if err = csvWriter.Error(); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	err = csvWriter.Error()
	return err
}


/*
	Helping functions
*/

// Function, that runs a process and checks limits
func (c *Init) run(runName string, RunArgs []string, RunAttr *os.ProcAttr,
		testNum int, resultChan chan Result) {
	// Setting timer
	startTime := time.Now()

	// Starting process
	pid, err := os.StartProcess(runName, RunArgs, RunAttr)

	if err != nil {
		resultChan <- Result {
			TestNumber:		testNum,
			ReturnedValue:	-1,
			Verdict:		"Cannot start run: " + err.Error(),
		}
		return
	}

	// If it take too long, kill
	done := make(chan string)
	go killTooLong(pid, c.RunLimits, done)

	// Scanning process info
	pState, err := pid.Wait()

	if err != nil {
		resultChan <- Result {
			TestNumber:		testNum,
			ReturnedValue:	-1,
			Verdict:		"Cannot wait for running end: " + err.Error(),
		}
		return
	}

	// Stopping timer
	endTime := time.Now()
	
	_ = <- done

	// Standart result
	result := Result {
		TestNumber:		testNum,
		ReturnedValue:	pState.ExitCode(),
		Verdict:		"OK",
	}

	// Process resources usage
	status := pState.SysUsage().(*syscall.Rusage)

	// Checking RTL
	if realTime := endTime.Sub(startTime); realTime > c.RunLimits.RTL && c.RunLimits.RTL != 0 {
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

	result.Time = time.Duration(status.Utime.Usec).Seconds()
	result.Memory = status.Maxrss

	// Result sending
	resultChan <- result
	return
}

// Function that kills process if it takes too long
func killTooLong(pid *os.Process, RunLimits Limits, done chan string) {
	// Wait ten times a real time limit
	for i := 0; i < 100; i++ {
		time.Sleep(10*RunLimits.RTL/100)
		// If process is still active, kill it
		if _, err := os.FindProcess(pid.Pid); err == nil {
			pid.Kill()
			done <- "Process took too long"
		} else {
			done <- "OK"
		}
	}
}

// Checks two files for equivalence
func checker(fileName1, fileName2 string) (res bool, err error) {
	file1, err := os.OpenFile(fileName1, os.O_RDONLY, 0766)
	if err != nil {
		return false, errors.New("checker: cannot open file: " + fileName1)
	}
	defer file1.Close()

	file2, err := os.OpenFile(fileName2, os.O_RDONLY, 0766)
	if err != nil {
		return false, errors.New("checker: cannot open file: " + fileName2)
	}
	defer file2.Close()

	firstScan := bufio.NewScanner(file1)
	secondScan := bufio.NewScanner(file2)

	for firstScan.Scan() {
		secondScan.Scan()
		if !bytes.Equal(firstScan.Bytes(), secondScan.Bytes()) {
			return false, nil
		}
	}

	return true, nil
}