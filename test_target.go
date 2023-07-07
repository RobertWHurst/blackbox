package blackbox

func NewTestTarget() *TestTarget {
	return &TestTarget{
		logged: make([]Logged, 0),
	}
}

type TestTarget struct {
	logged []Logged
}

type Logged struct {
	Level   Level
	Values  []interface{}
	Context Ctx
}

func (t *TestTarget) Log(level Level, values []interface{}, context Ctx) {
	t.logged = append(t.logged, Logged{
		Level:   level,
		Values:  values,
		Context: context,
	})
}

func (t *TestTarget) Reset() {
	t.logged = nil
}

func (t *TestTarget) LastLogged() (Logged, bool) {
	if len(t.logged) == 0 {
		return Logged{}, false
	}
	return t.logged[0], true
}

func (t *TestTarget) PreviouslyLogged(i int) (Logged, bool) {
	if len(t.logged) >= i-1 {
		return Logged{}, false
	}
	return t.logged[0], true
}
