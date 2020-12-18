package sql

// Where sql where
type Where struct {
	Operation string
	Field string
	Value string
}

// Join sql join
type Join struct {
	Table string
	FieldA string
	FieldB string
	Operation string
}

// RawUpdate 包含表達式及參數
type RawUpdate struct {
	Expression string
	Args []interface{}
}

// Value map[欄位]數值
type Value map[string]interface{}

// FilterCondition 過濾條件
type FilterCondition struct {
	Fields []string
	Functions []string
	TableName string
	Wheres []Where
	WhereRaws string
	UpdateRaws []RawUpdate
	Leftjoins []Join
	Args []interface{}
	Order string
	Offset string
	Limit string
	Group string
	Statement string
	Values Value
}

// CRUD 資料庫CRUD等方法
type CRUD interface {
	GetName() string

	// get all columns
	ShowColumns(table string) string

	// get tables of databases
	ShowTables() string

	// Insert
	Insert(condition *FilterCondition) string

	// Delete
	Delete(condition *FilterCondition) string

	// Update
	Update(condition *FilterCondition) string

	// Select
	Select(condition *FilterCondition) string
}

// GetCRUDByDriver 設置分隔符與取得CRUD(interface)
func GetCRUDByDriver(driver string) CRUD {
	switch driver {
	case "mysql":
		return mysql{
			delimiter: "`",
		}
	}
	panic("無資料庫引擎的CRUD(interface)")
}
