package service

import (
	"github.com/uoya/file-packer/compressor"
	"github.com/uoya/file-packer/fileutil"
)

type Name string

type Service interface {
	Name() Name
	Check(fileutil.File) ([]fileutil.File, error)
	Exec([]fileutil.File) error
}

type Option struct {
	Name             Name                      `json:"name" validate:"required,min=1"`
	BaseNameSuffix   string                    `json:"baseNameSuffix,omitempty" validate:"omitempty"`
	TargetExtensions []string                  `json:"targetExtensions" validate:"min=1"`
	Includes         int                       `json:"includes" validate:"min=1"`
	Compress         compressor.CompressOption `json:"compress" validate:"required,min=1"`
}

func New(name Name) (Service, error) {
	switch name {
	case "AdobeStock":
		return AdobeStock{}, nil
	}
	return nil, nil
}
