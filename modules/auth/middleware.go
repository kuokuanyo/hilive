package auth

import (
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	"net/http"
	"net/url"
	"text/template"

	"github.com/gin-gonic/gin"
)

var (
	// UserValue 紀錄UserModel
	userValue = make(map[string]interface{})
)

// Invoker 中間件驗證
type Invoker struct {
	prefix                 string
	authFailCallback       gin.HandlerFunc
	permissionDenyCallback gin.HandlerFunc
	conn                   db.Connection
}

// GetUserByMiddleware 取得middleware驗證後的user
func GetUserByMiddleware() models.UserModel {
	return userValue["user"].(models.UserModel)
}

// DefaultInvoker 預設Invoker
func DefaultInvoker(conn db.Connection) *Invoker {
	return &Invoker{
		prefix: config.Prefix(),
		authFailCallback: func(ctx *gin.Context) {
			if ctx.Request.URL.Path == "/admin"+config.GetLoginURL() {
				return
			}
			if ctx.Request.URL.Path == "/admin/logout" {
				_, err := template.New("").Parse(``)
				if err != nil {
					panic("模板發生錯誤")
				}
				ctx.Header("Location", "/admin"+config.GetLoginURL())
				ctx.Status(http.StatusFound)
				return
			}
			param := ""
			if ref := ctx.GetHeader("Referer"); ref != "" {
				param = "?ref=" + url.QueryEscape(ref)
			}

			u := "/admin" + config.GetLoginURL() + param
			_, err := ctx.Request.Cookie("session")
			referer := ctx.GetHeader("Referer")

			if (ctx.GetHeader("X-PJAX") == "" && ctx.Request.Method != "GET") ||
				err != nil || referer == "" {
				_, err := template.New("").Parse(``)
				if err != nil {
					panic("模板發生錯誤")
				}
				ctx.Header("Location", u)
				ctx.Status(http.StatusFound)
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
				tmpl, err := template.New("").Parse(h)
				if err != nil {
					panic("模板發生錯誤")
				}
				tmpl.Execute(ctx.Writer, nil)
				ctx.Status(http.StatusOK)
			}
		},
		permissionDenyCallback: func(ctx *gin.Context) {
			if ctx.GetHeader("X-PJAX") == "" && ctx.Request.Method != "GET" {
				ctx.JSON(http.StatusForbidden, map[string]interface{}{
					"code": http.StatusForbidden,
					"msg":  "permission denied",
				})
			} else {
				tmpl, err := template.New("").Parse(`<div class="missing-content">
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
			</style>`)
				if err != nil {
					panic("模板發生錯誤")
				}
				tmpl.Execute(ctx.Writer, nil)
				ctx.Status(http.StatusForbidden)
			}
		},
		conn: conn,
	}
}

// Middleware 驗證，判斷用戶是否有權限
func (invoker *Invoker) Middleware(conn db.Connection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, authOk, permissionOk := Filter(ctx, conn)
		if authOk && permissionOk {
			userValue["user"] = user
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
func Filter(ctx *gin.Context, conn db.Connection) (models.UserModel, bool, bool) {
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
