package orm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"xorm.io/xorm"
)

// DA 数据库操作对象
type DA[T IModel] struct {
	*xorm.Session
	xormEngine *xorm.Engine

	model    T
	hasOrder bool
	groupBy  string
}

// NewDA 创建数据库操作对象
func NewDA[T IModel](md T) (da *DA[T], err error) {
	xormEngine, err := instance.getXormEngine(md.DatabaseAlias())
	if err != nil {
		return
	}

	da = &DA[T]{
		xormEngine: xormEngine,
		model:      md,
	}

	da.Session = da.xormEngine.NewSession()
	return
}

// UpdateByModel 通过模型更新数据，没传参数时使用内置模型对象
func (da *DA[T]) UpdateByModel(imdOpt ...IModel) (int64, error) {
	var imd IModel
	if len(imdOpt) == 1 {
		imd = imdOpt[0]
	} else {
		imd = da.model
	}

	columns := imd.FieldNames(true)
	if len(columns) <= 0 {
		return 0, errors.New("IModel没有设置字段内容，无法更新")
	}

	return da.Session.Cols(columns...).Update(imd)
}

// UpdateByBuilder 使用Builder进行更新
func (da *DA[T]) UpdateByBuilder(builder *Builder) (rows int64, err error) {
	sqlResult, err := da.Exec(builder.Builder)
	if err != nil {
		return
	}

	rows, err = sqlResult.RowsAffected()
	return
}

// Desc 降序
func (da *DA[T]) Desc(colNames ...string) *xorm.Session {
	da.hasOrder = true
	return da.Session.Desc(colNames...)
}

// Asc 升序
func (da *DA[T]) Asc(colNames ...string) *xorm.Session {
	da.hasOrder = true
	return da.Session.Asc(colNames...)
}

// OrderBy 排序
func (da *DA[T]) OrderBy(order string) *xorm.Session {
	da.hasOrder = true
	return da.Session.OrderBy(order)
}

// GroupBy 分组
func (da *DA[T]) GroupBy(keys string) *xorm.Session {
	da.groupBy = keys
	return da.Session.GroupBy(keys)
}

// GetModelList 获取列表数据，items可指定接收查询结果的结构对象
func (da *DA[T]) GetModelList(page, pageSize int, itemsOpt ...[]T) (ml *ModelList[T], err error) {
	ml, err = da.GetModelListWithPageOffset(page, pageSize, 0, itemsOpt...)
	return
}

// GetModelListWithPageOffset 获取列表数据，items可指定接收查询结果的结构对象
func (da *DA[T]) GetModelListWithPageOffset(page, pageSize, pageOffset int, itemsOpt ...[]T) (ml *ModelList[T], err error) {
	if !da.hasOrder {
		err = errors.New("没有设置Orderby条件")
		return
	}

	ml = &ModelList[T]{
		Page:     page,
		PageSize: pageSize,
	}
	if ml.Page == 1 && pageOffset < 0 {
		pageOffset = 0
	}

	if len(itemsOpt) == 1 {
		ml.Items = itemsOpt[0]
	} else {
		ml.Items = make([]T, 0, pageSize)
	}

	// 查询总条数
	sessionCount := da.xormEngine.NewSession()
	defer sessionCount.Close()

	// 获取查询条件
	cond := da.Conds()
	if da.groupBy != "" {
		ml.Total, err = sessionCount.Select("count(DISTINCT " + da.groupBy + ")").Where(cond).Count(da.model)
	} else {
		ml.Total, err = sessionCount.Where(cond).Count(da.model)
	}
	if err != nil {
		return
	}
	ml.Total = ml.Total - int64(pageOffset)

	//没有数据，提前返回
	if ml.Total <= 0 {
		ml.Total = 0
		da.resetPagination(ml)
		return
	}

	da.Limit(da.getModelListLimit(ml, pageOffset))
	err = da.Table(da.model).Find(&ml.Items)
	return
}

// GetModelListByBuilder 获取列表数据，通过Builder构建查询，items可指定接收查询结果的结构对象
func (da *DA[T]) GetModelListByBuilder(builder *Builder, page, pageSize int, itemsOpt ...[]T) (ml *ModelList[T], err error) {
	ml, err = da.GetModelListByBuilderWithPageOffset(builder, page, pageSize, 0, itemsOpt...)
	return
}

