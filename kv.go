package total

type object []*Kv

func newObject(v ...*Kv) object {
	if v == nil || len(v) == 0 {
		return make(object, 0)
	}

	return v
}

type Kv struct {
	name  string
	value *value
}

func (p *Kv) getName() string {
	return p.name
}

func (p *Kv) getKind() int {
	return p.value.kind
}

func (p *Kv) getData() interface{} {
	return p.value.data
}

type value struct {
	kind int
	data interface{}
}

type values []*value
