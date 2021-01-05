package main

import (
	"hilive/engine"
	"hilive/modules/config"

	_ "github.com/go-sql-driver/mysql" // mysql引擎
)

func main() {
	cfg := config.Config{
		Database: config.Database{
			Host:       "35.194.236.160",
			Port:       "3306",
			User:       "yo",
			Pwd:        "yo123456",
			Name:       "hilive",
			MaxIdleCon: 50,
			MaxOpenCon: 150,
			Driver:     "mysql",
		},
		URLPrefix: "admin",
		IndexURL:  "/info/manager",
		LoginURL:  "/login",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
	}
	// InitDatabase 初始化資料庫引擎後將driver加入Engine.Services
	// InitRouter 初始化及設置api路由
	engine.DefaultEngine().InitDatabase(cfg).InitRouter().Gin.Run(":8080")
}
