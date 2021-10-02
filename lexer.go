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


type lexer struct {
	identRune          func(ch rune, i int) bool
	originalMode       uint
	originalWhitespace uint64

	total total
	err   error

	isLongText bool
	longText   string

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
	p.s.Mode = scanner.ScanRawStrings
	p.s.Init(strings.NewReader(s))

	p.identRune = func(ch rune, i int) bool {
		return ch == '>' || ch == '|' || ch == '<' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
	}

	p.s.IsIdentRune = p.identRune
	p.originalMode = p.s.Mode
	p.originalWhitespace = p.s.Whitespace
}

var word = regexp.MustCompile("[a-zA-Z_0-9]+")
var integer = regexp.MustCompile("^[0-9]+$")

// Lex (MANDATORY) satisfies yyLexer. Called every time the parser wants a new token
// which MUST be placed in its respective variable on yySymType (ch or part in this example)
func (p *lexer) Lex(lval *yySymType) int {
	if p.isLongText {
		p.longText = p.longTextCapture()
		p.isLongText = false
		lval.string = p.longText
		return TEXT
	} else if p.longText != "" {
		p.longText = ""
		return CLT
	}

	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		txt := p.s.TokenText()
		_, n := p.captureValue(txt, lval)
		if n != -1 {
			return n
		}

		switch txt {
		case "|>":
			p.prepareLongText()
			p.isLongText = true
			return OLT
		case "<|":
			p.restoreLexer()
			return CLT
		}

		switch tok {
		case '{':
			return OP
		case '}':
			return CL
		case ':':
			return COLON
		case '[':
			return OB
		case ']':
			return CB
		}

		return int(tok)
	}

	return 0
}

func (l *lexer) prepareLongText() {
	l.s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '|' || ch == '<'
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

func (p *lexer) captureValue(txt string, lval *yySymType) (interface{}, int) {
	if integer.MatchString(txt) {
		n, err := strconv.Atoi(txt)
		lval.integer = n
		if err != nil {
			fmt.Printf("error getting integer number from '%s': %s\n", txt, err.Error())
		}
		return n, INTEGER
	}

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

	lval.string = txt
	return txt, WORD
}

func (p *lexer) restoreLexer() {
	p.s.IsIdentRune = p.identRune
	p.s.Mode = p.originalMode
	p.s.Whitespace = p.originalWhitespace
}

// Error (MANDATORY) satisfies yyLexer.
func (l *lexer) Error(s string) {
	l.err = errors.New(s)
}
