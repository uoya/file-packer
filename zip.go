package main

import (
	"archive/zip"
	"io"
	"os"
)

func ZipFiles(sourceNames []string, baseName string) error {
	zipFile, err := os.Create(baseName + ".zip")
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

func addFileToZip(filename string, zipWriter *zip.Writer) error {
	fileToZip, err := os.Open(filename)
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

	header.Name = filename

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
