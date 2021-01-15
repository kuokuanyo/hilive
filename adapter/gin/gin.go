package gin

import (
	"bytes"
	"errors"
	"hilive/adapter"
	"hilive/context"
	"hilive/engine"
	"hilive/models"
	"hilive/modules/config"
	"hilive/plugins"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Gin 也符合adapter.WebFrameWork(interface)
type Gin struct {
	adapter.BaseAdapter
	// gin.Context(struct)為gin最重要的部分，允許在middleware傳遞變數(例如驗證請求、管理流程)
	ctx *gin.Context
	// app為框架中的實例，包含muxer,middleware ,configuration，藉由New() or Default()建立Engine
	app *gin.Engine
}

// init 建立引擎預設的配適器
func init() {
	engine.Register(new(Gin))
}

// -----下面為WebFrameWork方法-----start

// Name 回傳框架名稱，同時也是service(interface)
func (gins *Gin) Name() string {
	return "gin"
}

// Use 添加處理程序至api的路徑及方法中
func (gins *Gin) Use(app interface{}, plugs []plugins.Plugin) error {
	return gins.GetUse(app, plugs, gins)
}

// AddHandler 增加處理程序
func (gins *Gin) AddHandler(method, path string, handlers context.Handlers) {

	// Handle第三個參數為主要處理程序
	gins.app.Handle(strings.ToUpper(method), path, func(c *gin.Context) {
		// 預設Context(struct)
		ctx := context.NewContext(c.Request)

		// 將參數設置在url中
		for _, param := range c.Params {
			if c.Request.URL.RawQuery == "" {
				c.Request.URL.RawQuery += strings.Replace(param.Key, ":", "", -1) + "=" + param.Value
			} else {
				c.Request.URL.RawQuery += "&" + strings.Replace(param.Key, ":", "", -1) + "=" + param.Value
			}
		}

		// SetHandlers將參數(handlers)設置至至Context.Handlers
		// 執行迴圈Context.handlers[ctx.index](ctx)
		ctx.SetHandlers(handlers).Next()

		for key, head := range ctx.Response.Header {
			c.Header(key, head[0])
		}

		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			c.String(ctx.Response.StatusCode, buf.String())
		} else {
			c.Status(ctx.Response.StatusCode)
		}
	})
}

// User 先透過cookie值(session)取得用戶id，接著判斷用戶角色、權限及可用菜單
func (gins *Gin) User(ctx interface{}) (models.UserModel, bool) {
	return gins.GetUser(ctx, gins)
}

// SetApp 將參數轉換成gin.Engine型態設置至Gin.app
func (gins *Gin) SetApp(app interface{}) error {
	var (
		eng *gin.Engine
		ok  bool
	)
	// app.(*gin.Engine)將interface{}轉換為gin.Engine型態
	if eng, ok = app.(*gin.Engine); !ok {
		return errors.New("gin adapter SetApp: wrong parameter")
	}
	gins.app = eng
	return nil
}

// SetContext 將參數轉換成gin.Context設置至Gin.ctx
func (gins *Gin) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *gin.Context
		ok  bool
	)
	// 將contextInterface類別變成gin.Context(struct)
	if ctx, ok = contextInterface.(*gin.Context); !ok {
		panic("gin adapter SetContext: wrong parameter")
	}
	return &Gin{ctx: ctx}
}

// GetCookie 取得session裡設置的cookie
func (gins *Gin) GetCookie() (string, error) {
	// Cookie()回傳cookie(藉由參數裡的命名回傳的)
	return gins.ctx.Cookie(gins.CookieKey())
}

// Path return Gin.ctx.Request.URL.Path
func (gins *Gin) Path() string {
	return gins.ctx.Request.URL.Path
}

// Method return gins..ctx.Request.Method
func (gins *Gin) Method() string {
	return gins.ctx.Request.Method
}

// FormParam 解析參數(multipart/form-data裡的)
func (gins *Gin) FormParam() url.Values {
	_ = gins.ctx.Request.ParseMultipartForm(32 << 20)
	return gins.ctx.Request.PostForm
}

// IsPjax 設置標頭 X-PJAX = true
func (gins *Gin) IsPjax() bool {
	return gins.ctx.Request.Header.Get("X-PJAX") == "true"
}

// SetContentType return
func (gins *Gin) SetContentType() {
	return
}

// Redirect 重新導向至登入頁面(出現錯誤)
func (gins *Gin) Redirect() {
	gins.ctx.Redirect(302, config.Prefix()+config.GetLoginURL())
	gins.ctx.Abort()
}

// Write 將參數(body)寫入
func (gins *Gin) Write(body []byte) {
	gins.ctx.Data(http.StatusOK, gins.HTMLContentType(), body)
}

// DisableLog return panic
func (gins *Gin) DisableLog() { panic("not implement") }

// Static return panic
func (gins *Gin) Static(prefix, path string) { panic("not implement") }

// Run return panic
func (gins *Gin) Run() error { panic("not implement") }

// -----下面為WebFrameWork方法-----end
