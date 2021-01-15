package context

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Context gin.Context的簡化版本
type Context struct {
	Request   *http.Request
	Response  *http.Response
	UserValue map[string]interface{}
	index     int8
	handlers  Handlers
}

// App struct
type App struct {
	Requests    []Path     //Path包含url、method
	Handlers    HandlerMap // HandlerMap是 map[Path]Handlers，Handlers類別為[]Handler，Handler類別為func(ctx *Context)
	Middlewares Handlers
	Prefix      string
	Routers     RouterMap // RouterMap類別為map[string]Router，Router(struct)裡有methods、patten
	routeIndex  int
	routeANY    bool
}

// Path 包含url、method
type Path struct {
	URL    string
	Method string
}

// Handler 用於middleware
type Handler func(ctx *Context)

// Handlers is []Handler
type Handlers []Handler

// HandlerMap is map[Path]Handlers
type HandlerMap map[Path]Handlers

// RouterMap is map[string]Router
type RouterMap map[string]Router

// Router 包含方法([]string)、模式(url)
type Router struct {
	Methods []string
	Patten  string //url
}

// RouterGroup struct
type RouterGroup struct {
	app         *App
	Middlewares Handlers // 中間件
	Prefix      string
}

// NewContext 預設Context(struct)
func NewContext(req *http.Request) *Context {
	return &Context{
		Request:   req,
		UserValue: make(map[string]interface{}),
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		},
		index: -1,
	}
}

// NewApp 預設App(struct)
func NewApp() *App {
	return &App{
		Requests:    make([]Path, 0),
		Handlers:    make(HandlerMap),
		Prefix:      "/",
		Middlewares: make([]Handler, 0),
		routeIndex:  -1,
		Routers:     make(RouterMap),
	}
}

// User 取得目前登入的用戶(Context.UserValue["user"])
func (ctx *Context) User() interface{} {
	return ctx.UserValue["user"]
}

// SetHandlers 設置Context.Handlers
func (ctx *Context) SetHandlers(handlers Handlers) *Context {
	ctx.handlers = handlers
	return ctx
}

// Query 取得url中的參數
func (ctx *Context) Query(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

// Group 設置RouterGroup(struct)
func (app *App) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         app,
		Middlewares: append(app.Middlewares, middleware...),
		Prefix:      slash(prefix),
	}
}

// AppendReqAndResp stores the request info and handle into app.
// support the route parameter. The route parameter will be recognized as
// wildcard store into the RegUrl of Path struct. For example:
//
//         /user/:id      => /user/(.*)
//         /user/:id/info => /user/(.*?)/info
//
// The RegUrl will be used to recognize the incoming path and find
// the handler.
// AppendReqAndResp 在RouterGroup.app(struct)中新增Requests([]Path)路徑及方法、接著在該url中新增參數handler(Handler...)
func (g *RouterGroup) AppendReqAndResp(url, method string, handler []Handler) {
	g.app.Requests = append(g.app.Requests, Path{
		URL:    join(g.Prefix, url),
		Method: method,
	})
	g.app.routeIndex++

	var h = make([]Handler, len(g.Middlewares))
	copy(h, g.Middlewares)

	g.app.Handlers[Path{
		URL:    join(g.Prefix, url),
		Method: method,
	}] = append(h, handler...)
}

// POST 等於在AppendReqAndResp(url, "post", handler)
func (g *RouterGroup) POST(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "post", handler)
	return g
}

// GET 等於在AppendReqAndResp(url, "get", handler)
func (g *RouterGroup) GET(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "get", handler)
	return g
}

// DELETE 等於在AppendReqAndResp(url, "delete", handler)
func (g *RouterGroup) DELETE(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "delete", handler)
	return g
}

// PUT 等於在AppendReqAndResp(url, "put", handler)
func (g *RouterGroup) PUT(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "put", handler)
	return g
}

// OPTIONS 等於在AppendReqAndResp(url, "options", handler)
func (g *RouterGroup) OPTIONS(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "options", handler)
	return g
}

