package log_init

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ErrorLogConfig struct {
	LogDir    string `yaml:"dir"`
	FileName  string `yaml:"file_name"`
	LogLevel  string `yaml:"log_level"`
	LogExpire string `yaml:"log_expire"`
	LogPeriod string `yaml:"log_period"`
}

func (c *ErrorLogConfig) GetLogDir() string {
	logDir := c.LogDir
	if logDir == "" {
		//默认路径是可执行程序的上一层目录的 work/logs 根据系统自适应
		lastDir, err := GetLastAbsPathByExecutable()
		if err != nil {
			panic(err)
		}
		logDir = fmt.Sprintf("%s%swork%slog", lastDir, string(os.PathSeparator), string(os.PathSeparator))
	} else if !strings.HasPrefix(logDir, string(os.PathSeparator)) {
		//若目录配置不为绝对路径, 则路径为 上一层目录路径 + 配置的目录路径
		lastDir, err := GetLastAbsPathByExecutable()
		if err != nil {
			panic(err)
		}
		relativePathPrefix := fmt.Sprintf("..%s", string(os.PathSeparator))
		logDir = path.Join(lastDir, strings.TrimPrefix(logDir, relativePathPrefix))
	}
	dirPath, err := filepath.Abs(logDir)
	if err != nil {
		panic(err)
	}
	return dirPath
}

func (c *ErrorLogConfig) GetLogName() string {
	if c.FileName == "" {
		return "error.log"
	}
	return c.FileName
}
func (c *ErrorLogConfig) GetLogLevel() log.Level {
	l, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		l = log.InfoLevel
	}
	return l
}

func (c *ErrorLogConfig) GetLogExpire() time.Duration {
	if strings.HasSuffix(c.LogExpire, "h") {
		d, err := time.ParseDuration(c.LogExpire)
		if err != nil {
			return 7 * time.Hour
		}
		return d
	}
	if strings.HasSuffix(c.LogExpire, "d") {

		d, err := strconv.Atoi(strings.Split(c.LogExpire, "d")[0])
		if err != nil {
			return 7 * 24 * time.Hour
		}
		return time.Duration(d) * 24 * time.Hour
	}
	return 7 * 24 * time.Hour
}

func (c *ErrorLogConfig) GetLogPeriod() string {
	if c.LogPeriod == "" {
		return "day"
	}
	return c.LogPeriod
}

// GetLastAbsPathByExecutable 获取执行程序所在的上一层级目录的绝对路径
func GetLastAbsPathByExecutable() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	res, _ := filepath.EvalSymlinks(exePath)
	res = filepath.Dir(res) //获取程序所在目录
	res = filepath.Dir(res) //获取程序所在上一层级目录
	return res, nil
}
