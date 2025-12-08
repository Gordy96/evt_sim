//go:generate go run ../../internal/generation/getter_setter.go -struct=Parameters -chain -templates=../../internal/generation

package message

import "encoding/json"

type Parameters struct {
	params map[string]any
}

func (p *Parameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.params)
}
