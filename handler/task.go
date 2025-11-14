package handler

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/plouc/textree"
)

func (k Handler) TaskHandler(c echo.Context) error {
	// Convert local markdown files to HTML
	taskMD, err := os.ReadFile("tasks/" + c.Param("task") + "/README.md")
	if err != nil {
		slog.Error("os.ReadFile: " + err.Error())
		return c.NoContent(http.StatusNotFound)
	}
	var taskHTML bytes.Buffer
	if err := md.Convert(taskMD, &taskHTML); err != nil {
		slog.Error("md.Convert: " + err.Error())
		return err
	}

	// Replace occurrences of {{Code}} in the webpage with a custom file viewer
	htm := fmt.Sprintf(`
		<div class="row" style="margin-top: 3em; margin-bottom:3em">
			<div class="col-3">
				<div>%s</div>
			</div>

			<div class="col">
				<div class="tabs">
				%s
				</div>
			</div>
		</div>`, treeView(c.Param("task")), codeView(c.Param("task")))

	result := strings.ReplaceAll(taskHTML.String(), "{{Code}}", htm)

	data := k.page
	data.Body = result

	return c.Render(http.StatusOK, "index", data)
}

func treeView(name string) string {
	tree, err := textree.TreeFromDir("./tasks/" + name + "/code")
	if err != nil {
		panic(err)
	}

	var treebuf bytes.Buffer
	tree.Render(&treebuf, textree.NewRenderOptions())

	result := treebuf.String()
	result = strings.TrimSpace(result)
	result = strings.ReplaceAll(result, "\n", "<br />")

	return result
}

func codeView(name string) string {
	entries, err := os.ReadDir("./tasks/" + name + "/code")
	if err != nil {
		panic(err)
	}

	code := ""
	for _, e := range entries {
		content, err := os.ReadFile("./tasks/" + name + "/code/" + e.Name())
		if err != nil {
			panic(err)
		}

		code += fmt.Sprintf(`
			<input type="radio" name="tabs" id="tabone" checked="checked">
			<label for="tabone">%s</label>
			<div class="tab"><pre><code>%s</code></pre></div>`, e.Name(), string(content))
	}

	return code
}

func (k Handler) TaskList(c echo.Context) error {
	exercises, err := k.db.GetExercises(c.Request().Context())
	if err != nil {
		slog.Error("db.GetExercises: " + err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	result := make(map[string]string)
	for _, exercise := range exercises {
		result[exercise.Title] = exercise.ID
	}

	return c.JSON(http.StatusOK, result)
}
