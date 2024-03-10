package main

import (
	"os"
)

type AdobeStock struct{}

func (a AdobeStock) Name() ServiceName {
	return "AdobeStock"
}
func (a AdobeStock) Check(baseName FileBaseName) ([]FileName, error) {
	ext := []Extension{eps, jpg}
	var checked []FileName
	for _, ext := range ext {
		t := baseName.FullName(ext)
		_, err := os.Stat(string(t))
		if err != nil {
			return []FileName{}, err
		}
		checked = append(checked, t)
	}
	return checked, nil
}

func (a AdobeStock) Exec(sources []FileName) error {
	dstDir := DirectoryName(a.Name())
	if err := ZipFiles(sources, dstDir); err != nil {
		return err
	}
	return nil
}
