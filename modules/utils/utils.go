package utils

import "encoding/json"

// SetDefault 假設第一個參數 = 第二個參數回傳第三個參數，沒有的話回傳第一個參數
func SetDefault(value, condition, def string) string {
	if value == condition {
		return def
	}
	return value
}

// JSON 執行JSON編碼
func JSON(i interface{}) string{
	if i == nil {
		return ""
	}
	j, _ := json.Marshal(i)
	return string(j)
}

// InArray 是否在陣列中
func InArray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// AorB 判斷條件，true return a，false return b
func AorB(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
}