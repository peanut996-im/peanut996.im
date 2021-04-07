package file

import (
	"os"
	"path/filepath"
)

//GetAbsPath return a path string
func GetAbsPath(path string) (absFilePath string) {
	pwd, _ := os.Getwd()
	absFilePath = filepath.Join(pwd, path)
	return absFilePath
}
