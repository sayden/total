package main

import (
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	in := `user{
    key: value
    more_key: more_value
    list: [12 31 4]
}`

	v, err := Parse([]byte(in))
	assert.NoError(t, err)

	pp.Print(v)
}
