package total

import (
	"bytes"
	"regexp"
	"strings"
	"text/scanner"
	"unicode"
)

var openingLongText = regexp.MustCompile(`^>>`)

type scannerW struct {
	inner scanner.Scanner
}

func NewScanner(init string) *scannerW {
	inner := scanner.Scanner{}
	inner.Init(bytes.NewBufferString(init))
	return &scannerW{inner: inner}
}

func (sc *scannerW) Scan() (tok rune, tokText string) {
	tok = sc.inner.Scan()
	tokText = sc.inner.TokenText()
	tokText = strings.TrimSpace(tokText)

	// Check for the presence of a long text opening '>>'
	if len(tokText) > 2 && (tokText[0:2] == ">>" || tokText[0:3] == " >>") {
		sc.identLongText()
		tok = sc.inner.Scan()
		tokText2 := sc.inner.TokenText()
		tokText = strings.TrimSpace(tokText + "\n" + tokText2)
		tokText = strings.Replace(tokText, "<<", "", 1)
		tokText = strings.Replace(tokText, ">>", "", 1)
		tokText = strings.TrimSpace(tokText)
		sc.identLHS()
	}

	// Check for the beginning of a new block

	// Check for the beginning of a list

	return
}

func (sc *scannerW) identLongText() {
	sc.inner.Mode = 0
	sc.inner.Whitespace = scanner.GoWhitespace
	sc.inner.IsIdentRune = func(ch rune, i int) bool {
		return ch != '<'
	}
}

func (sc *scannerW) identDocumentName() {
	sc.inner.Mode = scanner.GoTokens
	sc.inner.Whitespace = scanner.GoWhitespace
	sc.inner.IsIdentRune = nil
}

// identLHS configures a scanner to scan correctly the left hand side
// of the document. LHS identifiers are limited to alphanumeric chars,
// underscores and dots.
func (sc *scannerW) identLHS() {
	sc.inner.Mode = scanner.GoTokens
	sc.inner.Whitespace = scanner.GoWhitespace
	sc.inner.IsIdentRune = func(ch rune, i int) bool {
		return ch == '_' || ch == '.' || unicode.IsDigit(ch) || unicode.IsLetter(ch)
	}
}

// identRHS configures a scanner to scan correctly the right hand Side
// of the document. RHS ends with newline except for the backstick version
// which requires a backstick before the newlinw:
// 		text with spaces\n
//		"quoted text"\n
//		99\n
//		89.9\n
// 		`backstick text with \n newlines`\n
func (sc *scannerW) identRHS() {
	sc.inner.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments
	sc.inner.Whitespace = 1 << '\t'
	//sc.inner.Whitespace = scanner.GoWhitespace
	sc.inner.IsIdentRune = func(ch rune, i int) bool {
		return ch != '\n' && ch != '\r'
	}
}
