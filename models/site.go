package models

import (
	"fmt"
	"hilive/modules/db"
	"hilive/modules/db/sql"
)

// SiteModel 為site資料表欄位
type SiteModel struct {
	// 記錄資料表名稱...等資訊
	Base

	ID          int64
	ConfigKey   string
	ConfigValue string
	CreatedAt   string
	UpdatedAt   string
}

// DefaultSiteModel 預設SiteModel
func DefaultSiteModel() SiteModel {
	return SiteModel{Base: Base{TableName: "site"}}
}

// SetConn 設置SiteModel.Base.Conn
func (s SiteModel) SetConn(conn db.Connection) SiteModel {
	s.Conn = conn
	return s
}

// Init 初始化Site資料表(插入或更新site資料表資料)
func (s SiteModel) Init(cfg map[string]string) {
	// 取得site資料表所有資料
	items, err := s.Base.Table(s.Base.TableName).All()
	if err != nil {
		panic("取得site資料表資料發生錯誤")
	}

	for key, value1 := range cfg {
		row := make([]map[string]interface{}, 0)
		for _, value2 := range items {
			keyStr := fmt.Sprintf("%v", value2["config_key"])
			if key == keyStr {
				row = append(row, value2)
			}
		}

		if len(row) == 0 {
			_, err := s.Base.Table(s.Base.TableName).Insert(sql.Value{
				"config_key":   key,
				"config_value": value1,
			})
			if err != nil {
				panic("插入資料至site資料表發生錯誤")
			}
		} else {
			if value1 != "" {
				_, err := s.Base.Table(s.Base.TableName).Where("config_key", "=", key).Update(sql.Value{
					"config_value": value1,
				})
				if err != nil {
					if fmt.Sprintf("%s", err) != "沒有影響任何資料" {
						panic("更新site資料表發生錯誤")
					}
				}
			}
		}
	}
}
