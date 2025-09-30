package blackbox

import (
	"runtime"
	"strings"
	"sync"
)

type Source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Target is an interface ment to be implemented by types that collect log
// data. blackbox ships with two of these: PrettyTarget and JSONTarget
type Target interface {
	Log(level Level, values []any, context Ctx, getSource func() *Source)
}

type targetSet struct {
	targets     []Target
	targetsLock sync.Mutex
}

func (t *targetSet) log(level Level, values []any, context Ctx, pc []uintptr) {
	for index, value := range values {
		if ctx, ok := value.(Ctx); ok {
			context = context.Extend(ctx)
			values = append(values[:index], values[index+1:]...)
		}
	}

	getSource := func() *Source {
		if len(pc) == 0 {
			return nil
		}

		frames := runtime.CallersFrames(pc)
		for {
			frame, more := frames.Next()
			if !more {
				break
			}
			funcPathChunks := strings.Split(frame.Function, "/")
			if len(funcPathChunks) == 0 {
				continue
			}
			packageAndFuncName := strings.Split(funcPathChunks[len(funcPathChunks)-1], ".")
			if len(packageAndFuncName) == 0 {
				continue
			}
			packageName := packageAndFuncName[0]
			if packageName != "blackbox" {
				return &Source{
					Function: frame.Function,
					File:     frame.File,
					Line:     frame.Line,
				}
			}
		}

		return nil
	}

	t.targetsLock.Lock()
	for _, target := range t.targets {
		target.Log(level, values, context, getSource)
	}
	t.targetsLock.Unlock()
}

func (t *targetSet) addTarget(target Target) {
	t.targetsLock.Lock()
	t.targets = append(t.targets, target)
	t.targetsLock.Unlock()
}
