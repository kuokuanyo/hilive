package config

import (
	"hilive/modules/utils"
	"html/template"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	globalCfg  = new(Config)
	count      uint32
	updateLock sync.Mutex
	lock       sync.Mutex
)

// Database 資料庫引擎設置
type Database struct {
	Host       string            `json:"host,omitempty" yaml:"host,omitempty" ini:"host,omitempty"`
	Port       string            `json:"port,omitempty" yaml:"port,omitempty" ini:"port,omitempty"`
	User       string            `json:"user,omitempty" yaml:"user,omitempty" ini:"user,omitempty"`
	Pwd        string            `json:"pwd,omitempty" yaml:"pwd,omitempty" ini:"pwd,omitempty"`
	Name       string            `json:"name,omitempty" yaml:"name,omitempty" ini:"name,omitempty"`
	MaxIdleCon int               `json:"max_idle_con,omitempty" yaml:"max_idle_con,omitempty" ini:"max_idle_con,omitempty"`
	MaxOpenCon int               `json:"max_open_con,omitempty" yaml:"max_open_con,omitempty" ini:"max_open_con,omitempty"`
	Driver     string            `json:"driver,omitempty" yaml:"driver,omitempty" ini:"driver,omitempty"`
	File       string            `json:"file,omitempty" yaml:"file,omitempty" ini:"file,omitempty"`
	Dsn        string            `json:"dsn,omitempty" yaml:"dsn,omitempty" ini:"dsn,omitempty"`
	Params     map[string]string `json:"params,omitempty" yaml:"params,omitempty" ini:"params,omitempty"`
}

// Service 用於Service之間的轉換
type Service struct {
	Config *Config
}

// Config 基本配置
type Config struct {
	// 支持多個資料庫連接
	Database Database `json:"database,omitempty" yaml:"database,omitempty" ini:"database,omitempty"`

	// The cookie domain used in the auth modules. see the session.go.
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty" ini:"domain,omitempty"`

	// The path where files will be stored into.
	Store Store `json:"store,omitempty" yaml:"store,omitempty" ini:"store,omitempty"`

	// The global url prefix.
	URLPrefix string `json:"prefix,omitempty" yaml:"prefix,omitempty" ini:"prefix,omitempty"`

	// The title of web page.
	Title string `json:"title,omitempty" yaml:"title,omitempty" ini:"title,omitempty"`

	// Login page title
	LoginTitle string `json:"login_title,omitempty" yaml:"login_title,omitempty" ini:"login_title,omitempty"`

	// Login page logo
	LoginLogo template.HTML `json:"login_logo,omitempty" yaml:"login_logo,omitempty" ini:"login_logo,omitempty"`

	// 側邊欄的logo
	Logo template.HTML `json:"logo,omitempty" yaml:"logo,omitempty" ini:"logo,omitempty"`

	// 側邊欄縮小後的logo
	MiniLogo template.HTML `json:"mini_logo,omitempty" yaml:"mini_logo,omitempty" ini:"mini_logo,omitempty"`

	// The url redirect to after login.
	IndexURL string `json:"index,omitempty" yaml:"index,omitempty" ini:"index,omitempty"`

	// Login page URL
	LoginURL string `json:"login_url,omitempty" yaml:"login_url,omitempty" ini:"login_url,omitempty"`

	// Signup page URL
	SignupURL string `json:"signup_URL,omitempty" yaml:"signup_URL,omitempty" ini:"signup_URL,omitempty"`

	// Menu url
	MenuURL string `json:"menu_URL,omitempty" yaml:"menu_URL,omitempty" ini:"menu_URL,omitempty"`

	// Menu new url
	MenuNewURL string `json:"menu_new_URL,omitempty" yaml:"menu_new_URL,omitempty" ini:"menu_new_URL,omitempty"`

	// Menu edit url
	MenuEditURL string `json:"menu_edit_URL,omitempty" yaml:"menu_edit_URL,omitempty" ini:"menu_edit_URL,omitempty"`

	// Menu delete url
	MenuDeleteURL string `json:"menu_delete_URL,omitempty" yaml:"menu_delete_URL,omitempty" ini:"menu_delete_URL,omitempty"`

	// ManagerURL manager url
	ManagerURL string `json:"manager_URL,omitempty" yaml:"manager_URL,omitempty" ini:"manager_URL,omitempty"`

	// ManagerNewURL manager url
	ManagerNewURL string `json:"manager_new_URL,omitempty" yaml:"manager_new_URL,omitempty" ini:"manager_new_URL,omitempty"`

	// ManagerEditURL manager url
	ManagerEditURL string `json:"manager_edit_URL,omitempty" yaml:"manager_edit_URL,omitempty" ini:"manager_edit_URL,omitempty"`

	// ManagerNewURLPost manager url
	ManagerNewURLPost string `json:"manager_new_URL_post,omitempty" yaml:"manager_new_URL_post,omitempty" ini:"manager_new_URL_post,omitempty"`

	// RolesURL roles url
	RolesURL string `json:"roles_URL,omitempty" yaml:"roles_URL,omitempty" ini:"roles_URL,omitempty"`

	// RolesNewURL manager url
	RolesNewURL string `json:"roles_new_URL,omitempty" yaml:"roles_new_URL,omitempty" ini:"roles_new_URL,omitempty"`

	// PermissionURL roles url
	PermissionURL string `json:"permission_URL,omitempty" yaml:"permission_URL,omitempty" ini:"permission_URL,omitempty"`

	// PermissionNewURL manager url
	PermissionNewURL string `json:"permission_new_URL,omitempty" yaml:"permission_new_URL,omitempty" ini:"permission_new_URL,omitempty"`

	// Assets visit link.
	AssetURL string `json:"asset_url,omitempty" yaml:"asset_url,omitempty" ini:"asset_url,omitempty"`

	// 使用者table名稱
	AuthUserTable string `json:"auth_user_table,omitempty" yaml:"auth_user_table,omitempty" ini:"auth_user_table,omitempty"`

	// Session valid time duration,units are seconds. Default 7200.
	SessionLifeTime int `json:"session_life_time,omitempty" yaml:"session_life_time,omitempty" ini:"session_life_time,omitempty"`

	// 不限制登入IP
	NoLimitLoginIP bool

	prefix string
}

