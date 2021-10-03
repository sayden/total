package total

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLhsRule(t *testing.T) {
	assert.True(t, lhsIdentifiers('a'))
	assert.True(t, lhsIdentifiers('Z'))
	assert.True(t, lhsIdentifiers('f'))
	assert.True(t, lhsIdentifiers('8'))
	assert.True(t, lhsIdentifiers('0'))
	assert.True(t, lhsIdentifiers('A'))
	assert.True(t, lhsIdentifiers('z'))
	assert.True(t, lhsIdentifiers('.'))
	assert.True(t, lhsIdentifiers('_'))

	assert.False(t, lhsIdentifiers('{'))
	assert.False(t, lhsIdentifiers('}'))
	assert.False(t, lhsIdentifiers('*'))
	assert.False(t, lhsIdentifiers('-'))
	assert.False(t, lhsIdentifiers('^'))
	assert.False(t, lhsIdentifiers('$'))
}

func TestScan(t *testing.T) {
	text := `user: {
		hello: world
		number: 99
		sentence: hello world
		long_text:>some stuff "going on" here<
		number2: 99.2
		inner: {
			hello:world
			long_text: > hello 
world <
		}
	}`

	s := newMyScanner([]byte(text))
	for {
		tok := s.Scan()
		if tok == "" {
			break
		}
		fmt.Println("Token found:", tok)
	}
}
