package controller

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/response"
	"hilive/views/login"
	"hilive/views/signup"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Auth 登入平台POST功能
func (h *Handler) Auth(ctx *gin.Context) {
	var (
		user models.UserModel
	)

	password := ctx.Request.FormValue("password")
	phone := ctx.Request.FormValue("phone")
	if password == "" || phone == "" {
		response.BadRequest(ctx, "密碼或手機號碼不能為空")
	}

	// 檢查登入用戶資訊並取得角色權限菜單
	user, ok := auth.Check(password, phone, h.Conn)
	if !ok {
		response.BadRequest(ctx, "登入失敗")
		return
	}

	// 設置cookie並加入header
	err := auth.SetCookie(ctx, user, h.Conn)
	if err != nil {
		response.BadRequest(ctx, "設置cookie發生錯誤")
		return
	}

	if ref := ctx.GetHeader("Referer"); ref != "" {
		if u, err := url.Parse(ref); err == nil {
			v := u.Query()
			if r := v.Get("ref"); r != "" {
				rr, _ := url.QueryUnescape(r)
				response.OkWithData(ctx, map[string]interface{}{
					"url": rr,
				})
				return
			}
		}
	}
	response.OkWithData(ctx, map[string]interface{}{
		"url": "/admin" + h.Config.IndexURL,
	})
	return
}

// ShowLogin 登入GET功能
func (h *Handler) ShowLogin(ctx *gin.Context) {
	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(login.LoginTmpl)
	if err != nil {
		panic("使用登入的模板發生錯誤")
	}

	if err := tmpl.Execute(ctx.Writer, struct {
		URLPrefix string
		Title     string
		Logo      template.HTML
		CdnURL    string
	}{
		URLPrefix: h.Config.AssertPrefix(),
		Title:     h.Config.LoginTitle,
		Logo:      h.Config.LoginLogo,
		CdnURL:    h.Config.AssetURL,
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		ctx.Status(http.StatusOK)
		panic("使用登入的模板發生錯誤")
	}
}

// Signup 註冊POST功能
func (h *Handler) Signup(ctx *gin.Context) {
	username := ctx.Request.FormValue("username")
	phone := ctx.Request.FormValue("phone")
	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")
	checkPassword := ctx.Request.FormValue("checkPassword")

	if username == "" || phone == "" || email == "" || password == "" || checkPassword == "" {
		response.BadRequest(ctx, "使用者名稱、手機號碼、信箱、密碼都不能為空")
		return
	}
	if !strings.Contains(phone[:2], "09") && len(phone) != 10 {
		response.BadRequest(ctx, "手機號碼錯誤，ex:09...")
		return
	}
	if password != checkPassword {
		response.BadRequest(ctx, "密碼不一致")
		return
	}
	if !strings.Contains(email, "@gmail") {
		response.BadRequest(ctx, "必須使用gmail信箱註冊")
		return
	}

	userPhone := models.DefaultUserModel("users").SetConn(h.Conn).FindByPhone(phone)
	if userPhone.ID != int64(0) {
		response.BadRequest(ctx, "電話號碼已被註冊過")
		return
	}
	UerEmail := models.DefaultUserModel("users").SetConn(h.Conn).FindByEmail(email)
	if UerEmail.ID != int64(0) {
		response.BadRequest(ctx, "信箱已被註冊過")
		return
	}

	// 加密
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.BadRequest(ctx, "密碼加密發生錯誤")
		return
	}

	// 新增註冊資料並增加角色權限
	user, err := models.DefaultUserModel("users").SetConn(h.Conn).
		AddUser(username, phone, email, string(hash[:]))
	if err != nil {
		response.BadRequest(ctx, "增加會員資料發生錯誤")
		return
	}
	_, addRoleErr := user.SetConn(h.Conn).AddRole("1")
	if addRoleErr != nil {
		response.BadRequest(ctx, "新增角色發生錯誤")
		return
	}
	_, addPermissionErr := user.SetConn(h.Conn).AddPermission("1")
	if addPermissionErr != nil {
		response.BadRequest(ctx, "新增權限發生錯誤")
		return
	}

	response.OkWithData(ctx, map[string]interface{}{
		"url": "/admin" + h.Config.LoginURL,
	})
}

// ShowSignup 註冊用戶GET功能
func (h *Handler) ShowSignup(ctx *gin.Context) {
	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(signup.SignupTmpl)
	if err != nil {
		panic("使用註冊用戶模板發生錯誤")
	}

	if err := tmpl.Execute(ctx.Writer, struct {
		URLPrefix string
		Logo      template.HTML
		CdnURL    string
	}{
		URLPrefix: h.Config.AssertPrefix(),
		Logo:      h.Config.LoginLogo,
		CdnURL:    h.Config.AssetURL,
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		ctx.Status(http.StatusOK)
		panic("使用註冊用戶模板發生錯誤")
	}
}
