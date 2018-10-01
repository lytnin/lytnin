package services

import (
	"log"

	"github.com/lytnin/lytnin"
	"github.com/urfave/cli"
)

// Hello service add a hello command to the main application
type Hello struct {
}

// Info returns information about the service
func (s *Hello) Info() interface{} {
	return "Hello command"
}

// Init initializes the service and registers it with the application
func (s *Hello) Init(a *lytnin.Application) {
	a.AddCommand(s.Command())
	a.AddService("hello", s)
}

// Close releases any resources used by the service
func (s *Hello) Close() {
}

// Greet outputs a greeting to the console.
func (s *Hello) Greet(ctx *cli.Context) error {
	log.Println("Hello!")
	return nil
}

// Command creates a cli command
func (s *Hello) Command() cli.Command {
	return cli.Command{
		Name:   "hello",
		Usage:  "greeting",
		Action: s.Greet,
	}
}
