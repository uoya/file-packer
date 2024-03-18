package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/uoya/file-packer/fileutil"
	"os"
)

var defaultConfig = &Config{
	Version:         version,
	Root:            root,
	MarkerExtension: ".ai",
	Services: []ServiceOption{
		{
			Name:             "AdobeStock",
			TargetExtensions: []string{".eps", ".png"},
			Includes:         2,
			Compress:         "zip",
		},
		{
			Name:             "PIXTA",
			TargetExtensions: []string{".eps", ".png", ".jpg"},
			Includes:         2,
			Compress:         "zip",
		}, {
			Name:             "ShutterStock",
			TargetExtensions: []string{".eps"},
			BaseNameSuffix:   "_ss",
			Includes:         1,
			Compress:         "none",
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
	Version         string          `json:"version" validate:"required,semver"`
	Root            string          `json:"root" validate:"required,min=1"`
	MarkerExtension string          `json:"markerExtension" validate:"required"`
	Services        []ServiceOption `json:"services" validate:"gt=0,dive,required"`
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
