package dataflow

import (
	"errors"
	"testing"
	"time"

	"github.com/efureev/go-dataflow/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	fnRunSuccessful := func(sleep time.Duration) NodeFunc {
		return func(ctx Context) (FlowData, error) {
			//time.Sleep(sleep * time.Second)
			list := ctx.Data().Raw.(map[string]interface{})
			nodeCount := list[`nodes`].(int)
			list[`nodes`] = nodeCount + 1
			//nodes := list[`nodesNames`].([]string)
			//nodes = append(nodes, )

			return NewFlowDataFromMap(list), nil
		}
	}

	fnRunFail := func(sleep time.Duration) NodeFunc {
		return func(ctx Context) (FlowData, error) {
			time.Sleep(sleep * time.Second)
			//println(`fn: `, sleep)
			return *ctx.Data(), errors.New(`test error!`)
		}
	}

	l := logger.New()
	l.SetLevel(logger.LogLevelError)

	d := NewFlowDataFromMap(map[string]interface{}{
		`email`: `furegin@test.com`,
		`name`:  `Fureev Eugene`,
		`nodes`: 0,
		//`nodesNames`: []string{},
	})

	initNode := NewEmptyNode(NewNodeConfig(`first`, fnRunSuccessful(1), l))
	initNode.
		AddNext(NewEmptyNode(NewNodeConfig(`second`, fnRunSuccessful(1), l))).
		AddNext(NewEmptyNode(NewNodeConfig(`third`, fnRunFail(0), l))).
		AddNext(NewEmptyNode(NewNodeConfig(`finish`, fnRunSuccessful(0), l)))

	f := New()
	f.AddInitNode(initNode)
	processedData := f.Run(d)

	expectedData := NewFlowDataFromMap(map[string]interface{}{
		`email`: `furegin@test.com`,
		`name`:  `Fureev Eugene`,
		`nodes`: 2,
		//`nodesNames`: []string{`first`, `second`},
	})

	assert.Equal(t, expectedData, processedData)
	assert.Equal(t, 3, f.history.Count())
	assert.True(t, f.history.HasErrors())
	assert.Equal(t, 4, f.Count())
	//result := f.Result()

	//spew.Dump(f.toTree())
	//spew.Dump(f.ToString())

	expectedTreeJson := `{"Node":{"Name":"first","Next":{"Name":"second","Next":{"Name":"third","Next":{"Name":"finish","Next":null}}}}}` + "\n"
	treeJson, err := f.ToString()
	assert.Nil(t, err)
	assert.Equal(t, expectedTreeJson, treeJson)
	f.Clear()
	assert.Nil(t, f.node)
	treeJsonEmpty, err := f.ToString()
	assert.Nil(t, err)
	assert.Equal(t, `{}`, treeJsonEmpty)
	assert.Equal(t, 0, f.Count())
	assert.Equal(t, 0, f.history.Count())
	assert.False(t, f.history.HasErrors())

	assert.Nil(t, f.Load(treeJson))
	assert.Equal(t, 4, f.Count())

}
