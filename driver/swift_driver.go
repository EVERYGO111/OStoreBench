package driver

import (
	"github.com/KGXarwen/COSB/utils"
	"github.com/ncw/swift"
	"math/rand"
)

type SwiftDriver struct {
	client swift.Connection
}

func NewSwiftDriver(username string, apikey string, authUrl string, domainName string, tenant string) *SwiftDriver {
	c := swift.Connection{
		UserName: username,
		ApiKey:   apikey,
		AuthUrl:  authUrl,
		Domain:   domainName,
		Tenant:   tenant, //project name in V3
	}
	if err := c.Authenticate(); err != nil {
		panic(err)
		return nil
	}

	return &SwiftDriver{
		client: c,
	}
}

func (d *SwiftDriver) Get(bucket string, key string) ([]byte, error) {
	data, err := d.client.ObjectGetBytes(bucket, key)
	if err == nil {
		return data, nil
	}
	return nil, err
}

func (d *SwiftDriver) Put(bucket string, fileName string, fileSize int64) (fileKey string, err error) {
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
func (d *SwiftDriver) Delete(bucket string, key string) error {
	if err := d.client.ObjectDelete(bucket, key); err != nil {
		return err
	}
	return nil
}
