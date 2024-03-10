package main

import (
	"os"
	"path"
)

type Pixta struct{}

func (p Pixta) Name() string {
	return string(pixta)
}

func (p Pixta) Check(baseName string) ([]string, error) {
	// eps必須もう1ファイルはpngがあればそちら、なければjpg
	var targets []string

	// eps 必須
	_, err := os.Stat(baseName + string(eps))
	if err != nil {
		return []string{}, err
	}
	targets = append(targets, baseName+string(eps))

	// png 優先, 存在しなければ jpg にフォールバック
	_, err = os.Stat(baseName + string(png))
	if err == nil {
		targets = append(targets, baseName+string(png))
	} else {
		_, err := os.Stat(baseName + string(jpg))
		if err != nil {
			return []string{}, err
		}
		targets = append(targets, baseName+string(jpg))
	}

	return targets, nil
}

func (p Pixta) Exec(targets []string, baseName string) error {
	var filenames []string
	// eps 必須
	filenames = append(filenames, baseName+string(eps))

	// png 優先, 存在しなければ jpg にフォールバック
	pngFile := baseName + string(png)
	_, err := os.Stat(pngFile)
	if err == nil {
		if err != nil {
			return err
		}
		filenames = append(filenames, pngFile)
	} else {
		filenames = append(filenames, baseName+string(jpg))
	}
	if err := ZipFiles(filenames, path.Join(p.Name(), baseName)); err != nil {
		return err
	}
	return nil
}
