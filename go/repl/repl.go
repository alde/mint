package repl

import (
	"bufio"
	"fmt"
	"io"

	"alde.nu/mint/evalutator"
	"alde.nu/mint/lexer"
	"alde.nu/mint/object"
	"alde.nu/mint/parser"
	"github.com/fatih/color"
)

const PROMPT = ">> "

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

func Start(in io.Reader, out io.Writer) {
	env := object.CreateEnvironment()
	scanner := bufio.NewScanner(in)
	num := 0
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		num += 1
		line := scanner.Text()
		l := lexer.Create(line)
		p := parser.Create(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evalutator.Eval(program, env)
		if evaluated != nil {
			index := fmt.Sprintf("%s%s%s ", green("["), yellow(num), green("]"))
			io.WriteString(out, index)
			io.WriteString(out, evaluated.Inspect())
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
