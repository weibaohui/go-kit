package filekit

import (
	"github.com/weibaohui/go-kit/strkit"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ParentDirectory(dirctory string) string {
	return strkit.Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}
func Pwd() string {
	return CurrentPath()
}
func CurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
