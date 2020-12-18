package table

import (
	"errors"
	"fmt"
	"hilive/modules/db"
	"hilive/modules/parameter"
	"hilive/modules/service"
	"hilive/modules/utils"
	"hilive/template/types"
)

// BaseTable 包含面板、表單資訊...等
type BaseTable struct {
	Informatoin      *types.InformationPanel
	Form             *types.FormPanel
	Detail           *types.InformationPanel
	CanAdd           bool
	EditAble         bool
	DeleteAble       bool
	PrimaryKey       PrimaryKey
	connectionDriver string
}

// PrimaryKey 紀錄主鍵及主鍵type
type PrimaryKey struct {
	Name string
	Type db.DatabaseType
}

// Table interface
type Table interface {
	// 設置主鍵並取得FormPanel
	GetFormPanel() *types.FormPanel

	// GetNewForm 處理並設置表單欄位細節資訊(允許增加的表單欄位)
	GetNewForm(services service.List) types.FormPanel

	// GetDataWithID 透過id取得資料並將值、預設值設置至BaseTable.Form.FormFields
	GetDataWithID(param parameter.Parameters, services service.List) (types.FormPanel, error)
}

// DefaultBaseTable 建立預設的BaseTable(同時也是Table(interface))
func DefaultBaseTable(cfgs ...ConfigTable) Table {
	var cfg ConfigTable
	if len(cfgs) > 0 && cfgs[0].PrimaryKey.Name != "" {
		cfg = cfgs[0]
	} else {
		cfg = DefaultConfig()
	}
	return &BaseTable{
		Informatoin:      types.DefaultInformationPanel(cfg.PrimaryKey.Name),
		Form:             types.DefaultFormPanel(),
		Detail:           types.DefaultInformationPanel(cfg.PrimaryKey.Name),
		CanAdd:           cfg.CanAdd,
		EditAble:         cfg.EditAble,
		DeleteAble:       cfg.DeleteAble,
		PrimaryKey:       cfg.PrimaryKey,
		connectionDriver: cfg.Driver,
	}
}

// -----BaseTable的所有Table方法-----start

// GetFormPanel 設置主鍵並取得FormPanel
func (base *BaseTable) GetFormPanel() *types.FormPanel {
	return base.Form.SetPrimaryKey(base.PrimaryKey.Name, base.PrimaryKey.Type)
}

// GetNewForm 處理並設置表單欄位細節資訊(允許增加的表單欄位)
func (base *BaseTable) GetNewForm(services service.List) types.FormPanel {
	return types.FormPanel{FieldList: base.Form.SetAllowAddValueOfField(services, base.GetSQLByService)}
}

// GetDataWithID 透過id取得資料並將值、預設值設置至BaseTable.Form.FormFields
func (base *BaseTable) GetDataWithID(param parameter.Parameters, services service.List) (types.FormPanel, error) {
	var (
		// FindPK 取得__pk的值(單個)
		id                     = param.FindPK()
		res                    map[string]interface{}
		args                   = []interface{}{id}
		fields, joins, groupBy = "", "", ""
		tableName              = base.GetFormPanel().Table
		pk                     = tableName + "." + base.PrimaryKey.Name
		queryStatement         = "select %s from %s" + " %s where " + pk + " = ? %s "
	)
	// 所有欄位
	columns, _ := base.getColumns(base.Form.Table, services)

	for _, field := range base.Form.FieldList {
		if field.Field != pk && utils.InArray(columns, field.Field) && !field.Joins.Valid() {
			fields += tableName + "." + field.Field + ","
		}
	}
	fields += pk
	queryCmd := fmt.Sprintf(queryStatement, fields, tableName, joins, groupBy)

	result, err := base.getConnectionByService(services).Query(queryCmd, args...)
	if err != nil {
		return types.FormPanel{Title: base.Form.Title, Description: base.Form.Description}, err
	}
	if len(result) == 0 {
		return types.FormPanel{Title: base.Form.Title, Description: base.Form.Description}, errors.New("錯誤的id")
	}
	res = result[0]

	var fieldList = base.Form.FieldWithValue(base.PrimaryKey.Name,
		id, columns, res, services, base.GetSQLByService)
	return types.FormPanel{
		FieldList:   fieldList,
		Title:       base.Form.Title,
		Description: base.Form.Description,
	}, nil
}

// -----BaseTable的所有Table方法-----end

// GetSQLByService 設置db.SQL(struct)的Connection、CRUD
func (base *BaseTable) GetSQLByService(services service.List) *db.SQL {
	if base.connectionDriver != "" {
		return db.SetConnectionAndCRUD(db.ConvertServiceToConnection(services.Get(base.connectionDriver)))
	}
	return nil
}

// getConnectionByService 取得Connection(interface)
func (base *BaseTable) getConnectionByService(services service.List) db.Connection {
	if base.connectionDriver != "" {
		return db.ConvertServiceToConnection(services.Get(base.connectionDriver))
	}
	return nil
}

// getColumns 取得所有欄位
func (base *BaseTable) getColumns(table string, services service.List) ([]string, bool) {
	var auto bool
	columnsModel, _ := base.GetSQLByService(services).Table(table).ShowColumns()
	columns := make([]string, len(columnsModel))
	for key, model := range columnsModel {
		columns[key] = model["Field"].(string)
		if columns[key] == base.PrimaryKey.Name { // 如果為主鍵
			if v, ok := model["Extra"].(string); ok {
				if v == "auto_increment" {
					auto = true
				}
			}
		}
	}
	return columns, auto
}
