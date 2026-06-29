package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/OJOMB/donkey/internal/evaluator"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/parser"
	"github.com/OJOMB/donkey/pkg/logs"
)

const donkeyASCII = `
                          /\          /\
                         ( \\        // )
                          \ \\      // /
                           \_\\||||//_/
                            \/ _  _ \
                           \/|(O)(O)|
                          \/ |      |
      ___________________\/  \      /
     //                //     |____|
    //                ||     /      \
   //|                \|     \ 0  0 /
  // \       )         V    / \____/
 //   \     /        (     /
""     \   /_________|  |_/
       /  /\   /     |  ||
      /  / /  /      \  ||
      | |  | |        | ||
      | |  | |        | ||
      |_|  |_|        |_||
       \_\  \_\        \_\\
`

const inPromptTemplate = "In [%d]: "
const outPromptTemplate = "Out [%d]: "

type Repl struct {
	in  io.Reader
	out io.Writer

	logger logs.Logger
}

func New(in io.Reader, out io.Writer, l logs.Logger) *Repl {
	if l == nil {
		l = logs.NewNullLogger()
	}

	return &Repl{in: in, out: out, logger: l.With("component", "repl")}
}

func (r *Repl) Start() {
	// create a new scanner to read input from the user
	scanner := bufio.NewScanner(r.in)

	_, _ = r.out.Write([]byte(donkeyASCII))
	_, _ = r.out.Write([]byte("Welcome to the Donkey programming language!\n"))

	// this is the global env
	var env *objects.Environment
	var inputLine = 1
	// loop indefinitely, reading input from the user and processing it until we encounter a SIGINT (Ctrl+C)
	for {
		if _, err := fmt.Fprintf(r.out, inPromptTemplate, inputLine); err != nil {
			r.logger.Error("failed to write input prompt", "error", err)
			return
		}

		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				r.logger.Error("failed to read input", "error", err)
			}

			return
		}

		line := scanner.Text()
		l := lexer.New(line, r.logger)

		p, err := parser.New(l, r.logger)
		if err != nil {
			r.logger.Error("failed to create parser", "error", err)
			return
		}

		program := p.ParseProgram()
		if len(p.Errors) != 0 {
			for _, err := range p.Errors {
				if _, err := fmt.Fprintf(r.out, outPromptTemplate, inputLine); err != nil {
					r.logger.Error("failed to write output prompt", "error", err)
					return
				}

				if _, err := fmt.Fprintf(r.out, "Parser error: %s\n", err); err != nil {
					r.logger.Error("failed to write parser error", "error", err)
					return
				}
			}

			continue
		}

		// now we have a valid AST, we can evaluate it and print the result
		e := evaluator.New(r.logger)
		env = objects.NewEnclosedEnvironment(env)
		result := e.Eval(program, env)
		if result != nil {
			if _, err := fmt.Fprintf(r.out, outPromptTemplate, inputLine); err != nil {
				r.logger.Error("failed to write output prompt", "error", err)
				return
			}

			if _, err := fmt.Fprintf(r.out, "%s\n", result.Inspect()); err != nil {
				r.logger.Error("failed to write evaluation result", "error", err)
				return
			}
		}

		inputLine++
	}
}
