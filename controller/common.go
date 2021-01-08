package controller

import (
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/service"
	"hilive/modules/table"
	"html/template"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	Config        *config.Config
	Conn          db.Connection
	Gin           *gin.Engine
	Services      service.List
	Alert         string
	TablelList map[string]func(conn db.Connection) table.Table
}

// URLRoute 模板需要使用的URL路徑
type URLRoute struct {
	URLPrefix   string
	IndexURL    string
	InfoURL     string
	NewURL      string
	EditURL     string
	DeleteURL   string
	SortURL     string
	PreviousURL string
}

// GetTable 取得table(面板資訊、表單資訊)
func (h *Handler) GetTable(ctx *gin.Context, prefix string) table.Table {
	return h.TablelList[prefix](h.Conn)
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

// isInfoURL 檢查url
func isInfoURL(s string) bool {
	reg, _ := regexp.Compile("(.*?)info/(.*?)$")
	sub := reg.FindStringSubmatch(s)
	return len(sub) > 2 && !strings.Contains(sub[2], "/")
}

// isNewURL 檢查url
func isNewURL(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/new")

	return reg.MatchString(s)
}

// isEditURL 檢查url
func isEditURL(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/edit")
	return reg.MatchString(s)
}
