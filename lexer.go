package main

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

type Total struct {
	docName string
	data    map[string]interface{}
}

type lexer struct {
	identRune          func(ch rune, i int) bool
	originalMode       uint
	originalWhitespace uint64

	total Total
	err   error

	s scanner.Scanner
}

func Parse(input []byte) (Total, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.total, l.err
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
	var err error
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		txt := p.s.TokenText()

		fmt.Printf("%s: '%s'\n", p.s.Position, txt)

		if txt == "|>" {
			longtext := p.longTextCapture()
			fmt.Printf("%s: Long text: '%s'\n", p.s.Position, longtext)
			lval.any = longtext
			return WORD
		}

		if integer.MatchString(txt) {
			lval.any, err = strconv.Atoi(txt)
			if err != nil {
				fmt.Printf("error getting integer number from '%s': %s\n", txt, err.Error())
			}
			return INTEGER
		}

		if word.MatchString(txt) {
			lval.any = txt
			return p.checkString(txt, lval)
		}

		switch txt {
		case "|>":
			return OLT
		case "<|":
			return CLT
		}

		switch tok {
		case '{':
			return OP
		case '}':
			return CL
		case ':':
			return COLON
		}

		return int(tok)
	}

	return 0
}

func (p *lexer) checkString(txt string, lval *yySymType) int {
	switch txt {
	case "null":
		lval.null = nil
		return NULL
	case "true":
		lval.ch = 1
		return TRUE
	case "false":
		lval.ch = 0
		return FALSE
	}
	lval.any = txt
	return WORD
}

func (p *lexer) longTextCapture() string {
	longtext := ""

	p.s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '|' || ch == '<'
	}
	p.s.Whitespace = 0

	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		txt := p.s.TokenText()

		if txt == "<|" {
			p.restore()
			break
		}

		longtext += txt
	}

	return longtext
}

func (p *lexer) restore() {
	p.s.IsIdentRune = p.identRune
	p.s.Mode = p.originalMode
	p.s.Whitespace = p.originalWhitespace
}

// Error (MANDATORY) satisfies yyLexer.
func (l *lexer) Error(s string) {
	l.err = errors.New(s)
}
