package lytnin

import (
	"os"
)

// Modules contains all registered modules
var Modules []Module

// Module is helps organise HTTP handlers and gives them access to
// the application struct
type Module interface {
	Name() string
	Init(app *Application)
}

// RegisterModule adds the module to the list of modules
func RegisterModule(m Module) {
	Modules = append(Modules, m)
}

// Service extends the application by providing additional features
// such as database access or html template rendering.
type Service interface {
	Init(a *Application)
	Info() interface{}
	Close()
}

// Configuration is the application specific configuration struct
type Configuration struct {
	Debug bool `env:"DEBUG"`
}

// Get retrieves configuration environment variables
func (c Configuration) Get(val string, replace ...string) string {
	var env string
	if env = os.Getenv(val); env == "" {
		if len(replace) > 0 {
			env = replace[0]
		}
	}
	return env
}

func init() {
	Modules = []Module{}
}
