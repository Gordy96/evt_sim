package simulation

type ParameterBag struct {
	params map[string]interface{}
}

func (p *ParameterBag) Param(name string) (any, bool) {
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
