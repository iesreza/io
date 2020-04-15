package gpath

import (
	"github.com/iesreza/io/lib/text"
	copy2 "github.com/otiai10/copy"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func MakePath(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func Parent(s string) string {
	list := text.SplitAny(s,"\\/")
	return strings.Join(list[0:len(list)-1],"/")
}


func WorkingDir() string {
	path, _ := os.Getwd()
	return path
}

func RSlash(path string) string {

	return strings.TrimRight(strings.TrimSpace(path), "/")
}

func IsDirExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func IsFileExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsDirEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}

func Stat(path string) *os.FileInfo {
	fileStat, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return &fileStat
}

type internalPathInfo struct {
	FileName  string
	Path      string
	Extension string
}

func PathInfo(path string) internalPathInfo {
	info := internalPathInfo{
		FileName:  filepath.Base(path),
		Path:      filepath.Dir(path),
		Extension: filepath.Ext(path),
	}
	return info
}

func CopyDir(src, dest string) error {
	return copy2.Copy(src, dest)
}

func CopyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func SymLink(src, dest string) error {
	return os.Link(src, dest)
}

func Remove(path string) error {
	if IsDir(path) {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func SafeFileContent(path string) []byte {
	data, _ := ioutil.ReadFile(path)
	return data
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
