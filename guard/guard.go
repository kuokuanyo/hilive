package guard

import (
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/service"
)

// Parameters 紀錄參數
var parameters = make(map[string]interface{})

// Guard struct
type Guard struct {
	Services service.List
	Conn     db.Connection
	Config   *config.Config
}
