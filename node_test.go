package dataflow

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/efureev/go-dataflow/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	fnRunSuccessful := func(sleep time.Duration) NodeFunc {
		return func(ctx Context) (FlowData, error) {
			//time.Sleep(sleep * time.Second)

			list := ctx.Data().Raw.(map[string]interface{})
			r := map[string]any{}
			for k, v := range list {
				r[k] = strings.ToUpper(v.(string))
			}

			return NewFlowDataFromMap(r), nil
		}
	}

	d := NewFlowDataFromMap(map[string]interface{}{
		`email`: `furegin@test.com`,
		`name`:  `Fureev Eugene`,
	})

	ch := make(chan NodeResult)
	l := logger.New()
	l.SetLevel(logger.LogLevelError)

	w := sync.WaitGroup{}
	var r *NodeResult

	go func(ch chan NodeResult) {
		res := <-ch

		r = &res
		w.Done()
	}(ch)

	w.Add(1)
	initNode := NewEmptyNode(NewNodeConfig(`first`, fnRunSuccessful(1), l))

	processesData := initNode.Run(ch, d)

	w.Wait()

	//spew.Dump(processesData)

	assert.NotNil(t, r)
	assert.Nil(t, r.err)
	assert.False(t, r.HasError())
	assert.True(t, r.finished)
	assert.True(t, r.Duration() > 0)
	assert.Equal(t, `first`, r.node.Name())

	expectedData := NewFlowDataFromMap(map[string]interface{}{
		`email`: `FUREGIN@TEST.COM`,
		`name`:  `FUREEV EUGENE`,
	})
	assert.Equal(t, expectedData, processesData)
}
