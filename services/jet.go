package services

import (
	"fmt"
	"io"
	"os"

	"github.com/CloudyKit/jet"
	"github.com/labstack/echo"
	"github.com/lytnin/lytnin"
)

// Renderer manages Jet Templates
type Renderer struct {
	baseDir string
	viewSet *jet.Set
}

// NewRenderer creates a new instance of Renderer
func NewRenderer(baseDir string) (*Renderer, error) {
	// check if baseDir exists
	fInfo, err := os.Lstat(baseDir)
	if err != nil {
		return nil, err
	}
	if fInfo.IsDir() == false {
		return nil, fmt.Errorf("%s is not a directory", baseDir)
	}

	rdr := Renderer{
		baseDir: baseDir,
		viewSet: jet.NewHTMLSet(baseDir),
	}

	return &rdr, nil
}

// Render implements echo.Render interface
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	t, err := r.viewSet.GetTemplate(name)
	if err != nil {
		return err
	}

	vars := make(jet.VarMap)
	// make sure data is a map
	m, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Please make sure template data is a Map")
	}
	// get values
	for k, v := range m {
		vars.Set(k, v)
	}

	return t.Execute(w, vars, c)
}

// HTMLRender service provides html template rendering to the application
type HTMLRender struct {
	BaseDir string
}

// Info returns information about the html renderer
func (s *HTMLRender) Info() interface{} {
	return s.BaseDir
}

// Init initializes the html render service and registers it with the application
func (s *HTMLRender) Init(a *lytnin.Application) {
	r, err := NewRenderer(s.BaseDir)
	checkErr(err)
	a.M.Renderer = r
	r.viewSet.SetDevelopmentMode(a.Config.Debug)

	a.AddService("jetrender", s)
}

// Close releases any resources used by the service
func (s *HTMLRender) Close() {
	// nothing to do
}
