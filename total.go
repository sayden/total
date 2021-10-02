package total

import "errors"

type total struct {
	docName string
	data    *value
}

func (t *total) getData() (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (t *total) mustGetData() interface{} {
	return nil
}

func (t *total) getRoot() *value {
	return t.data
}