// GetModelListByBuilderWithPageOffset 获取列表数据，通过Builder构建查询，items可指定接收查询结果的结构对象
func (da *DA[T]) GetModelListByBuilderWithPageOffset(builder *Builder, page, pageSize, pageOffset int, itemsOpt ...[]T) (ml *ModelList[T], err error) {
	// 取出参数临时置空，方便 count 符合条件的记录行数
	groupByStr := builder.GroupByStr
	orderByStr := builder.OrderByStr
	builder.GroupBy("")
	builder.OrderBy("")

	switch {
	case orderByStr == "":
		err = errors.New("没有设置Orderby条件")
	}
	if err != nil {
		return
	}

	ml = &ModelList[T]{
		Page:     page,
		PageSize: pageSize,
	}

	if len(itemsOpt) == 1 {
		ml.Items = itemsOpt[0]
	} else {
		ml.Items = make([]T, 0, pageSize)
	}

	conf, ok := instance.confs[da.model.DatabaseAlias()]
	if !ok || (conf.Type != "mssql" && conf.Type != "mysql") {
		err = fmt.Errorf("数据库 %s 的类型 %s 不支持 Builder 查询", conf.Alias, conf.Type)
		return
	}

	sql, args, err := builder.ToSQL()
	if err != nil {
		return
	}

	if sql == "" {
		err = errors.New("Builder不能为空")
		return
	}

	//重组sql语句，发起 count 查询
	countSQL := "select count("
	if groupByStr != "" {
		countSQL += "DISTINCT " + groupByStr
	} else {
		countSQL += "*"
	}
	countSQL += ") Total "

	tempSQL := strings.ToLower(sql)
	fromIndex := strings.Index(tempSQL, "from")
	countSQL = strings.Replace(sql, string(sql[0:fromIndex]), countSQL, 1)

	var sess = *da.Session
	countSession := &sess
	defer countSession.Close()
	_, err = countSession.SQL(countSQL, args...).Get(&ml.Total)
	if err != nil {
		return
	}
	ml.Total = ml.Total - int64(pageOffset)

	//没有数据，提前返回
	if ml.Total <= 0 {
		ml.Total = 0
		da.resetPagination(ml)
		return
	}

	selectIndex := strings.Index(tempSQL, "select")
	var selectSQL string

	_, rowStart := da.getModelListLimit(ml, pageOffset)

	if conf.Type == "mssql" {
		rowEnd := ml.Page*ml.PageSize + pageOffset

		selectSQL = "select * from (select" + sql[selectIndex+6:fromIndex]
		selectSQL += ",ROW_NUMBER() over (order by " + orderByStr + ") as row_num "
		selectSQL += sql[fromIndex:]

		if groupByStr != "" {
			selectSQL += " group by " + groupByStr
		}

		selectSQL += ") tb where row_num between "
		selectSQL += strconv.Itoa(rowStart+1) + " and " + strconv.Itoa(rowEnd)

	} else if conf.Type == "mysql" {
		selectSQL = "select " + sql[selectIndex+6:fromIndex] + " "
		selectSQL += sql[fromIndex:]

		if groupByStr != "" {
			selectSQL += " group by " + groupByStr
		}

		selectSQL += " order by " + orderByStr
		selectSQL += " limit " + strconv.Itoa(rowStart) + "," + strconv.Itoa(ml.PageSize)
	}

	err = da.Session.SQL(selectSQL, args...).Find(ml.Items)
	return
}

// getModelListLimit 获取分页计算
func (da *DA[T]) getModelListLimit(ml *ModelList[T], pageOffset int) (int, int) {
	da.resetPagination(ml)
	start := (ml.Page-1)*ml.PageSize + pageOffset
	return ml.PageSize, start
}

func (da *DA[T]) resetPagination(ml *ModelList[T]) {
	//检查页面长度
	if ml.PageSize <= 0 {
		ml.PageSize = 10
	}
	// } else if ml.PageSize > 500 {
	// 	ml.PageSize = 500
	// }

	//计算总页数
	cnt := int(ml.Total / int64(ml.PageSize))
	mod := int(ml.Total % int64(ml.PageSize))
	if mod > 0 {
		cnt++
	}

	//检查页面索引
	switch {
	case cnt == 0:
		ml.Page = 1
	case ml.Page > cnt:
		ml.Page = cnt
	case ml.Page <= 0:
		ml.Page = 1
	}

	//设置页面总页数
	ml.PageCount = cnt
}
