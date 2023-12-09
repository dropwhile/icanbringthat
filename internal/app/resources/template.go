package resources

import (
	"embed"
	"fmt"
	htmltemplate "html/template"
	"io"
	"io/fs"
	"maps"
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

func getTemplateFS(templatesDir string) (fs.FS, error) {
	var templateFS fs.FS
	if templatesDir == "embed" {
		var err error
		templateFS, err = fs.Sub(templateEmbedFS, "templates")
		if err != nil {
			return templateFS, err
		}
	} else {
		templateFS = os.DirFS(templatesDir)
	}
	return templateFS, nil
}

func ParseHtmlTemplates(templatesDir string) (TemplateMap, error) {
	templates := make(TemplateMap, 0)

	templateFS, err := getTemplateFS(templatesDir)
	if err != nil {
		return templates, err
	}

	nonViewHtmlTemplates, err := htmltemplate.New("").Funcs(templateFuncMap).ParseFS(
		templateFS,
		"html/layout/*.gohtml",
		"html/partial/*.gohtml",
	)
	if err != nil {
		return templates, err
	}

	viewSub, err := fs.Sub(templateFS, "html/view")
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
				templateFS, fmt.Sprintf("html/view/%s", name),
			)
			if err != nil {
				return err
			}
			templates[name] = t
		}
		return nil
	})
	return templates, err
}

func ParseTxtTemplates(templatesDir string) (TemplateMap, error) {
	templates := make(TemplateMap, 0)

	templateFS, err := getTemplateFS(templatesDir)
	if err != nil {
		return templates, err
	}

	viewSub, err := fs.Sub(templateFS, "txt")
	if err != nil {
		return templates, err
	}
	nonViewTxtTemplates := txttemplate.New("").Funcs(templateFuncMap)

	err = fs.WalkDir(viewSub, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ".gotxt" {
			name := filepath.Base(p)
			c, err := nonViewTxtTemplates.Clone()
			if err != nil {
				return err
			}
			t, err := c.New(name).Funcs(templateFuncMap).ParseFS(
				templateFS, fmt.Sprintf("txt/%s", name),
			)
			if err != nil {
				return err
			}
			templates[name] = t
		}
		return nil
	})
	return templates, err
}

func ParseTemplates(templatesDir string) (TemplateMap, error) {
	templates := make(TemplateMap, 0)

	htmlTemplates, err := ParseHtmlTemplates(templatesDir)
	if err != nil {
		return templates, err
	}
	maps.Copy(templates, htmlTemplates)

	txtTemplates, err := ParseTxtTemplates(templatesDir)
	if err != nil {
		return templates, err
	}
	maps.Copy(templates, txtTemplates)

	return templates, nil
}
