package main

import "os"

type ShutterStock struct{}

func (s ShutterStock) Name() string {
	return string(shutterStock)
}

func (s ShutterStock) Check(baseName string) ([]string, error) {
	// *_ss.eps が存在するかチェック
	targets := []string{baseName + "_ss.eps"}
	_, err := os.Stat(targets[0])
	if err != nil {
		return []string{}, err
	}
	return targets, nil
}

func (s ShutterStock) Exec(targets []string, baseName string) error {
	for _, t := range targets {
		_, err := CopyFile(t, s.Name())
		if err != nil {
			return err
		}
	}
	return nil
}
