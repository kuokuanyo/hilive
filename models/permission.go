package models

import (
	"hilive/modules/db"
	"strings"
)

// PermissionModel permissions table
type PermissionModel struct {
	Base
	ID         int64
	Name       string
	Slug       string
	HTTPMethod []string
	HTTPPath   []string
	CreatedAt  string
	UpdatedAt  string
}

// DefaultPermissionModel 預設PermissionModel
func DefaultPermissionModel() PermissionModel {
	return PermissionModel{Base: Base{TableName: "permissions"}}
}

// SetConn 設置connection
func (permission PermissionModel) SetConn(conn db.Connection) PermissionModel {
	permission.Base.Conn = conn
	return permission
}

// MapToPermissionModel 設置值至PermissionModel
func (permission PermissionModel) MapToPermissionModel(m map[string]interface{}) PermissionModel {
	permission.ID = m["id"].(int64)
	permission.Name, _ = m["name"].(string)
	permission.Slug, _ = m["slug"].(string)
	methods, _ := m["http_method"].(string)
	if methods != "" {
		permission.HTTPMethod = strings.Split(methods, ",")
	} else {
		permission.HTTPMethod = []string{""}
	}
	path, _ := m["http_path"].(string)
	permission.HTTPPath = strings.Split(path, "\n")
	permission.CreatedAt, _ = m["created_at"].(string)
	permission.UpdatedAt, _ = m["updated_at"].(string)
	return permission 
}

// IsSlugExist 檢查標誌是否已經存在
func (permission PermissionModel) IsSlugExist(slug string, id string) bool {
	if id == "" {
		check, _ := permission.Table(permission.TableName).Where("slug", "=", slug).First()
		return check != nil
	}
	check, _ := permission.Table(permission.TableName).
		Where("slug", "=", slug).
		Where("id", "!=", id).
		First()
	return check != nil
}
