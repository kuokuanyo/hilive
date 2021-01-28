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
	"html/template"
	"strconv"
)

// GetLotteryOtherPrize 取得轉盤遊戲謝謝參與獎頁面、表單資訊
func (s *SystemTable) GetLotteryOtherPrize(ctx *context.Context) (gametable Table) {
	// 建立BaseTable
	gametable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := gametable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("遊戲專屬ID", "game_id", db.Varchar).FieldFilterable()
	info.AddField("謝謝參與獎數量", "amount", db.Int)

	info.SetTable("activity_gamelottery_other_prize").SetTitle("幸運轉盤").SetDescription("謝謝參與獎管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_gamelottery_other_prize").WhereIn("id", ids).Delete()
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
	formList := gametable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("遊戲專屬ID", "game_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("遊戲辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("謝謝參與獎數量", "amount", db.Int, form.Number).
		SetFieldMust().SetFieldDefault("5")

	formList.SetTable("activity_gamelottery_other_prize").SetTitle("幸運轉盤").SetDescription("謝謝參與獎管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "amount") {
			return errors.New("遊戲ID、數量等欄位都不能為空")
		}

		if models.DefaultLotteryOtherModel().SetConn(s.conn).
			IsGameExist(values.Get("game_id"), "") {
			return errors.New("此遊戲已設置謝謝參與獎數量")
		}

		amount, _ := strconv.Atoi(values.Get("amount"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultLotteryOtherModel().SetTx(tx).SetConn(s.conn).AddLotteryOther(
				values.Get("game_id"), amount)
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
		if values.IsEmpty("game_id", "amount") {
			return errors.New("遊戲ID、數量等欄位都不能為空")
		}

		if models.DefaultLotteryOtherModel().SetConn(s.conn).
			IsGameExist(values.Get("game_id"), values.Get("id")) {
			return errors.New("此遊戲已設置謝謝參與獎數量")
		}

		amount, _ := strconv.Atoi(values.Get("amount"))
		model := models.GetLotteryOtherModelAndID("activity_gamelottery_other_prize", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateLotteryOther(values.Get("game_id"), amount)
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
