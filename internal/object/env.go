package object

type Env struct {
	store map[string]Object
}

func NewEnv() *Env {
	return &Env{store: make(map[string]Object)}
}

func (e *Env) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	return obj, ok
}

func (e *Env) Set(key string, val Object) Object {
	e.store[key] = val
	return val
}
