package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tiamxu/cactus/logic/repo"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/kit/sql"

	"github.com/koding/multiconfig"
	httpkit "github.com/tiamxu/kit/http"
)

var (
	cfg        *Config
	name       = "cmdb"
	configPath = "config/config.yaml"
)

// yaml文件内容映射到结构体
type Config struct {
	ENV      string                  `yaml:"env"`
	LogLevel string                  `yaml:"log_level"`
	HttpSrv  httpkit.GinServerConfig `yaml:"http_srv"`
	DB       *sql.Config             `yaml:"db" xml:"db" json:"db"`
}

// set log level
func (c *Config) Initial() (err error) {

	defer func() {
		if err == nil {
			log.Printf("config initialed, env: %s,name: %s", cfg.ENV, name)
		}
	}()
	//日志
	if level, err := logrus.ParseLevel(c.LogLevel); err != nil {
		return err
	} else {
		log.DefaultLogger().SetLevel(level)
	}

	if err = repo.Init(cfg.DB); err != nil {
		return fmt.Errorf("database initialization failed: %w", err)

	}

	return nil
}

// 读取配置文件
func loadConfig() {
	cfg = new(Config)

	// env := os.Getenv("ENV")
	env := "dev"

	switch env {
	case "dev":
		configPath = "config/config-dev.yaml"
	case "test":
		configPath = "config/config-test.yaml"
	case "prod":
		configPath = "config/config-prod.yaml"
	default:
		configPath = "config/config.yaml"
	}

	multiconfig.MustLoadWithPath(configPath, cfg)
}
