package blackbox

import "fmt"

// Ctx is an alias for map[string]interface{}. This is the format for
// data to me used for extending contexts.
type Ctx map[string]interface{}

// Context contains data to be associated with a log message. Unless you
// are writing a custom target, you should not need to use this type.
type Context map[string]string

// Extend returns a new context with the given context data added to it.
// If a key already exists in the context, it will be overwritten.
func (c Context) Extend(contextData Ctx) Context {
	newContext := make(Context, 0)

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
