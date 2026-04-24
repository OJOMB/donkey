package objects

import "fmt"

// Environment is a struct that represents the environment in which the Donkey programming language is evaluated.
// It contains a map of variable names to their corresponding objects, and a pointer to an outer environment (if any) for nested scopes.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new Environment with an empty store and no outer environment.
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), outer: nil}
}

// NewEnclosedEnvironment creates a new Environment that is enclosed within an outer environment.
// This is used for creating new scopes when evaluating functions or blocks.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	newEnv := NewEnvironment()
	if outer == nil {
		return newEnv
	}

	newEnv.outer = outer
	return newEnv
}

// Get retrieves the object associated with the given name from the environment.
// It first checks the current environment's store, and if not found, it checks the outer environment (if any).
// It returns the object and a boolean indicating whether the name was found in the environment.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}

	return obj, ok
}

// Bind assigns the given object to the given name in the environment's store or overwrites the existing value if the name already exists in the current environment. It returns the object that was bound to the name. If the name already exists in the current environment, it should log a warning and overwrite the existing value.
func (e *Environment) Bind(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Set updates the value of an existing variable in the environment. It first checks if the variable exists in the current environment's store, and if it does, it updates the value.
// If the variable does not exist, it recursively checks the outer environment (if any) for the variable. If the variable is found in an outer environment, it updates the value there. If the variable is not found in any environment, it returns an error indicating that the variable was not found.
func (e *Environment) Set(name string, val Object) (Object, error) {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val, nil
	}

	if e.outer != nil {
		return e.outer.Set(name, val)
	}

	return nil, fmt.Errorf("variable not found: %s", name)
}
