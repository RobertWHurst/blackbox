package blackbox

import "fmt"

type context map[string]string

func (c context) extend(contextData map[string]interface{}) context {
	newContext := make(context, 0)

	for key, value := range c {
		newContext[key] = value
	}

	for key, value := range contextData {
		if value != nil {
			newContext[key] = fmt.Sprintf("%+v", value)
		}
	}

	return newContext
}
