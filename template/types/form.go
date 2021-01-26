package types

import (
	"encoding/json"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/modules/service"
	"hilive/modules/utils"
	"hilive/template/form"
	"html/template"
	"strings"
)

// FormPostFunc 表單Post功能函式
type FormPostFunc func(values form2.Values) error

// PostType uint8
type PostType uint8

// FormPanel 表單面板
type FormPanel struct {
	FieldList         FormFields
	curFieldListIndex int // 欄位的index位置
	Table             string
	Title             string
	Description       string
	InsertFunc        FormPostFunc
	UpdateFunc        FormPostFunc
	primaryKey        primaryKey
}

// FormFields 紀錄所有表單欄位資訊
type FormFields []FormField

// FormField 表單欄位資訊
type FormField struct {
	Field                string          `json:"field"`
	TypeName             db.DatabaseType `json:"type_name"`
	Header               string          `json:"head"`
	FormType             form.Type       `json:"form_type"`
	Value                template.HTML   `json:"value"`
	Value2               string          `json:"value_2"` // 儲存檔案(ex:頭像)的資訊
	Placeholder          string          `json:"placeholder"`
	Editable             bool            `json:"editable"`      // 允許編輯
	NotAllowAdd          bool            `json:"not_allow_add"` // 不允許增加
	Must                 bool            `json:"must"`          // 該欄位必填
	Hide                 bool            `json:"hide"`
	Default              template.HTML   `json:"Default"`
	Joins                Joins           `json:"-"`
	FieldDisplay         FieldDisplay
	FieldOptions         FieldOptions
	FieldOptionFromTable FieldOptionFromTable
	OptionExt            template.JS   `json:"option_ext"`   // 不同欄位類型處理方式
	OptionExt2           template.JS   `json:"option_ext_2"` // 不同欄位類型處理方式
	HelpMsg              template.HTML `json:"help_msg"`     // 欄位提示訊息
}

// FieldOptions 紀錄所有表單欄位選單
type FieldOptions []FieldOption

// FieldOption 紀錄表單欄位選單
type FieldOption struct {
	Text          string        `json:"text"`
	Value         string        `json:"value"`
	TextHTML      template.HTML `json:"-"`
	Selected      bool          `json:"-"` // 選項是否被選擇
	SelectedLabel template.HTML `json:"-"` // 選項的label
}

// FieldOptionFromTable 紀錄表單欄位選單(選單名稱由資料表取得)
type FieldOptionFromTable struct {
	Table      string
	TextField  string
	ValueField string
}

// DefaultFormPanel 預設FormPanel
func DefaultFormPanel() *FormPanel {
	return &FormPanel{
		curFieldListIndex: -1,
	}
}

// SetPrimaryKey 設置主鍵
func (f *FormPanel) SetPrimaryKey(name string, typ db.DatabaseType) *FormPanel {
	f.primaryKey = primaryKey{Name: name, Type: typ}
	return f
}

// AddField 增加欄位資訊
func (f *FormPanel) AddField(header, field string, fieldType db.DatabaseType, formType form.Type) *FormPanel {
	f.FieldList = append(f.FieldList, FormField{
		Header:      header,
		Field:       field,
		TypeName:    fieldType,
		Editable:    true,
		Hide:        false,
		Placeholder: " 輸入 " + header,
		FormType:    formType,
		FieldDisplay: FieldDisplay{
			DisplayFunc: func(value FieldModel) interface{} {
				return value.Value
			},
		},
	})
	f.curFieldListIndex++

	// 不同欄位類型設置不同的處理方式
	op1, op2, js := formType.GetFieldOptions(field)
	f.FieldOptionExt(op1)
	f.FieldOptionExt2(op2)
	f.FieldOptionExtJS(js)

	// 如果表單欄位類型為選項(select)，設置DisplayFunc
	setDisplayFuncOfFormType(f, formType)
	return f
}

// FieldWithValue 取得欄位的值、預設值並設置至FormPanel.FormFields(帶有資料值)
func (f *FormPanel) FieldWithValue(pk, id string, columns []string, res map[string]interface{}, services service.List, sql func(services service.List) *db.SQL) FormFields {
	var (
		list  = make(FormFields, 0)
		hasPk = false
	)
	// 取得值將值更新至FormField
	for _, field := range f.FieldList {
		// if field.Editable {
			// 取得欄位的值
			dataValue := field.GetDataValue(columns, res[field.Field])
			// 將取得的欄位值放入FormField中
			list = append(list, *(field.UpdateValue(id, dataValue, res, sql(services))))
			// 判斷是否有主鍵
			if field.Field == pk {
				hasPk = true
			}
		// }
	}

	if !hasPk {
		list = append(list, FormField{
			Header:   pk,
			Field:    pk,
			Value:    template.HTML(id),
			FormType: form.Default,
			Hide:     true,
		})
	}
	return list
}

