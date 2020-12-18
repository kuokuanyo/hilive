package types

import (
	"fmt"
	"hilive/template/form"
	"html/template"
	"strings"
)

// FieldFilterFunc 欄位過濾函式
type FieldFilterFunc func(value FieldModel) interface{}

// FieldDisplay 欄位過濾function
type FieldDisplay struct {
	DisplayFunc FieldFilterFunc
}

// FieldModel 欄位的ID、value...等資訊
type FieldModel struct {
	// The primaryKey of the table.
	ID string
	// The value of the single query result.
	Value string
	// The current row data.
	Row map[string]interface{}
}

// setDisplayFuncOfFormType 如果表單欄位類型為選項(select)，設置DisplayFunc
func setDisplayFuncOfFormType(f *FormPanel, typ form.Type) {
	if typ.IsSelect() {
		f.FieldList[f.curFieldListIndex].FieldDisplay.DisplayFunc = func(value FieldModel) interface{} {
			return strings.Split(value.Value, ",")
		}
	}
}

// ToDisplayFunc 執行Display得Function，判斷是否為選單欄位並處理值
func (f FieldDisplay) ToDisplayFunc(value FieldModel) interface{} {
	// 執行function，如果是選單欄位結果會是空陣列，不然為空值
	val := f.DisplayFunc(value)
	if f.IsNotSelect(val) {
		return FieldModel{
			Row:   value.Row,
			Value: fmt.Sprintf("%v", val),
			ID:    value.ID,
		}
	}
	return val
}

// DisplayFuncToHTML 執行display function後轉換成hmtl
func (f FieldDisplay) DisplayFuncToHTML(value FieldModel) template.HTML {
	v := f.DisplayFunc(value)
	if h, ok := v.(template.HTML); ok {
		return h
	} else if s, ok := v.(string); ok {
		return template.HTML(s)
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return template.HTML(arr[0])
	} else if v != nil {
		return ""
	} else {
		return ""
	}
}

// IsNotSelect 判斷是否不是選單欄位
func (f FieldDisplay) IsNotSelect(v interface{}) bool {
	switch v.(type) {
	case template.HTML:
		return false
	case []string:
		return false
	case [][]string:
		return false
	default:
		return true
	}
}
