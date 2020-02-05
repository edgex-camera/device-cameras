package utils

import "os"

func MakeDirsIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		return os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	return nil
}