// SetTable 設置FormPanel.Table
func (f *FormPanel) SetTable(table string) *FormPanel {
	f.Table = table
	return f
}

// SetTitle 設置FormPanel.Title
func (f *FormPanel) SetTitle(title string) *FormPanel {
	f.Title = title
	return f
}

// SetDescription 設置FormPanel.Description
func (f *FormPanel) SetDescription(desc string) *FormPanel {
	f.Description = desc
	return f
}

// FieldNotAllowAdd 該欄位不允許增加
func (f *FormPanel) FieldNotAllowAdd() *FormPanel {
	f.FieldList[f.curFieldListIndex].NotAllowAdd = true
	return f
}

// FieldNotAllowEdit 該欄位不允許編輯
func (f *FormPanel) FieldNotAllowEdit() *FormPanel {
	f.FieldList[f.curFieldListIndex].Editable = false
	return f
}

// SetFieldMust 該欄位一定要填
func (f *FormPanel) SetFieldMust() *FormPanel {
	f.FieldList[f.curFieldListIndex].Must = true
	return f
}

// SetFieldOptions 設置Field.FieldOptions
func (f *FormPanel) SetFieldOptions(options FieldOptions) *FormPanel {
	f.FieldList[f.curFieldListIndex].FieldOptions = options
	return f
}

// SetInsertFunc 設置新增函式
func (f *FormPanel) SetInsertFunc(fn FormPostFunc) *FormPanel {
	f.InsertFunc = fn
	return f
}

// SetUpdateFunc 設置更新函式
func (f *FormPanel) SetUpdateFunc(fn FormPostFunc) *FormPanel {
	f.UpdateFunc = fn
	return f
}

// SetAllowAddValueOfField 處理並設置表單欄位細節資訊(允許增加的表單欄位)，回傳FormFields
func (f *FormPanel) SetAllowAddValueOfField(services service.List, sql ...func(services service.List) *db.SQL) FormFields {
	var list = make(FormFields, 0)
	for _, v := range f.FieldList {
		if !v.NotAllowAdd {
			v.Editable = true
			if len(sql) > 0 {
				v.Value = v.Default
				// UpdateValue 如果表單類型為選單，處理FormField.FieldOptions並設置已被選擇的選項、label 如果不是選單類型則設定FormField.Value
				list = append(list, *(v.UpdateValue("", string(v.Value), make(map[string]interface{}), sql[0](services))))
			} else {
				v.Value = v.Default
				list = append(list, *(v.UpdateValue("", string(v.Value), make(map[string]interface{}), nil)))
			}
		}
	}
	return list
}

// SetFieldOptionFromTable 設置Field.FieldOptionFromTable(選單名稱由資料表取得)
func (f *FormPanel) SetFieldOptionFromTable(table, textFieldName, valueFieldName string) *FormPanel {
	f.FieldList[f.curFieldListIndex].FieldOptionFromTable = FieldOptionFromTable{
		Table:      table,
		TextField:  textFieldName,
		ValueField: valueFieldName,
	}
	return f
}

// SetDisplayFunc 設置欄位過濾函式至DisplayFunc
func (f *FormPanel) SetDisplayFunc(filter FieldFilterFunc) *FormPanel {
	f.FieldList[f.curFieldListIndex].FieldDisplay.DisplayFunc = filter
	return f
}

// SetFieldDefault 設定預設值
func (f *FormPanel) SetFieldDefault(def string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Default = template.HTML(def)
	return f
}

// FieldOptionExt 處理不同欄位類型後設置至FormField.OptionExt
func (f *FormPanel) FieldOptionExt(m map[string]interface{}) *FormPanel {
	if m == nil {
		return f
	}
	if f.FieldList[f.curFieldListIndex].FormType == form.Code {
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(fmt.Sprintf(`
	theme = "%s";
	font_size = %s;
	language = "%s";
	options = %s;
`, m["theme"], m["font_size"], m["language"], m["options"]))
		return f
	}

	m = f.FieldList[f.curFieldListIndex].FormType.FixOptions(m)
	s, _ := json.Marshal(m)

	if f.FieldList[f.curFieldListIndex].OptionExt != template.JS("") {
		ss := string(f.FieldList[f.curFieldListIndex].OptionExt)
		ss = strings.Replace(ss, "}", "", strings.Count(ss, "}"))
		ss = strings.TrimRight(ss, " ")
		ss += ","
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(ss) + template.JS(strings.Replace(string(s), "{", "", 1))
	} else {
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(string(s))
	}
	return f
}

