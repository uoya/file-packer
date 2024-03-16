package service

import (
	"github.com/uoya/ImagePacker/fileutil"
	"os"
)

type AdobeStock struct{}

func (a AdobeStock) Name() Name {
	return "AdobeStock"
}
func (a AdobeStock) Check(baseName fileutil.FileBaseName) ([]fileutil.FileName, error) {
	ext := []fileutil.Extension{fileutil.Eps, fileutil.Jpg}
	var checked []fileutil.FileName
	for _, ext := range ext {
		t := baseName.FullName(ext)
		_, err := os.Stat(string(t))
		if err != nil {
			return []fileutil.FileName{}, err
		}
		checked = append(checked, t)
	}
	return checked, nil
}

func (a AdobeStock) Exec(sources []fileutil.FileName) error {
	dstDir := fileutil.DirectoryName(a.Name())
	if err := fileutil.ZipFiles(sources, dstDir); err != nil {
		return err
	}
	return nil
}
