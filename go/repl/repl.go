package repl

import (
	"bufio"
	"fmt"
	"io"

	"alde.nu/mint/compiler"
	"alde.nu/mint/lexer"
	"alde.nu/mint/parser"
	"alde.nu/mint/vm"
	"github.com/fatih/color"
)

const PROMPT = ">> "

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	num := 0
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "quit" {
			fmt.Fprint(out, "Exiting...\n")
			break
		}
		l := lexer.Create(line)
		p := parser.Create(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			io.WriteString(out, red(fmt.Sprintf("Woops! Compilation failed:\n%s\n", err)))
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			io.WriteString(out, red(fmt.Sprintf("Woops! Executing bytecode failed:\n%s\n", err)))
			continue
		}

		stackTop := machine.StackTop()

		if stackTop != nil {
			num += 1

			index := fmt.Sprintf("%s%s%s ", green("["), yellow(num), green("]"))
			io.WriteString(out, index)
			io.WriteString(out, stackTop.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, red("Parser Errors:\n"))
	for _, msg := range errors {
		err := fmt.Sprintf("\t%s\n", red(msg))
		io.WriteString(out, err)
	}
}
