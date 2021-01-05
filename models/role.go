package models

import (
	dbsql "database/sql"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
	"time"
)

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

// SetTx 設置Tx
func (role RoleModel) SetTx(tx *dbsql.Tx) RoleModel {
	role.Tx = tx
	return role
}

// GetRoleModelAndID 設置RoleModel與ID
func GetRoleModelAndID(id string) RoleModel {
	idInt, _ := strconv.Atoi(id)
	return RoleModel{Base: Base{TableName: "roles"}, ID: int64(idInt)}
}

// Update 更新角色資料
func (role RoleModel) Update(name, slug string) (int64, error) {
	return role.SetTx(role.Tx).Table(role.TableName).Where("id", "=", role.ID).
		Update(sql.Value{
			"name":       name,
			"slug":       slug,
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		})
}

// AddRole 新增角色
func (role RoleModel) AddRole(name, slug string) (RoleModel, error) {
	id, err := role.SetTx(role.Tx).Table(role.TableName).Insert(sql.Value{
		"name": name,
		"slug": slug,
	})

	role.ID = id
	role.Name = name
	role.Slug = slug
	return role, err
}

// AddPermission 新增權限
func (role RoleModel) AddPermission(id string) (int64, error) {
	checkPermission, _ := role.Table("role_permissions").
		Where("permission_id", "=", id).
		Where("role_id", "=", role.ID).First()
	if id != "" {
		if checkPermission == nil {
			return role.SetTx(role.Tx).Table("role_permissions").
				Insert(sql.Value{
					"permission_id": id,
					"role_id":       role.ID,
				})
		}
	}
	return 0, nil
}

// DeletePermissions 刪除角色所有權限
func (role RoleModel) DeletePermissions() error {
	return role.SetTx(role.Tx).Table("role_permissions").
		Where("role_id", "=", role.ID).Delete()
}

// IsSlugExist 檢查slug是否已經存在
func (role RoleModel) IsSlugExist(slug, id string) bool {
	if id == "" {
		check, _ := role.Table(role.TableName).Where("slug", "=", slug).First()
		return check != nil
	}
	check, _ := role.Table(role.TableName).
		Where("slug", "=", slug).
		Where("id", "!=", id).
		First()
	return check != nil
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
