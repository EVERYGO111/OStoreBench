package driver

// Error used for error message

// Client is an interface for different driver
type Client interface {
	SendRequest(data []byte) (bool, error)
}

// Driver interface for different storage backend driver
type Driver interface {
	Get(key string) ([]byte, error)
	Put(key string, value interface{}) error
	Delete(key string) error
}
