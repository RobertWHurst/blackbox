package blackbox

type TestTarget struct {
	logged []Logged
}

var _ Target = &TestTarget{}

func NewTestTarget() *TestTarget {
	return &TestTarget{
		logged: make([]Logged, 0),
	}
}

type Logged struct {
	LoggerID string
	Level    Level
	Values   []any
	Context  Ctx
	Source   *Source
}

func (t *TestTarget) Log(loggerID string, level Level, values []any, context Ctx, getSource func() *Source) {
	t.logged = append(t.logged, Logged{
		LoggerID: loggerID,
		Level:    level,
		Values:   values,
		Context:  context,
		Source:   getSource(),
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
