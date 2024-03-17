package service

import (
	"github.com/uoya/file-packer/compressor"
	"github.com/uoya/file-packer/fileutil"
	"os"
	"path"
)

type AdobeStock struct{}

func (a AdobeStock) Name() Name {
	return "AdobeStock"
}
func (a AdobeStock) Check(file fileutil.File) ([]fileutil.File, error) {
	ext := []fileutil.Extension{fileutil.Eps, fileutil.Jpg}
	var checked []fileutil.File
	for _, ext := range ext {
		fileName := file.Base().FullName(ext)
		_, err := os.Stat(path.Join(file.Root, string(fileName)))
		if err != nil {
			return []fileutil.File{}, err
		}
		checked = append(checked, fileutil.File{Name: fileName, Root: file.Root})
	}
	return checked, nil
}

func (a AdobeStock) Exec(sources []fileutil.File) error {
	dstDir := fileutil.DirectoryName(a.Name())
	compressor := compressor.NewCompressor(compressor.CompressZip)
	if err := compressor.Compress(sources, fileutil.DirectoryName(path.Join(sources[0].Root, string(dstDir)))); err != nil {
		return err
	}
	return nil
}
