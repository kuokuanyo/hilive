package types

import "hilive/modules/db"

// FieldList 所有欄位資訊
type FieldList []Field

// Field 欄位head、field、typename...等細節資訊
type Field struct {
	Header     string
	Field      string
	TypeName   db.DatabaseType
	Joins      Joins
	SortAble   bool
	EditAble   bool
	FilterAble bool
	Hide       bool
}

// InformationPanel 資訊面板
type InformationPanel struct {
	FieldList   FieldList
	Table       string
	Title       string
	Description string
	SortField   string
	Wheres      Wheres
	PageSizeList []int // 頁面顯示資料數
	DefaultPageSize int // 顯示頁數
}

// Wheres where的資訊
type Wheres []Where

// Where 條件資訊
type Where struct {
	Join      string
	Field     string
	Operation string
	Arg       interface{}
}

// Joins join其他表的資訊
type Joins []Join

// Join join資訊
type Join struct {
	Table     string
	Field     string
	JoinField string
	BaseTable string
}

// primaryKey 紀錄主鍵及主鍵type
type primaryKey struct {
	Type db.DatabaseType
	Name string
}

// DefaultInformationPanel 預設DefaultInformationPanel
func DefaultInformationPanel(pk string) *InformationPanel {
	return &InformationPanel{
		PageSizeList: []int{10, 20, 30, 50, 100},
		DefaultPageSize: 10,
		Wheres: make([]Where, 0),
		SortField: pk,
	}
}

// Valid 判斷是否有設置Joins
func (j Joins) Valid() bool {
	for i := 0; i < len(j); i++ {
		if j[i].Table != "" && j[i].Field != "" && j[i].JoinField != "" && j[i].BaseTable != "" {
			return true
		}
	}
	return false
}