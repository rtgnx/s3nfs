package s3

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/minio/minio-go/v7"
)

func NewS3() billy.Filesystem {
	return &S3{}
}

// Chroot implements billy.Filesystem.
func (*S3) Chroot(path string) (billy.Filesystem, error) {
	return nil, ErrNotImplemented
}

// Create implements billy.Filesystem.
func (*S3) Create(filename string) (billy.File, error) {
	return nil, ErrNotImplemented
}

// Join implements billy.Filesystem.
func (*S3) Join(elem ...string) string {
	return filepath.Join(append([]string{"/"}, elem...)...)
}

// Lstat implements billy.Filesystem.
func (*S3) Lstat(filename string) (fs.FileInfo, error) {
	return nil, ErrNotImplemented
}

// MkdirAll implements billy.Filesystem.
func (s3 *S3) MkdirAll(filename string, perm fs.FileMode) error {
	return ErrNotImplemented
}

// Open implements billy.Filesystem.
func (s3 *S3) Open(filename string) (billy.File, error) {
	filename = strings.TrimPrefix(filename, "/")
	_, err := s3.client.GetObject(context.Background(), s3.bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &File{s3: s3, s3Prefix: filename, locked: false}, nil
}

// OpenFile implements billy.Filesystem.
func (*S3) OpenFile(filename string, flag int, perm fs.FileMode) (billy.File, error) {
	return nil, ErrNotImplemented
}

// ReadDir implements billy.Filesystem.
func (s3 *S3) ReadDir(path string) ([]fs.FileInfo, error) {
	list := []fs.FileInfo{}
	pattern := filepath.Join(path, "*")
	for k, v := range s3.index {
		if k == "/" {
			continue
		}
		ok, err := filepath.Match(pattern, k)
		if err != nil {
			log.Print(err)
			continue
		}

		if ok {
			list = append(list, v)
		}

	}

	list = append(list, &FileInfo{isDir: true, name: ".", size: 4096})
	list = append(list, &FileInfo{isDir: true, name: "..", size: 4096})

	return list, nil
}

// Readlink implements billy.Filesystem.
func (*S3) Readlink(link string) (string, error) {
	return "", ErrNotImplemented
}

// Remove implements billy.Filesystem.
func (s3 *S3) Remove(filename string) error {
	return s3.client.RemoveObject(context.Background(), s3.bucket, filename, minio.RemoveObjectOptions{})
}

// Rename implements billy.Filesystem.
func (s3 *S3) Rename(oldpath string, newpath string) error { return ErrNotImplemented }

// Root implements billy.Filesystem.
func (*S3) Root() string { return "" }

// Stat implements billy.Filesystem.
func (s3 *S3) Stat(filename string) (fs.FileInfo, error) {
	log.Printf("STAT %s", filename)

	if filename == "" {
		filename = "/"
	}

	// if strings.Compare(filename, "/") == 0 || len(filename) == 0 {
	// 	return &FileInfo{isDir: true, name: ".", size: 4096}, nil
	// }

	stat, ok := s3.index[filename]
	if !ok {
		return stat, fmt.Errorf("no such file or directory")
	}
	return stat, nil
}

// Symlink implements billy.Filesystem.
func (*S3) Symlink(target string, link string) error {
	return ErrNotImplemented
}

// TempFile implements billy.Filesystem.
func (*S3) TempFile(dir string, prefix string) (billy.File, error) {
	return nil, ErrNotImplemented
}
