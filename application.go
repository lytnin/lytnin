package lytnin

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
)

// Application is a global struct that contains items
// that you want to make avaulable to your HTTP handlers
type Application struct {
	M        *echo.Echo
	Config   Configuration
	services map[string]Service
	Cli      *cli.App
}

// NewApplication creates a new application struct
func NewApplication() *Application {
	// read .env if any
	_ = godotenv.Load()

	// load application specific configuration
	c := Configuration{}
	err := env.Parse(&c)
	if err != nil {
		log.Fatal(err)
	}

	// create mux
	e := echo.New()

	// middleware
	e.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// static assets
	e.Static(
		"/public",
		c.Get("STATIC_DIR", "frontend/dist/static"),
	)

	// single page app
	e.File(
		"/app/*",
		c.Get("SPA_ENTRY_POINT", "frontend/dist/index.html"),
	)

	app := &Application{
		M:        e,
		Config:   c,
		services: map[string]Service{},
	}

	// setup cli
	app.Cli = cli.NewApp()
	app.Cli.Name = "Lytnin"
	app.Cli.Usage = "Web Application CLI"
	app.Cli.Version = "0.0.1"

	return app
}

func (a *Application) AddCommand(cmd cli.Command) {
	a.Cli.Commands = append(a.Cli.Commands, cmd)
}

// Start runs the application
func (a *Application) Start() {
	// initialize modules
	for _, mod := range Modules {
		mod.Init(a)
	}

	addr := ":" + a.Config.Get("PORT")
	a.M.Logger.Fatal(a.M.Start(addr))
}

// Close cleans up application resources
func (a *Application) Close() {
	for _, s := range a.services {
		s.Close()
	}
}

// AddService adds application dependencies/plugins
func (a *Application) AddService(name string, s Service) {
	a.services[name] = s
}

// GetService returns existing application dependencies/plugins
func (a *Application) GetService(name string) Service {
	var s Service
	s, _ = a.services[name]
	return s
}
