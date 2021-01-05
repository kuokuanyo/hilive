package table

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	form2 "hilive/modules/form"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strings"
	"time"
)

// GetPermissionFormPanel 取得權限表單資訊
func GetPermissionFormPanel(conn db.Connection) (permissionTable Table) {
	// DefaultBaseTable 建立預設的BaseTable(同時也是Table(interface))
	permissionTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))
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
	}).SetPostFilterFunc(func(model types.PostFieldModel) interface{} {
		return strings.Join(model.Value, ",")
	}).
		SetFieldHelpMsg(template.HTML("如果為空代表可以使用所有方法"))
	formList.AddField("可使用路徑", "http_path", db.Varchar, form.TextArea).
		SetPostFilterFunc(func(model types.PostFieldModel) interface{} {
			return strings.TrimSpace(model.Value.Value())
		}).
		SetFieldHelpMsg(template.HTML("路徑不包含前綴(/admin)並且一行設置一個路徑，若要輸入新路徑請換行"))
	formList.AddField("建立時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("更新時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.SetTable("permissions").SetTitle("權限").SetDescription("新增權限").
		SetPostValidatorFunc(func(values form2.Values) error {
			if values.IsEmpty("slug", "http_path", "name") {
				return errors.New("權限名稱、權限標誌、可使用路徑不能為空")
			}
			if models.DefaultPermissionModel().SetConn(conn).IsSlugExist(values.Get("slug"), values.Get("id")) {
				return errors.New("權限標誌已經存在")
			}
			return nil
		}).
		SetPostHookFunc(func(values form2.Values) error {
			_, err := db.SetConnectionAndCRUD(conn).Table("permissions").
				Where("id", "=", values.Get("id")).Update(sql.Value{
				"updated_at": time.Now().Format("2006-01-02 15:04:05"),
			})
			return err
		})
	return
}

// GetPermissionInfoPanel 取得角色資訊面板
func GetPermissionInfoPanel(conn db.Connection) (permissionTable Table) {
	permissionTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	info := permissionTable.GetInfo()
	// 增加面板欄位資訊
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("權限", "name", db.Varchar).FieldFilterable()
	info.AddField("標誌", "slug", db.Varchar).FieldFilterable()
	info.AddField("方法", "http_method", db.Varchar).
		FieldDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "" {
				return "All methods"
			}
			return value.Value
		})
	info.AddField("路徑", "http_path", db.Varchar).
		FieldDisplayFunc(func(model types.FieldModel) interface{} {
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
	return
}
