package dataflow

type flowTreeNode struct {
	Name string
	Next *flowTreeNode
	//Children []*flowTreeNode
}

func (n *flowTreeNode) AddNext(node *flowTreeNode) {
	n.Next = node
}

type FlowTree struct {
	Node *flowTreeNode
}

func (n *FlowTree) isEmpty() bool {
	return n == nil || n.Node == nil
}

/*
func (n *flowTreeNode) AddChildren(nodes ...*flowTreeNode) {
	for _, node := range nodes {
		n.Children = append(n.Children, node)
	}
}
*/

func flowTreeNodeFromNode(node Node) *flowTreeNode {
	if node == nil {
		return nil
	}

	return &flowTreeNode{
		Name: node.Name(),
	}
}
