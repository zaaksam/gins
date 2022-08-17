package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zaaksam/gins"
	"github.com/zaaksam/gins/example/model/qry"
	"github.com/zaaksam/gins/example/service"
)

// Demo 接口
var Demo demo

type demo struct{}

func init() {
	gins.AddInit(func(gs *gins.Server) {
		router := gs.Engine().Group("ops/tag")

		router.POST("list", Demo.List)
	})
}

// Get 详情
func (*demo) Get(ctx *gin.Context) {
	res := gins.GetAPIResponse(ctx)

	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		res.SetMsg("id格式错误")
		return
	}

	md, err := service.Demo.MustGet(ctx, id)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetData(md)
}

// Create 创建
func (*demo) Create(ctx *gin.Context) {
	res := gins.GetAPIResponse(ctx)

	qry := qry.NewDemo()
	err := parseModel(ctx, qry)
	if err != nil {
		res.SetError(err)
		return
	}

	md, err := service.Demo.Create(ctx, &qry.Demo)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetData(md)
}

// List 列表
func (*demo) List(ctx *gin.Context) {
	res := gins.GetAPIResponse(ctx)

	qry := qry.NewDemo()

	err := parseModel(ctx, qry)
	if err != nil {
		res.SetError(err)
		return
	}

	ml, err := service.Demo.List(ctx, qry)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetData(ml)
}
