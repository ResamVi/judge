package handler

import (
	"context"
	"fmt"
	"github.com/ResamVi/judge/db"
	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"net/http"
)

var (
	md = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
)

// Page contains all the data displayed on the current page
type Page struct {
	// List of exercises
	Exercises []MenuItem

	// The main content (derived from markdown files in exercises/)
	Body string

	// If the code editor should be shown or not
	Exercise bool
}

type MenuItem struct {
	Name string
	Link string
}

type Handler struct {
	db   *db.Queries
	page Page

	env string
}

func New(queries *db.Queries, env string) (*Handler, error) {
	exercises, err := getExercises(queries)
	if err != nil {
		return nil, fmt.Errorf("getting exercises: %w", err)
	}

	return &Handler{
		db: queries,
		page: Page{
			Exercises: exercises,
			Exercise:  false,
		},
		env: env,
	}, nil
}

// Username shows the name if the user is logged and if not a link to login
func (k Handler) Username(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return c.HTML(http.StatusOK, `<li class="float-right"><a href="/login">Anmelden</a></li>`)
	}

	return c.HTML(http.StatusOK, `<li class="float-right">Eingeloggt als <strong>`+cookie.Value+`</strong></li>`)
}

func getExercises(queries *db.Queries) ([]MenuItem, error) {
	exercises, err := queries.GetExercises(context.Background())
	if err != nil {
		return nil, fmt.Errorf("GetExercises: %w", err)
	}

	var items []MenuItem
	for _, ex := range exercises {
		items = append(items, MenuItem{
			Name: ex.Title,
			Link: "/exercises/" + ex.ID,
		})
	}

	return items, nil
}
