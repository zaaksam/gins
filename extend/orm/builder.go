package orm

import (
	"errors"

	"strings"

	"xorm.io/builder"
)

type Builder struct {
	*builder.Builder
	OrderByStr string
	GroupByStr string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func NewMySQLBuilder() *Builder {
	b := &Builder{}
	b.Builder = builder.MySQL()
	return b
}

func (b *Builder) Asc(colNames ...string) {
	if len(colNames) == 0 {
		return
	}

	if b.OrderByStr != "" {
		b.OrderByStr = ","
	}

	b.OrderByStr += strings.Join(colNames, " ASC,") + " ASC"

	b.Builder.OrderBy(b.OrderByStr)
}

func (b *Builder) Desc(colNames ...string) {
	if len(colNames) == 0 {
		return
	}

	if b.OrderByStr != "" {
		b.OrderByStr = ","
	}

	b.OrderByStr += strings.Join(colNames, " DESC,") + " DESC"

	b.Builder.OrderBy(b.OrderByStr)
}

func (b *Builder) OrderBy(orderBy string) {
	b.OrderByStr = orderBy

	b.Builder.OrderBy(b.OrderByStr)
}

func (b *Builder) GroupBy(groupBy string) {
	b.GroupByStr = groupBy

	b.Builder.GroupBy(b.GroupByStr)
}

func (b *Builder) getTableName(table interface{}) string {
	switch val := table.(type) {
	case IModel:
		return val.TableName()
	case string:
		return val
	default:
		panic(errors.New("table 只支持 IModel 接口和 string 类型"))
	}
}

//******** 以下为重写 xorm Builder 方法 ***********

//Select 重写Builder.Select方法，总是创建新对象
func (b *Builder) Select(cols ...string) *Builder {
	b.Builder = builder.Select(cols...)

	return b
}

//Insert 重写Builder.Insert方法，总是创建新对象
func (b *Builder) Insert(eq builder.Eq) *Builder {
	b.Builder = builder.Insert(eq)

	return b
}

//Update 重写Builder.Update方法，总是创建新对象
func (b *Builder) Update(updates ...builder.Cond) *Builder {
	b.Builder = builder.Update(updates...)

	return b
}

//Delete 重写Builder.Delete方法，总是创建新对象
func (b *Builder) Delete(conds ...builder.Cond) *Builder {
	b.Builder = builder.Delete(conds...)

	return b
}

//From 重写Builder.From方法，支持model传入
func (b *Builder) From(table interface{}) *Builder {
	b.Builder.From(b.getTableName(table))

	return b
}

//Into 重写Builder.Into方法，支持model传入
func (b *Builder) Into(table interface{}) *Builder {
	b.Builder.Into(b.getTableName(table))

	return b
}

//Join 重写Builder.Join方法，支持model传入
func (b *Builder) Join(joinType string, joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.Join(joinType, b.getTableName(joinTable), joinCond)

	return b
}

//InnerJoin 重写Builder.InnerJoin方法，支持model传入
func (b *Builder) InnerJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.InnerJoin(b.getTableName(joinTable), joinCond)

	return b
}

//LeftJoin 重写Builder.LeftJoin方法，支持model传入
func (b *Builder) LeftJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.LeftJoin(b.getTableName(joinTable), joinCond)

	return b
}

//RightJoin 重写Builder.RightJoin方法，支持model传入
func (b *Builder) RightJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.RightJoin(b.getTableName(joinTable), joinCond)

	return b
}

//CrossJoin 重写Builder.CrossJoin方法，支持model传入
func (b *Builder) CrossJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.CrossJoin(b.getTableName(joinTable), joinCond)

	return b
}

//FullJoin 重写Builder.FullJoin方法，支持model传入
func (b *Builder) FullJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.FullJoin(b.getTableName(joinTable), joinCond)

	return b
}

func (b *Builder) Where(cond builder.Cond) *Builder {
	b.Builder.Where(cond)

	return b
}

func (b *Builder) And(cond builder.Cond) *Builder {
	b.Builder.And(cond)

	return b
}

func (b *Builder) Or(cond builder.Cond) *Builder {
	b.Builder.Or(cond)

	return b
}
