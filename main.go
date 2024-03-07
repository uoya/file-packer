package main

import (
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	outputPath = "./処理済み"
	logFile    = "./process.log"
)

type ServiceName string

const (
	adobeStock   ServiceName = "AdobeStock"
	pixta        ServiceName = "PIXTA"
	imageMart    ServiceName = "イメージマート"
	shutterStock ServiceName = "ShutterStock"
)

type Extension string

const (
	ai  Extension = ".ai"
	png Extension = ".png"
	jpg Extension = ".jpg"
	eps Extension = ".eps"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logfile), nil))
	slog.SetDefault(logger)

	// 配置先フォルダ作成
	services := []Service{
		AdobeStock{},
		Pixta{},
		ImageMart{},
		ShutterStock{},
	}

	files, err := os.ReadDir("./")
	for _, file := range files {
		filename := file.Name()
		ext := filepath.Ext(file.Name())

		// ai ファイルのみ抽出
		if !file.IsDir() && ext == string(ai) {
			baseName := filename[:len(filename)-len(ext)]

			// チェック
			for _, service := range services {
				if err := service.Check(baseName); err != nil {
					slog.Error(err.Error(), "ステップ", "check", "対象", service.Name(), "ファイル", baseName)
					return 1
				}
			}

			// チェック通過したためフォルダ作成
			for _, service := range services {
				err = MkdirIfNotExists(service.Name())
				if err != nil {
					slog.Error("%s", err.Error())
					return 1
				}
			}
			// 実行
			for _, service := range services {
				err := service.Exec(baseName)
				if err != nil {
					slog.Error(err.Error(), "ステップ", "exec", "対象", service.Name(), "ファイル", baseName)
					return 1
				}
			}

			// 実行済みファイルを移動
		}
	}

	// 処理済みファイル格納フォルダ作成
	now := time.Now().Format("2006-01-02-15-04")
	nowDir := path.Join(outputPath, now)
	err = MkdirIfNotExists(nowDir)
	if err != nil {
		slog.Error("%s", err.Error())
		return 1
	}

	return 0
}

func zip(files []string, output string) (int, error) {
	return 0, nil
}

func MkdirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, 0777); err != nil {
			return err
		}
		return err
	}
	return nil
}

func CopyFile(src string, dstDir string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, nil
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dst := path.Join(dstDir, src)
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

type Service interface {
	Name() string
	Check(string) error
	Exec(string) error
}

type AdobeStock struct{}

func (a AdobeStock) Name() string {
	return string(adobeStock)
}

func (a AdobeStock) Check(baseName string) error {
	targets := []string{baseName + string(eps), baseName + string(jpg)}
	for _, target := range targets {
		_, err := os.Stat(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a AdobeStock) Exec(baseName string) error {
	extensions := []Extension{eps, jpg}
	for _, ext := range extensions {
		_, err := CopyFile(baseName+string(ext), a.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

type Pixta struct{}

func (p Pixta) Name() string {
	return string(pixta)
}

func (p Pixta) Check(baseName string) error {
	// eps 必須
	_, err := os.Stat(baseName + string(eps))
	if err != nil {
		return err
	}
	// png 優先, 存在しなければ jpg にフォールバック
	_, err = os.Stat(baseName + string(png))
	if err != nil {
		_, err := os.Stat(baseName + string(jpg))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Pixta) Exec(baseName string) error {
	// eps 必須
	_, err := os.Stat(baseName + string(eps))
	if err != nil {
		return err
	}
	_, err = CopyFile(baseName+string(eps), p.Name())
	if err != nil {
		return err
	}
	// png 優先, 存在しなければ jpg にフォールバック
	_, err = os.Stat(baseName + string(png))
	if err == nil {
		_, err = CopyFile(baseName+string(png), p.Name())
		if err != nil {
			return err
		}
	} else {
		_, err := os.Stat(baseName + string(jpg))
		if err != nil {
			return err
		}
		_, err = CopyFile(baseName+string(jpg), p.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

type ImageMart struct{}

func (i ImageMart) Name() string {
	return string(imageMart)
}

func (i ImageMart) Check(baseName string) error {
	extensions := []Extension{eps, jpg}
	for _, ext := range extensions {
		_, err := os.Stat(baseName + string(ext))
		if err != nil {
			return err
		}
	}
	return nil
}

func (i ImageMart) Exec(baseName string) error {
	extensions := []Extension{eps, jpg}
	for _, ext := range extensions {
		_, err := CopyFile(baseName+string(ext), i.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

type ShutterStock struct{}

func (s ShutterStock) Name() string {
	return string(shutterStock)
}

func (s ShutterStock) Check(baseName string) error {
	// *_ss.eps が存在するかチェック
	_, err := os.Stat(baseName + "_ss" + string(eps))
	if err != nil {
		return err
	}
	return nil
}

func (s ShutterStock) Exec(baseName string) error {
	_, err := CopyFile(baseName+"_ss"+string(eps), s.Name())
	if err != nil {
		return err
	}
	return nil
}