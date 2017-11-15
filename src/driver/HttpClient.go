package driver

import (
	"fmt"
	"net"
)

// HTTPClient for testing the distribution
type HTTPClient struct {
	host string
	port int
}

// NewHTTPClient creates a HTTPClient struct
func NewHTTPClient(host string, port int) *HTTPClient {
	return &HTTPClient{host: host, port: port}
}

// SendRequest impl
func (c *HTTPClient) SendRequest(data []byte) (bool, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return false, err
	}
	conn.Write(data)
	conn.Close()
	return true, nil
}
