package resources

import (
	"embed"
	"html/template"
	"io/fs"
	"os"
)

var (
	//go:embed templates
	templateEmbedFS embed.FS
	templateFS      fs.FS
	templates       *template.Template
)

func MustParseTemplates(templatesDir string) *template.Template {
	if templatesDir == "embed" {
		var err error
		templateFS, err = fs.Sub(templateEmbedFS, "templates")
		if err != nil {
			panic(err)
		}
	} else {
		templateFS = os.DirFS(templatesDir)
	}

	var err error
	templates, err = template.ParseFS(
		templateFS,
		"layout/*.gohtml",
		"view/*.gohtml",
		"partial/*.gohtml",
	)
	if err != nil {
		panic(err)
	}
	return templates
}
