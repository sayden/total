package total

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocument(t *testing.T) {
	yyErrorVerbose = true
	//yyDebug = 2

	v, err := parse([]byte(complexObject))
	assert.NoError(t, err)

	assert.Equal(t, "user", v.docName)
	root := v.data
	assert.Equal(t, OBJECT, root.kind)
	obj := root.data.(object)

	longText := obj[0]
	assert.Equal(t, "long_text", longText.name)
	assert.Equal(t, TEXT, longText.getKind())
	assert.Equal(t, `'this' is "maaaadnessss" \n `, longText.value.data.(string))

	number := obj[3]
	assert.Equal(t, INTEGER, number.getKind())
	assert.Equal(t, 12, number.getData().(int))

	innerBlock := obj[8]
	assert.Equal(t, OBJECT, innerBlock.getKind())
	assert.Equal(t, "inner_block", innerBlock.getName())

}

func TestListDocument(t *testing.T) {
	in := `user [
		{
			hello: world
		}
		{
			list: [1 2 3 4]
		}
]`

	yyErrorVerbose = true
	//yyDebug = 2

	v, err := parse([]byte(in))
	assert.NoError(t, err)
	//pp.Print(v)

	assert.Equal(t, "user", v.docName)
	blockInfo := v.data
	assert.Equal(t, LIST, blockInfo.kind, "doc has type list")
	blocksList := blockInfo.data.(values)
	assert.Equal(t, OBJECT, blocksList[0].kind, "first item in this doc has a list of key-values")
	assert.Equal(t, OBJECT, blocksList[1].kind, "second item in this doc has a list of key-values")

	helloBlock := blocksList[0].data.(object)
	assert.Equal(t, "hello", helloBlock[0].name)
	assert.Equal(t, WORD, helloBlock[0].value.kind)
	assert.Equal(t, "world", helloBlock[0].value.data.(string))

	listBlock := blocksList[1].data.(object)
	assert.Len(t, listBlock, 1)
	assert.Equal(t, "list", listBlock[0].name)
	assert.Equal(t, LIST, listBlock[0].value.kind)
	itemList := listBlock[0].value.data.(values)
	assert.Len(t, itemList, 4)
	assert.Equal(t, INTEGER, itemList[0].kind)
	assert.Equal(t, 1, itemList[0].data.(int))
	assert.Equal(t, INTEGER, itemList[1].kind)
	assert.Equal(t, 2, itemList[1].data.(int))
	assert.Equal(t, INTEGER, itemList[2].kind)
	assert.Equal(t, 3, itemList[2].data.(int))
	assert.Equal(t, INTEGER, itemList[3].kind)
	assert.Equal(t, 4, itemList[3].data.(int))

	in = `user {
		hello: world
	}`

	v, err = parse([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, "user", v.docName)
	blockInfo = v.data
	assert.Equal(t, OBJECT, blockInfo.kind, "doc has type object")
	simpleBlock := blockInfo.data.(object)
	assert.Equal(t, WORD, simpleBlock[0].value.kind)
}