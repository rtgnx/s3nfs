package noop

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

var (
	ErrNotImplemented = errors.New("err: not implemented")
)

type NoOP struct{}

// Chroot implements billy.Filesystem.
func (*NoOP) Chroot(path string) (billy.Filesystem, error) {
	return nil, ErrNotImplemented
}

// Create implements billy.Filesystem.
func (*NoOP) Create(filename string) (billy.File, error) {
	return nil, ErrNotImplemented
}

// Join implements billy.Filesystem.
func (*NoOP) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Lstat implements billy.Filesystem.
func (*NoOP) Lstat(filename string) (fs.FileInfo, error) {
	return nil, ErrNotImplemented
}

// MkdirAll implements billy.Filesystem.
func (*NoOP) MkdirAll(filename string, perm fs.FileMode) error {
	return ErrNotImplemented
}

// Open implements billy.Filesystem.
func (*NoOP) Open(filename string) (billy.File, error) {
	return nil, ErrNotImplemented
}

// OpenFile implements billy.Filesystem.
func (*NoOP) OpenFile(filename string, flag int, perm fs.FileMode) (billy.File, error) {
	return nil, ErrNotImplemented
}

// ReadDir implements billy.Filesystem.
func (*NoOP) ReadDir(path string) ([]fs.FileInfo, error) {
	return []fs.FileInfo{}, ErrNotImplemented
}

// Readlink implements billy.Filesystem.
func (*NoOP) Readlink(link string) (string, error) {
	return "", ErrNotImplemented
}

// Remove implements billy.Filesystem.
func (*NoOP) Remove(filename string) error {
	return ErrNotImplemented
}

// Rename implements billy.Filesystem.
func (*NoOP) Rename(oldpath string, newpath string) error {
	return ErrNotImplemented
}

// Root implements billy.Filesystem.
func (*NoOP) Root() string {
	return ""
}

// Stat implements billy.Filesystem.
func (*NoOP) Stat(filename string) (fs.FileInfo, error) {
	return nil, ErrNotImplemented
}

// Symlink implements billy.Filesystem.
func (*NoOP) Symlink(target string, link string) error {
	return ErrNotImplemented
}

// TempFile implements billy.Filesystem.
func (*NoOP) TempFile(dir string, prefix string) (billy.File, error) {
	return nil, ErrNotImplemented
}

func NewNoOP() billy.Filesystem {
	return &NoOP{}
}
