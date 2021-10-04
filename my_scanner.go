package total

import (
	"strings"
)

type myscanner struct {
	input []byte

	pos int

	err error
	total total

	identifier func(byte) bool
	separator  func(byte) bool

	currentToken []byte

	tokenStack []string
}

func newMyScanner(i []byte) *myscanner {
	return &myscanner{
		input:        i,
		currentToken: make([]byte, 0),
		identifier:   lhsIdentifiers,
		separator:    lhsSeparators,
		tokenStack:   make([]string, 0),
	}
}

func (m *myscanner) first() string {
	length := len(m.tokenStack)
	if length == 0 {
		return ""
	}

	first := m.tokenStack[0]
	m.tokenStack = m.tokenStack[1:]
	return first
}

func (m *myscanner) push(s string) {
	m.tokenStack = append(m.tokenStack, s)
}

func (m *myscanner) isEmpty() bool {
	return len(m.tokenStack) == 0
}

func (m *myscanner) Scan() string {
	for !m.isEmpty() {
		return m.first()
	}

	for ; m.pos < len(m.input); m.pos++ {
		tok := m.input[m.pos]

		//If any of the configured identifiers, append to current token
		if m.identifier(tok) {
			m.currentToken = append(m.currentToken, tok)

			// If next token is not an identifier, it might be an ignored char or a separator, but this token ends here
			if !m.identifier(m.input[m.pos+1]) {
				temp := string(m.currentToken)
				m.currentToken = make([]byte,0)
				m.pos++
				m.push(strings.TrimSpace(temp))
				return m.first()
			}

			continue
		}

		// When a separator is found, switch to RHS which must handle the potential exception of finding a long text
		if m.separator(tok) {
			if tok == ':' {
				// Add colon to stack
				m.push(":")
				m.pos++

				m.omitLeadingSpaces()

				tok = m.input[m.pos]

				switch tok {
				case '>':
					// Long text found.
					m.pos++

					// Long text. It requires a new "mode" to recognize 2 tokens as the final identifier
					for ; m.pos < len(m.input); m.pos++ {
						if string(m.input[m.pos:m.pos+2]) != "<\n" {
							m.currentToken = append(m.currentToken, m.input[m.pos])
						} else {
							temp := m.currentToken
							m.currentToken = m.currentToken[0:0]
							m.pos+=2
							m.push(strings.TrimSpace(string(temp)))
							m.push("\n")
							return m.first()
						}
					}
					break
				case '[':
					// List, can be anything
					m.separator = rhsSeparators
					m.identifier = lhsIdentifiers
					m.push("[")
					m.pos++
				case '{':
					// List, can be anything
					m.separator = rhsSeparators
					m.identifier = lhsIdentifiers
					m.push("{")
					m.pos++
				default:
					// Normal RHS
					m.separator = rhsSeparators
					m.identifier = rhsIdentifiers
				}

			} else if tok == ']' {
				m.push("]")
				m.pos++
				m.separator = lhsSeparators
				m.identifier = lhsIdentifiers
			} else if tok == '\n' {
				m.push("\n")
				m.pos++
				m.separator = lhsSeparators
				m.identifier = lhsIdentifiers
			} else if tok == '}' {
				m.push("}")
				m.pos++
				m.separator = lhsSeparators
				m.identifier = lhsIdentifiers
			} else {
				m.pos++
				m.push(string(tok))
			}

			return m.first()
		}

		// No identifier or separator found, omit
	}

	return ""
}

func (m *myscanner) omitLeadingSpaces() {
	for ; m.pos < len(m.input); m.pos++ {
		if m.input[m.pos] != ' ' {
			return
		}
	}
}

func lhsSeparators(r byte) bool {
	return r == '{' || r == '}' || r == '[' || r == ']' || r == ':' || r == '\n'
}

func rhsSeparators(r byte) bool {
	return r == '{' || r == '}' || r == '[' || r == ']' || r == '\n'
}

func lhsIdentifiers(r byte) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_'
}

func rhsIdentifiers(r byte) bool {
	return r != '\n'
}
