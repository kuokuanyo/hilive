package table

import (
	"database/sql"
	"errors"
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GetManagerPanel 取得用戶頁面、表單資訊
func GetManagerPanel(conn db.Connection) (managerTable Table) {
	// 建立BaseTable
	managerTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	info := managerTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("UserID", "userid", db.Varchar).FieldFilterable()
	info.AddField("用戶名稱", "username", db.Varchar).FieldFilterable()
	info.AddField("用戶照片", "picture", db.Varchar)
	info.AddField("電話號碼", "phone", db.Varchar).FieldFilterable()
	info.AddField("信箱", "email", db.Varchar).FieldFilterable()
	info.AddField("角色", "name", db.Varchar).FieldJoin(types.Join{
		JoinTable: "role_users",
		JoinField: "user_id",
		Field:     "id",
		BaseTable: "users",
	}).FieldJoin(types.Join{
		JoinTable: "roles",
		JoinField: "id",
		Field:     "role_id",
		BaseTable: "role_users",
	}).SetDisplayFunc(func(model types.FieldModel) interface{} {
		labels := template.HTML("")
		labelValues := strings.Split(model.Value, types.JoinFieldValueDelimiter)

		for key, label := range labelValues {
			if key == len(labelValues)-1 {
				labels += template.HTML(fmt.Sprintf(`<span class="label label-success" style="background-color: ;">%s</span>`, label))
			} else {
				labels += template.HTML(fmt.Sprintf(`<span class="label label-success" style="background-color: ;">%s</span>`, label) + "<br><br>")
			}
		}
		if labels == template.HTML("") {
			return "沒有角色"
		}
		return labels
	}).FieldFilterable()
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	info.SetTable("users").SetTitle("用戶").SetDescription("用戶管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("role_users").WhereIn("user_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_users資料表角色發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("user_permissions").WhereIn("user_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除user_permissions資料表權限發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("users").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除users資料表用戶發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})
	
	// 取得FormPanel
	formList := managerTable.GetFormPanel()

	// 增加欄位資訊
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("UserID", "userid", db.Varchar, form.Text).SetFieldMust().SetFieldHelpMsg(template.HTML("LINE ID"))
	formList.AddField("用戶名稱", "username", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("用戶照片", "picture", db.Varchar, form.Text)
	formList.AddField("電話號碼", "phone", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("用途: 登入")).SetFieldMust()
	formList.AddField("信箱", "email", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("角色", "role_id", db.Varchar, form.Select).
		SetFieldOptionFromTable("roles", "slug", "id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var roles []string
			if model.ID == "" {
				return roles
			}

			roleModel, _ := db.TableAndCleanData("role_users", conn).Select("role_id").
				Where("user_id", "=", model.ID).All()
			for _, v := range roleModel {
				roles = append(roles, strconv.FormatInt(v["role_id"].(int64), 10))
			}
			return roles
		}).SetFieldHelpMsg(template.HTML("沒有對應選項?") + link("/admin/info/roles/new", "立刻新增角色"))
	formList.AddField("權限", "permission_id", db.Varchar, form.Select).
		SetFieldOptionFromTable("permissions", "slug", "id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var permissions []string
			if model.ID == "" {
				return permissions
			}

			permissionModel, _ := db.TableAndCleanData("user_permissions", conn).
				Select("permission_id").Where("user_id", "=", model.ID).All()
			for _, v := range permissionModel {
				permissions = append(permissions, strconv.FormatInt(v["permission_id"].(int64), 10))
			}
			return permissions
		}).SetFieldHelpMsg(template.HTML("沒有對應選項?") + link("/admin/info/permission/new", "立刻新增權限"))
	formList.AddField("密碼", "password", db.Varchar, form.Password).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			return ""
		})
	formList.AddField("請確認密碼", "password_again", db.Varchar, form.Password).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			return ""
		})
	formList.AddField("更新時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("建立時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.SetTable("users").SetTitle("用戶").SetDescription("用戶管理")

	// 設置用戶更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("userid", "phone", "username", "email", "password") {
			return errors.New("userid、電話號碼、用戶名稱、電子信箱、密碼都不能為空值")
		}
		password := values.Get("password")
		if password != values.Get("password_again") {
			return errors.New("密碼不相符")
		}
		password = EncodePassword([]byte(values.Get("password")))

		user := models.GetUserModelAndID("users", values.Get("id")).SetConn(conn)

		_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新用戶資料
			_, err := user.SetTx(tx).Update(values.Get("userid"), values.Get("username"), values.Get("picture"), values.Get("phone"), values.Get("email"), password)
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("更新用戶發生錯誤，可能原因:用戶ID、手機號碼、電子信箱已被註冊"), nil
				}
			}
			// 刪除該ID用戶所有角色並新增新的角色
			err = user.SetTx(tx).DeleteRolesByID()
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("刪除該用戶角色發生錯誤"), nil
				}
			}
			for i := 0; i < len(values["role_id[]"]); i++ {
				_, err = user.SetTx(tx).AddRole(values["role_id[]"][i])
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("增加該用戶角色發生錯誤"), nil
					}
				}
			}

			// 刪除該ID用戶權限並新增新的權限
			err = user.SetTx(tx).DeletePermissionsByID()
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("刪除該用戶權限發生錯誤"), nil
				}
			}
			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, err = user.SetTx(tx).AddPermission(values["permission_id[]"][i])
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("增加該用戶權限發生錯誤"), nil
					}
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置用戶新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("phone", "username", "password") {
			return errors.New("電話號碼、用戶名稱、密碼都不能為空值")
		}
		password := values.Get("password")
		if password != values.Get("password_again") {
			return errors.New("密碼不相符")
		}
		password = EncodePassword([]byte(values.Get("password")))

		_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增用戶資料
			user, err := models.DefaultUserModel().SetTx(tx).SetConn(conn).AddUser(values.Get("userid"), values.Get("username"),
				values.Get("phone"), values.Get("email"), password)
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("新增用戶發生錯誤，可能原因:用戶ID、電話號碼、電子信箱已被註冊"), nil
				}
			}
			// 新增角色、權限
			for i := 0; i < len(values["role_id[]"]); i++ {
				_, err = user.SetTx(tx).AddRole(values["role_id[]"][i])
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("增加該用戶角色發生錯誤"), nil
					}
				}
			}
			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, err = user.SetTx(tx).AddPermission(values["permission_id[]"][i])
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("增加該用戶權限發生錯誤"), nil
					}
				}
			}
			return nil, nil
		})
		return txErr
	})
	return
}

// EncodePassword 加密
func EncodePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash[:])
}
