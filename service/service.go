package service

import "github.com/uoya/ImagePacker/fileutil"

type Name string
type Service interface {
	Name() Name
	Check(fileutil.FileBaseName) ([]fileutil.FileName, error)
	Exec([]fileutil.FileName) error
}
