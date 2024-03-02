package main

import (
	"fmt"
	"html/template"
	"io"
	"myapp/handlers"
	"myapp/models"
	"myapp/sessions"
	"os"

	"github.com/labstack/echo/v4"
)

var sm = models.NewSessionManager()

func main() {
	e := echo.New()
	e.Static("/css", "public/styles")
	e.Use(sessions.SessionMiddleware(sm))
	renderer := Template{templates: template.Must(template.ParseGlob("public/views/*.html"))}
	e.Renderer = renderer

	e.GET("/", handlers.Home_GET)
	e.GET("/login", handlers.Login_GET)
	e.POST("/login", handlers.Login_POST)
	e.POST("/logout", handlers.Logout_POST)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%n", os.Getenv("PORT"))))
}

type Template struct {
	templates *template.Template
}

// Render implements echo.Renderer.
func (r Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
