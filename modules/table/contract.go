package table

import (
	"database/sql"
	"errors"
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
)

// GetContractPanel 取得簽約牆頁面、表單資訊
func (s *SystemTable) GetContractPanel(ctx *context.Context) (contractTable Table) {
	// 建立BaseTable
	contractTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := contractTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("簽約牆標題", "contract_title", db.Varchar)
	info.AddField("大屏幕背景圖片", "contract_background", db.Varchar)
	info.AddField("簽名動畫區域大小", "signature_animation_size", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "小"
			} else if value.Value == "2" {
				return "中"
			}
			return "大"
		})
	info.AddField("簽名最終停留區域大小", "signature_area_size", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "小"
			} else if value.Value == "2" {
				return "中"
			}
			return "大"
		})
	info.AddField("使用設備", "mobile_device", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "平板"
			}
			return "手機"
		})
	info.AddField("設備方向", "device_direction", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "橫屏"
			}
			return "豎屏"
		})
	info.AddField("移動端背景圖片", "mobile_background", db.Varchar)
	info.AddField("建立時間", "create_time", db.Varchar)

	info.SetTable("activity_set_contract").SetTitle("簽約牆").SetDescription("簽約牆管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_contract").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 增加表單資訊欄位
	formList := contractTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust()
	formList.AddField("簽約牆標題", "contract_title", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("大屏幕背景圖片", "contract_background", db.Varchar, form.Text)
	formList.AddField("簽名動畫區域大小", "signature_animation_size", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "小"},
			{Value: "2", Text: "中"},
			{Value: "3", Text: "大"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_contract").Select("signature_animation_size").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["signature_animation_size"].(int64), 10))
			return stauts
		}).SetFieldDefault("2")
	formList.AddField("簽名最終停留區域大小", "signature_area_size", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "小"},
			{Value: "2", Text: "中"},
			{Value: "3", Text: "大"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_contract").Select("signature_area_size").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["signature_area_size"].(int64), 10))
			return stauts
		}).SetFieldDefault("2")
	formList.AddField("使用設備", "mobile_device", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "平板", Value: "1"},
			{Text: "手機", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_contract").Select("mobile_device").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["mobile_device"].(int64), 10))
			return status
		}).SetFieldDefault("0")
	formList.AddField("設備方向", "device_direction", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "橫屏", Value: "1"},
			{Text: "豎屏", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_contract").Select("device_direction").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["device_direction"].(int64), 10))
			return status
		}).SetFieldDefault("0")
	formList.AddField("移動端背景圖片", "mobile_background", db.Varchar, form.Text)

	formList.SetTable("activity_set_contract").SetTitle("簽約牆").SetDescription("簽約牆管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "contract_title", "signature_animation_size", "signature_area_size",
			"mobile_device", "device_direction") {
			return errors.New("活動ID、標題、設置等欄位都不能為空")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultContractModel().SetTx(tx).SetConn(s.conn).AddContract(
				values.Get("activity_id"), values.Get("contract_title"), values.Get("contract_background"),
				values.Get("signature_animation_size"), values.Get("signature_area_size"),
				values.Get("mobile_device"), values.Get("device_direction"), values.Get("mobile_background"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "contract_title", "signature_animation_size", "signature_area_size",
			"mobile_device", "device_direction") {
			return errors.New("活動ID、標題、設置等欄位都不能為空")
		}

		model := models.GetContractModelAndID("activity_set_contract", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateContract(values.Get("activity_id"), values.Get("contract_title"), values.Get("contract_background"),
				values.Get("signature_animation_size"), values.Get("signature_area_size"),
				values.Get("mobile_device"), values.Get("device_direction"), values.Get("mobile_background"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	return
}
