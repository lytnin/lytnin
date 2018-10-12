package services

import (
	"testing"

	"github.com/lytnin/lytnin"
)

func TestHelloService(t *testing.T) {
	app := lytnin.NewApplication()
	s1 := Hello{}
	s1.Init(app)

	s2 := app.GetService("hello")
	if s2 == nil {
		t.Fatal("Expected to get Hello Service but got nil")
	}

	hello, ok := s2.(*Hello)
	if !ok {
		t.Fatalf("Expected to successfully cast Service interface to Hello Type")
	}

	if err := hello.Greet(nil); err != nil {
		t.Fatalf("Expected Hello.Greet to succeed but instead got error: %s", err)
	}

}
