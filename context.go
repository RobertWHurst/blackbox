package blackbox

// Ctx is an alias for map[string]any. This is the format for
// data to me used for extending contexts.
type Ctx map[string]any

// Extend returns a new context with the given context data added to it.
// If a key already exists in the context, it will be overwritten.
func (c Ctx) Extend(contextData Ctx) Ctx {
	newContext := make(Ctx, 0)

	for key, value := range c {
		newContext[key] = value
	}

	for key, value := range contextData {
		if value != nil {
			newContext[key] = value
		}
	}

	return newContext
}
