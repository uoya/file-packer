package fileutil

import "os"

type DirectoryName string

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
