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

type Service interface {
	Name() string
	Check(string) ([]string, error)
	Exec([]string, string) error
}

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

	history := make(map[string]map[string][]string)
	for _, file := range files {
		filename := file.Name()
		ext := filepath.Ext(file.Name())

		// ai ファイルのみ抽出
		if !file.IsDir() && ext == string(ai) {
			baseName := filename[:len(filename)-len(ext)]

			// チェック
			var checked []string
			h := make(map[string][]string)
			h["original"] = []string{filename} // オリジナルの ai ファイル
			for _, service := range services {
				checked, err = service.Check(baseName)
				if err != nil {
					slog.Error(err.Error(), "ステップ", "check", "対象", service.Name(), "ファイル", baseName)
					return 1
				}
				h[service.Name()] = checked
			}
			history[baseName] = h

			// 実行
			for _, service := range services {
				err = MkdirIfNotExists(service.Name())
				err := service.Exec(history[baseName][service.Name()], baseName)
				if err != nil {
					slog.Error(err.Error(), "ステップ", "exec", "対象", service.Name(), "ファイル", baseName)
					return 1
				}
			}
		}
	}

	// 処理済みファイル格納フォルダ作成
	now := time.Now().Format("2006-01-02-15-04")
	nowDir := path.Join(outputPath, now)
	err = MkdirIfNotExists(nowDir)
	if err != nil {
		slog.Error(err.Error())
		return 1
	}

	// 実行済みファイルを移動
	// 実行済みファイルの重複を除去
	fileSet := NewSet()
	for _, h := range history {
		for _, v := range h {
			for _, vv := range v {
				fileSet.Add(vv)
			}
		}
	}
	for _, f := range fileSet.Values() {
		err = os.Rename(f, path.Join(nowDir, f))
		if err != nil {
			slog.Error(err.Error(), "ステップ", "move", "対象", f)
			return 1
		}
	}
	return 0
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
