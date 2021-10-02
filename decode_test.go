package total

import (
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"testing"
)

var complexObject = `user {
		long_text: |>'this' is "maaaadnessss" \n <|
		list_of_strings: [as asdf asdfa]
		key: value
		number: 12
		another_key: "another value"
		OneMore4: asda
		whatifitsnull: null
		another_number: 99
		inner_block {
			inner_key: inner_value
			list_of_numbers: [1 2 3]
			more_inner {
				hello: world
			}
			list_of_blocks: [
				{
					list_of_numbers: [4 5 6]
				}
				{
					hello: world
				}
			]
		}
	}`

func TestMarshalFromStruct(t *testing.T) {
	testStruct := struct {
		Key        string `json:"key"`
		AnotherKey string `json:"another_key"`
		Number     int    `json:"number"`
	}{
		Key:        "a_key",
		AnotherKey: "a value",
		Number:     99,
	}

	yyErrorVerbose = true
	byt, err := Marshal("user", testStruct)
	assert.NoError(t, err)
	pp.Println(string(byt))

	var m2 map[string]interface{}
	err = Unmarshal(byt, &m2)
	assert.NoError(t, err)

	pp.Print(m2)

	assert.Equal(t, "a_key", m2["key"])
	assert.Equal(t, "a value", m2["another_key"])
	assert.Equal(t, 99.0, m2["number"].(float64))

}

func TestMarshalFromMapStringInterface(t *testing.T) {
	var m map[string]interface{}
	err := Unmarshal([]byte(complexObject), &m)
	assert.NoError(t, err)

	totalByt, err := Marshal("user", m)
	assert.NoError(t, err)

	var m2 map[string]interface{}
	err = Unmarshal(totalByt, &m2)
	assert.NoError(t, err)

	assert.Equal(t, "value", m2["key"])
	assert.Equal(t, "world", m2["inner_block"].(map[string]interface{})["more_inner"].(map[string]interface{})["hello"])
}

func TestUnmarshalIntoMapStringInterface(t *testing.T) {
	var m map[string]interface{}
	err := Unmarshal([]byte(complexObject), &m)
	assert.NoError(t, err)

	assert.Equal(t, "value", m["key"])
	assert.Equal(t, "world", m["inner_block"].(map[string]interface{})["more_inner"].(map[string]interface{})["hello"])
}

func TestUnmarshalIntoStruct(t *testing.T) {
	testStruct := struct {
		Key        string `json:"key"`
		AnotherKey string `json:"another_key"`
		Number     int    `json:"number"`
	}{}
	err := Unmarshal([]byte(complexObject), &testStruct)
	assert.NoError(t, err)
	assert.Equal(t, "value", testStruct.Key)
	assert.Equal(t, 12, testStruct.Number)
}
