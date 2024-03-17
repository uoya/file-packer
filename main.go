package main

import (
	_ "embed"
	"errors"
	"github.com/uoya/file-packer/fileutil"
	"github.com/uoya/file-packer/service"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	version    = "0.1.0"
	outputPath = "処理済み"
	logFile    = "process.log"
	root       = "./work"
	configFile = "./config.json"
)

type ErrTitle string
type ErrMsg string

const (
	titleFinished ErrTitle = "処理完了"
	titleError    ErrTitle = "エラー"
	msgFinished   ErrMsg   = "処理が完了しました"
)

type History map[fileutil.FileBaseName]map[service.Name][]fileutil.File

func main() {
	if err := realMain(); err != nil {
		MessageBox(titleError, ErrMsg(err.Error()))
		os.Exit(1)
	}
	os.Exit(0)
}

func realMain() error {
	// ログ出力設定
	logPath := path.Join(root, logFile)
	logfile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logfile), &slog.HandlerOptions{
		AddSource: true,
	}))
	child := logger.With(slog.String("version", version))
	slog.SetDefault(child)

	conf, err := loadConf()
	if err != nil {
		slog.Error(err.Error(), "ステップ", "loadConf")
		return err
	}

	var services []service.Service
	for _, s := range conf.Services {
		srv, err := service.New(s.Name)
		if err != nil {
			slog.Error(err.Error(), "ステップ", "loadConf")
			return err
		}
		services = append(services, srv)
	}

	history := make(History)
	files, err := os.ReadDir(root)

	// 処理対象取得
	var markerFiles []fileutil.File
	for _, file := range files {

		if !file.IsDir() && filepath.Ext(file.Name()) == conf.MarkerExtension {
			f := fileutil.File{Name: fileutil.FileName(file.Name()), Root: root}
			markerFiles = append(markerFiles, f)

		}
	}

	// チェック
	for _, f := range markerFiles {

		h := make(map[service.Name][]fileutil.File)
		for _, service := range services {
			// 出力先フォルダがすでに存在している場合、フォルダ内のデータを確認
			serviceDir := path.Join(root, string(service.Name()))
			if _, err = os.Stat(serviceDir); err == nil {
				// フォルダが存在しているので中のファイルを確認
				items, err := os.ReadDir(serviceDir)
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
			checked, err := service.Check(f)
			if err != nil {
				slog.Error(err.Error(), "ステップ", "check", "対象", service.Name(), "ファイル", f.Base())
				return err
			}
			h[service.Name()] = checked
		}
		history[f.Base()] = h
	}

	// 実行
	for k, _ := range history {
		for _, service := range services {
			err = fileutil.MkdirIfNotExists(fileutil.DirectoryName(path.Join(root, string(service.Name()))))
			err := service.Exec(history[k][service.Name()])
			if err != nil {
				slog.Error(err.Error(), "ステップ", "exec", "対象", service.Name(), "ファイル", k)
				return err
			}
		}
	}

	// 処理済みファイル格納フォルダ作成
	now := time.Now().Format("2006-01-02-15-04")
	nowDir := fileutil.DirectoryName(path.Join(root, outputPath, now))
	err = fileutil.MkdirIfNotExists(nowDir)
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
				fileSet.Add(string(vv.Name))
			}
		}
	}
	for _, f := range fileSet.Values() {
		err = os.Rename(path.Join(root, f), path.Join(string(nowDir), f))
		if err != nil {
			slog.Error(err.Error(), "ステップ", "move", "対象", f)
			return err
		}
	}

	slog.Info(string(msgFinished))
	MessageBox(titleFinished, msgFinished)
	return nil
}
