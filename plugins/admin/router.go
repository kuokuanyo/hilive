package admin

import (
	"hilive/context"
	"hilive/modules/auth"
)

// initRouter 設置api功能路徑
func (admin *Admin) initRouter() *Admin {
	app := context.NewApp()
	route := app.Group("/admin")

	// 登入GET、POST
	route.GET(admin.config.LoginURL, admin.handler.ShowLogin)
	route.POST(admin.config.LoginURL, admin.handler.Auth)
	// 註冊GET、POST
	route.GET(admin.config.SignupURL, admin.handler.ShowSignup)
	route.POST(admin.config.SignupURL, admin.handler.Signup)

	authRoute := route.Group("/", auth.DefaultInvoker(admin.Base.Conn).Middleware(admin.Base.Conn))
	// 菜單
	authRoute.GET(admin.config.MenuURL, admin.handler.ShowMenu)
	authRoute.GET(admin.config.MenuNewURL, admin.handler.ShowNewMenu)
	authRoute.GET(admin.config.MenuEditURL, admin.handler.ShowEditMenu)
	authRoute.POST(admin.config.MenuNewURL, admin.guard.MenuNew, admin.handler.NewMenu)
	authRoute.POST(admin.config.MenuEditURL, admin.guard.MenuEdit, admin.handler.EditMenu)
	authRoute.POST(admin.config.MenuDeleteURL, admin.guard.MenuDelete, admin.handler.DeleteMenu)

	authPrefixRoute := route.Group("/", auth.DefaultInvoker(admin.Base.Conn).Middleware(admin.Base.Conn), admin.guard.CheckPrefix)
	// 使用者、角色、權限
	authPrefixRoute.GET("/info/:__prefix", admin.handler.ShowInfo)
	authPrefixRoute.GET("/info/:__prefix/new", admin.guard.ShowNewForm, admin.handler.ShowNewForm)
	authPrefixRoute.GET("/info/:__prefix/edit", admin.guard.ShowEditForm, admin.handler.ShowEditForm)
	authPrefixRoute.POST("/new/:__prefix", admin.guard.NewForm, admin.handler.NewForm)
	authPrefixRoute.POST("/edit/:__prefix", admin.guard.EditForm, admin.handler.EditForm)
	authPrefixRoute.POST("/delete/:__prefix", admin.guard.Delete, admin.handler.Delete)

	admin.App = app
	return admin
}
