package resources

import (
	"embed"
	"fmt"
	htmltemplate "html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	txttemplate "text/template"
)

//go:embed templates
var templateEmbedFS embed.FS

type (
	TemplateIf interface {
		Execute(wr io.Writer, data any) error
		ExecuteTemplate(io.Writer, string, any) error
	}
	TemplateMap map[string]TemplateIf
)

func (tm *TemplateMap) Get(name string) (TemplateIf, error) {
	if v, ok := (*tm)[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("template not found for name %s", name)
}

func ParseTemplates(templatesDir string) (TemplateMap, error) {
	templates := make(TemplateMap, 0)

	var templateFS fs.FS
	if templatesDir == "embed" {
		var err error
		templateFS, err = fs.Sub(templateEmbedFS, "templates")
		if err != nil {
			return templates, err
		}
	} else {
		templateFS = os.DirFS(templatesDir)
	}

	nonViewHtmlTemplates, err := htmltemplate.New("").Funcs(templateFuncMap).ParseFS(
		templateFS,
		"layout/*.gohtml",
		"partial/*.gohtml",
	)
	if err != nil {
		return templates, err
	}

	/* currently no inheritance for plain templates, uncomment if/when needed
	nonViewTxtTemplates, err := txttemplate.New("").Funcs(templateFuncMap).ParseFS(
		templateFS,
		"layout/*.gotxt",
		"partial/*.gotxt",
	)
	if err != nil {
		return templates, err
	}
	*/
	nonViewTxtTemplates := txttemplate.New("").Funcs(templateFuncMap)

	viewSub, err := fs.Sub(templateFS, "view")
	if err != nil {
		return templates, err
	}

	err = fs.WalkDir(viewSub, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ".gohtml" {
			name := filepath.Base(p)
			c, err := nonViewHtmlTemplates.Clone()
			if err != nil {
				return err
			}
			t, err := c.New(name).Funcs(templateFuncMap).ParseFS(
				templateFS, fmt.Sprintf("view/%s", name),
			)
			if err != nil {
				return err
			}
			templates[name] = t
		}
		if filepath.Ext(p) == ".gotxt" {
			name := filepath.Base(p)
			c, err := nonViewTxtTemplates.Clone()
			if err != nil {
				return err
			}
			t, err := c.New(name).Funcs(templateFuncMap).ParseFS(
				templateFS, fmt.Sprintf("view/%s", name),
			)
			if err != nil {
				return err
			}
			templates[name] = t
		}
		return nil
	})
	if err != nil {
		return templates, err
	}

	return templates, nil
}
