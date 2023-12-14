package resources

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