// HEAD 等於在AppendReqAndResp(url, "head", handler)
func (g *RouterGroup) HEAD(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "head", handler)
	return g
}

// Write 將狀態碼、標頭(header)及body寫入Context.Response
func (ctx *Context) Write(code int, header map[string]string, Body string) {
	ctx.Response.StatusCode = code
	for key, head := range header {
		ctx.AddHeader(key, head)
	}
	ctx.Response.Body = ioutil.NopCloser(strings.NewReader(Body))
}

// WriteString 將參數body保存至Context.response.Body中
func (ctx *Context) WriteString(body string) {
	ctx.Response.Body = ioutil.NopCloser(strings.NewReader(body))
}

// HTML 輸出HTML
func (ctx *Context) HTML(code int, body string) {
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetStatusCode(code)
	ctx.WriteString(body)
}

// JSON 轉換成JSON
func (ctx *Context) JSON(code int, Body map[string]interface{}) {
	ctx.Response.StatusCode = code
	ctx.SetContentType("application/json")
	// Marshal將struct轉成json
	BodyStr, err := json.Marshal(Body)
	if err != nil {
		panic(err)
	}
	ctx.Response.Body = ioutil.NopCloser(bytes.NewReader(BodyStr))
}

// Group 設置中間件(驗證)
func (g *RouterGroup) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         g.app,
		Middlewares: append(g.Middlewares, middleware...),
		Prefix:      join(slash(g.Prefix), slash(prefix)),
	}
}

// SetUserValue 將參數設定至Context.UserValue
func (ctx *Context) SetUserValue(key string, value interface{}) {
	ctx.UserValue[key] = value
}

// AddHeader 將參數添加header中(Context.Response.Header)
func (ctx *Context) AddHeader(key, value string) {
	ctx.Response.Header.Add(key, value)
}

// Headers 透過參數key獲得Header
func (ctx *Context) Headers(key string) string {
	return ctx.Request.Header.Get(key)
}

// SetCookie 設置cookie在response header Set-Cookie中
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		ctx.AddHeader("Set-Cookie", v)
	}
}

// Path return Request.url.path
func (ctx *Context) Path() string {
	return ctx.Request.URL.Path
}

// Method return ctx.Request.Method
func (ctx *Context) Method() string {
	return ctx.Request.Method
}

// Abort abort the context.
func (ctx *Context) Abort() {
	ctx.index = 63
}

// Next 執行迴圈Context.handlers[ctx.index](ctx)
func (ctx *Context) Next() {
	ctx.index++
	// Context.Handlers類別為[]Handler，Handler類別為func(ctx *Context)
	for s := int8(len(ctx.handlers)); ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// IsPjax 判斷header X-PJAX:true
func (ctx *Context) IsPjax() bool {
	return ctx.Headers("X-PJAX") == "true"
}

// SetContentType 將參數添加至Content-Type
func (ctx *Context) SetContentType(contentType string) {
	ctx.AddHeader("Content-Type", contentType)
}

// SetStatusCode 將參數設置至Context.Response.StatusCode
func (ctx *Context) SetStatusCode(code int) {
	ctx.Response.StatusCode = code
}

// slash fix the path which has wrong format problem.
//
// 	 ""      => "/"
// 	 "abc/"  => "/abc"
// 	 "/abc/" => "/abc"
// 	 "/abc"  => "/abc"
// 	 "/"     => "/"
//
// slash 處理斜線(路徑)
func slash(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix == "/" {
		return "/"
	}
	if prefix[0] != '/' {
		if prefix[len(prefix)-1] == '/' {
			return "/" + prefix[:len(prefix)-1]
		}
		return "/" + prefix
	}
	if prefix[len(prefix)-1] == '/' {
		return prefix[:len(prefix)-1]
	}
	return prefix
}

// join 路徑
func join(prefix, suffix string) string {
	if prefix == "/" {
		return suffix
	}
	if suffix == "/" {
		return prefix
	}
	return prefix + suffix
}
