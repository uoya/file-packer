package service

import (
	"github.com/uoya/ImagePacker/fileutil"
	"os"
)

type ImageMart struct{}

func (i ImageMart) Name() Name {
	return "イメージマート"
}

func (i ImageMart) Check(baseName fileutil.FileBaseName) ([]fileutil.FileName, error) {
	ext := []fileutil.Extension{fileutil.Eps, fileutil.Jpg}
	var checked []fileutil.FileName
	for _, ext := range ext {
		t := baseName.FullName(ext)
		_, err := os.Stat(string(baseName.FullName(ext)))
		if err != nil {
			return []fileutil.FileName{}, err
		}
		checked = append(checked, t)
	}
	return checked, nil
}

func (i ImageMart) Exec(sources []fileutil.FileName) error {
	for _, src := range sources {
		dstDir := fileutil.DirectoryName(i.Name())
		_, err := fileutil.CopyFile(src, dstDir)
		if err != nil {
			return err
		}
	}
	return nil
}
