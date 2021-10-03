package total

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentDocumentName(t *testing.T) {
	table := map[string]string{
		"user{":  "user",
		"user {": "user",
		"u":      "u",
		"u{":     "u",
		" user{": "user",
		"a_user": "a_user",
		"user5{": "user5",
	}
	for k, v := range table {
		s := NewScanner(k)
		s.identDocumentName()
		_, text := s.Scan()
		assert.Equal(t, v, text)
	}
}

func TestIdentLhsName(t *testing.T) {
	table := map[string]string{
		" hello :":    "hello",
		" hello: ":    "hello",
		" hello:":     "hello",
		"host.name :": "host.name",
	}
	for k, v := range table {
		s := NewScanner(k)
		s.identLHS()
		_, text := s.Scan()
		assert.Equal(t, v, text)
		_, text = s.Scan()
		assert.Equal(t, ":", text)
	}
}

func TestIdentRhsName(t *testing.T) {
	fmt.Println('a','z','A','Z','0','9')
	table := map[string]string{
		" {\n":                  "{",
		" [\n":                  "[",
		"hello\n":               "hello",
		"hello func world\n":    "hello func world",
		" hello\n":              "hello",
		"host\".\"name\n":       `host"."name`,
		">> some long text<<\n": "some long text",
		` >> some
spec.ially
long text<<\n`: `some
spec.ially
long text`,
	}
	for k, v := range table {
		s := NewScanner(k)
		s.identRHS()
		_, text := s.Scan()
		assert.Equal(t, v, text)
	}
}
