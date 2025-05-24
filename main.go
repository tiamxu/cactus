package main

import (
	"github.com/gin-gonic/gin"
	conf "github.com/tiamxu/cactus/config"
	"github.com/tiamxu/cactus/logic/routes"

	"github.com/tiamxu/kit/log"
)

func main() {
	cfg := conf.LoadConfig()
	if err := cfg.Initial(); err != nil {
		log.Fatalf("Config initialization failed: %v", err)
	}

	r := gin.Default()
	routes.InitRoutes(r)
	if err := r.Run(cfg.HttpSrv.Address); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}

}
