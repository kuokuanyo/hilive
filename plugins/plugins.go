package plugins

import (
	"hilive/context"
	"hilive/modules/db"
	"hilive/modules/service"
)

// Base 包含Plugin(interface)所有方法
type Base struct {
	App       *context.App
	Services  service.List
	Conn      db.Connection
	PlugName  string
	URLPrefix string
}

// Plugin 設置及配置路徑...等方法
type Plugin interface {
	GetHandler() context.HandlerMap
	InitPlugin(services service.List)
	Name() string
	Prefix() string
}

// -----plugin的所有方法-----start

// GetHandler 取得Base.App.Handlers(map[Path]Handlers)
func (b *Base) GetHandler() context.HandlerMap {
	return b.App.Handlers
}

// Name 也屬於service方法
func (b *Base) Name() string {
	return b.PlugName
}

// Prefix return Base.URLPrefix
func (b *Base) Prefix() string {
	return b.URLPrefix
}

// -----plugin的所有方法-----end

// InitBase 設置Base.Conn、Base.Services
func (b *Base) InitBase(srv service.List) {
	b.Services = srv
	b.Conn = db.GetConnectionFromService(b.Services)
}
