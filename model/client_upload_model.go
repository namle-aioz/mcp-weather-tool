package model

import (
	"io"
)

type UploadVideoClient struct {
	FileName string
	File     io.Reader
	FileSize int64
}

type MediaInfo struct {
	MediaID   string
	Name      string
	Size      int32
	Duration  float32
	CreatedAt string
}
