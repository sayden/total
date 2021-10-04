package total

import (
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
		A_hello: world
		number: 99
		sentence: hello world
		long_text:>some stuff "going on" here<
		number2: 99.2
		inner: {
			hello:world
			long_text: > hello
world <
		}
		list: [1 12 123]
		list2: [a as asdf]
		inner_block: [
			{
				hello3:world
				another_hello: another world
			}
		]
	}`

	expected := []string{"user", ":", "{", "\n", "A_hello", ":", "world", "\n", "number", ":", "99", "\n", "sentence", ":", "hello world", "\n",
		"long_text", ":", `some stuff "going on" here`, "\n", "number2", ":", "99.2", "\n", "inner", ":", "{", "\n", "hello", ":", "world", "\n",
		"long_text", ":", `hello
world`, "\n", "}", "\n", "list", ":", "[", "1", "12", "123", "]", "\n", "list2", ":", "[", "a", "as", "asdf", "]", "\n", "inner_block",
		":", "[","\n", "{", "\n", "hello3", ":", "world", "\n", "another_hello", ":", "another world", "\n", "}", "\n", "]", "\n", "}"}
	_ = expected

	s := newMyScanner([]byte(text))
	for i := 0; ; i++ {
		tok := s.Scan()
		if tok == "" {
			break
		}
		//fmt.Printf("Token found '%s'\n", tok)
		assert.Equal(t, expected[i], tok, i)
	}
}
