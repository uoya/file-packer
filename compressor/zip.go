package compressor

import (
	"archive/zip"
	"github.com/uoya/file-packer/fileutil"
	"io"
	"os"
	"path"
)

type ZipCompressor struct{}

func (z *ZipCompressor) Compress(sources []fileutil.File, dstDir fileutil.DirectoryName) error {
	if len(sources) == 0 {
		return nil
	}
	zipFile, err := os.Create(path.Join(string(dstDir), string(sources[0].Base())+".zip"))
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, f := range sources {
		if err := addFileToZip(f, zipWriter); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(file fileutil.File, zipWriter *zip.Writer) error {
	fileToZip, err := os.Open(file.StrPath())
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

	header.Name = string(file.Name)

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
