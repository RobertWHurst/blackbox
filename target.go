package blackbox

import "sync"

// Target is an interface ment to be implemented by types that collect log
// data. blackbox ships with two of these: PrettyTarget and JSONTarget
type Target interface {
	Log(level Level, values []interface{}, context Ctx)
}

type targetSet struct {
	targets     []Target
	targetsLock sync.Mutex
}

func (t *targetSet) log(level Level, values []interface{}, context Ctx) {
	for index, value := range values {
		if ctx, ok := value.(Ctx); ok {
			context = context.Extend(ctx)
			values = append(values[:index], values[index+1:]...)
		}
	}

	t.targetsLock.Lock()
	for _, target := range t.targets {
		target.Log(level, values, context)
	}
	t.targetsLock.Unlock()
}

func (t *targetSet) addTarget(target Target) {
	t.targetsLock.Lock()
	t.targets = append(t.targets, target)
	t.targetsLock.Unlock()
}
