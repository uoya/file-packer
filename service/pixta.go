package service

import (
	"github.com/uoya/ImagePacker/fileutil"
	"os"
)

type Pixta struct{}

func (p Pixta) Name() Name {
	return "PIXTA"
}

func (p Pixta) Check(baseName fileutil.FileBaseName) ([]fileutil.FileName, error) {
	// eps必須もう1ファイルはpngがあればそちら、なければjpg
	var checked []fileutil.FileName

	// eps 必須
	_, err := os.Stat(string(baseName.FullName(fileutil.Eps)))
	if err != nil {
		return []fileutil.FileName{}, err
	}
	checked = append(checked, baseName.FullName(fileutil.Eps))

	// png 優先, 存在しなければ jpg にフォールバック
	_, err = os.Stat(string(baseName.FullName(fileutil.Png)))
	if err == nil {
		checked = append(checked, baseName.FullName(fileutil.Png))
	} else {
		_, err := os.Stat(string(baseName.FullName(fileutil.Jpg)))
		if err != nil {
			return []fileutil.FileName{}, err
		}
		checked = append(checked, baseName.FullName(fileutil.Jpg))
	}

	return checked, nil
}

func (p Pixta) Exec(sources []fileutil.FileName) error {
	dstDir := fileutil.DirectoryName(p.Name())
	if err := fileutil.ZipFiles(sources, dstDir); err != nil {
		return err
	}
	return nil
}
