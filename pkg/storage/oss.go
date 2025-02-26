// pkg/storage/oss.go

package storage

import (
	"fmt"
	"mime/multipart"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSStorage struct {
    client     *oss.Client
    bucket     *oss.Bucket
    bucketName string
    baseURL    string
}

func NewOSSStorage(endpoint, accessKeyID, accessKeySecret, bucketName string) (*OSSStorage, error) {
    client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
    if err != nil {
        return nil, err
    }

    bucket, err := client.Bucket(bucketName)
    if err != nil {
        return nil, err
    }

    return &OSSStorage{
        client:     client,
        bucket:     bucket,
        bucketName: bucketName,
        baseURL:    fmt.Sprintf("https://%s.%s", bucketName, endpoint),
    }, nil
}

func (o *OSSStorage) Upload(file *multipart.FileHeader, directory string) (string, error) {
    // 打开源文件
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    // 生成OSS对象名
    ext := path.Ext(file.Filename)
    objectName := fmt.Sprintf("%s/%d%s", directory, time.Now().UnixNano(), ext)

    // 上传到OSS
    if err := o.bucket.PutObject(objectName, src); err != nil {
        return "", err
    }

    return fmt.Sprintf("%s/%s", o.baseURL, objectName), nil
}

func (o *OSSStorage) Delete(fileURL string) error {
    // 从URL中提取对象名
    objectName := path.Clean(fileURL[len(o.baseURL)+1:])

    // 删除对象
    if err := o.bucket.DeleteObject(objectName); err != nil {
        return err
    }

    return nil
}