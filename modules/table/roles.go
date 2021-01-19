package table

import (
	"database/sql"
	"errors"
	"fmt"
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
	"strings"
)

// GetRolesPanel 取得角色面板資訊、表單資訊
func (s *SystemTable) GetRolesPanel(ctx *context.Context) (rolesTable Table) {
	// DefaultBaseTable 建立預設的BaseTable(同時也是Table(interface))
	rolesTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	info := rolesTable.GetInfo()
	// 增加面板欄位資訊
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("角色", "name", db.Varchar).FieldFilterable()
	info.AddField("標誌", "slug", db.Varchar).FieldFilterable()
	info.AddField("權限", "name", db.Varchar).FieldJoin(types.Join{
		JoinTable: "role_permissions",
		JoinField: "role_id",
		Field:     "id",
		BaseTable: "roles",
	}).FieldJoin(types.Join{
		JoinTable: "permissions",
		JoinField: "id",
		Field:     "permission_id",
		BaseTable: "role_permissions",
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
			return "沒有權限"
		}
		return labels
	}).FieldFilterable()
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	// 資訊面板需要使用到刪除函式
	info.SetTable("roles").SetTitle("角色").SetDescription("角色管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)
			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("role_users").WhereIn("role_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_users資料表角色發生錯誤"), nil
					}
				}
				err = s.connection().SetTx(tx).
					Table("role_menu").WhereIn("role_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_menu資料表角色發生錯誤"), nil
					}
				}
				err = s.connection().SetTx(tx).
					Table("role_permissions").WhereIn("role_id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除role_permissions資料表角色發生錯誤"), nil
					}
				}
				err = s.connection().SetTx(tx).
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

	// 取得FormPanel
	formList := rolesTable.GetFormPanel()

	// 增加欄位資訊
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("角色名稱", "name", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("不能重複")).SetFieldMust()
	formList.AddField("角色標誌", "slug", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("不能重複")).SetFieldMust()
	formList.AddField("權限", "permission_id", db.Varchar, form.SelectBox).
		SetFieldOptionFromTable("permissions", "name", "id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var permissions = make([]string, 0)
			if model.ID == "" {
				return permissions
			}

			permissionModel, _ := s.table("role_permissions").
				Select("permission_id").Where("role_id", "=", model.ID).All()
			for _, v := range permissionModel {
				permissions = append(permissions, strconv.FormatInt(v["permission_id"].(int64), 10))
			}
			return permissions
		}).SetFieldHelpMsg(template.HTML("沒有對應選項?") + link("/admin/info/permission/new", "立刻新增權限"))
	formList.AddField("建立時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("更新時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.SetTable("roles").SetTitle("角色").SetDescription("角色管理")

	// 設置角色新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if models.DefaultRoleModel().SetConn(s.conn).IsSlugExist(values.Get("slug"), "") {
			return errors.New("角色標誌已經存在")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增角色
			role, err := models.DefaultRoleModel().SetTx(tx).SetConn(s.conn).AddRole(values.Get("name"), values.Get("slug"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("新增角色發生錯誤，可能原因:角色名稱已被註冊"), nil
				}
			}
			// 	新增權限
			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, err = role.SetTx(tx).AddPermission(values["permission_id[]"][i])
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("增加該角色權限發生錯誤"), nil
					}
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置角色更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if models.DefaultRoleModel().SetConn(s.conn).IsSlugExist(values.Get("slug"), values.Get("id")) {
			return errors.New("角色標誌已經存在")
		}
		// 設置RoleModel與ID
		role := models.GetRoleModelAndID(values.Get("id")).SetConn(s.conn)

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新角色
			_, err := role.SetTx(tx).Update(values.Get("name"), values.Get("slug"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("更新角色發生錯誤，可能原因:角色名稱已被註冊"), nil
				}
			}
			// 刪除該角色權限
			err = role.SetTx(tx).DeletePermissions()
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return errors.New("刪除該角色所有權限發生錯誤"), nil
				}
			}
			// 新增權限
			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, err = role.SetTx(tx).AddPermission(values["permission_id[]"][i])
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("增加該角色權限發生錯誤"), nil
					}
				}
			}
			return nil, nil
		})
		return txErr
	})
	return
}
