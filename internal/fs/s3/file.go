package s3

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/minio/minio-go/v7"
)

func NewFileInfo(info minio.ObjectInfo) fs.FileInfo {
	return &FileInfo{}
}

// IsDir implements fs.FileInfo.
func (fi *FileInfo) IsDir() bool {
	return fi.isDir
}

// ModTime implements fs.FileInfo.
func (fi *FileInfo) ModTime() time.Time {
	return time.Now()
}

// Mode implements fs.FileInfo.
func (fi *FileInfo) Mode() fs.FileMode {
	return fs.FileMode(fi.stat.Mode)
}

// Name implements fs.FileInfo.
func (fi *FileInfo) Name() string {
	return fi.name
}

// Size implements fs.FileInfo.
func (fi *FileInfo) Size() int64 {
	return fi.stat.Size
}

// Sys implements fs.FileInfo.
func (fi *FileInfo) Sys() any {
	return fi.stat
}

func NewFile() billy.File {
	return &File{}
}

// FILE

// Close implements billy.File.
func (f *File) Close() error {
	return f.fd.Close()
}

// Lock implements billy.File.
func (f *File) Lock() error {
	if _, loaded := f.s3.locks.LoadOrStore(f.s3Prefix, true); loaded {
		return os.ErrDeadlineExceeded
	}
	f.locked = true
	return nil
}

// Name implements billy.File.
func (f *File) Name() string {
	return filepath.Base(f.s3Prefix)
}

// Read implements billy.File.
func (f *File) Read(p []byte) (n int, err error) {
	return f.fd.Read(p)
}

// ReadAt implements billy.File.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return f.fd.ReadAt(p, off)
}

// Seek implements billy.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.Seek(offset, whence)
}

// Truncate implements billy.File.
func (f *File) Truncate(size int64) error {
	return ErrNotImplemented
}

// Unlock implements billy.File.
func (f *File) Unlock() error {
	if f.locked {
		f.s3.locks.Delete(f.s3Prefix)
		return nil
	}

	return errors.New("no lock acquired")
}

// Write implements billy.File.
func (*File) Write(p []byte) (n int, err error) {
	return 0, ErrNotImplemented
}
