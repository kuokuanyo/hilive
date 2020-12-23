package models

import (
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
	"time"
)

// MenuModel is menu model structure.
type MenuModel struct {
	Base
	ID        int64
	Title     string
	ParentID  int64
	Icon      string
	URL       string
	Header    string
	CreatedAt string
	UpdatedAt string
}

// DefaultMenuModel 預設MenuModel
func DefaultMenuModel() MenuModel {
	return MenuModel{Base: Base{TableName: "menu"}}
}

// SetConn 設置connection
func (menu MenuModel) SetConn(conn db.Connection) MenuModel {
	menu.Conn = conn
	return menu
}

// SetMenuModelByID 設置MenuModel
func SetMenuModelByID(id string) MenuModel {
	idInt, _ := strconv.Atoi(id)
	return MenuModel{Base: Base{TableName: "menu"}, ID: int64(idInt)}
}

// New 新增資料
func (menu MenuModel) New(title, icon, uri, header string, parentID, order int64) (MenuModel, error) {
	id, err := menu.Table(menu.TableName).Insert(sql.Value{
		"title":       title,
		"parent_id":   parentID,
		"icon":        icon,
		"url":         uri,
		"field_order": order,
		"header":      header,
	})

	menu.ID = id
	menu.Title = title
	menu.ParentID = parentID
	menu.Icon = icon
	menu.URL = uri
	menu.Header = header
	return menu, err
}

// Update 更新資料
func (menu MenuModel) Update(title, icon, url, header string, parentID int64) (int64, error) {
	return menu.Table(menu.Base.TableName).
		Where("id", "=", menu.ID).
		Update(sql.Value{
			"title":      title,
			"parent_id":  parentID,
			"icon":       icon,
			"url":        url,
			"header":     header,
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		})
}

// Delete 刪除資料
func (menu MenuModel) Delete() {
	// 刪除menu及role_menu資料表資料
	menu.Table(menu.Base.TableName).Where("id", "=", menu.ID).Delete()
	menu.Table("role_menu").Where("menu_id", "=", menu.ID).Delete()

	// 如果是其他菜單的父級也必須刪除
	items, _ := menu.Table(menu.Base.TableName).Where("parent_id", "=", menu.ID).All()
	if len(items) > 0 {
		ids := make([]interface{}, len(items))
		for i := 0; i < len(ids); i++ {
			ids[i] = items[i]["id"]
		}
		menu.Table("role_menu").WhereIn("menu_id", ids).Delete()
	}
	menu.Table(menu.Base.TableName).Where("parent_id", "=", menu.ID).Delete()
}

// AddRole 新建角色
func (menu MenuModel) AddRole(roleID string) (int64, error) {
	if roleID != "" {
		// 先檢查角色表裡是否有該菜單角色
		checkRole, _ := menu.Table("role_menu").
			Where("role_id", "=", roleID).
			Where("menu_id", "=", menu.ID).First()
		if checkRole == nil {
			return menu.Table("role_menu").Insert(sql.Value{
				"role_id": roleID,
				"menu_id": menu.ID,
			})
		}
	}
	return 0, nil
}

// DeleteRoles 刪除角色
func (menu MenuModel) DeleteRoles() error {
	return menu.Table("role_menu").
		Where("menu_id", "=", menu.ID).Delete()
}
