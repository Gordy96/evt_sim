package simulation

type CompositeNode interface {
	Node
	Children() []Node
	GetChild(id string) Node
}
