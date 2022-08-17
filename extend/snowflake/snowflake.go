package snowflake

import (
	"time"

	"github.com/zaaksam/gins/extend/logger"
)

var instance *sony

func init() {
	var err error
	instance, err = newSony()
	if err != nil {
		logger.Warnf("创建 snowflake 实例错误：%s", err)
	}
}

// NewID 创建新 ID
func NewID() (id uint64) {
	if instance == nil {
		// 确保ID总是创建成功
		id = uint64(time.Now().Unix())
		logger.Warnf("创建 snowflake id 出现问题：实例不存在")
		return
	}

	id, err := instance.NewID()
	if err != nil {
		// 确保ID总是创建成功
		id = uint64(time.Now().Unix())
		logger.Warnf("创建 snowflake id 出现问题：%s", err)
	}

	return
}
