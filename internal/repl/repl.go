package repl

import (
	"bufio"
	"io"

	"github.com/OJOMB/monkey/internal/lexer"
	"github.com/OJOMB/monkey/internal/tokens"
)

const Prompt = ">> "

type Repl struct {
	in  io.Reader
	out io.Writer
}

func New(in io.Reader, out io.Writer) *Repl {
	return &Repl{in: in, out: out}
}

func (r *Repl) Start() {
	scanner := bufio.NewScanner(r.in)

	if _, err := r.out.Write([]byte("Welcome to the Monkey programming language!\n")); err != nil {
		return
	}

	for {
		if _, err := r.out.Write([]byte(Prompt)); err != nil {
			return
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != tokens.TokenTypeEOF; tok = l.NextToken() {
			if _, err := r.out.Write([]byte(tok.Lexeme + "\n")); err != nil {
				return
			}
		}
	}
}
