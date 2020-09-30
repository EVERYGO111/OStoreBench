package driver

import (
	"github.com/KGXarwen/COSB/utils"
	"github.com/linxGnu/goseaweedfs"
	"math/rand"
	"time"
)

type WeedDriver struct {
	client *goseaweedfs.Seaweed
}

func NewWeeDriver(master string) *WeedDriver {
	return &WeedDriver{
		client: goseaweedfs.NewSeaweed("http", master, nil, 64*1024*1024, 2*time.Minute),
	}
}

func (w *WeedDriver) Get(bucket string, key string) ([]byte, error) {
	//kgx byteRead, err := w.client.GetFile(key)
	fileName, byteRead, err := w.client.DownloadFile(key, nil)
        if err != nil {
		return nil, err
	}
        _=fileName
	return byteRead, nil
}

func (w *WeedDriver) Put(bucket string, fileName string, fileSize int64) (fileKey string, err error) {
	_, fileId, err := w.client.Upload(utils.NewFakeReader(rand.Uint64(), fileSize), fileName, fileSize, "", "")
	if err != nil {
		return "", err
	}
	return fileId, nil
}

func (w *WeedDriver) Delete(bucket string, key string) error {
	if err := w.client.DeleteFile(key, nil); err != nil {
		return err
	}
	return nil
}
