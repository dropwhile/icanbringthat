package resources

import (
	"embed"
	"fmt"
	htmltemplate "html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
)

const embedMagicStr = "embed"

//go:embed templates
var templateEmbedFS embed.FS

type TContainer struct {
	tmap  TemplateMap
	tfs   fs.FS
	embed bool
}

func (tc *TContainer) Get(name string) (TemplateIf, error) {
	if tc.embed {
		return tc.tmap.Get(name)
	}

	baseTpl, viewSub, err := getBaseTpl(tc.tfs)
	if err != nil {
		return nil, err
	}
	c, err := baseTpl.Clone()
	if err != nil {
		return nil, err
	}
	t, err := c.New(name).
		Funcs(templateFuncMap).
		Funcs(sprig.FuncMap()).
		ParseFS(viewSub, name)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func getTemplateFS(templatesDir string) (fs.FS, error) {
	var templateFS fs.FS
	if templatesDir == embedMagicStr {
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

func getBaseTpl(tfs fs.FS) (*htmltemplate.Template, fs.FS, error) {
	var subdir string
	baseTpls, err := htmltemplate.New("").
		Funcs(templateFuncMap).
		Funcs(sprig.FuncMap()).
		ParseFS(
			tfs,
			"html/layout/*.gohtml",
			"html/partial/*.gohtml",
		)
	if err != nil {
		return nil, nil, err
	}
	subdir = "html/view"

	viewSub, err := fs.Sub(tfs, subdir)
	if err != nil {
		return nil, nil, err
	}
	return baseTpls, viewSub, nil
}

func ParseTemplates(templatesDir string) (*TContainer, error) {
	tfs, err := getTemplateFS(templatesDir)
	if err != nil {
		return nil, err
	}

	tc := &TContainer{
		tfs:   tfs,
		tmap:  make(TemplateMap, 0),
		embed: templatesDir == embedMagicStr,
	}

	nonViewHtmlTemplates, viewSub, err := getBaseTpl(tfs)
	if err != nil {
		return nil, err
	}
	err = fs.WalkDir(viewSub, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ".gohtml" {
			name := filepath.Base(p)
			c, err := nonViewHtmlTemplates.Clone()
			if err != nil {
				return err
			}
			t, err := c.New(name).Funcs(templateFuncMap).Funcs(sprig.FuncMap()).ParseFS(tfs, fmt.Sprintf("html/view/%s", name))
			if err != nil {
				return err
			}
			tc.tmap[name] = t
		}
		return nil
	})

	return tc, err
}
