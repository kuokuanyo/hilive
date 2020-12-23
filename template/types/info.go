package types

import (
	"hilive/modules/db"
	"hilive/modules/parameter"
	"hilive/modules/service"
	"hilive/modules/utils"
	"hilive/template/form"
	"html/template"
	"strconv"
	"strings"
)

// FieldList 所有欄位資訊
type FieldList []Field

// DeleteFunc 刪除函式
type DeleteFunc func(ids []string) error

// JoinFieldValueDelimiter 分隔符號
var JoinFieldValueDelimiter = utils.UUID(8)

// Field 欄位head、field、typename...等細節資訊
type Field struct {
	Header           string
	Field            string
	TypeName         db.DatabaseType
	Joins            Joins
	SortAble         bool
	EditAble         bool
	FilterAble       bool
	Hide             bool
	FieldDisplay     FieldDisplay
	FilterFormFields []FilterFormField // 可以過濾欄位的表單資訊
}

// InformationPanel 資訊面板
type InformationPanel struct {
	FieldList         FieldList
	curFieldListIndex int
	Table             string
	Title             string
	Description       string
	SortField         string
	Wheres            Wheres
	PageSizeList      []int // 頁面顯示資料數
	DefaultPageSize   int   // 顯示頁數
	primaryKey        primaryKey
	DeleteFunc        DeleteFunc
}

// FilterFormField 可以過濾欄位的表單資訊
type FilterFormField struct {
	FormType    form.Type
	Header      string
	Placeholder string
}

// TableInfo 資料表資訊
type TableInfo struct {
	Table      string
	PrimaryKey string
	Delimiter  string
	Driver     string
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
	JoinTable string
	Field     string
	JoinField string
	BaseTable string
}

// InfoList 顯示在界面上的所有資料
type InfoList []map[string]InfoItem

// InfoItem 顯示在界面上的資料
type InfoItem struct {
	Content template.HTML `json:"content"`
	Value   string        `json:"value"`
}

// primaryKey 紀錄主鍵及主鍵type
type primaryKey struct {
	Type db.DatabaseType
	Name string
}

// SetTable 設置資料表
func (i *InformationPanel) SetTable(table string) *InformationPanel {
	i.Table = table
	return i
}

// SetTitle 設置主題名稱
func (i *InformationPanel) SetTitle(title string) *InformationPanel {
	i.Title = title
	return i
}

// SetDescription 設置描述
func (i *InformationPanel) SetDescription(desc string) *InformationPanel {
	i.Description = desc
	return i
}

// SetPrimaryKey 設置主鍵至InformationPanel
func (i *InformationPanel) SetPrimaryKey(name string, typ db.DatabaseType) *InformationPanel {
	i.primaryKey = primaryKey{
		Name: name,
		Type: typ,
	}
	return i
}

// DefaultInformationPanel 預設DefaultInformationPanel
func DefaultInformationPanel(pk string) *InformationPanel {
	return &InformationPanel{
		PageSizeList:      []int{10, 20, 30, 50, 100},
		curFieldListIndex: -1,
		DefaultPageSize:   10,
		Wheres:            make([]Where, 0),
		SortField:         pk,
	}
}

// AddField 增加欄位
func (i *InformationPanel) AddField(header, field string, typeName db.DatabaseType) *InformationPanel {
	i.FieldList = append(i.FieldList, Field{
		Header:   header,
		Field:    field,
		TypeName: typeName,
		SortAble: false,
		Joins:    make(Joins, 0),
		EditAble: true,
		FieldDisplay: FieldDisplay{
			DisplayFunc: func(value FieldModel) interface{} {
				return value.Value
			},
		},
	})
	i.curFieldListIndex++
	return i
}

// GetFieldInformationAndJoinOrderAndFilterForm 取得欄位資訊、join的語法及table、可過濾欄位資訊
func (f FieldList) GetFieldInformationAndJoinOrderAndFilterForm(info TableInfo, params parameter.Parameters, columns []string, sql ...func(services service.List) *db.SQL) (
	FieldList, string, string, string, []string, []FormField) {
	var (
		fieldList  = make(FieldList, 0)
		fields     = ""
		joinFields = ""                   // ex: group_concat(roles.`name` separator 'CkN694kH') as roles_join_name,
		joins      = ""                   // join資料表語法，ex: left join `role_users` on role_users.`user_id` = users.`id` left join....
		joinTables = make([]string, 0)    // ex:{role_users roles}
		filterForm = make([]FormField, 0) // 可過濾欄位資訊
	)

	for _, field := range f {
		// 不是主鍵並且沒有join關聯會執行
		if field.Field != info.PrimaryKey && utils.InArray(columns, field.Field) &&
			!field.Joins.Valid() {
			fields += info.Table + "." + field.Field + ","
		}

		// 有join關聯的欄位會執行
		headField := field.Field
		if field.Joins.Valid() {
			// ex:roles_join_name
			headField = field.Joins.Last().JoinTable + "_join_" + field.Field

			joinFields += db.GetAggregationExpression(info.Driver, field.Joins.Last().JoinTable+"."+field.Field,
				headField, JoinFieldValueDelimiter) + ","

			for _, join := range field.Joins {
				if !utils.InArray(joinTables, join.JoinTable) {
					joinTables = append(joinTables, join.JoinTable)

					// ex: joins =  left join `role_users` on role_users.`user_id` = users.`id` left join....
					joins += " left join " + join.JoinTable + " on " +
						join.JoinTable + "." + join.JoinField + " = " +
						join.BaseTable + "." + join.Field
				}
			}
		}

		if field.FilterAble {
			filterForm = append(filterForm, field.GetFormFieldFromFilterFormFields(params, headField)...)
		}

		if field.Hide {
			continue
		}

		fieldList = append(fieldList, Field{
			Header:   field.Header,
			SortAble: field.SortAble,
			Field:    headField,
			Hide:     !utils.InArrayWithoutEmpty(params.Columns, headField),
			EditAble: field.EditAble,
		})
	}
	return fieldList, fields, joinFields, joins, joinTables, filterForm
}

