package utils

import (
	"encoding/json"

	"github.com/NebulousLabs/fastrand"
)

// SetDefault 假設第一個參數 = 第二個參數回傳第三個參數，沒有的話回傳第一個參數
func SetDefault(value, condition, def string) string {
	if value == condition {
		return def
	}
	return value
}

// JSON 執行JSON編碼
func JSON(i interface{}) string {
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

// InArrayWithoutEmpty 是否在陣列中
func InArrayWithoutEmpty(arr []string, str string) bool {
	if len(arr) == 0 {
		return true
	}
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

// Random 隨機數
func Random(strings []string) ([]string, error) {
	for i := len(strings) - 1; i > 0; i-- {
		num := fastrand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	str := make([]string, 0)
	for i := 0; i < len(strings); i++ {
		str = append(str, strings[i])
	}
	return str, nil
}

// UUID 設置uuid
func UUID(length int64) string {
	ele := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "v", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "Driver", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	ele, _ = Random(ele)
	uuid := ""
	var i int64
	for i = 0; i < length; i++ {
		uuid += ele[fastrand.Intn(59)]
	}
	return uuid
}
