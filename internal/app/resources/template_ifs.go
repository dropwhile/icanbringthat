package resources

import (
	"fmt"
	"io"
)

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
