package driver

// Error used for error message

// Client is an interface for different driver
type Client interface {
	SendRequest(data []byte) (bool, error)
}
