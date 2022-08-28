package fileutil

import (
	"io/ioutil"
	"os"
	"strings"
)

// MkdirAll 创建所有层级目录
func MkdirAll(path string) (err error) {
	err = os.MkdirAll(path, os.ModePerm)
	return
}

// ReadFile 读取文件
func ReadFile(name string) (body []byte, err error) {
	body, err = ioutil.ReadFile(name)
	return
}

// NewFile 创建 *os.File
func NewFile(name string, isAppend ...bool) (file *os.File, err error) {
	flag := os.O_RDWR | os.O_CREATE

	if len(isAppend) == 1 && isAppend[0] {
		flag = flag | os.O_APPEND
	}

	dirs := strings.Split(name, string(os.PathSeparator))
	dir := strings.Join(dirs[:len(dirs)-1], string(os.PathSeparator))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}

	file, err = os.OpenFile(name, flag, os.ModePerm)
	// file, err = os.OpenFile(name, flag, 0644)
	if err != nil {
		return
	}

	if len(isAppend) == 0 || !isAppend[0] {
		err = file.Truncate(0)
	}

	return
}

// WriteFile 写入文件
func WriteFile(name string, body []byte, isAppend ...bool) (err error) {
	file, err := NewFile(name, isAppend...)
	if err != nil {
		return
	}

	_, err = file.Write(body)
	return
}

// RemoveFile 删除文件
func RemoveFile(name string) (err error) {
	err = os.Remove(name)
	return
}
