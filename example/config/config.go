package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/zaaksam/gins"
	"github.com/zaaksam/gins/example/utils"
	"gopkg.in/yaml.v2"
)

// Instance 配置实例
var Instance config

type config struct {
	// Gins 配置
	Gins       *gins.Config `yaml:"-"`
	Task       bool         `yaml:"-"`
	RootPath   string       `yaml:"-"` // 程序所在目录
	ConfigPath string       `yaml:"-"` // config.yaml 所在目录

	LogLevel string `yaml:"logLevel"` // 默认：info
	Redis    struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	DBs []struct {
		Alias     string `yaml:"alias"`
		Type      string `yaml:"type"`
		Server    string `yaml:"server"`
		Port      int    `yaml:"port"`
		Database  string `yaml:"database"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		LogLevel  string `yaml:"logLevel"`
		IsShowSQL bool   `yaml:"isShowSQL"`
	} `yaml:"dbs"`
}

var _ = func() bool {
	testing.Init()
	return true
}()

func init() {
	Instance.Gins = &gins.Config{}

	// 从命令行参数加载配置
	flag.StringVar(&Instance.Gins.Name, "name", "", "应用名称")
	flag.StringVar(&Instance.Gins.Host, "host", "", "Host名称")
	flag.StringVar(&Instance.Gins.IP, "ip", "", "运行IP，默认：0.0.0.0")
	flag.StringVar(&Instance.Gins.BroadcastIP, "broadcast-ip", "", "广播IP，默认为：IP")
	flag.IntVar(&Instance.Gins.Port, "port", 8080, "运行端口，默认：8080")
	flag.IntVar(&Instance.Gins.BroadcastPort, "broadcast-port", 0, "广播端口，默认为：Port")
	flag.IntVar(&Instance.Gins.Timeout, "timeout", 30, "优雅退出时的超时秒数，默认：30秒")

	flag.BoolVar(&Instance.Gins.Debug, "debug", false, "调试模式，默认：false")
	flag.BoolVar(&Instance.Gins.Pprof, "pprof", false, "性能监控模式，默认：false")
	flag.BoolVar(&Instance.Task, "task", false, "是否开启后台任务，默认：false")

	flag.Parse()

	if Instance.Gins.IP == "" {
		Instance.Gins.IP = "0.0.0.0"
	}

	if Instance.Gins.Port <= 0 {
		Instance.Gins.Port = 8080
	}

	if Instance.Gins.Name == "" {
		Instance.Gins.Name = "gins_example"
	}

	Instance.RootPath, _ = os.Getwd()
}

// Init 初始化
func (conf *config) Init(version string) (err error) {
	if Instance.Gins.Debug {
		// 调试模式下，兼容单元测试模式的根目录路径识别
		_, currentFile, _, ok := runtime.Caller(0)
		if ok {
			currentFile = strings.ReplaceAll(currentFile, string(os.PathSeparator)+filepath.Join("config", "config.go"), "")
			if currentFile != Instance.RootPath {
				Instance.RootPath = currentFile
			}
		}
	}

	body, err := utils.ReadFile(filepath.Join(conf.RootPath, "config.yaml"))
	if err != nil {
		err = fmt.Errorf("读取 config.yaml 配置时出错：%s", err)
		return
	}

	err = yaml.Unmarshal(body, conf)
	if err != nil {
		err = fmt.Errorf("解析 config.yaml 配置时出错：%s", err)
		return
	}

	conf.Gins.LogLevel = conf.LogLevel
	conf.Gins.Version = version
	err = conf.Gins.Init()
	if err != nil {
		err = fmt.Errorf("初始化 gins 配置时出错：%s", err)
	}

	return
}
