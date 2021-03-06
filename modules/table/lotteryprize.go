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

// GetLotteryPrizePanel 取得幸運轉盤獎品頁面、表單資訊
func (s *SystemTable) GetLotteryPrizePanel(ctx *context.Context) (prizeTable Table) {
	// 建立BaseTable
	prizeTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := prizeTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("遊戲專屬ID", "game_id", db.Varchar).FieldFilterable()
	info.AddField("獎品名稱", "prize_name", db.Varchar)
	info.AddField("獎品類型", "prize_type", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "現金紅包"
			}
			return "實物獎品"
		})
	info.AddField("獎品照片", "picture", db.Varchar)
	info.AddField("中獎機率(%)", "percent", db.Int)
	info.AddField("獎品數量", "amount", db.Int)
	info.AddField("剩餘數量", "remain", db.Int)
	info.AddField("獎品價值", "price", db.Int)
	info.AddField("允許抽中", "allow_winning", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "允許"
			}
			return "禁止"
		})
	info.AddField("兌獎方式", "redeem_method", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "現場兌獎"
			}
			return "郵寄兌獎"
		})
	info.AddField("兌獎密碼", "redeem_password", db.Varchar)

	info.SetTable("activity_gamelottery_prize").SetTitle("幸運轉盤").SetDescription("獎品管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_gamelottery_prize").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 取得FormPanel
	formList := prizeTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("遊戲專屬ID", "game_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("遊戲辨別ID")).SetFieldMust().
		FieldNotAllowEdit()
	formList.AddField("獎品名稱", "prize_name", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("獎品類型", "prize_type", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "現金紅包"},
			{Value: "0", Text: "實物獎品"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_gamelottery_prize").
				Select("prize_type").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["prize_type"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("獎品照片", "picture", db.Varchar, form.Text)
	formList.AddField("中獎機率(%)", "percent", db.Int, form.Number).
		SetFieldMust().SetFieldDefault("20").
		SetFieldHelpMsg(template.HTML("數值請設置0~100區間內"))
	formList.AddField("獎品數量", "amount", db.Int, form.Number).
		SetFieldMust().SetFieldDefault("5")
	formList.AddField("剩餘數量", "remain", db.Int, form.Number).
		SetFieldMust().SetFieldDefault("5").FieldNotAllowAdd().
		SetFieldHelpMsg(template.HTML("不可大於獎品數量"))
	formList.AddField("產品價值", "price", db.Int, form.Number).
		SetFieldMust().SetFieldDefault("1000")
	formList.AddField("允許抽中", "allow_winning", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "允許"},
			{Value: "0", Text: "禁止"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_gamelottery_prize").
				Select("allow_winning").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["allow_winning"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("兌獎方式", "redeem_method", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "現場兌獎"},
			{Value: "0", Text: "郵寄兌獎"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_gamelottery_prize").
				Select("redeem_method").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["redeem_method"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("兌獎密碼", "redeem_password", db.Varchar, form.Text)

	formList.SetTable("activity_gamelottery_prize").SetTitle("幸運轉盤").SetDescription("獎品管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_name", "prize_type", "percent",
			"amount", "allow_winning", "redeem_method", "price") {
			return errors.New("遊戲ID、獎品名稱、類型、兌獎等欄位都不能為空")
		}

		percent, _ := strconv.Atoi(values.Get("percent"))
		amount, _ := strconv.Atoi(values.Get("amount"))
		price, _ := strconv.Atoi(values.Get("price"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultLotteryPrizeModel().SetTx(tx).SetConn(s.conn).AddLotteryPrize(
				values.Get("game_id"), values.Get("prize_name"), values.Get("prize_type"),
				values.Get("allow_winning"), values.Get("redeem_method"), values.Get("redeem_password"),
				values.Get("picture"), percent, amount, price)
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
		if values.IsEmpty("game_id", "prize_name", "prize_type", "percent",
			"amount", "allow_winning", "redeem_method", "price", "remain") {
			return errors.New("遊戲ID、獎品名稱、類型、兌獎等欄位都不能為空")
		}

		percent, _ := strconv.Atoi(values.Get("percent"))
		amount, _ := strconv.Atoi(values.Get("amount"))
		remain, _ := strconv.Atoi(values.Get("remain"))
		price, _ := strconv.Atoi(values.Get("price"))
		model := models.GetLotteryPrizeModelAndID("activity_gamelottery_prize", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateLotteryPrize(values.Get("game_id"), values.Get("prize_name"), values.Get("prize_type"),
				values.Get("allow_winning"), values.Get("redeem_method"), values.Get("redeem_password"),
				values.Get("picture"), percent, amount, price, remain)
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
