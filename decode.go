package total

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	tot, err := unmarshalTotal(data)
	if err != nil {
		return err
	}

	bs := traverseToJSON(tot)
	return json.Unmarshal(bs, v)
}

func Marshal(n string, v interface{}) ([]byte, error) {
	bs := make([]byte, 0)
	buf := bytes.NewBuffer(bs)

	buf.WriteString(n)

	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Array:
		traverseSliceMsi(v, buf)
		break
	case reflect.Map:
		traverseMsi(v, buf)
		break
	case reflect.Struct:
		//TODO for the love of god Mario, improve this
		j, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err = json.Unmarshal(j, &m); err != nil {
			return nil, err
		}

		traverseMsi(m, buf)
		break
	default:
		return nil, errors.New("invalid value")
	}

	return buf.Bytes(), nil
}

func marshalTotal(t *total) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func unmarshalTotal(data []byte) (*total, error) {
	return parse(data)
}

func traverseMsi(i interface{}, buf *bytes.Buffer) {
	buf.WriteByte('{')
	buf.WriteByte('\n')
	defer buf.WriteByte('}')
	m := i.(map[string]interface{})

	pos := 0
	for k, v := range m {
		if v == nil {
			buf.WriteString(fmt.Sprintf("%s:null", k))
		} else {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				buf.WriteString(fmt.Sprintf("%s:", k))
				traverseSliceMsi(v, buf)
			case reflect.Array:
				buf.WriteString(fmt.Sprintf("%s:", k))
				traverseSliceMsi(v, buf)
			case reflect.Map:
				buf.WriteString(fmt.Sprintf("%s", k))
				traverseMsi(v, buf)
			case reflect.String:
				buf.WriteString(fmt.Sprintf("%s:", k))
				processText(v, buf)
				break
			case reflect.Bool:
				buf.WriteString(fmt.Sprintf("%s:%v", k, v))
				break
			case reflect.Int:
				buf.WriteString(fmt.Sprintf("%s:%v", k, v))
				break
			case reflect.Float64:
				buf.WriteString(fmt.Sprintf("%s:%v", k, v))
				break
			default:
				fmt.Printf("type not found %#v\n", v)
			}
		}

		//if pos != len(m)-1 {
		buf.WriteByte('\n')
		//}

		pos++
	}
}

func processText(i interface{}, buf *bytes.Buffer) {
	s := i.(string)

	if strings.Contains(s, "\"") {
		//consider a long text
		buf.WriteString("|>")
		defer buf.WriteString("<|")
	}

	buf.WriteString(s)
}

func traverseSliceMsi(i interface{}, buf *bytes.Buffer) {
	buf.WriteByte('[')
	defer buf.WriteByte(']')

	c := i.([]interface{})
	for i, val := range c {

		switch reflect.TypeOf(val).Kind() {
		case reflect.Slice:
			traverseSliceMsi(val, buf)
			continue
		case reflect.Map:
			traverseMsi(val, buf)
			continue
		default:
			//not possible
		}

		buf.WriteString(fmt.Sprintf("%v", val))
		if i != len(c)-1 {
			buf.WriteString(fmt.Sprintf(" "))
		}
	}
}

func traverseToJSON(t *total) []byte {
	// TODO Do something with name. Now it's being ignored
	//t.docName

	bs := make([]byte, 0)
	buf := bytes.NewBuffer(bs)

	// Data is either a block or a list of blocks
	if t.data.kind == OBJECT {
		//traverseObject(t.data.data.(object))
		traverseObjectToJSON(t.data.data.(object), buf)
		return buf.Bytes()
	}

	//traverseList(t.data.data.(values))
	traverseListToJSON(t.data.data.(values), buf)

	return buf.Bytes()
}

func traverseDocName(s string, buf *bytes.Buffer) {
	buf.WriteString(s)
}

func traverseValueToJSON(v *value, buf *bytes.Buffer) error {
	switch v.kind {
	case WORD:
		buf.WriteString(fmt.Sprintf(`"%s"`, v.data.(string)))
	case OBJECT:
		traverseObjectToJSON(v.data.(object), buf)
	case FLOAT:
		buf.WriteString(fmt.Sprintf("%f", v.data.(float64)))
	case INTEGER:
		buf.WriteString(fmt.Sprintf("%d", v.data.(int)))
	case TEXT:
		buf.WriteString(fmt.Sprintf(`"%s"`, escapeText(v.data.(string))))
	case BOOLEAN:
		buf.WriteString(fmt.Sprintf("%s", v.data.(string)))
	case NULLTYPE:
		buf.WriteString(fmt.Sprintf("null"))
	case LIST:
		traverseListToJSON(v.data.(values), buf)
	}
	return nil
}

func escapeText(s string) string {
	//return strings.Replace(s, `"`, `\"`, -1)
	return strings.Replace(strings.Replace(s, `\`, `\\`, -1), `"`, `\"`, -1)
}

func traverseKvI(kv *keyValue, buf *bytes.Buffer) error {
	buf.WriteString(fmt.Sprintf(`"%s":`, kv.name))
	traverseValueToJSON(kv.value, buf)
	return nil
}

func traverseObjectToJSON(v object, buf *bytes.Buffer) error {
	buf.WriteByte('{')
	defer buf.WriteByte('}')

	for i, kv := range v {
		traverseKvI(kv, buf)

		//is last? don't add comma
		if i != len(v)-1 {
			buf.WriteByte(',')
		}
	}
	return nil
}

func traverseListToJSON(l values, buf *bytes.Buffer) error {
	buf.WriteByte('[')

	for i, value := range l {
		traverseValueToJSON(value, buf)

		//is last? don't add comma
		if i != len(l)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	return nil
}

func traverseValue(v *value) {
	fmt.Printf("%s", v)
	switch v.kind {
	case WORD:
		fmt.Printf(", %s\n", v.data.(string))
	case OBJECT:
		traverseObject(v.data.(object))
	case FLOAT:
		fmt.Printf(": %f\n", v.data.(float64))
	case INTEGER:
		fmt.Printf(", Value: %d\n", v.data.(int))
	case TEXT:
		fmt.Printf(", Value: %s\n", v.data.(string))
	case BOOLEAN:
		fmt.Printf(", Value: %s\n", v.data.(string))
	case NULLTYPE:
		fmt.Printf(", Value: NULL\n")
	case LIST:
		traverseList(v.data.(values))
	}
}

func traverseKv(kv *keyValue) {
	fmt.Printf("Name: %s, ", kv.name)
	traverseValue(kv.value)
}

func traverseObject(v object) {
	fmt.Printf("\n")
	for _, kv := range v {
		traverseKv(kv)
	}
}

func traverseList(l values) {
	fmt.Printf("\n")
	for _, value := range l {
		traverseValue(value)
	}
}

func tokName(i int) string {
	switch i {
	case WORD:
		return "WORD"
	case OBJECT:
		return "OBJECT"
	case VALUE:
		return "VALUE"
	case FLOAT:
		return "FLOAT"
	case INTEGER:
		return "INTEGER"
	case TEXT:
		return "TEXT"
	case BOOLEAN:
		return "BOOLEAN"
	case NULLTYPE:
		return "NULLTYPE"
	case LIST:
		return "LIST"
	default:
		return "UNKNOWN"
	}
}
