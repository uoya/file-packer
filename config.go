package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/uoya/file-packer/fileutil"
	"os"
)

const (
	version    = "0.1.0"
	outputPath = "処理済み"
	logFile    = "process.log"
	configFile = "config.json"
)

var defaultConfig = &Config{
	Version:         "0.1.0",
	WorkDir:         "work",
	OutputPath:      "処理済み",
	MarkerExtension: ".ai",
	Services: []ServiceOption{
		{
			Name:             "AdobeStock･ShutterStock",
			TargetExtensions: []string{".eps"},
			Includes:         1,
			BaseNameSuffix:   "_lg",
			Compress:         "none",
		},
		{
			Name:             "PIXTA",
			TargetExtensions: []string{".eps", ".png", ".jpg"},
			Includes:         2,
			Compress:         "zip",
		}, {
			Name:             "イメージマート",
			TargetExtensions: []string{".eps", ".jpg"},
			Includes:         2,
			Compress:         "none",
		},
	},
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type Config struct {
	Version         string          `json:"version" validate:"required,semver"` // config.jsonの版を識別するためのバージョン番号
	WorkDir         string          `json:"workDir" validate:"required,min=1"` // 作業用フォルダ名
	OutputPath      string          `json:"outputPath" validate:"required,min=1"` // 処理済みファイルの移動先フォルダ名
	MarkerExtension string          `json:"markerExtension" validate:"required"` // 処理の基準とするファイルの拡張子
	Services        []ServiceOption `json:"services" validate:"gt=0,dive,required"` // アップロード先のサービスに応じたパラメータ
}

func loadConf() (*Config, error) {

	validate = validator.New(validator.WithRequiredStructEnabled())

	content, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err := fileutil.CreateFileIfNotExists(configFile, string(content)); err != nil {
		return nil, err
	}
	var conf Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
