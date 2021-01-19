package menu

import (
	"hilive/models"
	"hilive/modules/db"
	"regexp"
	"strconv"
)

// Item 為資料表menu的欄位
type Item struct {
	Name         string
	ID           string
	URL          string
	Icon         string
	Header       string
	Active       string
	ChildrenList []Item // 放子選單
}

// Menu 紀錄menu表資訊
type Menu struct {
	List     []Item
	Options  []map[string]string // 紀錄menu的title、ID
	MaxOrder int64
}

// GetMenuInformation 透過user取得menu資料表資訊
func GetMenuInformation(user models.UserModel, conn db.Connection) *Menu {
	var (
		menus      []map[string]interface{}
		menuOption = make([]map[string]string, 0)
	)

	user.GetUserRoles().GetUserMenus()
	// 判斷是否為超級管理員
	if user.IsSuperAdmin() {
		menus, _ = db.TableAndCleanData("menu", conn).
			Where("id", ">", 0).OrderBy("field_order", "asc").All()
	} else {
		var ids []interface{}
		for i := 0; i < len(user.MenuIDs); i++ {
			ids = append(ids, user.MenuIDs[i])
		}
		menus, _ = db.TableAndCleanData("menu", conn).
			WhereIn("id", ids).OrderBy("field_order", "asc").All()
	}

	for i := 0; i < len(menus); i++ {
		title := menus[i]["title"].(string)
		menuOption = append(menuOption, map[string]string{
			"id":    strconv.FormatInt(menus[i]["id"].(int64), 10),
			"title": title,
		})
	}

	// 將map轉換成Item，第二個參數設為0是因為只取沒有子選單的menu，並將子選單放置ChildrenList
	menuList := MapConvertToMenuItem(menus, 0)
	return &Menu{
		List:     menuList,
		Options:  menuOption,
		MaxOrder: menus[len(menus)-1]["parent_id"].(int64),
	}
}

// MapConvertToMenuItem 將map轉換成Item(menu資料表欄位)
func MapConvertToMenuItem(menus []map[string]interface{}, parentID int64) []Item {
	items := make([]Item, 0)

	for j := 0; j < len(menus); j++ {
		if parentID == menus[j]["parent_id"].(int64) {
			title := menus[j]["title"].(string)
			header, _ := menus[j]["header"].(string)

			child := Item{
				Name:   title,
				ID:     strconv.FormatInt(menus[j]["id"].(int64), 10),
				URL:    menus[j]["url"].(string),
				Icon:   menus[j]["icon"].(string),
				Header: header,
				Active: "",
				// 將子選單放置ChildrenList
				ChildrenList: MapConvertToMenuItem(menus, menus[j]["id"].(int64)),
			}
			items = append(items, child)
		}
	}
	return items
}

// SetActiveClass 設定側邊欄menu的展開功能
func (menu *Menu) SetActiveClass(path string) *Menu {
	reg, _ := regexp.Compile(`\?(.*)`)
	path = reg.ReplaceAllString(path, "")
	for i := 0; i < len(menu.List); i++ {
		menu.List[i].Active = "active"
	}

	// for i := 0; i < len(menu.List); i++ {
	// 	if menu.List[i].URL == path && len(menu.List[i].ChildrenList) == 0 {
	// 		menu.List[i].Active = "active"
	// 		return menu
	// 	}
	// 	for j := 0; j < len(menu.List[i].ChildrenList); j++ {
	// 		if menu.List[i].ChildrenList[j].URL == path {
	// 			menu.List[i].Active = "active"
	// 			menu.List[i].ChildrenList[j].Active = "active"
	// 			return menu
	// 		}
	// 		menu.List[i].Active = ""
	// 		menu.List[i].ChildrenList[j].Active = ""
	// 	}
	// 	if strings.Contains(path, menu.List[i].URL) {
	// 		menu.List[i].Active = "active"
	// 	}
	// }
	return menu
}
