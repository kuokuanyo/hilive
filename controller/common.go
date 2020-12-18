package controller

import (
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/service"
	"html/template"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	Config *config.Config
	Conn   db.Connection
	Gin    *gin.Engine
	Services service.List
}

// DefaultFuncMap 模板需要使用的函式
var DefaultFuncMap = template.FuncMap{
	"link": func(cdnUrl, prefixUrl, assetsUrl string) string {
		if cdnUrl == "" {
			return prefixUrl + assetsUrl
		}
		return cdnUrl + assetsUrl
	},
	"isLinkURL": func(s string) bool {
		return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
	},
}
