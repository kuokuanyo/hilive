package table

import (
	"database/sql"
	"errors"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/template/types"
	"html/template"
	"strings"
)

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
			_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
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
