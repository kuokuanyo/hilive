package models

import (
	dbsql "database/sql"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
	"strings"
	"time"
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

// GetPermissionModelAndID 設置PermissionModel與ID
func GetPermissionModelAndID(id string) PermissionModel {
	idInt, _ := strconv.Atoi(id)
	return PermissionModel{Base: Base{TableName: "permissions"}, ID: int64(idInt)}
}

// SetConn 設置connection
func (permission PermissionModel) SetConn(conn db.Connection) PermissionModel {
	permission.Base.Conn = conn
	return permission
}

// SetTx 設置Tx
func (permission PermissionModel) SetTx(tx *dbsql.Tx) PermissionModel {
	permission.Base.Tx = tx
	return permission
}

// AddPermission 新增權限
func (permission PermissionModel) AddPermission(name, slug, methods, path string) (PermissionModel, error) {
	id, err := permission.SetTx(permission.Tx).Table(permission.TableName).Insert(sql.Value{
		"name":        name,
		"slug":        slug,
		"http_method": methods,
		"http_path":   path,
	})

	permission.ID = id
	permission.Name = name
	permission.Slug = slug
	if methods != "" {
		permission.HTTPMethod = strings.Split(methods, ",")
	} else {
		permission.HTTPMethod = []string{""}
	}
	permission.HTTPPath = strings.Split(path, "\n")
	return permission, err
}

// Update 更新權限資料
func (permission PermissionModel) Update(name, slug, methods, path string) (int64, error) {
	return permission.SetTx(permission.Tx).Table(permission.TableName).Where("id", "=", permission.ID).
		Update(sql.Value{
			"name":        name,
			"slug":        slug,
			"http_method": methods,
			"http_path":   path,
			"updated_at":  time.Now().Format("2006-01-02 15:04:05"),
		})
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
