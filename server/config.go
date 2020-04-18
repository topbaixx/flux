package server

import (
	"github.com/bytepowered/flux"
	"github.com/bytepowered/flux/ext"
	"github.com/bytepowered/flux/logger"
	"github.com/bytepowered/flux/pkg"
	"os"
	"sync"
)

const (
	EnvKeyApplicationConfigPath  = "APP_CONF_FILE"
	DefaultApplicationConfigPath = "conf.d/application.toml"
)

var (
	_globals    flux.Config
	_globalOnce = new(sync.Once)
)

// GlobalConfig 返回全局配置
func GlobalConfig() flux.Config {
	return _globals
}

func LoadConfig(outpath string) flux.Config {
	_globalOnce.Do(func() {
		configPath := DefaultApplicationConfigPath
		// 1. Env配置
		if envpath := os.Getenv(EnvKeyApplicationConfigPath); envpath != "" {
			configPath = envpath
		} else if outpath != "" {
			// 2. 外部配置
			configPath = outpath
		}
		logger.Infof("Using config, path: %s", configPath)
		if data, err := pkg.LoadTomlFile(configPath); nil != err {
			logger.Panicf("Config not found: %s", configPath)
		} else {
			_globals = ext.NewMapConfig(data)
		}
	})
	return _globals
}