// Store 儲存文件
type Store struct {
	Path   string `json:"path,omitempty" yaml:"path,omitempty" ini:"path,omitempty"`
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty" ini:"prefix,omitempty"`
}

// SetGlobalConfig 設置Config(struct)title、theme、登入url、前綴url...資訊，如果參數cfg(struct)有些數值為空值，設置預設值
// *****此函數設置全局變數globalCfg
func SetGlobalConfig(cfg Config) *Config {
	// 避免與更新數值同時執行，因此使用互斥鎖
	// 互斥鎖
	lock.Lock()
	defer lock.Unlock()

	// 不能設置config兩次
	// LoadUint32 atomically loads *addr.
	if atomic.LoadUint32(&count) != 0 {
		panic("can not set config twice")
	}
	atomic.StoreUint32(&count, 1)

	// SetDefault 如果參數cfg(struct)數值為空值，設置預設值
	cfg = SetDefault(cfg)

	// global url前綴
	if cfg.URLPrefix == "" {
		cfg.prefix = "/"
	} else if cfg.URLPrefix[0] != '/' {
		cfg.prefix = "/" + cfg.URLPrefix
	} else {
		cfg.prefix = cfg.URLPrefix
	}

	globalCfg = &cfg
	return globalCfg
}

// SetDefault 如果參數cfg(struct)數值為空值，設置預設值
func SetDefault(cfg Config) Config {
	// SetDefault假設第一個參數 = 第二個參數回傳第三個參數，沒有的話回傳第一個參數
	cfg.Title = utils.SetDefault(cfg.Title, "", "晶橙資訊")
	cfg.LoginTitle = utils.SetDefault(cfg.LoginTitle, "", "晶橙")
	cfg.Logo = template.HTML(utils.SetDefault(string(cfg.Logo), "", "<b>晶橙</b>"))
	cfg.MiniLogo = template.HTML(utils.SetDefault(string(cfg.MiniLogo), "", "<b>晶橙</b>"))
	cfg.IndexURL = utils.SetDefault(cfg.IndexURL, "", "/info/manager")
	cfg.LoginURL = utils.SetDefault(cfg.LoginURL, "", "/login")
	cfg.SignupURL = utils.SetDefault(cfg.SignupURL, "", "/signup")
	cfg.AuthUserTable = utils.SetDefault(cfg.AuthUserTable, "", "users")
	cfg.MenuURL = utils.SetDefault(cfg.MenuURL, "", "/menu")
	cfg.MenuEditURL = utils.SetDefault(cfg.MenuEditURL, "", "/menu/edit")
	cfg.MenuDeleteURL = utils.SetDefault(cfg.MenuDeleteURL, "", "/menu/delete")
	cfg.MenuNewURL = utils.SetDefault(cfg.MenuNewURL, "", "/menu/new")
	cfg.ManagerURL = utils.SetDefault(cfg.ManagerURL, "", "/info/manager")
	cfg.ManagerNewURL = utils.SetDefault(cfg.ManagerNewURL, "", "/info/manager/new")
	cfg.ManagerEditURL = utils.SetDefault(cfg.ManagerEditURL, "", "/info/manager/edit")
	cfg.ManagerNewURLPost = utils.SetDefault(cfg.ManagerNewURLPost, "", "/new/manager")
	cfg.RolesURL = utils.SetDefault(cfg.RolesURL, "", "/info/roles")
	cfg.RolesNewURL = utils.SetDefault(cfg.RolesNewURL, "", "/info/roles/new")
	cfg.PermissionURL = utils.SetDefault(cfg.PermissionURL, "", "/info/permission")
	cfg.PermissionNewURL = utils.SetDefault(cfg.PermissionNewURL, "", "/info/permission/new")

	// cookie時效
	if cfg.SessionLifeTime == 0 {
		// default two hours
		cfg.SessionLifeTime = 7200
	}
	return cfg
}

