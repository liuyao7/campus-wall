// pkg/storage/local.go

package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"
)

type LocalStorage struct {
    basePath   string // 基础存储路径
    baseURL    string // 基础访问URL
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
    return &LocalStorage{
        basePath: basePath,
        baseURL:  baseURL,
    }
}

func (l *LocalStorage) Upload(file *multipart.FileHeader, directory string) (string, error) {
    // 创建目录
    uploadDir := filepath.Join(l.basePath, directory)
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create directory: %w", err)
    }

    // 生成文件名
    ext := path.Ext(file.Filename)
    fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    filePath := filepath.Join(uploadDir, fileName)

    // 打开源文件
    src, err := file.Open()
    if err != nil {
        return "", fmt.Errorf("failed to open file: %w", err)
    }
    defer src.Close()

    // 创建目标文件
    dst, err := os.Create(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %w", err)
    }
    defer dst.Close()

    // 复制文件内容
    if _, err = io.Copy(dst, src); err != nil {
        return "", fmt.Errorf("failed to copy file: %w", err)
    }

    // 返回可访问的URL
    return fmt.Sprintf("%s/%s/%s", l.baseURL, directory, fileName), nil
}

func (l *LocalStorage) Delete(fileURL string) error {
    // 从URL中提取文件路径
    relativePath := path.Clean(fileURL[len(l.baseURL):])
    filePath := filepath.Join(l.basePath, relativePath)

    // 删除文件
    if err := os.Remove(filePath); err != nil {
        return fmt.Errorf("failed to delete file: %w", err)
    }

    return nil
}