package main

import (
	"errors"
	"github.com/uoya/file-packer/compressor"
	"github.com/uoya/file-packer/fileutil"
	"os"
	"path/filepath"
)

type ServiceName string

type Service struct {
	Name             ServiceName
	BaseNameSuffix   string
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
		BaseNameSuffix:   opt.BaseNameSuffix,
		Compress:         opt.Compress,
		Includes:         opt.Includes,
		TargetExtensions: opt.TargetExtensions,
	}
}

func (s Service) Check(file fileutil.File) ([]fileutil.File, error) {
	// 再起で存在確認をしないといけない
	//TODO このままだと PIXTA と shutterstockに対応できない
	var checked []fileutil.File
	var errs []error
	for _, ext := range s.TargetExtensions {
		if len(checked) >= s.Includes {
			break
		}
		fileName := file.Base().Suffix(s.BaseNameSuffix).FullName(ext)
		_, err := os.Stat(filepath.Join(file.WorkDir, string(fileName)))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		checked = append(checked, fileutil.File{Name: fileName, WorkDir: file.WorkDir})
	}
	if len(checked) < s.Includes {
		return checked, errors.Join(errs...)
	}
	return checked, nil
}

func (s Service) Exec(sources []fileutil.File) error {
	dstDir := filepath.Join(sources[0].WorkDir, string(s.Name))
	compressor := compressor.NewCompressor(s.Compress)
	// TODO dstPath にしたほうが良いかな
	if err := compressor.Compress(sources, fileutil.DirectoryName(dstDir)); err != nil {
		return err
	}
	return nil
}
