package dataflow

type (
	Context interface {
		Data() *FlowData
		Node() nodeContext
		Logger() *ProxyLogger
	}

	nodeContext struct {
		name string
	}

	context struct {
		data   *FlowData
		node   nodeContext
		logger *ProxyLogger
	}
)

func (c context) Data() *FlowData {
	return c.data
}

func (c context) Node() nodeContext {
	return c.node
}

func (c context) Logger() *ProxyLogger {
	return c.logger
}
