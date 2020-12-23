package table

import (
	"database/sql"
	"errors"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
	"strings"
)

// GetManagerInfoPanel 取得使用者資訊面板
func GetManagerInfoPanel(conn db.Connection) (managerTable Table) {
	// 建立BaseTable
	managerTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	info := managerTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("用戶ID", "userid", db.Varchar).FieldFilterable()
	info.AddField("用戶名稱", "username", db.Varchar).FieldFilterable()
	info.AddField("電話", "phone", db.Varchar).FieldFilterable()
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
	}).FieldDisplayFunc(func(model types.FieldModel) interface{} {
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
	info.AddField("用戶照片", "picture", db.Varchar)
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	info.SetTable("users").SetTitle("用戶").SetDescription("用戶管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := db.SetConnectionAndCRUD(conn).WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := db.SetConnectionAndCRUD(conn).WithTx(tx).
					Table("role_users").WhereIn("user_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_users資料表角色發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).WithTx(tx).
					Table("user_permissions").WhereIn("user_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除user_permissions資料表權限發生錯誤"), nil
					}
				}
				err = db.SetConnectionAndCRUD(conn).WithTx(tx).
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

	return managerTable
}

// GetMenuFormPanel 取得menu表單面板資訊
func GetMenuFormPanel(conn db.Connection) (menuTable Table) {
	// 建立BaseTable
	menuTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 父級菜單選單
	var parentIDOptions = types.FieldOptions{
		{
			Text:  "ROOT",
			Value: "0",
		},
	}

	// 處理父級菜單欄位
	// 取得父級菜單資料(parent_id=0)
	allMenus, _ := db.TableAndCleanData("menu", conn).
		Where("parent_id", "=", 0).Select("id", "title").
		OrderBy("field_order", "asc").All()

	// 紀錄父級菜單ID
	allMenuIDs := make([]interface{}, len(allMenus))
	if len(allMenuIDs) > 0 {
		for i := 0; i < len(allMenus); i++ {
			allMenuIDs[i] = allMenus[i]["id"]
		}

		// 取得父級底下的菜單
		secondLevelMenus, _ := db.TableAndCleanData("menu", conn).
			WhereIn("parent_id", allMenuIDs).Select("id", "title", "parent_id").All()

		for i := 0; i < len(allMenus); i++ {
			// 新增父級的選項名稱
			parentIDOptions = append(parentIDOptions, types.FieldOption{
				TextHTML: "&nbsp;&nbsp;┝  " + template.HTML(allMenus[i]["title"].(string)),
				Value:    strconv.Itoa(int(allMenus[i]["id"].(int64))),
			})

			// 取得父級底下的菜單並加入
			for j := 0; j < len(secondLevelMenus); j++ {
				if secondLevelMenus[j]["parent_id"].(int64) == allMenus[i]["id"].(int64) {
					parentIDOptions = append(parentIDOptions, types.FieldOption{
						TextHTML: "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;┝  " +
							template.HTML(secondLevelMenus[j]["title"].(string)),
						Value: strconv.Itoa(int(secondLevelMenus[j]["id"].(int64))),
					})
				}
			}
		}
	}

	// 取得types.FormPanel
	formList := menuTable.GetFormPanel()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("父級菜單", "parent_id", db.Int, form.Text).
		// SetFieldOptions 設置Field.FieldOptions
		// SetDisplayFunc 設置欄位過濾函式至DisplayFunc
		SetFieldOptions(parentIDOptions).
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var menuItem []string
			if model.ID == "" {
				return menuItem
			}

			menuModel, _ := db.TableAndCleanData("menu", conn).
				Select("parent_id").FindByID(model.ID)
			menuItem = append(menuItem, strconv.FormatInt(menuModel["parent_id"].(int64), 10))
			return menuItem
		})
	formList.AddField("菜單名稱", "title", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("標頭名稱", "header", db.Varchar, form.Text)
	formList.AddField("圖標", "icon", db.Varchar, form.Text)
	formList.AddField("路徑", "url", db.Varchar, form.Text)
	formList.AddField("角色", "roles", db.Int, form.Text).
		// SetFieldOptionFromTable 設置Field.FieldOptionFromTable(選單名稱由資料表取得)
		SetFieldOptionFromTable("roles", "slug", "id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var roles []string
			if model.ID == "" {
				return roles
			}

			roleModel, _ := db.TableAndCleanData("role_menu", conn).
				Select("role_id").Where("menu_id", "=", model.ID).All()
			for _, m := range roleModel {
				roles = append(roles, strconv.FormatInt(m["role_id"].(int64), 10))
			}
			return roles
		})
	formList.AddField("更新時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("建立時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.SetTable("menu").SetTitle("菜單").SetDescription("菜單處理")
	return
}

// interfaces 將[]string轉換成[]interface{}
func interfaces(arr []string) []interface{} {
	var iarr = make([]interface{}, len(arr))
	for key, v := range arr {
		iarr[key] = v
	}
	return iarr
}
