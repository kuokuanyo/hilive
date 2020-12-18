package models

import "hilive/modules/db"

// RoleModel roles table
type RoleModel struct {
	Base
	ID        int64
	Name      string
	Slug      string
	CreatedAt string
	UpdatedAt string
}

// DefaultRoleModel 預設RoleModel
func DefaultRoleModel() RoleModel {
	return RoleModel{Base: Base{TableName: "roles"}}
}

// SetConn 設置connection
func (role RoleModel) SetConn(conn db.Connection) RoleModel {
	role.Conn = conn
	return role
}

// MapToRoleModel 設置值至RoleModel
func (role RoleModel) MapToRoleModel(m map[string]interface{}) RoleModel {
	role.ID = m["id"].(int64)
	role.Name, _ = m["name"].(string)
	role.Slug, _ = m["slug"].(string)
	role.CreatedAt, _ = m["created_at"].(string)
	role.UpdatedAt, _ = m["updated_at"].(string)
	return role
}