package utils

import (
	"path/filepath"
	"os"
	"strings"
	"io"
)

const TMPD = "/tmp/tmpgo_ins"
const DATE_DIR = "2006_01_02"
const DATE_FILE = "2006_01_02T150405"

/**
 * 获取当前文件执行的路径
 */
func GetCurrPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}

/**
 * 截取目录
 */
func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

/**
 * 获取上级目录
 */
func GetParentDirectory(dir string) string {
	return substr(dir, 0, strings.LastIndex(dir, "/"))
}

/**
 * 获取目录文件名
 */
func GetFileName(dir string) string {
	return substr(dir, strings.LastIndex(dir, "/") + 1, len(dir))
}

/**
 * 截断左右空格
 */
func TrimLeftRightSpace(str string) string {
	return strings.TrimRight(strings.TrimLeft(str, " "), " ")
}

/*
Copyfile 拷贝文件
 */
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

/**
 * 列出目录下所有文件, 不递归
 */
func GetAllFiles(dirs string) ([]os.FileInfo, error) {
	files := make([]os.FileInfo, 0)
	dir, err := os.Open(dirs)
	if err != nil {
		return files, err
	}
	fs, err := dir.Readdir(0)
	if err != nil {
		return files, err
	}
	for _, f := range fs {
		if !f.IsDir() {
			files = append(files, f)
		}
	}
	return files, nil
}


// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions perm;
// otherwise WriteFile truncates it before writing.
func WriteFileA(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_RDWR, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func MkdirConsolesPath(base string) error {
	start := base + "/start"
	stop := base + "/stop"
	restart := base + "/restart"

	var err error
	err = os.MkdirAll(start, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(stop, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(restart, 0755)
	if err != nil {
		return err
	}
	return nil
}

