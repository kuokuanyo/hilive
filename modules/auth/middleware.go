package auth

import (
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	"net/http"
	"net/url"
)

// MiddlewareCallback is type of callback function.
type MiddlewareCallback func(ctx *context.Context)

// Invoker 中間件驗證
type Invoker struct {
	prefix                 string
	authFailCallback       MiddlewareCallback
	permissionDenyCallback MiddlewareCallback
	conn                   db.Connection
}

// DefaultInvoker 預設Invoker
func DefaultInvoker(conn db.Connection) *Invoker {
	return &Invoker{
		prefix: config.Prefix(),
		authFailCallback: func(ctx *context.Context) {
			if ctx.Request.URL.Path == "/admin"+config.GetLoginURL() {
				return
			}
			if ctx.Request.URL.Path == "/admin/logout" {
				ctx.Write(302, map[string]string{
					"Location": config.Prefix() + config.GetLoginURL(),
				}, ``)
				return
			}
			param := ""
			if ref := ctx.Headers("Referer"); ref != "" {
				param = "?ref=" + url.QueryEscape(ref)
			}

			u := "/admin" + config.GetLoginURL() + param
			_, err := ctx.Request.Cookie("session")
			referer := ctx.Headers("Referer")

			if (ctx.Headers("X-PJAX") == "" && ctx.Request.Method != "GET") ||
				err != nil || referer == "" {
				ctx.Write(302, map[string]string{
					"Location": u,
				}, ``)
			} else {
				msg := "登入逾時，請重新登入"
				h := `<script>
				if (typeof(swal) === "function") {
					swal({
						type: "info",
						title: "login info",
						text: "` + msg + `",
						showCancelButton: false,
						confirmButtonColor: "#3c8dbc",
						confirmButtonText: '` + "got it" + `',
					})
					setTimeout(function(){ location.href = "` + u + `"; }, 2000);
				} else {
					alert("` + msg + `")
					location.href = "` + u + `"
				}
			</script>`
				ctx.HTML(http.StatusOK, h)
			}
		},
		permissionDenyCallback: func(ctx *context.Context) {
			if ctx.Headers("X-PJAX") == "" && ctx.Request.Method != "GET" {
				ctx.JSON(http.StatusForbidden, map[string]interface{}{
					"code": http.StatusForbidden,
					"msg":  "permission denied",
				})
			} else {
				h := `<div class="missing-content">
				<div class="missing-content-title">403</div>
				<div class="missing-content-title-subtitle">Sorry, you don't have access to this page.</div>
			</div>
			
			<style>
			.missing-content {
				padding: 48px 32px;
			}
			.missing-content-title {
				color: rgba(0,0,0,.85);
				font-size: 54px;
				line-height: 1.8;
				text-align: center;
			}
			.missing-content-title-subtitle {
				color: rgba(0,0,0,.45);
				font-size: 18px;
				line-height: 1.6;
				text-align: center;
			}
			</style>`
				ctx.HTML(http.StatusForbidden, h)
			}
		},
		conn: conn,
	}
}

// Middleware 驗證，判斷用戶是否有權限
func (invoker *Invoker) Middleware(conn db.Connection) context.Handler {
	return func(ctx *context.Context) {
		user, authOk, permissionOk := Filter(ctx, conn)
		if authOk && permissionOk {
			ctx.SetUserValue("user", user)
			ctx.Next()
			return
		}
		if !authOk {
			invoker.authFailCallback(ctx)
			ctx.Abort()
			return
		}
		if !permissionOk {
			invoker.permissionDenyCallback(ctx)
			ctx.Abort()
			return
		}
	}
}

// Filter 透過用戶id取得角色權限菜單，並判斷用戶使否有權限訪問該頁面
func Filter(ctx *context.Context, conn db.Connection) (models.UserModel, bool, bool) {
	var (
		user = models.UserModel{Base: models.Base{TableName: "users"}}
		id   float64
		ok   bool
	)

	// InitSession 初始化Session並取得session資料表的cookie_values欄位(ex:{"user_id":1})
	ses, err := InitSession(ctx, conn)
	if err != nil {
		return user, false, false
	}

	// 取得session資料表cookie_values欄位值
	if id, ok = ses.Values["user_id"].(float64); !ok {
		return user, false, false
	}

	// 取得用戶角色權限菜單，以及是否有可以訪問的menu
	user, ok = GetUserByID(int64(id), conn)
	if !ok {
		return user, false, false
	}
	return user, true,
		user.CheckPermissionByURLMethod(ctx.Request.URL.String(), ctx.Request.Method, ctx.Request.PostForm)
}

// GetUserByID 透過id取得用戶角色權限菜單，以及是否有可以訪問的menu
func GetUserByID(id int64, conn db.Connection) (user models.UserModel, ok bool) {
	var superAdmin bool
	user = models.DefaultUserModel().SetConn(conn).FindByID(id)
	if user.ID == int64(0) {
		ok = false
		return
	}

	// 取得角色權限菜單
	user = user.GetUserRoles().GetUserPermissions().GetUserMenus()

	// 判斷是否為超級管理員
	for _, permission := range user.Permissions {
		if len(permission.HTTPPath) > 0 && permission.HTTPPath[0] == "*" && permission.HTTPMethod[0] == "" {
			superAdmin = true
			break
		}
		superAdmin = false
	}
	if len(user.MenuIDs) != 0 || superAdmin {
		ok = true
	}
	return
}

// GetCurUser 先透過cookie值(session)取得用戶id，接著判斷用戶角色、權限及可用菜單
func GetCurUser(sesKey string, conn db.Connection) (user models.UserModel, ok bool) {
	if sesKey == "" {
		ok = false
		return
	}

	// 取得session資料表的cookie_values[key]的值(id)，如果沒有則回傳-1
	id := GetUserID(sesKey, conn)
	if id == -1 {
		ok = false
		return
	}
	// GetCurUserByID 透過參數(id)取得role、permission及可用menu
	return GetCurUserByID(id, conn)
}

// GetCurUserByID 透過參數(id)取得role、permission及可用menu
func GetCurUserByID(id int64, conn db.Connection) (user models.UserModel, ok bool) {
	// 透過參數(id)取得UserModel(struct)
	user = models.DefaultUserModel().SetConn(conn).FindByID(id)
	if user.ID == int64(0) {
		ok = false
		return
	}

	// 取得角色、權限及可使用菜單
	user = user.GetUserRoles().GetUserPermissions().GetUserMenus()
	// 檢查用戶是否有可訪問的menu
	ok = len(user.MenuIDs) != 0 || user.IsSuperAdmin()
	return
}

// GetUserID 取得session資料表的cookie_values[key]的值(id)，如果沒有則回傳-1
func GetUserID(sesKey string, conn db.Connection) int64 {
	// GetSessionByKey 取得session資料表的cookie_values[key]的值(id)
	id, err := GetSessionByKey(sesKey, "user_id", conn)
	if err != nil {
		return -1
	}
	if idFloat64, ok := id.(float64); ok {
		return int64(idFloat64)
	}
	return -1
}
