package auth

import (
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/service"
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// CSRFToken is type of a csrf token list.
type CSRFToken []string

// TokenService struct
type TokenService struct {
	Tokens CSRFToken //[]string
	lock   sync.Mutex
}

// ***************初始化*******************
// 將token_csrf_helper加入services(map[string]Generator)
func init() {
	service.Register("token_csrf_helper", func() (service.Service, error) {
		return &TokenService{
			Tokens: make(CSRFToken, 0),
		}, nil
	})
}

// Auth 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
func Auth(ctx *context.Context) models.UserModel {
	return ctx.User().(models.UserModel)
}

// GetTokenServiceByService 藉由Service取得TokenService
func GetTokenServiceByService(s interface{}) *TokenService {
	if srv, ok := s.(*TokenService); ok {
		return srv
	}
	panic("錯誤的Service")
}

// CSRFToken的Service方法-----start

// Name Service方法
func (s *TokenService) Name() string {
	return "token_csrf_helper"
}

// AddToken 新增Token
func (s *TokenService) AddToken() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	tokenStr, err := uuid.NewV4()
	if err != nil {
		panic("產生uuid發生錯誤")
	}
	s.Tokens = append(s.Tokens, tokenStr.String())
	return tokenStr.String()
}

// CSRFToken的Service方法-----end

// ConvertInterfaceToTokenService 將interface轉換TokenService
func ConvertInterfaceToTokenService(s interface{}) *TokenService {
	if srv, ok := s.(*TokenService); ok {
		return srv
	}
	panic("interface轉換TokenService發生錯誤")
}

// Check 檢查登入資訊並取得用戶角色權限可用menu
func Check(password string, phone string, conn db.Connection) (user models.UserModel, ok bool) {
	user = models.DefaultUserModel().SetConn(conn).FindByPhone(phone)
	// 檢查是否為空
	if user.ID == int64(0) {
		ok = false
	} else {
		if comparePassword(password, user.Password) {
			ok = true
			// 取得role、permission、menu
			user = user.GetUserRoles().GetUserPermissions().GetUserMenus()
			// 更新密碼
			user.UpdatePassword(EncodePassword([]byte(password)))
		} else {
			ok = false
		}
	}
	return
}

// CheckToken 檢查是否存在token
func (s *TokenService) CheckToken(CheckToken string) bool {
	for i := 0; i < len(s.Tokens); i++ {
		if s.Tokens[i] == CheckToken {
			s.Tokens = append((s.Tokens)[:i], (s.Tokens)[i+1:]...)
			return true
		}
	}
	return false
}

// SetCookie 設置cookie並加入header
func SetCookie(ctx *context.Context, user models.UserModel, conn db.Connection) error {
	ses, err := InitSession(ctx, conn)
	if err != nil {
		return err
	}
	ses.Values["user_id"] = user.ID
	if err := ses.Driver.Update(ses.Sid, ses.Values); err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:  ses.Cookie,
		Value: ses.Sid,
		// 回傳globalCfg.SessionLifeTime
		MaxAge: config.GetSessionLifeTime(),
		// cookie存在時間
		Expires:  time.Now().Add(ses.Expires),
		HttpOnly: true,
		Path:     "/",
	}
	if config.GetDomain() != "" {
		cookie.Domain = config.GetDomain()
	}

	// 設置cookie(struct)在response header Set-Cookie中
	ses.Context.SetCookie(&cookie)
	return nil
}

// comparePassword 檢查密碼是否相符
func comparePassword(comPwd, pwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(comPwd))
	return err == nil
}

// EncodePassword 加密
func EncodePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash[:])
}
