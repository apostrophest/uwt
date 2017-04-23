package uwt

// An Environment represents a key-value mapping of environment vars
type Environment struct {
	env map[string]string
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{make(map[string]string)}
}

// Get returns the value of the specified env variable
func (e *Environment) Get(varname string) string {
	return e.env[varname]
}

// Set sets the value of the specified env variable
func (e *Environment) Set(varname string, value string) {
	e.env[varname] = value
}
