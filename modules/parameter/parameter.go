package parameter

import (
	"hilive/modules/utils"
	"net/url"
	"strconv"
	"strings"
)

// keys url參數
var keys = []string{"__page", "__pageSize", "__sort", "__columns", "__prefix"}

// operators 運算符號
var operators = map[string]string{
	"like": "like",
	"gr":   ">",
	"gq":   ">=",
	"eq":   "=",
	"ne":   "!=",
	"le":   "<",
	"lq":   "<=",
	"free": "free",
}

// Parameters 紀錄頁面資訊，頁數及頁數size...等
type Parameters struct {
	Page      string
	PageSize  string
	SortField string
	SortType  string
	Columns   []string
	Fields    map[string][]string
	URLPath   string
}

// DefaultParameters 預設Parameters(struct)
func DefaultParameters() Parameters {
	return Parameters{
		Page:     "1",
		PageSize: "10",
		Fields:   make(map[string][]string),
	}
}

// SetPKs 將參數(多個)join成string加入Parameters.Fields["__pk"]
func (param Parameters) SetPKs(id ...string) Parameters {
	param.Fields["__pk"] = []string{strings.Join(id, ",")}
	return param
}

// SetPage 設置Parameters.Page
func (param Parameters) SetPage(page string) Parameters {
	param.Page = page
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

// FindPKs 取得__pk的值(多個)
func (param Parameters) FindPKs() []string {
	value, ok := param.Fields["__pk"]
	if ok && len(value) > 0 {
		return strings.Split(value[0], ",")
	}
	return []string{}
}

// GetParamFromURL 解析URL後設置頁面資訊
func GetParamFromURL(urlStr string, defaultPageSize int) Parameters {
	// 解析url
	u, err := url.Parse(urlStr)
	if err != nil {
		return DefaultParameters()
	}
	return GetParam(u, defaultPageSize)
}

// GetParam 設置頁面資訊
func GetParam(u *url.URL, defaultPageSize int) Parameters {
	// Query從url取得設定參數
	// ex: map[__columns:[id,username,name,goadmin_roles_goadmin_join_name,created_at,updated_at] __page:[1] __pageSize:[10]  __sort:[id] __sort_type:[desc] ...]
	values := u.Query()

	primaryKey := "id"
	SortType := "asc"

	page := getDefault(values, "__page", "1")
	pageSize := getDefault(values, "__pageSize", strconv.Itoa(defaultPageSize))
	sortField := getDefault(values, "__sort", primaryKey)
	sortType := getDefault(values, "__sort_type", SortType)

	columns := getDefault(values, "__columns", "")
	columnsArr := make([]string, 0)
	if columns != "" {
		columns, _ = url.QueryUnescape(columns)
		columnsArr = strings.Split(columns, ",")
	}

	// 將除了keys以外(其他過濾條件)的參數加入fields中
	fields := make(map[string][]string)
	for key, value := range values {
		if !utils.InArray(keys, key) && len(value) > 0 && value[0] != "" {
			if key == "__sort_type" {
				if value[0] != "desc" && value[0] != "asc" {
					fields[key] = []string{"asc"}
				}
			} else {
				if strings.Contains(key, "__operator__") &&
					values.Get(strings.Replace(key, "__operator__", "", -1)) == "" {
					continue
				}
				fields[strings.Replace(key, "[]", "", -1)] = value
			}
		}
	}

	return Parameters{
		Page:      page,
		PageSize:  pageSize,
		URLPath:   u.Path,
		SortField: sortField,
		SortType:  sortType,
		Fields:    fields,
		Columns:   columnsArr,
	}
}

// WhereStatement 處理過濾的where語法
func (param Parameters) WhereStatement(wheres, table string, whereArgs []interface{}, columns, existKeys []string) (string, []interface{}, []string) {
	for key, value := range param.Fields {
		// 運算符號
		var op string
		if strings.Contains(key, "_end") {
			key = strings.Replace(key, "_end", "", -1)
			op = "<="
		} else if strings.Contains(key, "_start") {
			key = strings.Replace(key, "_start", "", -1)
			op = ">="
		} else if len(value) > 1 {
			op = "in"
		} else if !strings.Contains(key, "__operator__") {
			op = "="
		}

		if utils.InArray(columns, key) {
			// 判斷運算符號
			if op == "in" {
				qmark := ""
				for range value {
					qmark += "?,"
				}
				wheres += table + "." + key + " " + op + " (" + qmark[:len(qmark)-1] + ") and "
			} else {
				wheres += table + "." + key + " " + op + " ? and "
			}

			if op == "like" && !strings.Contains(value[0], "%") {
				whereArgs = append(whereArgs, "%"+value[0]+"%")
			} else {
				for _, v := range value {
					whereArgs = append(whereArgs, v)
				}
			}
		} else {
			keys := strings.Split(key, "_join_") // 刪選角色欄位會用到
			if len(keys) > 1 {
				if op == "in" {
					qmark := ""
					for range value {
						qmark += "?,"
					}
					wheres += keys[0] + "." + keys[1] + " " + op + " (" + qmark[:len(qmark)-1] + ") and "
				} else {
					wheres += keys[0] + "." + keys[1] + " " + op + " ? and "
				}
				if op == "like" && !strings.Contains(value[0], "%") {
					whereArgs = append(whereArgs, "%"+value[0]+"%")
				} else {
					for _, v := range value {
						whereArgs = append(whereArgs, v)
					}
				}
			}
		}
		existKeys = append(existKeys, key)
	}
	if len(wheres) > 3 {
		wheres = wheres[:len(wheres)-4]
	}
	return wheres, whereArgs, existKeys
}

// getDefault 判斷url是否有設置key參數，如果沒有則回傳def(預設值)
func getDefault(values url.Values, key, def string) string {
	value := values.Get(key)
	if value == "" {
		return def
	}
	return value
}

// GetFieldValue 取得Parameters.Fields[field]的值，若沒有則回傳""
func (param Parameters) GetFieldValue(field string) string {
	value, ok := param.Fields[field]
	if ok && len(value) > 0 {
		return value[0]
	}
	return ""
}

// GetRouteParamStr 將url.value{}處理成url後的參數
// ex: ?__page=1&__pageSize=10&__sort=id&__sort_type=desc
func (param Parameters) GetRouteParamStr() string {
	p := param.GetFixedParam()
	p.Add("__page", param.Page)
	return "?" + p.Encode()
}

// GetFixedParam 將sort、page相關資訊設置至url.values{}
func (param Parameters) GetFixedParam() url.Values {
	p := url.Values{}
	p.Add("__sort", param.SortField)
	p.Add("__pageSize", param.PageSize)
	p.Add("__sort_type", param.SortType)
	if len(param.Columns) > 0 {
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		p[key] = value
	}
	return p
}

// GetFixedParamWithoutSort 處理url參數(不包含sort)
func (param Parameters) GetFixedParamWithoutSort() string {
	p := url.Values{}
	p.Add("__pageSize", param.PageSize)
	for key, value := range param.Fields {
		p[key] = value
	}
	if len(param.Columns) > 0 {
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	return "&" + p.Encode()
}

// GetLastPageRouteParam 取得上一頁路徑參數
func (param Parameters) GetLastPageRouteParam() string {
	p := param.GetFixedParam()
	pageInt, _ := strconv.Atoi(param.Page)
	p.Add("__page", strconv.Itoa(pageInt-1))
	return "?" + p.Encode()
}

// GetNextPageRouteParam 取得下一頁路徑參數
func (param Parameters) GetNextPageRouteParam() string {
	p := param.GetFixedParam()
	pageInt, _ := strconv.Atoi(param.Page)
	p.Add("__page", strconv.Itoa(pageInt+1))
	return "?" + p.Encode()
}

// URL 設置Page後回傳url參數
func (param Parameters) URL(page string) string {
	return param.URLPath + param.SetPage(page).GetRouteParamStr()
}

// GetRouteParamWithoutPageSize 取得url參數路徑(沒有pagesize)
func (param Parameters) GetRouteParamWithoutPageSize(page string) string {
	p := url.Values{}
	p.Add("__sort", param.SortField)
	p.Add("__page", page)
	p.Add("__sort_type", param.SortType)
	if len(param.Columns) > 0 {
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		p[key] = value
	}
	return "?" + p.Encode()
}

// DeleteField 刪除Parameters.Fields[參數]
func (param Parameters) DeleteField(field string) Parameters {
	delete(param.Fields, field)
	return param
}

// DeletePK 刪除Parameters.Fields[__pk]
func (param Parameters) DeletePK() Parameters {
	delete(param.Fields, "__pk")
	return param
}
