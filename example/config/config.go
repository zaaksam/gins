package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zaaksam/gins"
	"github.com/zaaksam/gins/extend/fileutil"
	"github.com/zaaksam/gins/extend/logger"
	"gopkg.in/yaml.v2"
)

// Instance 配置实例
var Instance config

type config struct {
	// Gins 配置
	Gins     *gins.Config `yaml:"-"`
	Task     bool         `yaml:"-"`
	ExecPath string       `yaml:"-"` // 运行程序的目录
	RootPath string       `yaml:"-"` // 程序所在目录

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

	Instance.ExecPath, _ = os.Getwd()

	executablePath, _ := os.Executable()
	executablePaths := strings.Split(executablePath, string(os.PathSeparator))
	Instance.RootPath = strings.Join(executablePaths[:len(executablePaths)-1], string(os.PathSeparator))

	if Instance.Gins.Debug {
		logger.Infof("execPath: %s", Instance.ExecPath)
		logger.Infof("rootPath: %s", Instance.RootPath)
	}
}

// Init 初始化
func (conf *config) Init(version string) (err error) {
	confName := "config.yaml"
	confPath := filepath.Join(conf.RootPath, confName)

	if conf.Gins.Debug {
		// FIXME: 目前判断比较粗糙，最多 5 层目录检查
		// 调试模式下，兼容单元测试模式的根目录路径识别
		tempExecPath := conf.ExecPath
		tempConfPath := confPath
		for i := 0; i < 5; i++ {
			_, e := os.Stat(tempConfPath)
			if e == nil {
				// 文件存在
				confPath = tempConfPath
				break
			}

			tempExecPath = filepath.Join(tempExecPath, "..")
			tempConfPath = filepath.Join(tempExecPath, confName)
		}
	}

	body, err := fileutil.ReadFile(confPath)
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
