package sensu

type Processor interface {
	Start() error
	Close()
}
