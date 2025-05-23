package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/conf"
	"github.com/tiamxu/cactus/logic/routes"

	"github.com/tiamxu/kit/log"
)

var cfg *conf.Config

func init() {
	conf.LoadConfig()
	if err := cfg.Initial(); err != nil {
		log.Fatalf("Config initialization failed: %v", err)
	}
}
func main() {
	r := gin.Default()

	// 注册路由
	routes.InitRoutes(r)

	// 启动服务
	if err := r.Run(cfg.HttpSrv.Address); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}

}
