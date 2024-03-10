package main

import "os"

type ImageMart struct{}

func (i ImageMart) Name() string {
	return string(imageMart)
}

func (i ImageMart) Check(baseName string) ([]string, error) {
	ext := []Extension{eps, jpg}
	var targets []string
	for _, ext := range ext {
		t := baseName + string(ext)
		_, err := os.Stat(baseName + string(ext))
		if err != nil {
			return []string{}, err
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func (i ImageMart) Exec(targets []string, baseName string) error {
	for _, t := range targets {
		_, err := CopyFile(t, i.Name())
		if err != nil {
			return err
		}
	}
	return nil
}