// GetFormFieldFromFilterFormFields 藉由可過濾欄位資訊(FilterFormField)取得[]FormField
func (f Field) GetFormFieldFromFilterFormFields(params parameter.Parameters, headField string) []FormField {
	var (
		filterFormFields = make([]FormField, 0)
		keySuffix, value string
	)

	for index, filter := range f.FilterFormFields {
		if index > 0 {
			keySuffix = "__index__" + strconv.Itoa(index)
		}

		// 判斷可過濾欄位type
		if filter.FormType.IsRange() {
		} else if filter.FormType.IsMultiSelect() {
		} else {
			// GetFieldValue 取得Parameters.Fields[field]的值，若沒有則回傳""
			value = params.GetFieldValue(headField + keySuffix)
		}

		field := &FormField{
			Field:       headField + keySuffix,
			TypeName:    f.TypeName,
			Header:      filter.Header,
			FormType:    filter.FormType,
			Value:       template.HTML(value),
			Value2:      "",
			Placeholder: filter.Placeholder,
			Editable:    true,
			OptionExt:   template.JS(""),
			OptionExt2:  template.JS(""),
		}
		filterFormFields = append(filterFormFields, *field)
	}
	return filterFormFields
}

// WhereStatement 處理where後回傳
func (whs Wheres) WhereStatement(wheres string, whereArgs []interface{}, existKeys, columns []string) (string, []interface{}) {
	pwheres := ""

	for k, wh := range whs {
		whFieldArr := strings.Split(wh.Field, ".")
		whField := ""
		whTable := ""

		if len(whFieldArr) > 1 {
			whField = whFieldArr[1]
			whTable = whFieldArr[0]
		} else {
			whField = whFieldArr[0]
		}

		if utils.InArray(existKeys, whField) {
			continue
		}

		// TODO: support like operation and join table
		if utils.InArray(columns, whField) {

			joinMark := ""
			if k != len(whs)-1 {
				joinMark = whs[k+1].Join
			}

			if whTable != "" {
				pwheres += whTable + "." + whField + " " + wh.Operation + " ? " + joinMark + " "
			} else {
				pwheres += whField + " " + wh.Operation + " ? " + joinMark + " "
			}
			whereArgs = append(whereArgs, wh.Arg)
		}
	}
	if wheres != "" && pwheres != "" {
		wheres += " and "
	}
	return wheres + pwheres, whereArgs
}

// GetFieldByFieldName 透過參數取得Field資訊
func (f FieldList) GetFieldByFieldName(name string) Field {
	for _, field := range f {
		if field.Field == name {
			return field
		}
		if JoinField(field.Joins.Last().JoinTable, field.Field) == name {
			return field
		}
	}
	return Field{}
}

// FieldSortable 欄位可以排序
func (i *InformationPanel) FieldSortable() *InformationPanel {
	i.FieldList[i.curFieldListIndex].SortAble = true
	return i
}

// FieldFilterable 欄位可以過濾
func (i *InformationPanel) FieldFilterable() *InformationPanel {
	i.FieldList[i.curFieldListIndex].FilterAble = true
	i.FieldList[i.curFieldListIndex].FilterFormFields = append(i.FieldList[i.curFieldListIndex].FilterFormFields,
		FilterFormField{
			FormType:    form.Text,
			Header:      i.FieldList[i.curFieldListIndex].Header,
			Placeholder: "輸入 " + i.FieldList[i.curFieldListIndex].Header,
		})
	return i
}

// FieldJoin 欄位有關聯其他表
func (i *InformationPanel) FieldJoin(join Join) *InformationPanel {
	i.FieldList[i.curFieldListIndex].Joins = append(i.FieldList[i.curFieldListIndex].Joins, join)
	return i
}

// FieldDisplayFunc 設置display function
func (i *InformationPanel) FieldDisplayFunc(filter FieldFilterFunc) *InformationPanel {
	i.FieldList[i.curFieldListIndex].FieldDisplay.DisplayFunc = filter
	return i
}

// SetDeleteFunc 設置刪除函式
func (i *InformationPanel) SetDeleteFunc(fn DeleteFunc) *InformationPanel {
	i.DeleteFunc = fn
	return i
}

// Valid 判斷是否有設置Joins
func (j Joins) Valid() bool {
	for i := 0; i < len(j); i++ {
		if j[i].JoinTable != "" && j[i].Field != "" && j[i].JoinField != "" && j[i].BaseTable != "" {
			return true
		}
	}
	return false
}

// Last 回傳Joins裡的最後一個Join
func (j Joins) Last() Join {
	if len(j) > 0 {
		return j[len(j)-1]
	}
	return Join{}
}

// JoinField return table_join_field
func JoinField(table, field string) string {
	return table + "_join_" + field
}