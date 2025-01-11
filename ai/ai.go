package ai

type Ai interface {
	Init() error
	GetMessage() (string, error)
}
