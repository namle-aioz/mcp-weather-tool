package model

import (
	"io"
)

type UploadVideoClient struct {
	FileName string
	File     io.Reader
	FileSize int64
}
