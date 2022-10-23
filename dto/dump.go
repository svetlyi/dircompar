package dto

type File struct {
	Name string
	Hash string
	Size int64
}

type Dump struct {
	Files []File
	Path  string
}
