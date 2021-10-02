package total

import "errors"

type Total struct {
	docName string
	data    interface{}
}

func (t *Total)getData()(interface{}, error){
	return nil, errors.New("not implemented")
}

func (t *Total)mustGetData()interface{}{
	return nil
}

func (t *Total) getRoot() object {
	return t.data.(object)
}

func (t *Total) getList() values {
	return t.data.(values)
}
