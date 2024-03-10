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
	root       = "./"
)

type ServiceName string
type Service interface {
	Name() ServiceName
	Check(FileBaseName) ([]FileName, error)
	Exec([]FileName) error
}

type History map[FileBaseName]map[ServiceName][]FileName

func main() {
	os.Exit(realMain())
}

func realMain() int {
	// ログ出力設定
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logfile), nil))
	slog.SetDefault(logger)

	services := []Service{
		AdobeStock{},
		Pixta{},
		ImageMart{},
		ShutterStock{},
	}

	history := make(History)
	files, err := os.ReadDir(root)
	for _, file := range files {
		filename := FileName(file.Name())

		// ai ファイルのみ抽出
		if !file.IsDir() && filepath.Ext(file.Name()) == string(ai) {
			baseName := filename.Base()

			// チェック
			h := make(map[ServiceName][]FileName)
			h["original"] = []FileName{filename} // オリジナルの ai ファイル
			for _, service := range services {
				checked, err := service.Check(baseName)
				if err != nil {
					slog.Error(err.Error(), "ステップ", "check", "対象", service.Name(), "ファイル", baseName)
					return 1
				}
				h[service.Name()] = checked
			}
			history[baseName] = h

			// 実行
			for _, service := range services {
				err = MkdirIfNotExists(DirectoryName(service.Name()))
				err := service.Exec(history[baseName][service.Name()])
				if err != nil {
					slog.Error(err.Error(), "ステップ", "exec", "対象", service.Name(), "ファイル", baseName)
					return 1
				}
			}
		}
	}

	// 処理済みファイル格納フォルダ作成
	now := time.Now().Format("2006-01-02-15-04")
	nowDir := DirectoryName(path.Join(outputPath, now))
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
				fileSet.Add(string(vv))
			}
		}
	}
	for _, f := range fileSet.Values() {
		err = os.Rename(f, path.Join(string(nowDir), f))
		if err != nil {
			slog.Error(err.Error(), "ステップ", "move", "対象", f)
			return 1
		}
	}
	slog.Info("処理が完了しました")
	return 0
}
