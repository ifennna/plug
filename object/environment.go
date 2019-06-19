package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (env *Environment) Get(name string) (Object, bool) {
	object, ok := env.store[name]
	if !ok && env.outer != nil {
		// check the parent scope
		object, ok = env.outer.Get(name)
	}
	return object, ok
}

func (env *Environment) Set(name string, value Object) Object {
	env.store[name] = value
	return value
}
