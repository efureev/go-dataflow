package dataflow

import "time"

type NodeResult struct {
	start, finish time.Time
	err           *NodeError
	finished      bool
	node          Node
	data          *FlowData
}

func (r *NodeResult) fail(err error) {
	r.err = newNodeError(r.node, err)
	r.done()
}

func (r *NodeResult) failOnData(err error, data *FlowData) {
	r.data = data
	r.fail(err)
}

func (r *NodeResult) done() {
	r.finished = true
	r.finish = time.Now()
}

func (r NodeResult) HasError() bool {
	return r.err != nil
}

func (r NodeResult) Duration() time.Duration {
	return r.finish.Sub(r.start)
}

func newNodeResult(node Node) NodeResult {
	return NodeResult{
		start:    node.GetStartTime(),
		err:      nil,
		finished: false,
		node:     node,
	}
}
