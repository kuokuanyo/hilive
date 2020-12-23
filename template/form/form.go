package form

import "html/template"

// Type uint8
type Type uint8

const (
	// Default Default
	Default Type = iota
	// Text Text
	Text
	// SelectSingle SelectSingle
	SelectSingle
	// Select Select
	Select
	// IconPicker IconPicker
	IconPicker
	// SelectBox SelectBox
	SelectBox
	// File File
	File
	// Multifile Multifile
	Multifile
	// Password Password
	Password
	// RichText RichText
	RichText
	// Datetime Datetime
	Datetime
	// DatetimeRange DatetimeRange
	DatetimeRange
	// Radio Radio
	Radio
	// Checkbox Checkbox
	Checkbox
	// CheckboxStacked CheckboxStacked
	CheckboxStacked
	// CheckboxSingle CheckboxSingle
	CheckboxSingle
	// Email Email
	Email
	// Date Date
	Date
	// DateRange DateRange
	DateRange
	// URL URL
	URL
	// IP IP
	IP
	// Color Color
	Color
	// Array Array
	Array
	// Currency Currency
	Currency
	// Rate Rate
	Rate
	// Number Number
	Number
	// Table Table
	Table
	// NumberRange NumberRange
	NumberRange
	// TextArea TextArea
	TextArea
	// Custom Custom
	Custom
	// Switch Switch
	Switch
	// Code Code
	Code
	// Slider Slider
	Slider
)

// GetFieldOptions 不同欄位類型設置不同處理方式
func (t Type) GetFieldOptions(field string) (map[string]interface{}, map[string]interface{}, template.JS) {
	switch t {
	case File, Multifile:
		return map[string]interface{}{
			"overwriteInitial":     true,
			"initialPreviewAsData": true,
			"browseLabel":          "瀏覽",
			"showRemove":           false,
			"previewClass":         "preview-" + field,
			"showUpload":           false,
			"allowedFileTypes":     []string{"image"},
		}, nil, ""
	case Slider:
		return map[string]interface{}{
			"type":     "single",
			"prettify": false,
			"hasGrid":  true,
			"max":      100,
			"min":      1,
			"step":     1,
			"postfix":  "",
		}, nil, ""
	case DatetimeRange:
		format := "YYYY-MM-DD HH:mm:ss"
		if t == DateRange {
			format = "YYYY-MM-DD"
		}
		m := map[string]interface{}{
			"format": format,
		}
		m1 := map[string]interface{}{
			"format":     format,
			"useCurrent": false,
		}
		return m, m1, ""
	case Datetime:
		format := "YYYY-MM-DD HH:mm:ss"
		if t == Date {
			format = "YYYY-MM-DD"
		}
		m := map[string]interface{}{
			"format":           format,
			"allowInputToggle": true,
		}
		return m, nil, ""
	case Date:
		format := "YYYY-MM-DD HH:mm:ss"
		if t == Date {
			format = "YYYY-MM-DD"
		}
		m := map[string]interface{}{
			"format":           format,
			"allowInputToggle": true,
		}
		return m, nil, ""
	case DateRange:
		format := "YYYY-MM-DD HH:mm:ss"
		if t == DateRange {
			format = "YYYY-MM-DD"
		}
		m := map[string]interface{}{
			"format": format,
		}
		m1 := map[string]interface{}{
			"format":     format,
			"useCurrent": false,
		}
		return m, m1, ""
	case Code:
		return nil, nil, `
	theme = "monokai";
	font_size = 14;
	language = "html";
	options = {useWorker: false};
`
	}
	return nil, nil, ""
}

// FixOptions 判斷欄位類型後處理，回傳map
func (t Type) FixOptions(m map[string]interface{}) map[string]interface{} {
	switch t {
	case Slider:
		if _, ok := m["type"]; !ok {
			m["type"] = "single"
		}
		if _, ok := m["prettify"]; !ok {
			m["prettify"] = false
		}
		if _, ok := m["hasGrid"]; !ok {
			m["hasGrid"] = true
		}
		return m
	}
	return m
}

// IsFile 表單類型是否為檔案
func (t Type) IsFile() bool {
	return t == File || t == Multifile
}

// IsSelect 判斷type是否為select
func (t Type) IsSelect() bool {
	return t == Select || t == SelectSingle || t == SelectBox || t == Radio || t == Switch ||
		t == Checkbox || t == CheckboxStacked || t == CheckboxSingle
}

// SelectedLabel 判斷類別後加入html陣列
func (t Type) SelectedLabel() []template.HTML {
	if t == Select || t == SelectSingle || t == SelectBox {
		return []template.HTML{"selected", ""}
	}
	if t == Radio || t == Switch || t == Checkbox || t == CheckboxStacked || t == CheckboxSingle {
		return []template.HTML{"checked", ""}
	}
	return []template.HTML{"", ""}
}

// IsRange 是否值為範圍
func (t Type) IsRange() bool {
	return t == DatetimeRange || t == NumberRange
}

// IsMultiSelect 是否有多個選擇
func (t Type) IsMultiSelect() bool {
	return t == Select || t == SelectBox || t == Checkbox || t == CheckboxStacked
}

// String 將type轉換成string
func (t Type) String() string {
	switch t {
	case Default:
		return "default"
	case Text:
		return "text"
	case SelectSingle:
		return "select_single"
	case Select:
		return "select"
	case IconPicker:
		return "iconpicker"
	case SelectBox:
		return "selectbox"
	case File:
		return "file"
	case Table:
		return "table"
	case Multifile:
		return "multi_file"
	case Password:
		return "password"
	case RichText:
		return "richtext"
	case Rate:
		return "rate"
	case Checkbox:
		return "checkbox"
	case CheckboxStacked:
		return "checkbox_stacked"
	case CheckboxSingle:
		return "checkbox_single"
	case Date:
		return "datetime"
	case DateRange:
		return "datetime_range"
	case Datetime:
		return "datetime"
	case DatetimeRange:
		return "datetime_range"
	case Radio:
		return "radio"
	case Slider:
		return "slider"
	case Array:
		return "array"
	case Email:
		return "email"
	case URL:
		return "url"
	case IP:
		return "ip"
	case Color:
		return "color"
	case Currency:
		return "currency"
	case Number:
		return "number"
	case NumberRange:
		return "number_range"
	case TextArea:
		return "textarea"
	case Custom:
		return "custom"
	case Switch:
		return "switch"
	case Code:
		return "code"
	default:
		panic("wrong form type")
	}
}
