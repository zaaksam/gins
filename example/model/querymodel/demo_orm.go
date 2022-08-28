package querymodel

import "github.com/zaaksam/gins/extend/orm"

// NewDemo 创建数据对象
func NewDemo() *Demo {
	md := &Demo{}

	imd := &md.Model

	md.ID = orm.NewField[uint64](imd, "id", "id", true)
	md.User = orm.NewField[string](imd, "user", "user", false)
	md.Pswd = orm.NewField[string](imd, "pswd", "pswd", false)
	md.Status = orm.NewField[int](imd, "status", "status", false)
	md.Created = orm.NewField[int64](imd, "created", "created", false)
	md.Updated = orm.NewField[int64](imd, "updated", "updated", false)

	return md
}
