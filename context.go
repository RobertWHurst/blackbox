package blackbox

import "fmt"

type context struct {
	data map[string]string
}

func (c *context) extend(contextData map[string]interface{}) context {
	newContext := context{
		data: make(map[string]string, 0),
	}

	for key, value := range c.data {
		newContext.data[key] = value
	}

	for key, value := range contextData {
		if value != nil {
			newContext.data[key] = fmt.Sprintf("%+v", value)
		}
	}

	return newContext
}
