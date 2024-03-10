package main

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Extension string

const (
	ai  Extension = ".ai"
	png Extension = ".png"
	jpg Extension = ".jpg"
	eps Extension = ".eps"
)

type FileName string
type FileBaseName string
type DirectoryName string

func (f FileName) Base() FileBaseName {
	strName := string(f)
	return FileBaseName(strings.TrimSuffix(strName, filepath.Ext(strName)))
}

func (b FileBaseName) Suffix(suffix string) FileBaseName {
	return FileBaseName(string(b) + suffix)
}
func (b FileBaseName) FullName(ext Extension) FileName {
	return FileName(string(b) + string(ext))
}
func MkdirIfNotExists(dirName DirectoryName) error {
	strDirName := string(dirName)
	if _, err := os.Stat(strDirName); os.IsNotExist(err) {
		if err = os.MkdirAll(strDirName, 0777); err != nil {
			return err
		}
		return err
	}
	return nil
}

func CopyFile(src FileName, dstDir DirectoryName) (int64, error) {
	strSrc := string(src)
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

	dst := path.Join(string(dstDir), strSrc)
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
