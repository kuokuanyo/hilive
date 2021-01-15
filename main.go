package main

import (
	"hilive/engine"
	"hilive/modules/config"
	"log"
	"net/http"

	_ "hilive/adapter/gin" // 框架引擎

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // mysql引擎
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	r := gin.Default()
	router := r.Group("/admin")
	// 設置靜態檔案
	router.Static("/assets", "./assets")

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
	// Use 設置Plugin、Admin、Guard、Handler等資訊
	engine.DefaultEngine().InitDatabase(cfg).Use(r)

	// r.Run(":8080")
	log.Fatal(http.Serve(autocert.NewListener("hilive.com.tw"), r))
}
