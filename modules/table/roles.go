package table

import (
	"database/sql"
	"errors"
	"hilive/modules/config"
	"hilive/modules/db"
)

// GetRolesInfoPanel 取得角色資訊面板
func GetRolesInfoPanel(conn db.Connection) (rolesTable Table) {
	rolesTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	info := rolesTable.GetInfo()
	// 增加面板欄位資訊
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("角色", "name", db.Varchar).FieldFilterable()
	info.AddField("標誌", "slug", db.Varchar).FieldFilterable()
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	// 資訊面板需要使用到刪除函式
	info.SetTable("roles").SetTitle("角色").SetDescription("角色管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)
			_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("role_users").WhereIn("role_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_users資料表角色發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("role_menu").WhereIn("role_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_menu資料表角色發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("role_permissions").WhereIn("role_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_permissions資料表角色發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).SetTx(tx).
					Table("roles").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除roles資料表角色發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})
	return
}
