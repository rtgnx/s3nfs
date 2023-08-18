package nfsd

import (
	"github.com/go-git/go-billy/v5"
	"github.com/rtgnx/s3nfs/internal/fs/s3"
)

type Config struct {
	Global  string
	Regions map[string]s3.S3Config
}

func (cfg Config) Open() (billy.Filesystem, error) { return nil, nil }
