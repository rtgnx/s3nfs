package s3

import (
	"errors"
	"sync"
	"syscall"

	"github.com/minio/minio-go/v7"
)

var (
	ErrNotImplemented = errors.New("feature not implemented")
)

type File struct {
	s3       *S3
	s3Prefix string
	locked   bool
	fd       *minio.Object
}

type FileInfo struct {
	stat  syscall.Stat_t
	isDir bool
	name  string
	size  int64
}

type S3 struct {
	client        *minio.Client
	bucket        string
	locks         sync.Map
	index         map[string]*FileInfo
	localCache    bool
	localCacheDir string
}
