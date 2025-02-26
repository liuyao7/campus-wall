// pkg/storage/storage.go

package storage

import (
	"mime/multipart"
)

// Storage 存储接口
type Storage interface {
    // Upload 上传文件
    Upload(file *multipart.FileHeader, directory string) (string, error)
    // Delete 删除文件
    Delete(fileURL string) error
}

// FileInfo 文件信息
type FileInfo struct {
    URL      string `json:"url"`
    FileName string `json:"file_name"`
    Size     int64  `json:"size"`
    Type     string `json:"type"`
}