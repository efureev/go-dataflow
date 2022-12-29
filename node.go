package dataflow

import (
	"time"
)

type NodeFunc func(ctx Context) (FlowData, error)

var emptyNodeFunc NodeFunc = func(ctx Context) (FlowData, error) { return FlowData{}, nil }

type Node interface {
	Run(doneCh chan NodeResult, data FlowData) FlowData
	Name() string
	GetStartTime() time.Time
	Next() Node
	Count() int
	AddNext(node Node) Node
}

type NodeConfig struct {
	Name   string
	Fn     NodeFunc
	logger Logger
}

type NodeOption func(*EmptyNode)

func NewNodeConfig(name string, fn NodeFunc, l Logger) NodeConfig {
	return NodeConfig{
		Name:   name,
		Fn:     fn,
		logger: l,
	}
}

type EmptyNode struct {
	start    time.Time
	finished bool
	next     Node
	name     string

	fn NodeFunc

	doneCh chan NodeResult
	logger *ProxyLogger
}

func (en EmptyNode) GetStartTime() time.Time {
	return en.start
}

func (en EmptyNode) Name() string {
	return en.name
}

func (en *EmptyNode) Run(doneCh chan NodeResult, data FlowData) FlowData {
	en.logger.Tracef("node [%s]: start Run\n", en.name)

	en.doneCh = doneCh
	en.start = time.Now()

	ctx := context{
		node:   nodeContext{name: en.name},
		logger: en.logger,
		data:   &data,
	}

	processedData, err := en.fn(ctx)
	en.logger.Tracef("node [%s]: executed Run\n", en.name)
	if err != nil {
		en.throwError(err, data)
		return data
	}

	en.done()
	en.logger.Tracef("node [%s]: finish Run\n", en.name)

	totalProcessedData := en.runNext(processedData)

	return totalProcessedData
}

func (en *EmptyNode) Next() Node {
	return en.next
}

func (en *EmptyNode) HasNext() bool {
	return en.next != nil
}

func (en *EmptyNode) AddNext(node Node) Node {
	en.next = node

	return node
}

func (en *EmptyNode) runNext(data FlowData) FlowData {
	if en.HasNext() {
		return en.Next().Run(en.doneCh, data)
	}

	return data
}

func (en *EmptyNode) done() {
	r := newNodeResult(en)
	r.done()

	en.doneCh <- r
}

func (en *EmptyNode) throwError(err error, data FlowData) {
	r := newNodeResult(en)
	r.failOnData(err, &data)

	en.doneCh <- r
}

func (en *EmptyNode) WithFinished() {
	en.finished = true
}

func (en EmptyNode) Count() int {

	if en.Next() == nil {
		return 0
	}

	var curNode Node
	curNode = en.Next()

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

func NewEmptyNode(nc NodeConfig, opts ...NodeOption) *EmptyNode {
	node := &EmptyNode{
		finished: false,
		name:     nc.Name,
		fn:       nc.Fn,
		logger: &ProxyLogger{
			client: nc.logger,
		},
	}

	for _, opt := range opts {
		opt(node)
	}

	return node
}

func WithName(name string) NodeOption {
	return func(n *EmptyNode) { n.name = name }
}

func WithFinished() NodeOption {
	return func(n *EmptyNode) { n.finished = true }
}

func WithDoneChan(ch chan NodeResult) NodeOption {
	return func(n *EmptyNode) { n.doneCh = ch }
}

func WithLogger(logger Logger) NodeOption {
	return func(n *EmptyNode) { n.logger.client = logger }
}
