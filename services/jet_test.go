package services

import (
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/labstack/echo"

	"github.com/lytnin/lytnin"
)

func TestHTMLRenderService(t *testing.T) {
	app := lytnin.NewApplication()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Expected to get current working directory but got error instead: %s", err)
	}

	s1 := HTMLRender{BaseDir: path.Join(dir, "test_fixtures")}
	s1.Init(app)

	s2 := app.GetService("jetrender")
	if s2 == nil {
		t.Fatal("Expected to get HTMLRender Service but got nil")
	}

	_, ok := s2.(*HTMLRender)
	if !ok {
		t.Fatalf("Expected to successfully cast Service interface to HTMLRender Type")
	}

	h := func(c echo.Context) error {
		// call index.jet from test fixtures
		return c.Render(200, "index", map[string]interface{}{})
	}

	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := app.M.NewContext(req, rec)

	user := struct {
		Username string
	}{"admin"}
	c.Set("user", user)

	err = h(c)
	if err != nil {
		t.Fatalf("Expected handler to return no error but instead got: %s", err)
	}

	b := rec.Body.String()
	val := "<p>admin</p>"
	if b != val {
		t.Fatalf("Expected '%s' but instead go '%s'", val, b)
	}

}
