package services

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/lytnin/lytnin"
)

// Renderer manages a pongo2 TemplateSet
type Renderer struct {
	baseDir   string
	templates *template.Template
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
		baseDir:   baseDir,
		templates: template.Must(template.ParseGlob(filepath.Join(baseDir, "*.html"))),
	}

	return &rdr, nil
}

// Render implements echo.Render interface
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

// HTMLRender service provides html template rendering to the application
type HTMLRender struct {
	BaseDir string
}

// Info returns information about the html renderer
func (s *HTMLRender) Info() interface{} {
	return s.BaseDir
}

// Init initializes the html rendere service and registers it with the application
func (s *HTMLRender) Init(a *lytnin.Application) {
	r, err := NewRenderer(s.BaseDir)
	checkErr(err)
	a.M.Renderer = r

	a.AddService("htmlrender", s)
}

// Close releases any resources used by the service
func (s *HTMLRender) Close() {
	// nothing to do
}
