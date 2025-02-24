package search

type Search interface {
	Search(query string) (any, error)
}
