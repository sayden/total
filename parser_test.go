package main

import (
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	in := `user {
    key: value
	number: 12
	another_key: another_value
	OneMore4: asda
	whatifitsnull: null
	another_number: 99
}`

	v, err := Parse([]byte(in))
	assert.NoError(t, err)

	pp.Print(v)


}
