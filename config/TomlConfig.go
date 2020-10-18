package config

import (
	"github.com/pelletier/go-toml"
	"stu-manager/logger"
)

type Config struct {
	mRootTree *toml.Tree
}

var sConfig *Config = nil

func GetConfig() *Config {
	if sConfig == nil {
		logger.Error.Fatalln("must call InitConfig first")
		return nil
	}
	return sConfig
}

func InitConfig(configPath string) {
	rootTree, err := toml.LoadFile(configPath)
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
	sConfig = &Config{mRootTree: rootTree}
}

func (config Config) GetRootTree() *toml.Tree {
	return config.mRootTree
}
