package util

import (
	"os"
)

func GetFilesAndDirs(dirPath string) (files []string, dirs []string, err error) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}

	pathSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPath+pathSep+fi.Name())
			sonFiles, sonDirs, err := GetFilesAndDirs(dirPath + pathSep + fi.Name())
			if err != nil {
				return nil, nil, err
			}
			files = append(files, sonFiles...)
			dirs = append(dirs, sonDirs...)
		} else {
			files = append(files, dirPath+pathSep+fi.Name())
		}
	}

	return files, dirs, nil
}
