package paginator

import (
	"hilive/modules/parameter"
	"html/template"
	"math"
	"strconv"
)

// Paginator 分頁器
type Paginator struct {
	Total         string // 總資料數
	URL           string
	PageSizeList  []string
	PreviousClass string              // 如沒上一頁則PreviousClass=disabled，否則為空
	PreviousURL   string              // 前一頁的url參數，例如在第二頁時回傳第一頁的url參數
	Pages         []map[string]string // 每一分頁的資訊，包刮page、active、issplit、url...
	NextClass     string              // 如沒有下一頁=disables，否則為空
	NextURL       string              // 下頁的url參數
	Option        map[string]template.HTML
}

// GetPaginatorInformation 取得頁面的分頁資訊
func GetPaginatorInformation(size int, params parameter.Parameters) Paginator {
	paginator := Paginator{}

	pageInt, _ := strconv.Atoi(params.Page)
	pageSizeInt, _ := strconv.Atoi(params.PageSize)

	// 判斷總共頁數
	totalPage := int(math.Ceil(float64(size) / float64(pageSizeInt)))

	// 取得第一頁的url(不包含pagesize)
	paginator.URL = params.URLPath + params.GetRouteParamWithoutPageSize("1")
	paginator.Total = strconv.Itoa(size)

	// 如果第一頁，則PreviousClass設置disabled
	if pageInt == 1 {
		paginator.PreviousClass = "disabled"
		paginator.PreviousURL = params.URLPath
	} else {
		paginator.PreviousClass = ""
		paginator.PreviousURL = params.URLPath + params.GetLastPageRouteParam()
	}
	// 如果在最後一頁，則NextClass設置disabled
	if pageInt == totalPage {
		paginator.NextClass = "disabled"
		paginator.NextURL = params.URLPath
	} else {
		paginator.NextClass = ""
		paginator.NextURL = params.URLPath + params.GetNextPageRouteParam()
	}

	// 處理頁面的顯示方式
	paginator.Pages = []map[string]string{}
	if totalPage < 10 {
		var pagesArr []map[string]string
		for i := 1; i < totalPage+1; i++ {
			if i == pageInt {
				pagesArr = append(pagesArr, map[string]string{
					"page":    params.Page,
					"active":  "active",
					"isSplit": "0",
					"url":     params.URL(params.Page),
				})
			} else {
				page := strconv.Itoa(i)
				pagesArr = append(pagesArr, map[string]string{
					"page":    page,
					"active":  "",
					"isSplit": "0",
					"url":     params.URL(page),
				})
			}
		}
		paginator.Pages = pagesArr
	} else {
		var pagesArr []map[string]string
		if pageInt < 6 {
			for i := 1; i < totalPage+1; i++ {

				if i == pageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    params.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     params.URL(params.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     params.URL(page),
					})
				}

				if i == 6 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     params.URL("6"),
					})
					i = totalPage - 1
				}
			}
		} else if pageInt < totalPage-4 {
			for i := 1; i < totalPage+1; i++ {

				if i == pageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    params.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     params.URL(params.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     params.URL(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     params.URL("2"),
					})
					if pageInt < 7 {
						i = 5
					} else {
						i = pageInt - 2
					}
				}

				if pageInt < 7 {
					if i == pageInt+5 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     params.URL(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				} else {
					if i == pageInt+3 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     params.URL(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				}
			}
		} else {
			for i := 1; i < totalPage+1; i++ {

				if i == pageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    params.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     params.URL(params.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     params.URL(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     params.URL("2"),
					})
					i = totalPage - 4
				}
			}
		}
		paginator.Pages = pagesArr
	}

	return paginator
}
