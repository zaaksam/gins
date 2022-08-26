package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaaksam/gins/errors"
	"github.com/zaaksam/gins/example/model"
	"github.com/zaaksam/gins/example/model/querymodel"
	"github.com/zaaksam/gins/extend/orm"
	"github.com/zaaksam/gins/extend/snowflake"
)

// Demo 服务
var Demo demo

type demo struct{}

// MustGet 必须获取详情
func (*demo) MustGet(ctx *gin.Context, id uint64) (md *model.Demo, err error) {
	msg := ""
	if id == 0 {
		msg = "id不能为空"
	}
	if msg != "" {
		err = errors.NewAPIError(msg)
		return
	}

	md, err = Demo.mustGet(nil, id)
	return
}

func (*demo) mustGet(da *orm.DA[*model.Demo], id uint64) (md *model.Demo, err error) {
	md, err = Demo.get(da, id)
	if err != nil {
		return
	}

	if md == nil {
		err = errors.NewAPIError("样例不存在")
	}
	return
}

// Get 获取详情
func (*demo) Get(ctx *gin.Context, id uint64) (md *model.Demo, err error) {
	msg := ""
	if id == 0 {
		msg = "id不能为空"
	}
	if msg != "" {
		err = errors.NewAPIError(msg)
		return
	}

	md, err = Demo.get(nil, id)
	return
}

// Get 获取
func (*demo) get(da *orm.DA[*model.Demo], id uint64) (md *model.Demo, err error) {
	md = model.NewDemo()
	if da == nil {
		da, err = orm.NewDA(md)
		if err != nil {
			err = errors.NewAPIErrorWrap(err)
			return
		}
		defer da.Close()
	}

	has, err := da.Where(md.ID.FieldName()+"=?", id).Get(md)
	if err != nil {
		err = errors.NewAPIErrorWrap(err)
		return
	}
	if !has {
		md = nil
	}
	return
}

// Create 创建
func (*demo) Create(ctx *gin.Context, data *model.Demo) (md *model.Demo, err error) {
	msg := ""
	if data.User.Val() == "" {
		msg = "User不能为空"
	} else if data.Pswd.Val() == "" {
		msg = "Pswd不能为空"
	} else if !data.Status.IsSetVal() {
		msg = "Status不能为空"
	}
	if msg != "" {
		err = errors.NewAPIError(msg)
		return
	}

	md = model.NewDemo()
	da, err := orm.NewDA(md)
	if err != nil {
		err = errors.NewAPIErrorWrap(err)
		return
	}
	defer da.Close()

	md, err = Demo.create(da, data)
	return
}

func (*demo) create(da *orm.DA[*model.Demo], data *model.Demo) (md *model.Demo, err error) {
	md = model.NewDemo()
	if da == nil {
		da, err = orm.NewDA(md)
		if err != nil {
			err = errors.NewAPIErrorWrap(err)
			return
		}
		defer da.Close()
	}

	ux := time.Now().Unix()
	md.ID.SetVal(snowflake.NewID())
	md.User.SetVal(data.User.Val())
	md.Pswd.SetVal(data.Pswd.Val())
	md.Status.SetVal(data.Status.Val())
	md.Created.SetVal(ux)
	md.Updated.SetVal(ux)

	rows, err := da.Insert(md)
	if err != nil {
		err = errors.NewAPIErrorWrap(err)
		return
	}

	if rows <= 0 {
		err = errors.NewAPIError("创建样例记录失败")
	}
	return
}

// Update 更新
func (*demo) Update(da *orm.DA[*model.Demo], md *model.Demo) (err error) {
	if da == nil {
		da, err = orm.NewDA(md)
		if err != nil {
			err = errors.NewAPIErrorWrap(err)
			return
		}
		defer da.Close()
	}

	var rows int64
	md.Updated.SetVal(time.Now().Unix())

	da.Where(md.ID.FieldName()+"=?", md.ID)
	rows, err = da.UpdateByModel(md)
	if err != nil {
		err = errors.NewAPIErrorWrap(err)
		return
	}
	if rows <= 0 {
		err = errors.NewAPIError("更新样例时失败")
	}
	return
}

// Del 删除
func (*demo) Del(da *orm.DA[*model.Demo], id uint64) (err error) {
	md := model.NewDemo()

	if da == nil {
		da, err = orm.NewDA(md)
		if err != nil {
			err = errors.NewAPIErrorWrap(err)
			return
		}
		defer da.Close()
	}

	var rows int64

	da.Where(md.ID.FieldName()+"=?", id)
	rows, err = da.Delete(md)
	if err != nil {
		err = errors.NewAPIErrorWrap(err)
		return
	}
	if rows <= 0 {
		err = errors.NewAPIError("删除样例时失败")
	}
	return
}

// List 列表
func (*demo) List(ctx *gin.Context, qry *querymodel.Demo) (ml *orm.ModelList[*model.Demo], err error) {
	md := model.NewDemo()
	da, err := orm.NewDA(md)
	if err != nil {
		err = errors.NewAPIErrorWrap(err)
		return
	}
	defer da.Close()

	if qry.User.IsSetVal() {
		da.And(fmt.Sprintf("%s=?", md.User.FieldName()), qry.User)
	}

	if qry.Status.IsSetVal() {
		da.And(fmt.Sprintf("%s=?", md.Status.FieldName()), qry.Status)
	}

	if qry.Page <= 0 {
		qry.Page = 1
	}

	if qry.PageSize <= 0 {
		qry.PageSize = 10
	}

	da.Desc(md.ID.FieldName())
	ml, err = da.GetModelList(qry.Page, qry.PageSize)

	if err != nil {
		err = errors.NewAPIErrorWrap(err)
	}

	return
}
