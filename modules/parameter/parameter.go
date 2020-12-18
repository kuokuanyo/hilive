package parameter

import "strings"

// Parameters 紀錄頁面資訊，頁數及頁數size...等
type Parameters struct {
	Page      string
	PageSize  string
	SortField string
	SortType  string
	Columns   []string
	Fields    map[string][]string
}

// DefaultParameters 預設Parameters(struct)
func DefaultParameters() Parameters {
	return Parameters{
		Page:     "1",
		PageSize: "10",
		Fields:   make(map[string][]string),
	}
}

// SetFieldPKByJoinParam 將參數(多個)join成string加入Parameters.Fields["__pk"]
func (param Parameters) SetFieldPKByJoinParam(id ...string) Parameters {
	param.Fields["__pk"] = []string{strings.Join(id, ",")}
	return param
}

// FindPK 取得__pk的值(單個)
func (param Parameters) FindPK() string {
	value, ok := param.Fields["__pk"]
	if ok && len(value) > 0 {
		return strings.Split(value[0], ",")[0]
	}
	return ""
}
