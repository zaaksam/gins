package utils

import "io/ioutil"

// ReadFile 读取文件
func ReadFile(name string) (body []byte, err error) {
	body, err = ioutil.ReadFile(name)
	return
}
