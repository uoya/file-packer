package main

import (
	"github.com/uoya/file-packer/compressor"
	"github.com/uoya/file-packer/fileutil"
	"os"
	"path"
)

type ServiceName string

type Service struct {
	Name             ServiceName
	Compress         compressor.CompressOption
	Includes         int
	TargetExtensions []string
}

type ServiceOption struct {
	Name             ServiceName               `json:"name" validate:"required,min=1"`
	BaseNameSuffix   string                    `json:"baseNameSuffix,omitempty" validate:"omitempty"`
	TargetExtensions []string                  `json:"targetExtensions" validate:"min=1"`
	Includes         int                       `json:"includes" validate:"min=1"`
	Compress         compressor.CompressOption `json:"compress" validate:"required,min=1"`
}

func NewService(opt ServiceOption) Service {
	return Service{
		Name:             opt.Name,
		Compress:         opt.Compress,
		Includes:         opt.Includes,
		TargetExtensions: opt.TargetExtensions,
	}
}

func (s Service) Check(file fileutil.File) ([]fileutil.File, error) {
	// 再起で存在確認をしないといけない
	//TODO このままだと PIXTA と shutterstockに対応できない
	var checked []fileutil.File
	for _, ext := range s.TargetExtensions {
		fileName := file.Base().FullName(ext)
		_, err := os.Stat(path.Join(file.Root, string(fileName)))
		if err != nil {
			return []fileutil.File{}, err
		}
		checked = append(checked, fileutil.File{Name: fileName, Root: file.Root})
	}
	return checked, nil
}

func (s Service) Exec(sources []fileutil.File) error {
	dstDir := fileutil.DirectoryName(s.Name)
	compressor := compressor.NewCompressor(s.Compress)
	// TODO dstPath にしたほうが良いかな
	if err := compressor.Compress(sources, dstDir); err != nil {
		return err
	}
	return nil
}