// -----ConfigService的Service方法-----start

// Name 為Service(interface)
func (c *Service) Name() string {
	return "config"
}

// -----ConfigService的Service方法的Service方法-----end

// Update 更新Config
func (c *Config) Update(m map[string]string) error {
	// 互斥鎖
	updateLock.Lock()
	defer updateLock.Unlock()

	c.Domain = m["domain"]
	c.Title = m["title"]
	c.Logo = template.HTML(m["logo"])
	c.MiniLogo = template.HTML(m["mini_logo"])
	c.LoginTitle = m["login_title"]
	return nil
}

// ToMap 設置Config(放在Site資料表中)
func (c *Config) ToMap() map[string]string {
	var m = make(map[string]string, 0)
	m["database"] = utils.JSON(c.Database)
	m["domain"] = c.Domain
	m["url_prefix"] = c.URLPrefix
	m["title"] = c.Title
	m["logo"] = string(c.Logo)
	m["mini_logo"] = string(c.MiniLogo)
	m["index_url"] = c.IndexURL
	m["login_url"] = c.LoginURL
	m["signup_url"] = c.SignupURL
	m["login_title"] = c.LoginTitle
	m["auth_user_table"] = c.AuthUserTable
	m["no_limit_login_ip"] = strconv.FormatBool(c.NoLimitLoginIP)
	m["login_logo"] = string(c.LoginLogo)
	m["session_life_time"] = strconv.Itoa(c.SessionLifeTime)
	m["menu_url"] = c.MenuURL
	m["menu_edit_url"] = c.MenuEditURL
	m["menu_delete_url"] = c.MenuDeleteURL
	m["menu_new_url"] = c.MenuNewURL
	m["manager_url"] = c.ManagerURL
	m["manager_new_url"] = c.ManagerNewURL
	m["manager_edit_url"] = c.ManagerEditURL
	m["manager_new_url_post"] = c.ManagerNewURLPost
	m["roles_url"] = c.RolesURL
	m["roles_new_url"] = c.RolesNewURL
	m["permission_url"] = c.PermissionURL
	m["permission_new_url"] = c.PermissionNewURL
	m["store"] = c.Store.JSON()
	return m
}

// AssertPrefix 處理config.prefix
func (c *Config) AssertPrefix() string {
	if c.prefix == "/" {
		return ""
	}
	return c.prefix
}

// ParamStr 設置Database.Params
func (d Database) ParamStr() {
	if d.Params == nil {
		d.Params = make(map[string]string)
	}
	if d.Driver == "mysql" {
		if _, ok := d.Params["charset"]; !ok {
			d.Params["charset"] = "utf8mb4"
		}
	}
}

// URL 處理URL
func (s Store) URL(suffix string) string {
	if len(suffix) > 4 && suffix[:4] == "http" {
		return suffix
	}
	if s.Prefix == "" {
		if suffix[0] == '/' {
			return suffix
		}
		return "/" + suffix
	}
	if s.Prefix[0] == '/' {
		if suffix[0] == '/' {
			return s.Prefix + suffix
		}
		return s.Prefix + "/" + suffix
	}
	if suffix[0] == '/' {
		if len(s.Prefix) > 4 && s.Prefix[:4] == "http" {
			return s.Prefix + suffix
		}
		return "/" + s.Prefix + suffix
	}
	if len(s.Prefix) > 4 && s.Prefix[:4] == "http" {
		return s.Prefix + "/" + suffix
	}
	return "/" + s.Prefix + "/" + suffix
}

// JSON 將Store(struct)JSON編碼
func (s Store) JSON() string {
	if s.Path == "" && s.Prefix == "" {
		return ""
	}
	return utils.JSON(s)
}

// Prefix globalCfg.prefix
func Prefix() string {
	return globalCfg.prefix
}

// GetDatabaseDriver 取得全局變數globalCfg的driver
func GetDatabaseDriver() string {
	return globalCfg.Database.Driver
}

// GetSessionLifeTime 取得session時間
func GetSessionLifeTime() int {
	return globalCfg.SessionLifeTime
}

// GetNoLimitLoginIP return globalCfg.NoLimitLoginIP
func GetNoLimitLoginIP() bool {
	return globalCfg.NoLimitLoginIP
}

// GetDomain return globalCfg.Domain
func GetDomain() string {
	return globalCfg.Domain
}

// GetLoginURL return globalCfg.LoginUrl
func GetLoginURL() string {
	return globalCfg.LoginURL
}

// IndexGetIndexURL return globalCfg.IndexURL
func IndexGetIndexURL() string {
	return globalCfg.IndexURL
}

// GetStore return globalCfg.Store
func GetStore() Store {
	return globalCfg.Store
}
