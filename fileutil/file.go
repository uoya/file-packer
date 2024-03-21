package fileutil

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type FilePath string
type FileName string
type FileBaseName string

type File struct {
	Name    FileName
	WorkDir string
}

func (f File) Base() FileBaseName {
	strName := string(f.Name)
	return FileBaseName(strings.TrimSuffix(strName, filepath.Ext(strName)))
}

func (b FileBaseName) Suffix(suffix string) FileBaseName {
	return FileBaseName(string(b) + suffix)
}
func (b FileBaseName) FullName(ext string) FileName {
	return FileName(string(b) + string(ext))
}

func (f File) Path() FilePath {
	strName := string(f.Name)
	return FilePath(filepath.Join(f.WorkDir, strName))
}

func (f File) StrPath() string {
	return string(f.Path())
}

func CopyFile(srcPath FilePath, dstDir DirectoryName) (int64, error) {
	strSrc := string(srcPath)
	sourceFileStat, err := os.Stat(strSrc)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, nil
	}

	source, err := os.Open(strSrc)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dst := filepath.Join(string(dstDir), filepath.Base(strSrc))
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CreateFileIfNotExists(path FilePath, content string) error {
	srcPath := string(path)
	_, err := os.Stat(srcPath)
	if os.IsNotExist(err) {
		file, err := os.Create(srcPath)
		slog.Info("file created.", "path", path)
		if err != nil {
			return err
		}
		if _, err := file.WriteString(content); err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}
