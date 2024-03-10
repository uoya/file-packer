package main

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func ZipFiles(sourceNames []FileName, dstDir DirectoryName) error {
	if len(sourceNames) == 0 {
		return nil
	}
	zipFile, err := os.Create(path.Join(string(dstDir), string(sourceNames[0].Base())+".zip"))
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range sourceNames {
		if err := addFileToZip(file, zipWriter); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(filename FileName, zipWriter *zip.Writer) error {
	fileToZip, err := os.Open(string(filename))
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = string(filename)

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
