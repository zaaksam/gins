package model

import "github.com/zaaksam/gins/extend/orm"

// Demo 模型
type Demo struct {
	orm.Model `xorm:"-"`

	ID      uint64 `json:"id" xorm:"bigint pk 'id'"`
	User    string `json:"user" xorm:"varchar(50) not null 'user'"`
	Pswd    string `json:"pswd" xorm:"varchar(100) not null 'pswd'"`
	Status  int    `json:"status" xorm:"tinyint not null 'status'"`
	Created int64  `json:"created" xorm:"bigint not null 'cteated'"`
	Updated int64  `json:"updated" xorm:"bigint not null 'updated'"`
}

func (*Demo) TableName() string {
	return "demo"
}

func (*Demo) DatabaseAlias() string {
	return "gins_db"
}
