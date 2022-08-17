package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaaksam/gins/errors"
	"github.com/zaaksam/gins/example/model"
	"github.com/zaaksam/gins/example/model/qry"
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

	has, err := da.Where(md.GetIDMarkKey()+"=?", id).Get(md)
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
	if data.User == "" {
		msg = "User不能为空"
	} else if data.Pswd == "" {
		msg = "Pswd不能为空"
	} else if !data.IsStatusMark() {
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
	md.SetID(snowflake.NewID())
	md.SetUser(data.User)
	md.SetPswd(data.Pswd)
	md.SetStatus(data.Status)
	md.SetCreated(ux)
	md.SetUpdated(ux)

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
	md.SetUpdated(time.Now().Unix())

	da.Where(md.GetIDMarkKey()+"=?", md.ID)
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

	da.Where(md.GetIDMarkKey()+"=?", id)
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
func (*demo) List(ctx *gin.Context, qry *qry.Demo) (ml *orm.ModelList[*model.Demo], err error) {
	md := model.NewDemo()
	da, err := orm.NewDA(md)
	if err != nil {
		err = errors.NewAPIErrorWrap(err)
		return
	}
	defer da.Close()

	if qry.IsUserMark() {
		da.Where(fmt.Sprintf("%s=?", md.GetUserMarkKey()), qry.User)
	}

	if qry.IsStatusMark() {
		da.Where(fmt.Sprintf("%s=?", md.GetStatusMarkKey()), qry.Status)
	}

	if qry.Page <= 0 {
		qry.SetPage(1)
	}

	if qry.PageSize <= 0 {
		qry.SetPageSize(10)
	}

	da.Desc(md.GetIDMarkKey())
	ml, err = da.GetModelList(qry.Page, qry.PageSize)

	if err != nil {
		err = errors.NewAPIErrorWrap(err)
	}

	return
}
