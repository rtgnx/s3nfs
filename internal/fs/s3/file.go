package s3

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/minio/minio-go/v7"
)

const DefaultACL = 0540

func NewFileInfo(info minio.ObjectInfo) fs.FileInfo {
	return &FileInfo{
		name:    filepath.Base(strings.TrimPrefix(info.Key, "/")),
		size:    info.Size,
		isDir:   false,
		lastMod: info.LastModified,
		stat: syscall.Stat_t{
			Size: info.Size,
		},
	}
}

func NewDirInfo(prefix string) fs.FileInfo {
	return &FileInfo{
		name:  filepath.Base(strings.TrimPrefix(prefix, "/")),
		size:  4096,
		isDir: true,
	}
}

// IsDir implements fs.FileInfo.
func (fi *FileInfo) IsDir() bool {
	return fi.isDir
}

// ModTime implements fs.FileInfo.
func (fi *FileInfo) ModTime() time.Time {
	return fi.lastMod
}

// Mode implements fs.FileInfo.
func (fi *FileInfo) Mode() fs.FileMode {
	return DefaultACL
}

// Name implements fs.FileInfo.
func (fi *FileInfo) Name() string {
	return fi.name
}

// Size implements fs.FileInfo.
func (fi *FileInfo) Size() int64 {
	return fi.size
}

// Sys implements fs.FileInfo.
func (fi *FileInfo) Sys() any {
	return fi.stat
}

// FILE

// Close implements billy.File.
func (f *File) Close() error {
	return f.fd.Close()
}

// Lock implements billy.File.
func (f *File) Lock() error {
	log.Printf("LOCK: %s", f.s3Prefix)
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
	log.Printf("UNLOCK: %s", f.s3Prefix)
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
