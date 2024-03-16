package service

import (
	"github.com/uoya/ImagePacker/fileutil"
	"os"
)

type ShutterStock struct{}

func (s ShutterStock) Name() Name {
	return "ShutterStock"
}

func (s ShutterStock) Check(baseName fileutil.FileBaseName) ([]fileutil.FileName, error) {
	// *_ss.eps が存在するかチェック
	checked := []fileutil.FileName{baseName.Suffix("_ss").FullName(fileutil.Eps)}
	_, err := os.Stat(string(checked[0]))
	if err != nil {
		return []fileutil.FileName{}, err
	}
	return checked, nil
}

func (s ShutterStock) Exec(sources []fileutil.FileName) error {
	for _, src := range sources {
		dstDir := fileutil.DirectoryName(s.Name())
		_, err := fileutil.CopyFile(src, dstDir)
		if err != nil {
			return err
		}
	}
	return nil
}
