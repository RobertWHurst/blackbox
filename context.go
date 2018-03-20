package blackbox

import "fmt"

type context struct {
	data map[string]string
}

func (c *context) extend(contextData map[string]interface{}) context {
	mergedContext := context{
		data: c.data,
	}
	for key, value := range contextData {
		if value == nil {
			delete(mergedContext.data, key)
		} else {
			mergedContext.data[key] = fmt.Sprintf("%+v", value)
		}
	}
	return mergedContext
}
