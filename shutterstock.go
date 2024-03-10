package main

import "os"

type ShutterStock struct{}

func (s ShutterStock) Name() ServiceName {
	return "ShutterStock"
}

func (s ShutterStock) Check(baseName FileBaseName) ([]FileName, error) {
	// *_ss.eps が存在するかチェック
	checked := []FileName{baseName.Suffix("_ss").FullName(eps)}
	_, err := os.Stat(string(checked[0]))
	if err != nil {
		return []FileName{}, err
	}
	return checked, nil
}

func (s ShutterStock) Exec(sources []FileName) error {
	for _, src := range sources {
		dstDir := DirectoryName(s.Name())
		_, err := CopyFile(src, dstDir)
		if err != nil {
			return err
		}
	}
	return nil
}
