package table

import (
	"database/sql"
	dbsql "database/sql"
	"errors"
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strings"
)

// GetPermissionPanel 取得權限資訊面板、表單資訊
func GetPermissionPanel(conn db.Connection) (permissionTable Table) {
	// DefaultBaseTable 建立預設的BaseTable(同時也是Table(interface))
	permissionTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	info := permissionTable.GetInfo()
	// 增加面板欄位資訊
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("權限", "name", db.Varchar).FieldFilterable()
	info.AddField("標誌", "slug", db.Varchar).FieldFilterable()
	info.AddField("方法", "http_method", db.Varchar).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "" {
				return "All methods"
			}
			return value.Value
		})
	info.AddField("路徑", "http_path", db.Varchar).
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			pathArr := strings.Split(model.Value, "\n")
			res := ""
			for i := 0; i < len(pathArr); i++ {
				if i == len(pathArr)-1 {
					res += string(template.HTML(fmt.Sprintf(`<span class="label label-success" style="background-color: ;">%s</span>`, pathArr[i])))
				} else {
					res += string(template.HTML(fmt.Sprintf(`<span class="label label-success" style="background-color: ;">%s</span>`, pathArr[i]) + "<br><br>"))
				}
			}
			return res
		})
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	// 資訊面板需要使用到刪除函式
	info.SetTable("permissions").SetTitle("權限").SetDescription("權限管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)
			_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *dbsql.Tx) (e error, i map[string]interface{}) {
				err := db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("role_permissions").WhereIn("permission_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_permissions資料表角色發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("user_permissions").WhereIn("permission_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除user_permissions資料表角色發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("permissions").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除permissions資料表角色發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 取得FormPanel
	formList := permissionTable.GetFormPanel()

	// 增加欄位資訊
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("權限名稱", "name", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("不能重複")).SetFieldMust()
	formList.AddField("權限標誌", "slug", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("不能重複")).SetFieldMust()
	formList.AddField("可使用方法", "http_method", db.Varchar, form.Select).
		SetFieldOptions(types.FieldOptions{
			{Value: "GET", Text: "GET"},
			{Value: "POST", Text: "POST"},
			{Value: "PUT", Text: "PUT"},
			{Value: "DELETE", Text: "DELETE"},
		}).SetDisplayFunc(func(model types.FieldModel) interface{} {
		return strings.Split(model.Value, ",")
	}).
		SetFieldHelpMsg(template.HTML("如果為空代表可以使用所有方法"))
	formList.AddField("可使用路徑", "http_path", db.Varchar, form.TextArea).
		SetFieldHelpMsg(template.HTML("路徑不包含前綴(/admin)並且一行設置一個路徑，若要輸入新路徑請換行"))
	formList.AddField("建立時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("更新時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.SetTable("permissions").SetTitle("權限").SetDescription("權限管理")

	// 設置權限新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("slug", "name", "http_path") {
			return errors.New("權限名稱、權限標誌、可使用路徑不能為空")
		}
		if models.DefaultPermissionModel().SetConn(conn).IsSlugExist(values.Get("slug"), "") {
			return errors.New("權限標誌已經存在")
		}

		method := strings.Join(values["http_method[]"], ",")
		path := strings.TrimSpace(values.Get("http_path"))
		_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增權限資料
			_, err := models.DefaultPermissionModel().SetTx(tx).SetConn(conn).AddPermission(
				values.Get("name"), values.Get("slug"), method, path)
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("新增權限發生錯誤，可能原因:權限名稱已被註冊"), nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置權限更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if models.DefaultPermissionModel().SetConn(conn).IsSlugExist(values.Get("slug"), values.Get("id")) {
			return errors.New("權限標誌已經存在")
		}

		// 設置RoleModel與ID
		permission := models.GetPermissionModelAndID(values.Get("id")).SetConn(conn)

		method := strings.Join(values["http_method[]"], ",")
		path := strings.TrimSpace(values.Get("http_path"))
		_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新角色
			_, err := permission.SetTx(tx).Update(values.Get("name"), values.Get("slug"), method, path)
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("更新權限發生錯誤，可能原因:權限名稱已被註冊"), nil
				}
			}
			return nil, nil
		})
		return txErr
	})
	return
}
