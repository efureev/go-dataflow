package dataflow

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Marshaller interface {
	Load(str string) (*FlowTree, error)
	ToString(tree FlowTree) (string, error)
}

type JsonMarshaller struct {
}

func (m *JsonMarshaller) Load(str string) (*FlowTree, error) {
	jsonStream := strings.NewReader(str)

	decoder := json.NewDecoder(jsonStream)
	var tree FlowTree
	err := decoder.Decode(&tree)

	return &tree, err
}

func (m *JsonMarshaller) ToString(tree FlowTree) (string, error) {
	if tree.isEmpty() {
		return `{}`, nil
	}

	buf := new(bytes.Buffer)
	bufEncoder := json.NewEncoder(buf)

	err := bufEncoder.Encode(tree)

	return buf.String(), err
}
