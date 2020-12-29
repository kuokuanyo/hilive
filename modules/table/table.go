package table

import (
	"errors"
	"fmt"
	"hilive/modules/db"
	"hilive/modules/paginator"
	"hilive/modules/parameter"
	"hilive/modules/service"
	"hilive/modules/utils"
	"hilive/template/types"
	"html/template"
	"strconv"
)

// BaseTable 包含面板、表單資訊...等
type BaseTable struct {
	Informatoin      *types.InformationPanel
	Form             *types.FormPanel
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

// PanelInfo 頁面資訊
type PanelInfo struct {
	FieldList      types.FieldList  `json:"fieldlist"`        // 介面上的欄位資訊，是否可編輯、編輯選項、是否隱藏...等資訊
	InfoList       types.InfoList   `json:"info_list"`        // 每一筆資料的資訊
	FilterFormData types.FormFields `json:"filter_form_data"` // 可以篩選條件的欄位表單資訊
	Paginator      paginator.Paginator
	PrimaryKey     string
	Title          string `json:"title"`
	Description    string `json:"description"`
}

// FormInfo 表單資訊
type FormInfo struct {
	FieldList   types.FormFields `json:"field_list"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
}

// Table interface
type Table interface {
	// GetPrimaryKey 取得主鍵
	GetPrimaryKey() PrimaryKey

	// 設置主鍵取得InformationPanel
	GetInfo() *types.InformationPanel

	// GetForm 取得表單資訊
	GetForm() *types.FormPanel

	// 設置主鍵並取得FormPanel
	GetFormPanel() *types.FormPanel

	// GetNewForm 處理並設置表單欄位細節資訊(允許增加的表單欄位)
	GetNewForm(services service.List) FormInfo

	// GetDataWithID 透過id取得資料並將值、預設值設置至BaseTable.Form.FormFields
	GetDataWithID(param parameter.Parameters, services service.List) (FormInfo, error)

	// GetData 從資料庫取得頁面需要顯示的資料，回傳每一筆資料資訊、欄位資訊、可過濾欄位資訊、分頁資訊...等
	GetData(params parameter.Parameters, services service.List) (PanelInfo, error)
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
		CanAdd:           cfg.CanAdd,
		EditAble:         cfg.EditAble,
		DeleteAble:       cfg.DeleteAble,
		PrimaryKey:       cfg.PrimaryKey,
		connectionDriver: cfg.Driver,
	}
}

// -----BaseTable的所有Table方法-----start

// GetPrimaryKey 取得主鍵
func (base *BaseTable) GetPrimaryKey() PrimaryKey {
	return base.PrimaryKey
}

// GetInfo 設置主鍵取得InformationPanel
func (base *BaseTable) GetInfo() *types.InformationPanel {
	return base.Informatoin.SetPrimaryKey(base.PrimaryKey.Name, base.PrimaryKey.Type)
}

// GetForm 設置主鍵取得FormPanel
func (base *BaseTable) GetForm() *types.FormPanel {
	return base.Form.SetPrimaryKey(base.PrimaryKey.Name, base.PrimaryKey.Type)
}

// GetFormPanel 設置主鍵並取得FormPanel
func (base *BaseTable) GetFormPanel() *types.FormPanel {
	return base.Form.SetPrimaryKey(base.PrimaryKey.Name, base.PrimaryKey.Type)
}

// GetNewForm 處理並設置表單欄位細節資訊(新增資料的表單欄位)
func (base *BaseTable) GetNewForm(services service.List) FormInfo {
	return FormInfo{FieldList: base.Form.SetAllowAddValueOfField(services, base.GetSQLByService)}
}

// GetDataWithID 透過id取得資料並將值、預設值設置至BaseTable.Form.FormFields
func (base *BaseTable) GetDataWithID(param parameter.Parameters, services service.List) (FormInfo, error) {
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
		return FormInfo{Title: base.Form.Title, Description: base.Form.Description}, err
	}
	if len(result) == 0 {
		return FormInfo{Title: base.Form.Title, Description: base.Form.Description}, errors.New("錯誤的id")
	}
	res = result[0]

	var fieldList = base.Form.FieldWithValue(base.PrimaryKey.Name,
		id, columns, res, services, base.GetSQLByService)
	return FormInfo{
		FieldList:   fieldList,
		Title:       base.Form.Title,
		Description: base.Form.Description,
	}, nil
}

// GetData 從資料庫取得頁面需要顯示的資料，回傳每一筆資料資訊、欄位資訊、可過濾欄位資訊...等
func (base *BaseTable) GetData(params parameter.Parameters, services service.List) (PanelInfo, error) {
	return base.getDataFromDatabase(params, services)
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

// getDataFromDatabase 從資料庫取得頁面需要顯示的資料，回傳每一筆資料資訊、欄位資訊、可過濾欄位資訊...等
func (base *BaseTable) getDataFromDatabase(params parameter.Parameters, services service.List) (PanelInfo, error) {
	var (
		connection     = base.getConnectionByService(services)
		ids            = params.FindPKs() // FindPKs 取得__pk的值(多個)
		countStatement string
		queryStatement string
		primaryKey     = base.Informatoin.Table + "." + base.PrimaryKey.Name // 主鍵
		wheres         = ""
		args           = make([]interface{}, 0)
		whereArgs      = make([]interface{}, 0)
		existKeys      = make([]string, 0)
		size           int
	)

	if len(ids) > 0 {
		queryStatement = "select %s from %s%s where " + primaryKey + " in (%s) %s ORDER BY %s.%s %s"
		countStatement = "select count(*) from %s %s where " + primaryKey + " in (%s)"
	} else {
		queryStatement = "select %s from %s%s %s %s order by %s.%s %s LIMIT ? OFFSET ?"
		countStatement = "select count(*) from %s %s %s"
	}

	// 取得所有欄位
	columns, _ := base.getColumns(base.Informatoin.Table, services)

	fieldList, fields, joinFields, joins, joinTables, filterForm := base.getFieldInformationAndJoinOrderAndFilterForm(params, columns)

	// 加上主鍵
	fields += primaryKey
	// 所有欄位
	allFields := fields

	if joinFields != "" {
		// 加上join其他表的欄位(ex: group_concat(roles.`name` separator 'CkN694kH') as roles_join_name,)
		allFields += "," + joinFields[:len(joinFields)-1]
	}

	if len(ids) > 0 {
		for _, value := range ids {
			if value != "" {
				wheres += "?,"
				args = append(args, value)
			}
		}
		wheres = wheres[:len(wheres)-1]
	} else {
		wheres, whereArgs = base.Informatoin.Wheres.WhereStatement(wheres, whereArgs, existKeys, columns)
		if wheres != "" {
			wheres = " where " + wheres
		}

		if connection.Name() == "mysql" {
			pageSizeInt, _ := strconv.Atoi(params.PageSize)
			pageInt, _ := strconv.Atoi(params.Page)
			args = append(args, pageSizeInt, (pageInt-1)*(pageSizeInt))
		}
	}

	groupBy := ""
	if len(joinTables) > 0 {
		if connection.Name() == "mysql" {
			groupBy = " GROUP BY " + primaryKey
		}
	}

	// sql 語法
	queryCmd := fmt.Sprintf(queryStatement, allFields, base.Informatoin.Table, joins, wheres, groupBy,
		base.Informatoin.Table, params.SortField, params.SortType)

	res, err := connection.Query(queryCmd, args...)
	if err != nil {
		return PanelInfo{}, err
	}

	// 頁面上顯示的所有資料
	infoList := make([]map[string]types.InfoItem, 0)
	for i := 0; i < len(res); i++ {
		infoList = append(infoList, base.getTemplateDataModel(res[i], params, columns))
	}

	// 計算資料數
	if len(ids) == 0 {
		countCmd := fmt.Sprintf(countStatement, base.Informatoin.Table, joins, wheres)

		total, err := connection.Query(countCmd, whereArgs...)
		if err != nil {
			return PanelInfo{}, err
		}
		if base.connectionDriver == "mysql" {
			size = int(total[0]["count(*)"].(int64))
		}
	}

	// 設置Paginator.option(在被選中顯示資料筆數的地方加上select)
	paginator := paginator.GetPaginatorInformation(size, params)
	paginator.PageSizeList = base.Informatoin.GetPageSizeList()
	paginator.Option = make(map[string]template.HTML, len(paginator.PageSizeList))
	for i := 0; i < len(paginator.PageSizeList); i++ {
		paginator.Option[paginator.PageSizeList[i]] = template.HTML("")
	}
	paginator.Option[params.PageSize] = template.HTML("select")

	return PanelInfo{
		InfoList:       infoList,
		FieldList:      fieldList,
		Paginator:      paginator,
		PrimaryKey:     base.PrimaryKey.Name,
		Title:          base.Informatoin.Title,
		FilterFormData: filterForm,
		Description:    base.Informatoin.Description,
	}, nil
}

// getTemplateDataModel 取得並處理模板的每一筆資料
func (base *BaseTable) getTemplateDataModel(res map[string]interface{}, params parameter.Parameters, columns []string) map[string]types.InfoItem {
	var templateDataModel = make(map[string]types.InfoItem)
	headField := ""

	// 取得id
	primaryKeyValue := db.GetValueFromDatabaseType(base.PrimaryKey.Type, res[base.PrimaryKey.Name])

	for _, field := range base.Informatoin.FieldList {
		headField = field.Field

		// 如果有關聯其他表
		if field.Joins.Valid() {
			// ex: roles_join_name
			headField = field.Joins.Last().JoinTable + "_join_" + field.Field
		}

		if field.Hide {
			continue
		}
		if !utils.InArrayWithoutEmpty(params.Columns, headField) {
			continue
		}

		typeName := field.TypeName
		if field.Joins.Valid() {
			typeName = db.Varchar
		}

		// 每個欄位的值
		combineValue := db.GetValueFromDatabaseType(typeName, res[headField]).String()

		var value interface{}
		if len(columns) == 0 || utils.InArray(columns, headField) || field.Joins.Valid() {
			value = field.FieldDisplay.DisplayFunc(types.FieldModel{
				ID:    primaryKeyValue.String(),
				Value: combineValue,
				Row:   res,
			})
		} else {
			value = field.FieldDisplay.DisplayFunc(types.FieldModel{
				ID:    primaryKeyValue.String(),
				Value: "",
				Row:   res,
			})
		}

		if valueStr, ok := value.(string); ok {
			templateDataModel[headField] = types.InfoItem{
				Content: template.HTML(valueStr),
				Value:   combineValue,
			}
		} else {
			// 角色欄位會執行
			templateDataModel[headField] = types.InfoItem{
				Content: value.(template.HTML),
				Value:   combineValue,
			}
		}
	}

	// 不管有沒有顯示id(主鍵)欄位，都要加上id的欄位資訊
	primaryKeyField := base.Informatoin.FieldList.GetFieldByFieldName(base.PrimaryKey.Name)
	value := primaryKeyField.FieldDisplay.DisplayFunc(types.FieldModel{
		ID:    primaryKeyValue.String(),
		Value: primaryKeyValue.String(),
		Row:   res,
	})
	if valueStr, ok := value.(string); ok {
		templateDataModel[base.PrimaryKey.Name] = types.InfoItem{
			Content: template.HTML(valueStr),
			Value:   primaryKeyValue.String(),
		}
	} else {
		// 角色欄位會執行
		templateDataModel[base.PrimaryKey.Name] = types.InfoItem{
			Content: value.(template.HTML),
			Value:   primaryKeyValue.String(),
		}
	}

	return templateDataModel
}

// getFieldInformationAndJoinOrderAndFilterForm 取得欄位資訊、join的語法及table、可過濾欄位資訊
func (base *BaseTable) getFieldInformationAndJoinOrderAndFilterForm(params parameter.Parameters, columns []string) (types.FieldList,
	string, string, string, []string, []types.FormField) {
	return base.Informatoin.FieldList.GetFieldInformationAndJoinOrderAndFilterForm(types.TableInfo{
		Table:      base.Informatoin.Table,
		Delimiter:  base.getDelimiter(),
		Driver:     base.connectionDriver,
		PrimaryKey: base.PrimaryKey.Name,
	}, params, columns)
}

// getDelimiter 取得分隔符號
func (base *BaseTable) getDelimiter() string {
	if base.connectionDriver == "mysql" {
		return "'"
	}
	return ""
}
