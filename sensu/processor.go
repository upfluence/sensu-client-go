package sensu

type Processor interface {
	SetClient(c *Client) error
	Start() error
	Close()
}
