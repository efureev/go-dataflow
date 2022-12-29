package dataflow

import "sync"

type NodeResultHistory struct {
	m    sync.RWMutex
	list []NodeResult
}

// Clear History
func (h *NodeResultHistory) Clear() {
	h.list = []NodeResult{}
}

func (h *NodeResultHistory) HasErrors() bool {
	return len(h.Errors()) > 0
}

func (h *NodeResultHistory) Errors() []*NodeError {
	h.m.RLock()
	defer h.m.RUnlock()

	var list []*NodeError

	for _, r := range h.list {
		if r.HasError() {
			list = append(list, r.err)
		}
	}

	return list
}

func (h *NodeResultHistory) Add(result NodeResult) {
	h.m.Lock()
	h.list = append(h.list, result)
	h.m.Unlock()
}

func (h *NodeResultHistory) Count() int {
	h.m.RLock()
	defer h.m.RUnlock()

	return len(h.list)
}

func (h *NodeResultHistory) IsSuccessful() bool {
	return !h.HasErrors()
}

func newNodeResultHistory() *NodeResultHistory {
	return &NodeResultHistory{}
}
