package model

import "github.com/zaaksam/gins/extend/orm"

func init() {
	// 注册 db
	orm.Register(&Demo{})
}

// NewDemo 创建数据对象
func NewDemo() *Demo {
	md := &Demo{}

	imd := &md.Model

	md.ID = orm.NewField[uint64](imd, "id")
	md.User = orm.NewField[string](imd, "user")
	md.Pswd = orm.NewField[string](imd, "pswd")
	md.Status = orm.NewField[int](imd, "status")
	md.Created = orm.NewField[int64](imd, "created")
	md.Updated = orm.NewField[int64](imd, "updated")

	return md
}
