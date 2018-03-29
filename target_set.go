package blackbox

type targetSet struct {
	targets []Target
}

func (t *targetSet) log(level Level, values []interface{}, context context) {
	for _, target := range t.targets {
		target.Log(level, values, context)
	}
}

func (t *targetSet) addTarget(target Target) {
	t.targets = append(t.targets, target)
}
