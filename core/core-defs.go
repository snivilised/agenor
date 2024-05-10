package core

// core contains universal definitions and handles cross cutting concerns
// try to keep to a minimum to reduce rippling changes

type TraverseResult interface {
}

type DuffResult struct {
}

type Client func(node *Node) error
