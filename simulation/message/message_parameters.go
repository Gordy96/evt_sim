//go:generate go run ../../internal/generation/getter_setter.go -struct=Parameters -chain -templates=../../internal/generation

package message

type Parameters struct {
	params map[string]any
}
