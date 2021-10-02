package total

import "fmt"

type object []*keyValue

func newObject(v ...*keyValue) object {
	if v == nil || len(v) == 0 {
		return make(object, 0)
	}

	return v
}

type keyValue struct {
	name  string
	value *value
}

func (p *keyValue) getName() string {
	return p.name
}

func (p *keyValue) getKind() int {
	return p.value.kind
}

func (p *keyValue) getData() interface{} {
	return p.value.data
}

func (p *keyValue) String() string {
	return fmt.Sprintf("Name: %s, ValueObj: %s", p.name, p.value)
}

type value struct {
	kind int
	data interface{}
}

type values []*value

// String produces a debug message without newline of the 'kind'
func (v *value) String() string {
	return fmt.Sprintf("Kind: %s", tokName(v.kind))
}
