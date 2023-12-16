package resources

import (
	"embed"
	"fmt"
	htmltemplate "html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	txttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

const embedMagicStr = "embed"

//go:embed templates
var templateEmbedFS embed.FS

type TemplateIf interface {
	Execute(wr io.Writer, data any) error
	ExecuteTemplate(io.Writer, string, any) error
}

type TemplateMap map[string]TemplateIf

func (tm *TemplateMap) Get(name string) (TemplateIf, error) {
	if v, ok := (*tm)[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("template not found for name %s", name)
}

type ttype int

const (
	tmapIdxHtml ttype = iota
	tmapIdxTxt
)

type TContainer struct {
	tmaps [2]TemplateMap
	tfs   fs.FS
	embed bool
}

func (tc *TContainer) Get(name string) (TemplateIf, error) {
	ttype := tmapIdxHtml
	if strings.HasSuffix(name, ".gotxt") {
		ttype = tmapIdxTxt
	}

	if tc.embed {
		return tc.tmaps[ttype].Get(name)
	}

	switch ttype {
	case tmapIdxHtml:
		baseTpl, viewSub, err := getBaseHtmlTpl(tc.tfs)
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
	case tmapIdxTxt:
		baseTpl, viewSub, err := getBaseTxtTpl(tc.tfs)
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
	default:
		return nil, fmt.Errorf("unknown template type")
	}
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

func getBaseTpl(tfs fs.FS, tt ttype) (any, fs.FS, error) {
	var t any
	var subdir string
	switch tt {
	case tmapIdxHtml:
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
		t = baseTpls
		subdir = "html/view"
	case tmapIdxTxt:
		baseTpls := txttemplate.New("").
			Funcs(templateFuncMap).
			Funcs(sprig.FuncMap())
		t = baseTpls
		subdir = "txt"
	default:
		return nil, nil, fmt.Errorf("unknown type")
	}

	viewSub, err := fs.Sub(tfs, subdir)
	if err != nil {
		return nil, nil, err
	}
	return t, viewSub, nil
}

func getBaseHtmlTpl(tfs fs.FS) (*htmltemplate.Template, fs.FS, error) {
	tpl, sub, err := getBaseTpl(tfs, tmapIdxHtml)
	if err != nil {
		return nil, nil, err
	}
	t := tpl.(*htmltemplate.Template)
	return t, sub, err
}

func getBaseTxtTpl(tfs fs.FS) (*txttemplate.Template, fs.FS, error) {
	tpl, sub, err := getBaseTpl(tfs, tmapIdxTxt)
	if err != nil {
		return nil, nil, err
	}
	t := tpl.(*txttemplate.Template)
	return t, sub, err
}

func ParseHtmlTemplates(tfs fs.FS) (TemplateMap, error) {
	templates := make(TemplateMap, 0)

	nonViewHtmlTemplates, viewSub, err := getBaseHtmlTpl(tfs)
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
			t, err := c.New(name).
				Funcs(templateFuncMap).
				Funcs(sprig.FuncMap()).
				ParseFS(
					tfs, fmt.Sprintf("html/view/%s", name),
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

func ParseTxtTemplates(tfs fs.FS) (TemplateMap, error) {
	templates := make(TemplateMap, 0)

	nonViewTxtTemplates, viewSub, err := getBaseTxtTpl(tfs)
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(viewSub, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ".gotxt" {
			name := filepath.Base(p)
			c, err := nonViewTxtTemplates.Clone()
			if err != nil {
				return err
			}
			t, err := c.New(name).
				Funcs(templateFuncMap).
				Funcs(sprig.FuncMap()).
				ParseFS(
					tfs, fmt.Sprintf("txt/%s", name),
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

func ParseTemplates(templatesDir string) (*TContainer, error) {
	tfs, err := getTemplateFS(templatesDir)
	if err != nil {
		return nil, err
	}

	tc := &TContainer{
		tfs:   tfs,
		embed: templatesDir == embedMagicStr,
	}

	htmlTemplates, err := ParseHtmlTemplates(tfs)
	if err != nil {
		return nil, err
	}
	tc.tmaps[tmapIdxHtml] = htmlTemplates

	txtTemplates, err := ParseTxtTemplates(tfs)
	if err != nil {
		return nil, err
	}
	tc.tmaps[tmapIdxTxt] = txtTemplates

	return tc, nil
}
