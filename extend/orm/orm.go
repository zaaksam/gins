package orm

var (
	isInit     bool
	instance   *Engine
	dbRegistry map[string]any // map[alias]struct{}
	dbDisable  map[string]any // map[alias]struct{}
)

func init() {
	dbRegistry = make(map[string]any)
	dbDisable = make(map[string]any)
}

// Init 初始化
func Init(confs []*Config) (err error) {
	if isInit {
		return
	}

	instance, err = newEngine(confs)
	if err != nil {
		return
	}

	isInit = true
	return
}

// Register 注册数据库
func Register(md IModel) {
	dbRegistry[md.DatabaseAlias()] = struct{}{}
}

// Disable 关闭数据库（不创建连接）
func Disable(mds ...IModel) {
	for _, md := range mds {
		dbDisable[md.DatabaseAlias()] = struct{}{}
	}
}
