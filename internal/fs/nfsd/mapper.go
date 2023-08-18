package nfsd

import (
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/rtgnx/s3nfs/internal/fs/noop"
	"github.com/rtgnx/s3nfs/internal/fs/s3"
)

var (
	ErrShareNotFound = func(share string) error { return fmt.Errorf("share not found: %s", share) }
)

/*
  Mapper FS maps shares to their corresponding filesystems

  Query path => /share/some/file/file.txt
*/

type MapperFS struct {
	noop.NoOP
	shares map[string]billy.Filesystem
}

func NewMapperFS(cfg Config) billy.Filesystem {
	fs := &MapperFS{shares: make(map[string]billy.Filesystem)}
	var err error
	for k, v := range cfg.Regions {
		if fs.shares[k], err = v.Open(); err != nil {
			log.Println(err)
			delete(fs.shares, k)
			continue
		}
	}
	return fs
}

func (m *MapperFS) Resolve(filename string) (lfs billy.Filesystem, p string, err error) {
	for k, v := range m.shares {
		if strings.HasPrefix(filename, k) {
			return v, strings.TrimPrefix(filename, k), nil
		}
	}

	return nil, "", ErrShareNotFound(filename)
}

// Open implements billy.Filesystem.
func (m *MapperFS) Open(filename string) (billy.File, error) {
	fs, rpath, err := m.Resolve(filename)
	if err != nil {
		return nil, err
	}
	return fs.Open(rpath)
}

// OpenFile implements billy.Filesystem.
func (m *MapperFS) OpenFile(filename string, flag int, perm fs.FileMode) (billy.File, error) {
	return m.Open(filename)
}

// ReadDir implements billy.Filesystem.
func (m *MapperFS) ReadDir(path string) ([]fs.FileInfo, error) {

	if len(path) == 0 {
		dirs := []fs.FileInfo{}
		for k, _ := range m.shares {
			dirs = append(dirs, s3.NewDirInfo(k))
		}
		return dirs, nil
	}

	if share, ok := m.shares[strings.TrimPrefix(path, "/")]; ok {
		return share.ReadDir("/")
	}

	fs, rpath, err := m.Resolve(path)
	if err != nil {
		return nil, err
	}

	return fs.ReadDir(rpath)
}

// Stat implements billy.Filesystem.
func (m *MapperFS) Stat(filename string) (fs.FileInfo, error) {
	if len(filename) == 0 || filename == "/" {
		log.Println("MAPPER_STAT_ROOT ", filename)
		return s3.NewDirInfo("/"), nil
	}

	if _, ok := m.shares[strings.TrimPrefix(filename, "/")]; ok {
		return s3.NewDirInfo(strings.TrimPrefix(filename, "/")), nil
	}

	fs, rpath, err := m.Resolve(filename)
	if err != nil {
		return nil, err
	}
	log.Printf("MAPPER_STAT share: %s ", rpath)
	return fs.Stat(rpath)
}
