package total
//go:generate go run golang.org/x/tools/cmd/goyacc -l -o gen_parser.go parser.y

import (
	"errors"
	"regexp"
	"strconv"
)

var word = regexp.MustCompile("^[a-zA-Z_0-9,\\.@ ]+$")
var integer = regexp.MustCompile("^[0-9\\.]+$")

func parse(input []byte) (*total, error) {
	l := newMyScanner(input)
	_ = yyParse(l)
	return &l.total, l.err
}

func (p *myscanner) Lex(lval *yySymType) int {
	tok := p.Scan()

	switch tok {
	case "\n":
		return NL
	case "{":
		return OP
	case "}":
		return CL
	case ":":
		return COLON
	case "[":
		return OB
	case "]":
		return CB
	}

	if integer.MatchString(tok) {
		lval.integer, _ = strconv.Atoi(tok)
		return INTEGER
	}

	if len(tok)==0{
		return 0
	}

	//if word.MatchString(tok) {
		lval.string = tok
		return WORD
	//}

	//return 0
}

// Error (MANDATORY) satisfies yyLexer.
func (l *myscanner) Error(s string) {
	l.err = errors.New(s)
}
