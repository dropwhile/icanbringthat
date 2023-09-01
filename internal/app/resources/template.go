// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package resources

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
)

//go:embed templates
var templateEmbedFS embed.FS

type TExecuter interface {
	Execute(wr io.Writer, data any) error
	ExecuteTemplate(io.Writer, string, any) error
}

type TGetter interface {
	Get(string) (TExecuter, error)
}

type TemplateMap map[string]TExecuter

func (tm *TemplateMap) Get(name string) (TExecuter, error) {
	if v, ok := (*tm)[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("template not found for name %s", name)
}

type TContainer struct {
	tmap  TemplateMap
	tfs   fs.FS
	embed bool
}

func (tc *TContainer) Get(name string) (TExecuter, error) {
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
	return c.New(name).
		Funcs(templateFuncMap).
		Funcs(sprig.FuncMap()).
		ParseFS(viewSub, name)
}

func getTemplateFS(loc Location) (fs.FS, error) {
	var templateFS fs.FS
	switch loc {
	case Embed:
		var err error
		templateFS, err = fs.Sub(templateEmbedFS, "templates")
		if err != nil {
			return templateFS, err
		}
	case Filesystem:
		sdir := "./internal/app/resources/templates/"
		templateFS = os.DirFS(sdir)
	}
	return templateFS, nil
}

func getBaseTpl(tfs fs.FS) (*template.Template, fs.FS, error) {
	var subdir string
	baseTpls, err := template.New("").
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

func ParseTemplates(loc Location) (*TContainer, error) {
	tfs, err := getTemplateFS(loc)
	if err != nil {
		return nil, err
	}

	tc := &TContainer{
		tfs:   tfs,
		tmap:  make(TemplateMap, 0),
		embed: loc == Embed,
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
			t, err := c.New(name).
				Funcs(templateFuncMap).
				Funcs(sprig.FuncMap()).
				ParseFS(tfs, fmt.Sprintf("html/view/%s", name))
			if err != nil {
				return err
			}
			tc.tmap[name] = t
		}
		return nil
	})

	return tc, err
}
