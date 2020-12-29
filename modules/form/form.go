package form

// Values map[string][]string
type Values map[string][]string

// Get 取得Values[key][0]
func (v Values) Get(key string) string {
	if len(v[key]) > 0 {
		return v[key][0]
	}
	return ""
}

// IsEmpty 判斷是否為空
func (v Values) IsEmpty(key ...string) bool {
	for _, k := range key {
		if v.Get(k) == "" {
			return true
		}
	}
	return false
}