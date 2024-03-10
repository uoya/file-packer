package main

import "os"

type ImageMart struct{}

func (i ImageMart) Name() ServiceName {
	return "イメージマート"
}

func (i ImageMart) Check(baseName FileBaseName) ([]FileName, error) {
	ext := []Extension{eps, jpg}
	var checked []FileName
	for _, ext := range ext {
		t := baseName.FullName(ext)
		_, err := os.Stat(string(baseName.FullName(ext)))
		if err != nil {
			return []FileName{}, err
		}
		checked = append(checked, t)
	}
	return checked, nil
}

func (i ImageMart) Exec(sources []FileName) error {
	for _, src := range sources {
		dstDir := DirectoryName(i.Name())
		_, err := CopyFile(src, dstDir)
		if err != nil {
			return err
		}
	}
	return nil
}
