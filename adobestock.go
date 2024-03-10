package main

import (
	"os"
	"path"
)

type AdobeStock struct{}

func (a AdobeStock) Name() string {
	return string(adobeStock)
}

func (a AdobeStock) Check(baseName string) ([]string, error) {
	ext := []Extension{eps, jpg}
	var targets []string
	for _, ext := range ext {
		t := baseName + string(ext)
		_, err := os.Stat(t)
		if err != nil {
			return []string{}, err
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func (a AdobeStock) Exec(targets []string, baseName string) error {
	if err := ZipFiles(targets, path.Join(a.Name(), baseName)); err != nil {
		return err
	}
	return nil
}
