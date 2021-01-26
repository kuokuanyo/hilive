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

// GetDanmuPanel 取得彈幕的頁面、表單資訊
func (s *SystemTable) GetDanmuPanel(ctx *context.Context) (danmuTable Table) {
	// 建立BaseTable
	danmuTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := danmuTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("彈幕循環", "danmu_loop", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("顯示位置", "position", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "頂部"
			} else if value.Value == "2" {
				return "中部"
			}
			return "底部"
		})
	info.AddField("顯示暱稱", "display_user", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("彈幕大小", "danmu_size", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "小"
			} else if value.Value == "2" {
				return "中"
			}
			return "大"
		})
	info.AddField("彈幕速度", "danmu_speed", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "慢"
			} else if value.Value == "2" {
				return "中"
			}
			return "快"
		})
	info.AddField("彈幕密度", "danmu_density", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "低"
			} else if value.Value == "2" {
				return "中"
			}
			return "高"
		})
	info.AddField("彈幕不透明度", "danmu_opacity", db.Int).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "低"
			} else if value.Value == "2" {
				return "中"
			}
			return "高"
		})

	info.SetTable("activity_set_danmu").SetTitle("彈幕").SetDescription("彈幕管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_danmu").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 增加表單欄位資訊
	formList := danmuTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("彈幕循環", "danmu_loop", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_danmu").Select("danmu_loop").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["danmu_loop"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("顯示位置", "position", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "頂部"},
			{Value: "2", Text: "中部"},
			{Value: "3", Text: "底部"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_danmu").Select("position").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["position"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("顯示暱稱", "display_user", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_danmu").Select("display_user").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["display_user"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("彈幕大小", "danmu_size", db.Int, form.Radio).
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

			statusModel, _ := s.table("activity_set_danmu").Select("danmu_size").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["danmu_size"].(int64), 10))
			return stauts
		}).SetFieldDefault("2")
	formList.AddField("彈幕速度", "danmu_speed", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "慢"},
			{Value: "2", Text: "中"},
			{Value: "3", Text: "快"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_danmu").Select("danmu_speed").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["danmu_speed"].(int64), 10))
			return stauts
		}).SetFieldDefault("2")
	formList.AddField("彈幕密度", "danmu_density", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "低"},
			{Value: "2", Text: "中"},
			{Value: "3", Text: "高"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_danmu").Select("danmu_density").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["danmu_density"].(int64), 10))
			return stauts
		}).SetFieldDefault("2")
	formList.AddField("彈幕不透明度", "danmu_opacity", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "低"},
			{Value: "2", Text: "中"},
			{Value: "3", Text: "高"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_danmu").Select("danmu_opacity").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["danmu_opacity"].(int64), 10))
			return stauts
		}).SetFieldDefault("2")

	formList.SetTable("activity_set_danmu").SetTitle("彈幕").SetDescription("彈幕管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "danmu_loop", "position", "display_user", "danmu_size",
			"danmu_speed", "danmu_density", "danmu_opacity") {
			return errors.New("活動ID、彈幕資訊等欄位都不能為空")
		}

		if models.DefaultDanmuModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過彈幕的基礎設定")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultDanmuModel().SetTx(tx).SetConn(s.conn).AddDanmu(
				values.Get("activity_id"), values.Get("danmu_loop"),
				values.Get("position"), values.Get("display_user"),
				values.Get("danmu_size"), values.Get("danmu_speed"),
				values.Get("danmu_density"), values.Get("danmu_opacity"))
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
		if values.IsEmpty("activity_id", "danmu_loop", "position", "display_user", "danmu_size",
			"danmu_speed", "danmu_density", "danmu_opacity") {
			return errors.New("活動ID、彈幕資訊等欄位都不能為空")
		}

		danmuModel := models.GetDanmuModelAndID("activity_set_danmu", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := danmuModel.SetTx(tx).UpdateDanmu(
				values.Get("activity_id"), values.Get("danmu_loop"),
				values.Get("position"), values.Get("display_user"),
				values.Get("danmu_size"), values.Get("danmu_speed"),
				values.Get("danmu_density"), values.Get("danmu_opacity"))
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
