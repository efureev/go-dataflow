package dataflow

import (
	"sync"
	"time"
)

type Flow struct {
	start, finish time.Time
	running       bool
	node          Node

	logger  *ProxyLogger
	history *NodeResultHistory

	marshaller Marshaller

	doneCh chan NodeResult
	wg     sync.WaitGroup
}

func New() (flow *Flow) {
	flow = &Flow{
		running:    false,
		logger:     &ProxyLogger{},
		doneCh:     make(chan NodeResult),
		history:    newNodeResultHistory(),
		marshaller: &JsonMarshaller{},
	}

	return flow
}

func (f *Flow) initLogging() {
	for res := range f.doneCh {
		f.history.Add(res)

		status := `success`
		if res.HasError() {
			childrenCount := res.node.Count()
			f.wg.Add(-childrenCount)

			status = `error`
		}

		f.logger.Logf("[%s] Result from %s with result: %s\n", status, res.node.Name(), res.Duration())

		f.wg.Done()
	}
}

func (f *Flow) init() {
	f.wg.Add(f.Count())

	go f.initLogging()
}

func (f *Flow) startRun() {
	f.start = time.Now()
	f.logger.Tracef(`flow: start: %s`, f.start)
	f.init()

	f.running = true
}

func (f *Flow) stopRun() {
	f.running = false
	f.finish = time.Now()

	close(f.doneCh)
	f.logger.Tracef(`flow: finish: %s`, f.finish)
}

func (f *Flow) Run(data FlowData) FlowData {
	f.startRun()
	processedData := f.node.Run(f.doneCh, data)

	f.wg.Wait()
	f.stopRun()

	return processedData
}

func (f *Flow) AddInitNode(node Node) {
	f.node = node
}

func (f *Flow) Count() int {
	if f.node == nil {
		return 0
	}

	var curNode Node
	curNode = f.node

	count := 1
	for curNode != nil {
		n := curNode.Next()
		if n != nil {
			count++
		}
		curNode = n
	}

	return count
}

func (f *Flow) Result() FlowResult {
	return FlowResult{
		start:    f.start,
		finish:   f.finish,
		finished: !f.running,
		errors:   f.history.Errors(),
	}
}

// Load DataFlow from custom format
func (f *Flow) fromTree(tree *FlowTree) error {
	f.Clear()

	if tree == nil || tree.Node == nil {
		return nil
	}

	curFlowNode := tree.Node

	initNode := NewEmptyNode(NewNodeConfig(curFlowNode.Name, emptyNodeFunc, f.logger))
	curNode := initNode

	for curFlowNode != nil {
		n := curFlowNode.Next
		if n != nil {
			nextNode := NewEmptyNode(NewNodeConfig(n.Name, emptyNodeFunc, f.logger))
			curNode.AddNext(nextNode)
			curNode = nextNode
		}
		curFlowNode = n
	}
	f.node = initNode

	return nil
}

// Load DataFlow from custom format
func (f *Flow) toTree() FlowTree {
	f.throwOnRunning()

	tree := FlowTree{}
	if f.node == nil {
		return tree
	}

	curNode := f.node
	tree.Node = flowTreeNodeFromNode(f.node)
	curFlowNode := tree.Node

	for curNode != nil {
		n := curNode.Next()
		if n != nil {
			cfNode := flowTreeNodeFromNode(n)
			curFlowNode.AddNext(cfNode)
			curFlowNode = cfNode
		}
		curNode = n
	}

	return tree
}

// Clear Flow from Nodes
func (f *Flow) Clear() {
	f.throwOnRunning()

	f.node = nil
	f.history.Clear()
}

// ToString DataFlow to custom Format to store in DB
func (f *Flow) ToString() (string, error) {
	return f.marshaller.ToString(f.toTree())
}

// Load DataFlow from custom format
func (f *Flow) Load(str string) error {
	f.throwOnRunning()

	tree, err := f.marshaller.Load(str)
	if err != nil {
		return err
	}

	return f.fromTree(tree)
}

func (f *Flow) SetMarshaller(marshaller Marshaller) {
	f.marshaller = marshaller
}

func (f *Flow) throwOnRunning() {
	if f.running {
		panic(`Flow should be finish`)
	}
}
