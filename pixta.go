package main

import (
	"os"
)

type Pixta struct{}

func (p Pixta) Name() ServiceName {
	return "PIXTA"
}

func (p Pixta) Check(baseName FileBaseName) ([]FileName, error) {
	// eps必須もう1ファイルはpngがあればそちら、なければjpg
	var checked []FileName

	// eps 必須
	_, err := os.Stat(string(baseName.FullName(eps)))
	if err != nil {
		return []FileName{}, err
	}
	checked = append(checked, baseName.FullName(eps))

	// png 優先, 存在しなければ jpg にフォールバック
	_, err = os.Stat(string(baseName.FullName(png)))
	if err == nil {
		checked = append(checked, baseName.FullName(png))
	} else {
		_, err := os.Stat(string(baseName.FullName(jpg)))
		if err != nil {
			return []FileName{}, err
		}
		checked = append(checked, baseName.FullName(jpg))
	}

	return checked, nil
}

func (p Pixta) Exec(sources []FileName) error {
	dstDir := DirectoryName(p.Name())
	if err := ZipFiles(sources, dstDir); err != nil {
		return err
	}
	return nil
}
