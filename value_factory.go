package main

type Value struct {
	kind int
	data interface{}
}

type Values []*Value

type Pair struct {
	name  string
	value *Value
}

type Object []*Pair

func newObject(v ...*Pair) Object {
	if v == nil || len(v) == 0 {
		return make(Object, 0)
	}

	return v
}
