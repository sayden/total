package total

//go:generate go run golang.org/x/tools/cmd/goyacc -l -o gen_parser.go parser.y

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

var word = regexp.MustCompile("^[a-zA-Z_0-9,\\.@ ]+$")
var integer = regexp.MustCompile("^[0-9]+$")

type lexer struct {
	identRune          func(ch rune, i int) bool
	originalMode       uint
	originalWhitespace uint64

	total total
	err   error

	ignoreNewline bool

	isLongText bool
	isWord     bool
	longText   string

	tokenStack []int

	s scanner.Scanner
}

func parse(input []byte) (*total, error) {
	l := newLex(input)
	_ = yyParse(l)
	return &l.total, l.err
}

func newLex(input []byte) *lexer {
	p := lexer{}
	p.init(string(input))

	return &p
}

func (p *lexer) init(s string) {
	p.s.Mode = scanner.ScanStrings
	p.s.Init(strings.NewReader(s))
	p.s.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	p.ignoreNewline = true
	p.tokenStack = make([]int, 0)

	p.identRune = func(ch rune, i int) bool {
		return ch == '>' || ch == '|' || ch == '<' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
	}

	p.s.IsIdentRune = p.identRune
	p.originalMode = p.s.Mode
	p.originalWhitespace = p.s.Whitespace
}

func (p *lexer) popStack() int {
	last := p.tokenStack[len(p.tokenStack)-1]
	p.tokenStack = p.tokenStack[0 : len(p.tokenStack)-1]
	return last
}

func (p *lexer) push(n int) {
	p.tokenStack = append(p.tokenStack, n)
}

// Lex (MANDATORY) satisfies yyLexer. Called every time the parser wants a new token
// which MUST be placed in its respective variable on yySymType (ch or part in this example)
func (p *lexer) Lex(lval *yySymType) int {
	//if len(p.tokenStack )> 0 {
	//	return p.popStack()
	//}
	//
	//if p.isLongText {
	//	p.longText = p.longTextCapture()
	//	p.isLongText = false
	//	lval.string = p.longText
	//	return TEXT
	//} else if p.longText != "" {
	//	p.longText = ""
	//	return CLT
	//}
	//
	//if p.isWord {
	//	p.isWord = false
	//	p.ignoreNewline = true
	//	p.restoreLexer()
	//	return NL
	//}

	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		txt := p.s.TokenText()
		//if txt == "\n" && p.ignoreNewline {
		//	continue
		//}
		//
		//switch txt {
		//case "|>":
		//	p.prepareLongText()
		//	p.isLongText = true
		//	return OLT
		//case "<|":
		//	p.restoreLexer()
		//	return CLT
		//}

		switch tok {
		case '\n':
			if p.ignoreNewline {
				continue
			}
			p.restoreLexer()
			return NL
		case '{':
			return OP
		case '}':
			return CL
		case ':':
			p.prepareValueCapture()
			return COLON
		case '[':
			return OB
		case ']':
			return CB
		}

		_, n := p.captureValue(txt, lval)
		// -1 represents an unknown token
		if n != -1 {
			return n
		}

		return int(tok)
	}

	return 0
}

func (l *lexer) prepareValueCapture() {
	l.ignoreNewline = false
	l.s.Whitespace = scanner.GoWhitespace
	l.s.IsIdentRune = func(ch rune, i int) bool {
		return ch != '\n'
	}

	l.s.Whitespace = 0
}

func (l *lexer) prepareLongText() {
	l.s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '|' || ch == '<'
	}

	l.s.Whitespace = 0
}

func (l *lexer) prepareWord() {
	l.ignoreNewline = false
	l.s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '\n'
	}

	l.s.Whitespace = 0
}

func (p *lexer) longTextCapture() string {
	longtext := ""
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		txt := p.s.TokenText()

		if txt == "<|" {
			p.restoreLexer()
			break
		}

		longtext += txt
	}

	return longtext
}

func (p *lexer) wordCapture(s string) string {
	word := s
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		txt := p.s.TokenText()

		if tok == '\n' {
			break
		}

		word += txt
	}

	return word
}

func (p *lexer) captureList(lval *yySymType) int {
	list := make(values, 0)

	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		if tok == ']' {
			lval.list = list
			return LIST
		}

		v, n := p.captureValue(p.s.TokenText(), lval)

		list = append(list, &value{kind: n, data: v})
	}

	lval.list = list
	return LIST
}

// captureValue attempts to return a number or a string and the representative token.
// It returns -1 if none is found
func (p *lexer) captureValue(txt string, lval *yySymType) (interface{}, int) {
	txt = strings.TrimSpace(txt)
	txt = strings.TrimFunc(txt, func(r rune) bool {
		return r == '"'
	})

	// check if it's an integer and parse it
	if integer.MatchString(txt) {
		n, err := strconv.Atoi(txt)
		lval.integer = n
		if err != nil {
			fmt.Printf("error getting integer number from '%s': %s\n", txt, err.Error())
		}
		return n, INTEGER
	}

	//check if it's a string and parse it
	if word.MatchString(txt) {
		return p.checkString(txt, lval)
	}

	return nil, -1
}

func (p *lexer) checkString(txt string, lval *yySymType) (interface{}, int) {
	switch txt {
	case "null":
		lval.nulltype = nil
		return nil, NULLTYPE
	case "true":
		lval.boolean = true
		return true, BOOLEAN
	case "false":
		lval.boolean = false
		return false, BOOLEAN
	}


	lval.string = strings.TrimSpace(txt)
	return txt, WORD
}

func (p *lexer) restoreLexer() {
	p.ignoreNewline = true
	p.s.IsIdentRune = p.identRune
	p.s.Mode = p.originalMode
	p.s.Whitespace = p.originalWhitespace
}

// Error (MANDATORY) satisfies yyLexer.
func (l *lexer) Error(s string) {
	l.err = errors.New(s)
}
