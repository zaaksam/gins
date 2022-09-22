package querymodel

import (
	"github.com/zaaksam/gins/example/model"
	"github.com/zaaksam/gins/extend/orm"
)

// NewDemo 创建数据对象
func NewDemo() *Demo {
	md := &Demo{
		Demo: model.Demo{
			ID: &orm.Field[uint64]{},

			User: &orm.Field[string]{},

			Pswd: &orm.Field[string]{},

			Status: &orm.Field[int]{},

			Created: &orm.Field[int64]{},

			Updated: &orm.Field[int64]{},
		},
	}

	md.FieldReset()

	return md
}

// FieldReset 重设字段配置
func (md *Demo) FieldReset() {
	if md == nil {
		return
	}

	imd := &md.Model

	md.ID.Reset(imd, "id", "id", true)

	md.User.Reset(imd, "user", "user", false)

	md.Pswd.Reset(imd, "pswd", "pswd", false)

	md.Status.Reset(imd, "status", "status", false)

	md.Created.Reset(imd, "cteated", "created", false)

	md.Updated.Reset(imd, "updated", "updated", false)

}
