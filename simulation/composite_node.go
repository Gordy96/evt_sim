package simulation

type CompositeNode interface {
	Node
	Children() []Node
}
