package engine

import (
	"hilive/controller"
	"hilive/guard"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/service"

	"github.com/gin-gonic/gin"
)

var (
	// Router gin引擎
	r *gin.Engine
)

func init() {
	// Default returns an Engine instance with the Logger and Recovery middleware already attached.
	r = gin.Default()
}

// Engine 核心組件，有PluginList及Adapter兩個屬性
type Engine struct {
	config   *config.Config
	Services service.List // 儲存資料庫引擎、Config(struct)
	Gin      *gin.Engine
	handler  controller.Handler
	guard    guard.Guard
}

// DefaultEngine 預設Engine(struct)
func DefaultEngine() *Engine {
	return &Engine{
		Services: service.GetServices(),
		Gin:      r,
	}
}

// InitDatabase 初始化資料庫引擎後將driver加入Engine.Services
func (eng *Engine) InitDatabase(cfg config.Config) *Engine {
	// 設置Config(struct)title、theme、登入url、前綴url...資訊，如果config數值為空值則設置預設值
	// ***此函式處理全局變數globalCfg
	eng.config = config.SetGlobalConfig(cfg)

	// GetConnectionByDriver藉由資料庫引擎(mysql、mssql...)取得Connection(interface)
	// InitDB初始化資料庫
	eng.Services.Add(cfg.Database.Driver, db.GetConnectionByDriver(cfg.Database.Driver).InitDB(cfg.Database))

	// 增加token的service
	eng.Services.Add("token_csrf_helper", &auth.TokenService{
		Tokens: make(auth.CSRFToken, 0),
	})

	eng.handler = controller.Handler{
		Config:   eng.config,
		Conn:     eng.GetConnectionByDriver(),
		Gin:      r,
		Services: eng.Services,
	}

	eng.guard = guard.Guard{
		Conn:     eng.GetConnectionByDriver(),
		Services: eng.Services,
		Config:   eng.config,
	}

	// 設置TableName、Conn
	// GetConnectionByDriver 透過資料庫引擎轉換Connection(interface)
	site := models.DefaultSiteModel().SetConn(eng.GetConnectionByDriver())

	// ToMap 將Config轉換為map[string][string](放在Site資料表中)
	// 初始化site資料表，插入或更新site資料表
	site.Init(eng.config.ToMap())

	var m = make(map[string]string, 0)
	items, err := site.Base.Table(site.Base.TableName).All()
	if err != nil {
		panic("取得site資料表資料發生錯誤")
	}
	for _, item := range items {
		m[item["config_key"].(string)] = item["config_value"].(string)
	}

	// 更新Config
	eng.config.Update(m)
	return eng
}

// GetConnectionByDriver 透過資料庫引擎取得Connection(interface)
func (eng *Engine) GetConnectionByDriver() db.Connection {
	// ***Engine.Services儲存資料庫引擎
	// ***資料庫引擎是Service也是Connection
	// Get 取得匹配的資料庫引擎
	// ConvertServiceToConnection 將Service轉換Connection
	return db.ConvertServiceToConnection(eng.Services.Get(eng.config.Database.Driver))
}

// InitRouter 設置路由
func (eng *Engine) InitRouter() *Engine {
	router := eng.Gin.Group("/admin")
	router.Static("/assets", "./assets")
	// 登入GET、POST
	router.GET(eng.config.LoginURL, eng.handler.ShowLogin)
	router.POST(eng.config.LoginURL, eng.handler.Auth)
	// 註冊GET、POST
	router.GET(eng.config.SignupURL, eng.handler.ShowSignup)
	router.POST(eng.config.SignupURL, eng.handler.Signup)

	authRoute := router.Use(auth.DefaultInvoker(eng.handler.Conn).Middleware(eng.handler.Conn))
	// 菜單
	authRoute.GET(eng.config.MenuURL, eng.handler.ShowMenu)
	authRoute.GET(eng.config.MenuNewURL, eng.handler.ShowNewMenu)
	authRoute.GET(eng.config.MenuEditURL, eng.handler.ShowEditMenu)
	authRoute.POST(eng.config.MenuNewURL, eng.guard.MenuNew, eng.handler.NewMenu)
	authRoute.POST(eng.config.MenuEditURL, eng.guard.MenuEdit, eng.handler.EditMenu)
	authRoute.POST(eng.config.MenuDeleteURL, eng.guard.MenuDelete, eng.handler.DeleteMenu)

	// 使用者
	authRoute.GET(eng.config.ManagerURL, eng.handler.ShowManegerInfo)
	authRoute.GET(eng.config.ManagerNewURL, eng.guard.ShowManagerNewForm, eng.handler.ShowManagerNewForm)

	// 角色
	authRoute.GET(eng.config.RolesURL, eng.handler.ShowRolesInfo)

	// 權限
	authRoute.GET(eng.config.PermissionURL, eng.handler.ShowPermissionInfo)

	return eng
}
