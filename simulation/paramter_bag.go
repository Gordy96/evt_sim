package simulation

type ParameterBag struct {
	params map[string]any
}

func (p *ParameterBag) GetParam(name string) (any, bool) {
	if p.params == nil {
		return nil, false
	}

	ret, ok := p.params[name]
	return ret, ok
}

func (p *ParameterBag) SetParam(name string, value any) {
	if p.params == nil {
		p.params = make(map[string]any)
	}

	p.params[name] = value
}

func (p *ParameterBag) RemoveParam(name string) {
	if p.params != nil {
		delete(p.params, name)
	}
}
