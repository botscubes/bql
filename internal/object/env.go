package object

type Env struct {
	store map[string]Object
	outer *Env
}

func NewEnv() *Env {
	return &Env{store: make(map[string]Object)}
}

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func (e *Env) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}
	return obj, ok
}

func (e *Env) Set(key string, val Object) Object {
	e.store[key] = val
	return val
}
