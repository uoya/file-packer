package main

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	version    = "0.1.0"
	outputPath = "./処理済み"
	logFile    = "./process.log"
	root       = "./"
)

const (
	titleFinished = "処理完了"
	msgFinished   = "処理が完了しました"
	titleError    = "エラー"
)

type ServiceName string
type Service interface {
	Name() ServiceName
	Check(FileBaseName) ([]FileName, error)
	Exec([]FileName) error
}

type History map[FileBaseName]map[ServiceName][]FileName

func main() {
	if err := realMain(); err != nil {
		MessageBox(titleError, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func realMain() error {
	// ログ出力設定
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logfile), &slog.HandlerOptions{
		AddSource: true,
	}))
	child := logger.With(slog.String("version", version))
	slog.SetDefault(child)

	services := []Service{
		AdobeStock{},
		Pixta{},
		ImageMart{},
		ShutterStock{},
	}

	history := make(History)
	files, err := os.ReadDir(root)

	// チェック
	for _, file := range files {
		filename := FileName(file.Name())

		// ai ファイルのみ抽出
		if !file.IsDir() && filepath.Ext(file.Name()) == string(ai) {
			baseName := filename.Base()

			h := make(map[ServiceName][]FileName)
			h["original"] = []FileName{filename} // オリジナルの ai ファイル
			for _, service := range services {
				// 出力先フォルダがすでに存在している場合、フォルダ内のデータを確認
				if _, err = os.Stat(string(service.Name())); err == nil {
					// フォルダが存在しているので中のファイルを確認
					items, err := os.ReadDir(string(service.Name()))
					if err != nil {
						slog.Error(err.Error(), "ステップ", "check", "対象", service.Name())
						return err
					}
					if len(items) > 0 {
						message := "出力先フォルダ内にファイルが存在します"
						slog.Error(message, "対象", service.Name())
						return errors.New(message)
					}
				} else if !os.IsNotExist(err) {
					// フォルダが存在しない、以外のエラーは異常として扱う
					slog.Error(err.Error(), "ステップ", "check", "対象", service.Name())
					return err
				}

				// 必要なファイルの存在確認
				checked, err := service.Check(baseName)
				if err != nil {
					slog.Error(err.Error(), "ステップ", "check", "対象", service.Name(), "ファイル", baseName)
					return err
				}
				h[service.Name()] = checked
			}
			history[baseName] = h
		}
	}

	// 実行
	for k, _ := range history {
		for _, service := range services {
			err = MkdirIfNotExists(DirectoryName(service.Name()))
			err := service.Exec(history[k][service.Name()])
			if err != nil {
				slog.Error(err.Error(), "ステップ", "exec", "対象", service.Name(), "ファイル", k)
				return err
			}
		}
	}

	// 処理済みファイル格納フォルダ作成
	now := time.Now().Format("2006-01-02-15-04")
	nowDir := DirectoryName(path.Join(outputPath, now))
	err = MkdirIfNotExists(nowDir)
	if err != nil {
		slog.Error(err.Error())
		return err
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
			return err
		}
	}

	slog.Info(msgFinished)
	MessageBox(titleFinished, msgFinished)
	return nil
}