// FieldOptionExt2 處理不同欄位類型後設置至FormField.OptionExt2
func (f *FormPanel) FieldOptionExt2(m map[string]interface{}) *FormPanel {
	if m == nil {
		return f
	}

	m = f.FieldList[f.curFieldListIndex].FormType.FixOptions(m)
	s, _ := json.Marshal(m)

	if f.FieldList[f.curFieldListIndex].OptionExt2 != template.JS("") {
		ss := string(f.FieldList[f.curFieldListIndex].OptionExt2)
		ss = strings.Replace(ss, "}", "", strings.Count(ss, "}"))
		ss = strings.TrimRight(ss, " ")
		ss += ","
		f.FieldList[f.curFieldListIndex].OptionExt2 = template.JS(ss) + template.JS(strings.Replace(string(s), "{", "", 1))
	} else {
		f.FieldList[f.curFieldListIndex].OptionExt2 = template.JS(string(s))
	}
	return f
}

// FieldOptionExtJS 處理不同欄位類型後設置至FormField.OptionExt
func (f *FormPanel) FieldOptionExtJS(js template.JS) *FormPanel {
	if js != template.JS("") {
		f.FieldList[f.curFieldListIndex].OptionExt = js
	}
	return f
}

// UpdateValue 如果表單類型為選單，處理FormField.FieldOptions並設置已被選擇的選項、label
// 如果不是選單類型則設定FormField.Value
func (f *FormField) UpdateValue(id, val string, res map[string]interface{}, s *db.SQL) *FormField {
	m := FieldModel{
		ID:    id,
		Value: val,
		Row:   res,
	}

	// 判斷表單類型是否是選單
	if f.FormType.IsSelect() {
		// setFieldOptionFromTableBySQL 處理表單類型為選單的欄位並且選單名稱是由資料表取得的(ex:角色欄位)，設置FormField.FieldOptions
		f.setFieldOptionFromTableBySQL(s)
		// SetSelected 設置已被選擇的選項及label
		f.FieldOptions.SetSelected(f.FieldDisplay.ToDisplayFunc(m), f.FormType.SelectedLabel())
	} else {
		f.Value = f.FieldDisplay.DisplayFuncToHTML(m)
		// 判斷表單類型是否為檔案，ex:頭像欄位
		if f.FormType.IsFile() {
			if f.Value != template.HTML("") {
				f.Value2 = config.GetStore().URL(string(f.Value))
			}
		}
	}
	return f
}

// setFieldOptionFromTableBySQL 處理表單類型為選單的欄位並且選單名稱是由資料表取得的(ex:角色欄位)，設置FormField.FieldOptions(選項)
func (f *FormField) setFieldOptionFromTableBySQL(s *db.SQL) {
	if s != nil && f.FieldOptionFromTable.Table != "" && len(f.FieldOptions) == 0 {
		queryRes, err := s.Table(f.FieldOptionFromTable.Table).
			Select(f.FieldOptionFromTable.ValueField, f.FieldOptionFromTable.TextField).All()
		if err == nil {
			for _, item := range queryRes {
				f.FieldOptions = append(f.FieldOptions, FieldOption{
					Value: fmt.Sprintf("%v", item[f.FieldOptionFromTable.ValueField]), // ex: id值
					Text:  fmt.Sprintf("%v", item[f.FieldOptionFromTable.TextField]),  // ex: slug的值
				})
			}
		}
	}
}

// GetDataValue 取得欄位的值
func (f *FormField) GetDataValue(columns []string, v interface{}) string {
	return utils.AorB(utils.InArray(columns, f.Field),
		string(db.GetValueFromDatabaseType(f.TypeName, v)), "")
}

// SetFieldHelpMsg 設置欄位提示訊息
func (f *FormPanel) SetFieldHelpMsg(s template.HTML) *FormPanel {
	f.FieldList[f.curFieldListIndex].HelpMsg = s
	return f
}

// SetSelected 設置已被選擇的選項及label
func (f FieldOptions) SetSelected(val interface{}, labels []template.HTML) FieldOptions {
	if valArr, ok := val.([]string); ok {
		for k := range f {
			text := f[k].Text
			if text == "" {
				text = string(f[k].TextHTML)
			}
			f[k].Selected = utils.InArray(valArr, f[k].Value) || utils.InArray(valArr, text)
			if f[k].Selected {
				f[k].SelectedLabel = labels[0]
			} else {
				f[k].SelectedLabel = labels[1]
			}
		}
	} else {
		for k := range f {
			text := f[k].Text
			if text == "" {
				text = string(f[k].TextHTML)
			}
			f[k].Selected = f[k].Value == val || text == val
			if f[k].Selected {
				f[k].SelectedLabel = labels[0]
			} else {
				f[k].SelectedLabel = labels[1]
			}
		}
	}
	return f
}
