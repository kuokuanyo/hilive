package adapter

import (
	"hilive/context"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/db"
	"hilive/plugins"
	"net/url"
)

// WebFrameWork 功能都設定在框架中(使用/adapter/gin/gin.go框架)
type WebFrameWork interface {

	// Use 添加處理程序至api的路徑及方法中
	Use(app interface{}, plugins []plugins.Plugin) error

	// User 先透過cookie值(session)取得用戶id，接著判斷用戶角色、權限及可用菜單
	User(ctx interface{}) (models.UserModel, bool)

	// AddHandler 添加處理程序
	AddHandler(method, path string, handlers context.Handlers)

	DisableLog()

	Static(prefix, path string)

	Run() error

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
	Write(body []byte)
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

// GetUse 添加處理程序至api的路徑及方法中
func (base *BaseAdapter) GetUse(app interface{}, plugin []plugins.Plugin, wf WebFrameWork) error {
	// SetApp 將參數轉換成gin.Engine型態設置至Gin.app
	if err := wf.SetApp(app); err != nil {
		return err
	}

	// plugin is interface
	for _, plug := range plugin {
		// GetHandler 取得Base.App.Handlers(map[Path]Handlers)
		for path, handlers := range plug.GetHandler() {
			// AddHandler 添加處理程序
			wf.AddHandler(path.Method, path.URL, handlers)
		}
	}
	return nil
}

// GetUser 先透過cookie值(session)取得用戶id，接著判斷用戶角色、權限及可用菜單
func (base *BaseAdapter) GetUser(ctx interface{}, wf WebFrameWork) (models.UserModel, bool) {
	// SetContext 將參數轉換成gin.Context設置至Gin.ctx
	// 取得cookie
	cookie, err := wf.SetContext(ctx).GetCookie()
	if err != nil {
		return models.UserModel{}, false
	}
	// wf.GetConnection()回傳BaseAdapter.db(interface)
	// 透過cookie、conn可以得到角色、權限以及可使用菜單
	user, exist := auth.GetCurUser(cookie, wf.GetConnection())

	// 設置UserModel.Conn = nil後回傳UserModel
	return user.ReleaseConn(), exist
}
