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

const (
	tmapIdxHtml = iota
	tmapIdxTxt
)

type TContainer struct {
	tmaps [2]TemplateMap
	tfs   fs.FS
	embed bool
}

type TGetter interface {
	Get(string) (TemplateIf, error)
}

type anonGetter struct {
	tpls TemplateMap
}

func (ag *anonGetter) Get(name string) (TemplateIf, error) {
	return ag.tpls.Get(name)
}

func MockTContainer(tplm TemplateMap) TGetter {
	return &anonGetter{tplm}
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

func getBaseHtmlTpl(tfs fs.FS) (*htmltemplate.Template, fs.FS, error) {
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

	viewSub, err := fs.Sub(tfs, "html/view")
	if err != nil {
		return nil, nil, err
	}
	return baseTpls, viewSub, nil
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

func getBaseTxtTpl(tfs fs.FS) (*txttemplate.Template, fs.FS, error) {
	baseTpls := txttemplate.New("").
		Funcs(templateFuncMap).
		Funcs(sprig.FuncMap())

	viewSub, err := fs.Sub(tfs, "txt")
	if err != nil {
		return nil, nil, err
	}
	return baseTpls, viewSub, nil
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
