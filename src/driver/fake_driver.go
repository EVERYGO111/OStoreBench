package driver

import "fmt"

type FakeDriver struct {
	client *HTTPClient
}

func NewFakeDriver(host string, port int) *FakeDriver {
	return &FakeDriver{client: NewHTTPClient(host, port)}
}

func (fd *FakeDriver) Get(key string) ([]byte, error) {
	data := fmt.Sprintf("Get %s", key)
	if res, err := fd.client.SendRequest([]byte(data)); res {
		return []byte(data), nil
	} else {
		return nil, err
	}
}
func (fd *FakeDriver) Put(fileName string, fileSize int64) (fileKey string, err error) {
	data := fmt.Sprintf("Put %s %d", fileName, fileSize)
	if res, err := fd.client.SendRequest([]byte(data)); res {
		return fileName, nil
	} else {
		return "", err
	}
}
func (fd *FakeDriver) Delete(key string) error {
	data := fmt.Sprintf("Delete %s", key)
	if res, err := fd.client.SendRequest([]byte(data)); res {
		return nil
	} else {
		return err
	}
}
