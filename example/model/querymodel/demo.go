package querymodel

import "github.com/zaaksam/gins/example/model"

// Demo 参数接收结构体
type Demo struct {
	Paging
	model.Demo
}
