package models

import (
	dbsql "database/sql"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// UserModel user資料表欄位
type UserModel struct {
	Base `json:"-"`

	ID          int64             `json:"id"`
	UserID      string            `json:"userid"`
	UserName    string            `json:"user_name"`
	Phone       string            `json:"phone"`
	Email       string            `json:"email"`
	Password    string            `json:"password"`
	Picture     string            `json:"Picture"`
	Permissions []PermissionModel `json:"permissions"`
	MenuIDs     []int64           `json:"menu_ids"`
	Roles       []RoleModel       `json:"role"`
	Level       string            `json:"level"`
	LevelName   string            `json:"level_name"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// DefaultUserModel 預設UserModel
func DefaultUserModel(tablename string) UserModel {
	return UserModel{Base: Base{TableName: tablename}}
}

// GetUserModelAndID 設置UserModel與ID
func GetUserModelAndID(id, tablename string) UserModel {
	idInt, _ := strconv.Atoi(id)
	return UserModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (user UserModel) SetConn(conn db.Connection) UserModel {
	user.Conn = conn
	return user
}

// SetTx 設置Tx
func (user UserModel) SetTx(tx *dbsql.Tx) UserModel {
	user.Base.Tx = tx
	return user
}

// FindByID 透過ID尋找資料
func (user UserModel) FindByID(id interface{}) UserModel {
	item, _ := user.Table(user.Base.TableName).Where("id", "=", id).First()
	return user.MapToUserModel(item)
}

// FindByPhone 透過電話號碼查詢資料
func (user UserModel) FindByPhone(phone interface{}) UserModel {
	item, _ := user.Base.Table(user.Base.TableName).Where("phone", "=", phone).First()
	return user.MapToUserModel(item)
}

// FindByEmail 透過信箱查詢資料
func (user UserModel) FindByEmail(email interface{}) UserModel {
	item, _ := user.Base.Table(user.Base.TableName).
		Where("email", "=", email).First()
	return user.MapToUserModel(item)
}

// Update 更新用戶資料
func (user UserModel) Update(username, phone, email, password string) (int64, error) {
	fieldValues := sql.Value{
		"username":   username,
		"phone":      phone,
		"email":      email,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}

	if password != "" {
		fieldValues["password"] = password
	}
	return user.SetTx(user.Tx).Table(user.Base.TableName).
		Where("id", "=", user.ID).Update(fieldValues)
}

// DeleteRolesByID 刪除該ID用戶所有角色
func (user UserModel) DeleteRolesByID() error {
	return user.Base.Table("role_users").Where("user_id", "=", user.ID).Delete()
}

// DeletePermissionsByID 刪除該ID用戶所有權限
func (user UserModel) DeletePermissionsByID() error {
	return user.SetTx(user.Tx).Table("user_permissions").
		Where("user_id", "=", user.ID).Delete()
}

// AddUser 增加會員資料
func (user UserModel) AddUser(userid, username, phone, email, password string) (UserModel, error) {
	id, err := user.SetTx(user.Base.Tx).Table(user.TableName).Insert(sql.Value{
		"userid":   userid,
		"username": username,
		"phone":    phone,
		"email":    email,
		"password": password,
	})

	user.ID = id
	user.UserID = userid
	user.UserName = username
	user.Phone = phone
	user.Email = email
	user.Password = password
	return user, err
}

// AddRole 增加角色
func (user UserModel) AddRole(id string) (int64, error) {
	// 檢查是否設置過角色
	checkRole, _ := user.Table("role_users").
		Where("role_id", "=", id).
		Where("user_id", "=", user.ID).First()
	if id != "" {
		if checkRole == nil {
			user.SetTx(user.Base.Tx).Table("role_users").
				Insert(sql.Value{
					"role_id": id,
					"user_id": user.ID,
				})
		}
	}
	return 0, nil
}

// AddPermission 增加權限
func (user UserModel) AddPermission(id string) (int64, error) {
	checkPermission, _ := user.Table("user_permissions").
		Where("permission_id", "=", id).
		Where("user_id", "=", user.ID).First()
	if id != "" {
		if checkPermission == nil {
			user.SetTx(user.Base.Tx).Table("user_permissions").
				Insert(sql.Value{
					"permission_id": id,
					"user_id":       user.ID,
				})
		}
	}
	return 0, nil
}

// MapToUserModel 將值設置至UserModel
func (user UserModel) MapToUserModel(m map[string]interface{}) UserModel {
	user.ID, _ = m["id"].(int64)
	user.UserID, _ = m["userid"].(string)
	user.UserName, _ = m["username"].(string)
	user.Phone, _ = m["phone"].(string)
	user.Email, _ = m["email"].(string)
	user.Password, _ = m["password"].(string)
	user.Picture, _ = m["picture"].(string)
	user.CreatedAt, _ = m["created_at"].(string)
	user.UpdatedAt, _ = m["updated_at"].(string)
	return user
}

// GetUserRoles 取得用戶角色
func (user UserModel) GetUserRoles() UserModel {
	roleModel, _ := user.Base.Table("role_users").
		LeftJoin("roles", "roles.id", "=", "role_users.role_id").
		Where("user_id", "=", user.ID).
		Select("roles.id", "roles.name", "roles.slug",
			"roles.created_at", "roles.updated_at").All()

	for _, role := range roleModel {
		user.Roles = append(user.Roles, DefaultRoleModel().MapToRoleModel(role))
	}

	if len(user.Roles) > 0 {
		user.Level = user.Roles[0].Slug
		user.LevelName = user.Roles[0].Name
	}

	return user
}

// GetUserPermissions 取得用戶權限
func (user UserModel) GetUserPermissions() UserModel {
	var permissions = make([]map[string]interface{}, 0)

	roleIDs := user.GetAllRoleID()
	// *****
	// permission會依照user_id以及role_id取得不同的權限，因此需要做下列兩次判斷
	// *****
	// 利用role取得的權限
	if len(roleIDs) > 0 {
		permissions, _ = user.Base.Table("role_permissions").
			LeftJoin("permissions", "permissions.id", "=", "role_permissions.permission_id").
			WhereIn("role_id", roleIDs).
			Select("permissions.http_method", "permissions.http_path",
				"permissions.id", "permissions.name", "permissions.slug",
				"permissions.created_at", "permissions.updated_at").All()
	}
	// 利用user取得的權限
	userPermissions, _ := user.Base.Table("user_permissions").
		LeftJoin("permissions", "permissions.id", "=", "user_permissions.permission_id").
		Where("user_id", "=", user.ID).
		Select("permissions.http_method", "permissions.http_path",
			"permissions.id", "permissions.name", "permissions.slug",
			"permissions.created_at", "permissions.updated_at").All()

	permissions = append(permissions, userPermissions...)

	// 加入權限
	for i := 0; i < len(permissions); i++ {
		exist := false
		for j := 0; j < len(user.Permissions); j++ {
			if user.Permissions[j].ID == permissions[i]["id"] {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		user.Permissions = append(user.Permissions,
			DefaultPermissionModel().MapToPermissionModel(permissions[i]))
	}
	return user
}

// GetUserMenus 取得用戶可用menu
func (user UserModel) GetUserMenus() UserModel {
	var (
		menuIDsModel []map[string]interface{}
	)

	// 判斷是否為超級管理員
	if user.IsSuperAdmin() {
		menuIDsModel, _ = user.Base.Table("role_menu").
			LeftJoin("menu", "menu.id", "=", "role_menu.menu_id").
			Select("menu_id", "parent_id").All()
	} else {
		rolesID := user.GetAllRoleID()
		if len(rolesID) > 0 {
			menuIDsModel, _ = user.Base.Table("role_menu").
				LeftJoin("menu", "menu.id", "=", "role_menu.menu_id").
				WhereIn("role_menu.role_id", rolesID).
				Select("menu_id", "parent_id").All()
		}
	}

	var menuIDs []int64
	// 加入menu_id
	for _, mid := range menuIDsModel {
		menuIDs = append(menuIDs, mid["menu_id"].(int64))
	}
	user.MenuIDs = menuIDs
	return user
}

// GetAllRoleID 取得用戶所有角色的role_id
func (user UserModel) GetAllRoleID() []interface{} {
	var ids = make([]interface{}, len(user.Roles))

	for key, role := range user.Roles {
		ids[key] = role.ID
	}
	return ids
}

// UpdatePassword 更新密碼
func (user UserModel) UpdatePassword(password string) UserModel {
	user.Base.Table(user.Base.TableName).
		Where("id", "=", user.ID).
		Update(sql.Value{
			"password": password,
		})
	user.Password = password
	return user
}

// GetCheckPermissionByURLMethod 檢查用戶權限，如果沒有權限回傳""
func (user UserModel) GetCheckPermissionByURLMethod(path, method string) string {
	// 檢查權限
	if !user.CheckPermissionByURLMethod(path, method, url.Values{}) {
		return ""
	}
	return path
}

// CheckPermissionByURLMethod 透過url、method判斷用戶是否有權限訪問該頁面
func (user UserModel) CheckPermissionByURLMethod(path, method string, formParams url.Values) bool {
	// 判斷是否為超級管理員
	if user.IsSuperAdmin() {
		return true
	}

	logoutCheck, _ := regexp.Compile("/admin/logout" + "(.*?)")
	if logoutCheck.MatchString(path) {
		return true
	}

	if path == "" {
		return false
	}
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	path = strings.Replace(path, "__edit_pk", "id", -1)
	path = strings.Replace(path, "__detail_pk", "id", -1)

	// GetURLParam 回傳處理後的url參數及url
	path, params := GetURLParam(path)
	for key, value := range formParams {
		if len(value) > 0 {
			params.Add(key, value[0])
		}
	}

	for _, p := range user.Permissions {
		if p.HTTPMethod[0] == "" || methodInSlice(p.HTTPMethod, method) {
			if p.HTTPPath[0] == "*" {
				return true
			}

			for i := 0; i < len(p.HTTPPath); i++ {
				matchPath := "/admin" + strings.TrimSpace(p.HTTPPath[i])
				matchPath, matchParam := GetURLParam(matchPath)
				if matchPath == path {
					// checkParam 檢查url與url參數是否符合
					if checkParam(params, matchParam) {
						return true
					}
				}

				reg, err := regexp.Compile(matchPath)
				if err != nil {
					continue
				}

				if reg.FindString(path) == path {
					// checkParam 檢查url與url參數是否符合
					if checkParam(params, matchParam) {
						return true
					}
				}
			}
		}
	}

	return false
}

// GetURLParam 回傳處理後的url參數及url
func GetURLParam(u string) (string, url.Values) {
	m := make(url.Values)
	urr := strings.Split(u, "?")
	if len(urr) > 1 {
		// 處理url的參數
		m, _ = url.ParseQuery(urr[1])
	}
	return urr[0], m
}

// methodInSlice 判斷string是否在slice裡
func methodInSlice(arr []string, str string) bool {
	for i := 0; i < len(arr); i++ {
		if strings.EqualFold(arr[i], str) {
			return true
		}
	}
	return false
}

// checkParam 檢查url與url參數是否符合
func checkParam(src, comp url.Values) bool {
	if len(comp) == 0 {
		return true
	}
	if len(src) == 0 {
		return false
	}
	for key, value := range comp {
		v, find := src[key]
		if !find {
			return false
		}
		if len(value) == 0 {
			continue
		}
		if len(v) == 0 {
			return false
		}
		for i := 0; i < len(v); i++ {
			if v[i] == value[i] {
				continue
			} else {
				return false
			}
		}
	}
	return true
}

// IsSuperAdmin 判斷是否為超級管理員
func (user UserModel) IsSuperAdmin() bool {
	for _, permission := range user.Permissions {
		if len(permission.HTTPPath) > 0 && permission.HTTPPath[0] == "*" && permission.HTTPMethod[0] == "" {
			return true
		}
	}
	return false
}
