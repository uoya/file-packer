package compressor

import (
	"encoding/json"
	"fmt"
	"github.com/uoya/file-packer/fileutil"
)

type Compressor interface {
	Compress([]fileutil.File, fileutil.DirectoryName) error
}

func NewCompressor(opt CompressOption) Compressor {
	switch opt {
	case CompressZip:
		return &ZipCompressor{}
	default:
		return nil
	}
}

type CompressOption string

const (
	CompressNone CompressOption = "none"
	CompressZip  CompressOption = "zip"
)

var validCompressors = []CompressOption{
	CompressNone,
	CompressZip,
}

func (c *CompressOption) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	for _, validCompressor := range validCompressors {
		if CompressOption(str) == validCompressor {
			*c = validCompressor
			return nil
		}
	}

	return fmt.Errorf("invalid compressor value: %s", str)
}
