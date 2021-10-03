package total

import "strings"

type myscanner struct {
	input []byte

	pos int

	identifier func(byte) bool
	separator  func(byte) bool

	currentToken []byte

	longTextMode bool
}

func newMyScanner(i []byte) *myscanner {
	return &myscanner{
		input:        i,
		currentToken: make([]byte, 0),
		identifier:   lhsIdentifiers,
		separator:    lhsSeparators,
	}
}

func (m *myscanner) Scan() string {
	for ; m.pos < len(m.input); m.pos++ {
		tok := m.input[m.pos]
		//fmt.Printf("Tok: '%d' '%s'\n", tok, string(tok))

		//If any of the configured identifiers, append to current token
		if m.identifier(tok) {
			m.currentToken = append(m.currentToken, tok)

			// If next token is not an identifier, it might be an ignored char or a separator, but this token ends here
			if !m.identifier(m.input[m.pos+1]) && !m.longTextMode {
				temp := m.currentToken
				m.currentToken = m.currentToken[0:0]
				m.pos++
				return strings.TrimSpace(string(temp))
			} else if m.longTextMode && m.input[m.pos+1] == '>' && m.input[m.pos+2] == '\n' {
				// End long text mode here
				temp := m.currentToken
				m.currentToken = m.currentToken[0:0]
				m.pos += 2
				return strings.TrimSpace(string(temp))
			}

			continue
		}

		if m.separator(tok) {
			// fmt.Printf("Separator: '%d'\n", tok). When a separator is found, switch to RHS rules RHS rules must handle
			// the potential exception of finding a long text
			if tok == ':' {
				if string(m.input[m.pos+1:m.pos+2]) == ">" || string(m.input[m.pos+1:m.pos+3]) == " >" {
					// Long text. It requires a new "mode" to recognize 2 tokens
					// as the final identifier
					m.separator = longTextSeparator
					m.identifier = longTextIdentifiers
					m.longTextMode = true
				} else {
					// Normal RHS
					m.separator = rhsSeparators
					m.identifier = rhsIdentifiers
				}
			} else if tok == '\n' {
				// RHS finished
				m.separator = lhsSeparators
				m.identifier = lhsIdentifiers
			}
			m.pos++
			return string(tok)
		}
	}

	return ""
}

func lhsSeparators(r byte) bool {
	return r == '{' || r == '}' || r == '[' || r == ']' || r == ':'
}

func rhsSeparators(r byte) bool {
	return r == '{' || r == '}' || r == '[' || r == ']' || r == '\n'
}

func longTextSeparator(r byte) bool {
	return r == '<'
}

func longTextIdentifiers(r byte) bool {
	return r != '<'
}

func lhsIdentifiers(r byte) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_'
}

func rhsIdentifiers(r byte) bool {
	return r != '\n'
}
