package main

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

//func main() {
//	testText := "hello world |> sasdfasdf asdf asdasdf om\n newline<| And life goes on { hello }"
//
//	p := lexer{}
//	p.init(testText)
//	p.lex()
//}

type lexer struct {
	identRune          func(ch rune, i int) bool
	originalMode       uint
	originalWhitespace uint64

	scl   Scl
	err   error

	s scanner.Scanner
}

func (p *lexer) init(s string) {
	p.s.Mode = scanner.ScanRawStrings
	p.s.Init(strings.NewReader(s))

	p.identRune = func(ch rune, i int) bool {
		return ch == '>' || ch == '|' || ch == '<' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
	}

	p.s.IsIdentRune = p.identRune
	p.originalMode = p.s.Mode
	p.originalWhitespace = p.s.Whitespace
}

func (p *lexer) Lex(lval *yySymType) int {

	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		txt := p.s.TokenText()

		if txt == "|>" {
			longtext := p.longTextCapture()
			fmt.Printf("%s: Long text: '%s'\n", p.s.Position, longtext)
			continue
		}

		fmt.Printf("%s: '%s'\n", p.s.Position, p.s.TokenText())
	}
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

func (p *lexer)restore(){
	p.s.IsIdentRune = p.identRune
	p.s.Mode = p.originalMode
	p.s.Whitespace = p.originalWhitespace
}