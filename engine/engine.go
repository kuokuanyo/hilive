package engine

import (
	"hilive/adapter"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/service"
	"hilive/plugins"
	"hilive/plugins/admin"

	"github.com/gin-gonic/gin"
)

var (
	// gin引擎
	r              *gin.Engine
	defaultAdapter adapter.WebFrameWork
)

// Engine 核心組件，有PluginList及Adapter兩個屬性
type Engine struct {
	Adapter    adapter.WebFrameWork
	config     *config.Config
	Services   service.List // 儲存資料庫引擎、Config(struct)
	PluginList []plugins.Plugin
}

// DefaultEngine 預設Engine(struct)
func DefaultEngine() *Engine {
	return &Engine{
		Adapter:  defaultAdapter,
		Services: service.GetServices(),
	}
}

// Register 建立引擎預設的配適器
func Register(ada adapter.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}

// InitDatabase 初始化資料庫引擎後將driver加入Engine.Services
func (eng *Engine) InitDatabase(cfg config.Config) *Engine {
	if eng.Adapter == nil {
		panic("adapter is nil, import the default adapter or use AddAdapter method add the adapter")
	}

	// 設置Config(struct)title、theme、登入url、前綴url...資訊，如果config數值為空值則設置預設值
	// ***此函式處理全局變數globalCfg
	eng.config = config.SetGlobalConfig(cfg)

	// GetConnectionByDriver藉由資料庫引擎(mysql、mssql...)取得Connection(interface)
	// InitDB初始化資料庫
	eng.Services.Add(cfg.Database.Driver, db.GetConnectionByDriver(cfg.Database.Driver).InitDB(cfg.Database))

	if defaultAdapter == nil {
		panic("adapter is nil")
	}
	return eng
}

// FindPluginByName 尋找與參數符合的plugin(interface)，如果有回傳Plugin,true，反之nil, false
func (eng *Engine) FindPluginByName(name string) (plugins.Plugin, bool) {
	for _, plug := range eng.PluginList {
		if plug.Name() == name {
			return plug, true
		}
	}
	return nil, false
}

// Use 設置Plugin、Admin、Guard、Handler等資訊
func (eng *Engine) Use(router interface{}) error {
	if eng.Adapter == nil {
		panic("adapter is nil, import the default adapter or use AddAdapter method add the adapter")
	}

	// FindPluginByName 尋找與參數符合的plugin(interface)，如果有回傳Plugin,true，反之nil, false
	_, exist := eng.FindPluginByName("admin")
	if !exist {
		eng.PluginList = append(eng.PluginList, admin.NewAdmin())
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

	// 將config加入services中
	eng.Services.Add("config", config.ConvertConfigToService(eng.config))

	// 取得匹配的Services然後轉換成Connection(interface)
	defaultConnection := db.GetConnectionFromService(eng.Services)

	// SetConnection為WebFrameWork(interface)的方法
	//設定連線
	defaultAdapter.SetConnection(defaultConnection)
	eng.Adapter.SetConnection(defaultConnection)

	// 執行初始化plugin
	for i := range eng.PluginList {
		eng.PluginList[i].InitPlugin(eng.Services)
	}

	return eng.Adapter.Use(router, eng.PluginList)
}

// GetConnectionByDriver 透過資料庫引擎取得Connection(interface)
func (eng *Engine) GetConnectionByDriver() db.Connection {
	// ***Engine.Services儲存資料庫引擎
	// ***資料庫引擎是Service也是Connection
	// Get 取得匹配的資料庫引擎
	// ConvertServiceToConnection 將Service轉換Connection
	return db.ConvertServiceToConnection(eng.Services.Get(eng.config.Database.Driver))
}
