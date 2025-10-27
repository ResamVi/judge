package main

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ResamVi/judge/db"
	"github.com/ResamVi/judge/handler"
	"github.com/ResamVi/judge/migrate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rdbell/echo-pretty-logger"
)

//go:generate go tool sqlc generate

func main() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres:postgres@localhost:5432/mydb?sslmode=disable"
	}

	err := migrate.DB(url)
	if err != nil {
		panic(err)
	}

	queries, err := db.Init(url)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir("submissions", 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	e := echo.New()
	e.Use(prettylogger.Logger)
	e.Use(middleware.Recover())
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("www/index.html")),
	}

	h, err := handler.New(queries, os.Getenv("ENV"))
	if err != nil {
		panic(err)
	}

	e.GET("/", h.Homepage)
	e.GET("/login", h.LoginView)
	e.GET("/register", h.RegisterView)

	e.GET("/tasks/:task", h.TaskHandler)
	e.GET("/tasks/:task/code", CodeHandler)

	e.POST("/submission", h.Submit)
	e.POST("/login", h.Login)
	e.POST("/register", h.Register)
	e.POST("/username", h.Username)
	e.POST("/validate/name", h.ValidateUsername)
	e.POST("/validate/password", h.ValidatePassword)
	e.POST("/validate/confirm", h.ValidateConfirmation)

	e.Static("/tasks", "tasks")
	e.Static("/www", "www")

	e.Logger.Fatal(e.Start(":8080"))
}

func CodeHandler(c echo.Context) error {
	var b []byte
	buf := bytes.NewBuffer(b)
	w := zip.NewWriter(buf)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		strippedPath := strings.TrimPrefix(path, "tasks/"+c.Param("task")+"/code/")

		f, err := w.Create(strippedPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk("tasks/"+c.Param("task")+"/code", walker)
	if err != nil {
		panic(err)
	}
	w.Close()

	return c.Blob(http.StatusOK, "application/zip", buf.Bytes())
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
