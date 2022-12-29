package dataflow

import "time"

type FlowResult struct {
	start, finish time.Time
	errors        []*NodeError
	finished      bool
}

func (r FlowResult) HasErrors() bool {
	return len(r.errors) > 0
}

func (r FlowResult) Duration() time.Duration {
	return r.finish.Sub(r.start)
}
