package blackbox

// Target is an interface ment to be implemented by types that collect log
// data. blackbox ships with two of these: PrettyTarget and JSONTarget
type Target interface {
	Log(level Level, values []interface{}, context Ctx)
}

type targetSet struct {
	targets []Target
}

func (t *targetSet) log(level Level, values []interface{}, context Ctx) {
	for _, target := range t.targets {
		target.Log(level, values, context)
	}
}

func (t *targetSet) addTarget(target Target) {
	t.targets = append(t.targets, target)
}
