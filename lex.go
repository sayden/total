package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"text/scanner"
)

//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go parser.y

// Scl is the type of the parser result
type Scl struct {
	docName string
	data    map[string]interface{}
}

// Parse (optional, called from test) implements yyParser, parses the input and returs a Scl.
func Parse(input []byte) (Scl, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.scl, l.err
}

func newLex(input []byte) *lexer {
	p := lexer{}
	p.init(string(input))

	return &p
}

// MyHandmadeLexer (MANDATORY, any name works but it's referenced on parser.y:5)
// implements yyLexer Lex(*yySymType) and Error(string)
//type MyHandmadeLexer struct {
//	input []byte
//	scl   Scl
//	err   error
//
//	lastPos   int
//	lookAhead int
//
//	myScanner *scanner.Scanner
//}

var word = regexp.MustCompile("[a-zA-Z_]+")
var integer = regexp.MustCompile("[0-9]+")

// Lex (MANDATORY) satisfies yyLexer. Called every time the parser wants a new token
// which MUST be placed in its respective variable on yySymType (ch or part in this example)
func (l *MyHandmadeLexer) Lex(lval *yySymType) int {
	
	//lval.str = ""
	//
	//for tok := l.myScanner.Scan(); tok != scanner.EOF; tok = l.myScanner.Scan() {
	//	txt := l.myScanner.TokenText()
	//	fmt.Printf("%v, %s\n", l.myScanner.Pos(), txt)
	//	if word.MatchString(txt) {
	//		lval.any = txt
	//		return WORD
	//	}
	//
	//	if integer.MatchString(txt){
	//		lval.integer,_ = strconv.Atoi(txt)
	//		return INTEGER
	//	}
	//
	//	switch tok {
	//	case '{':
	//		return OP
	//	case '}':
	//		return CL
	//	case ':':
	//		return COLON
	//	}
	//
	//	return int(tok)
	//}

	return 0
}

// Error (MANDATORY) satisfies yyLexer.
func (l *MyHandmadeLexer) Error(s string) {
	l.err = errors.New(s)
}
