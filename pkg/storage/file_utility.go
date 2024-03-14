package storage

import (
	"os"
	"path/filepath"
)

func CheckIfFileOrFolderExistsAndNotEmpty(path string) (bool, error) {

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if fileInfo.Size() <= 0 {
		return false, nil
	}

	return true, nil
}

func CreateFile(path string) (*os.File, error) {

	err := os.MkdirAll(filepath.Dir(path), 0755)
	//fmt.Println("test===", filepath.Dir(path), " end of path")
	if err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}
