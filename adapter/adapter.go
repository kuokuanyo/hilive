package adapter

import (
	"hilive/modules/db"
	"net/url"
)

// WebFrameWork 功能都設定在框架中(使用/adapter/gin/gin.go框架)
type WebFrameWork interface {

	// 輔助功能
	SetApp(app interface{}) error
	SetContext(ctx interface{}) WebFrameWork
	SetConnection(db.Connection)
	GetConnection() db.Connection
	GetCookie() (string, error)
	Path() string
	Method() string
	FormParam() url.Values
	IsPjax() bool
	Redirect()
	SetContentType()
	CookieKey() string
	HTMLContentType() string
}

// BaseAdapter 是db.Connection(interface)
type BaseAdapter struct {
	db db.Connection
}

// -----下面為WebFrameWork方法-----start

// SetConnection 將參數(conn)設置至BaseAdapter.db
func (base *BaseAdapter) SetConnection(conn db.Connection) {
	base.db = conn
}

// GetConnection 回傳BaseAdapter.db
func (base *BaseAdapter) GetConnection() db.Connection {
	return base.db
}

// CookieKey return "session"
func (base *BaseAdapter) CookieKey() string {
	return "session"
}

// HTMLContentType return "text/html; charset=utf-8"
func (base *BaseAdapter) HTMLContentType() string {
	return "text/html; charset=utf-8"
}

// -----下面為WebFrameWork方法-----end
