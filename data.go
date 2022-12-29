package dataflow

type FlowData struct {
	Raw any
}

func NewFlowDataFromMap(data map[string]interface{}) FlowData {
	return FlowData{
		Raw: data,
	}
}
