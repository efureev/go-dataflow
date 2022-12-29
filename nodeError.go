package dataflow

type NodeError struct {
	Node
	Err error
}

func newNodeError(node Node, err error) *NodeError {
	return &NodeError{
		Node: node,
		Err:  err,
	}
}
