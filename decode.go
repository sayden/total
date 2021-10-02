package total

import "errors"

func UnmarshalTotal(data []byte) (*Total, error) {
	return parse(data)
}

func Unmarshal(data []byte, v interface{}) (name string, err error) {
	return "", errors.New("not implemented")
}

func MarshalTotal(t *Total) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func Marshal(n string, v interface{}) ([]byte, error){
	return nil, errors.New("not implemented")
}

func traverse(t *Total){
	// Do something with name
	//t.docName

	// Data is either a block or a list of blocks
	//t.data
}