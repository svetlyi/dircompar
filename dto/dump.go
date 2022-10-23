package dto

import (
	"path/filepath"
	"strings"
)

type File struct {
	Name string
	Hash string
	Size int64
}

func (f File) GetCleanUnixName() string {
	return strings.ReplaceAll(filepath.Clean(f.Name), "\\", "/")
}

type Dump struct {
	Files []File
	Path  string
}
