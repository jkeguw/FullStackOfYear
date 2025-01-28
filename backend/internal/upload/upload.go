package upload

import (
	"fmt"
	"mime/multipart"
	"path"
)

func CheckFile(file *multipart.FileHeader) error {
	// 检查大小
	if file.Size > DefaultConfig.MaxSize<<20 {
		return fmt.Errorf("file exceeds the maximum size：%dMB", DefaultConfig.MaxSize)
	}

	// 检查类型
	ext := path.Ext(file.Filename)
	for _, allowType := range DefaultConfig.AllowTypes {
		if ext == allowType {
			return nil
		}
	}
	return fmt.Errorf("unsupported file type：%s", ext)
}
