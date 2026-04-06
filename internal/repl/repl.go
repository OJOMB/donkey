package repl

import (
	"bufio"
	"io"

	"github.com/OJOMB/donkey/internal/lexer"
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

const Prompt = ">> "

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

	// loop indefinitely, reading input from the user and processing it until we encounter a SIGINT (Ctrl+C)
	for {
		if _, err := r.out.Write([]byte(Prompt)); err != nil {
			r.logger.Error("failed to write prompt", "error", err)
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
				if _, err := r.out.Write([]byte("parser error: " + err + "\n")); err != nil {
					r.logger.Error("failed to write parser error", "error", err)
					return
				}
			}

			continue
		}

		if _, err := r.out.Write([]byte(program.String() + "\n")); err != nil {
			r.logger.Error("failed to write program output", "error", err)
			return
		}
	}
}
