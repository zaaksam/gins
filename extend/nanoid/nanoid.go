package nanoid

import (
	"fmt"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/zaaksam/gins/extend/logger"
)

var generator func() string

func init() {
	var err error
	generator, err = nanoid.CustomASCII("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	if err != nil {
		generator = nil
		logger.Warnf("创建 nanoid 实例错误：%s", err)
		return
	}
}

// New 创建
func New() string {
	if generator == nil {
		// 确保总是创建成功
		return fmt.Sprintf("%d", time.Now().Unix())
	}

	return generator()
}
