package driver

import (
	"github.com/KDF5000/COSB/utils"
	"github.com/ncw/swift"
	"math/rand"
)

type CephDriver struct {
	client swift.Connection
}

func NewCephDriver(username string, apikey string, authUrl string) *CephDriver {
	c := swift.Connection{
		UserName: username,
		ApiKey:   apikey,
		AuthUrl:  authUrl,
	}
	if err := c.Authenticate(); err != nil {
		panic(err)
		return nil
	}

	return &CephDriver{
		client: c,
	}
}

func (d *CephDriver) Get(bucket string, key string) ([]byte, error) {
	data, err := d.client.ObjectGetBytes(bucket, key)
	if err == nil {
		return data, nil
	}
	return nil, err
}

func (d *CephDriver) Put(bucket string, fileName string, fileSize int64) (fileKey string, err error) {
	reader := utils.NewFakeReader(rand.Uint64(), fileSize)
	bytes := make([]byte, fileSize)
	if _, err := reader.Read(bytes); err != nil {
		return "", err
	}
	if err := d.client.ObjectPutBytes(bucket, fileName, bytes, ""); err != nil {
		return "", err
	}
	return fileName, nil
}
func (d *CephDriver) Delete(bucket string, key string) error {
	if err := d.client.ObjectDelete(bucket, key); err != nil {
		return err
	}
	return nil
}
