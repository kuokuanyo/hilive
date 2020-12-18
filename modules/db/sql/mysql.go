package sql

// mysql 分隔符號
type mysql struct {
	delimiter string
}

// -----設置CRUD的所有方法-----start

// GetName return "mysql"
func (mysql) GetName() string {
	return "mysql"
}

// ShowColumns return "show columns in table"
func (mysql) ShowColumns(table string) string {
	return "show columns in " + table
}

// ShowTables return "show tables"
func (mysql) ShowTables() string {
	return "show tables"
}

// Insert 處理insert命令
func (mysql mysql) Insert(condition *FilterCondition) string {
	fields := " ("
	// 添加問號(ex: ?,?,...)
	questionMark := "("

	// 處理sql命令
	for key, value := range condition.Values {
		fields += key + ","
		questionMark += "?,"
		condition.Args = append(condition.Args, value)
	}
	fields = fields[:len(fields)-1] + ")"
	questionMark = questionMark[:len(questionMark)-1] + ")"

	condition.Statement = "insert into " + condition.TableName + fields + " values " + questionMark
	return condition.Statement
}

// Delete 處理delete命令
func (mysql mysql) Delete(condition *FilterCondition) string {
	wheres := " where "
	if len(condition.Wheres) == 0 {
		if condition.WhereRaws != "" {
			wheres += condition.WhereRaws
		} else {
			wheres = ""
		}
	} else {
		for _, where := range condition.Wheres {
			wheres += where.Field + " " +
				where.Operation + " " + where.Value + " and "
		}
		if condition.WhereRaws != "" {
			wheres += condition.WhereRaws + " and "
		}
		wheres = wheres[:len(wheres)-5]
	}

	condition.Statement = "delete from " + condition.TableName + wheres
	return condition.Statement
}

// Update 處理update命令
func (mysql mysql) Update(condition *FilterCondition) string {
	fields := ""
	wheres := " where "
	args := make([]interface{}, 0)

	if len(condition.Values) != 0 {
		for field, value := range condition.Values {
			fields += field + " = ?, "
			args = append(args, value)
		}

		if len(condition.UpdateRaws) == 0 {
			fields = fields[:len(fields)-2]
		} else {
			for i := 0; i < len(condition.UpdateRaws); i++ {
				if i == len(condition.UpdateRaws)-1 {
					fields += condition.UpdateRaws[i].Expression + " "
				} else {
					fields += condition.UpdateRaws[i].Expression + ","
				}
				args = append(args, condition.UpdateRaws[i].Args...)
			}
		}

		condition.Args = append(args, condition.Args...)
	} else {
		if len(condition.UpdateRaws) == 0 {
			panic("資料表更新資料發生錯誤，必須設置參數")
		} else {
			for i := 0; i < len(condition.UpdateRaws); i++ {
				if i == len(condition.UpdateRaws)-1 {
					fields += condition.UpdateRaws[i].Expression + " "
				} else {
					fields += condition.UpdateRaws[i].Expression + ","
				}
				args = append(args, condition.UpdateRaws[i].Args...)
			}
		}
		condition.Args = append(args, condition.Args...)
	}

	if len(condition.Wheres) == 0 {
		if condition.WhereRaws != "" {
			wheres += condition.WhereRaws
		} else {
			wheres = ""
		}
	} else {
		for _, where := range condition.Wheres {
			wheres += where.Field + " " +
				where.Operation + " " + where.Value + " and "
		}

		if condition.WhereRaws != "" {
			wheres += condition.WhereRaws + " and "
		}
		wheres = wheres[:len(wheres)-5]
	}

	condition.Statement = "update " + condition.TableName + " set " +
		fields + wheres
	return condition.Statement
}

// Select 處理查詢命令
func (mysql mysql) Select(condition *FilterCondition) string {
	var fields, joins, group, order, limit, offset string
	wheres := " where "

	if len(condition.Fields) == 0 {
		fields = "* "
	} else {
		if len(condition.Leftjoins) == 0 {
			for i, field := range condition.Fields {
				if condition.Functions[i] != "" {
					fields += condition.Functions[i] + "(" +
						field + "),"
				} else {
					fields += field + ","
				}
			}
		} else {
			for _, field := range condition.Fields {
				fields += field + ","
			}
		}
	}
	fields = fields[:len(fields)-1]

	if len(condition.Leftjoins) != 0 {
		for _, join := range condition.Leftjoins {
			joins += " left join " + join.Table + " on " + join.FieldA +
				" " + join.Operation + " " + join.FieldB + " "
		}
	}

	if len(condition.Wheres) == 0 {
		if condition.WhereRaws != "" {
			wheres += condition.WhereRaws
		} else {
			wheres = ""
		}
	} else {
		for _, where := range condition.Wheres {
			wheres += where.Field + " " +
				where.Operation + " " + where.Value + " and "
		}

		if condition.WhereRaws != "" {
			wheres += condition.WhereRaws + " and "
		}
		wheres = wheres[:len(wheres)-5]
	}

	if condition.Group != "" {
		group = " group by " + condition.Group + " "
	}
	if condition.Order != "" {
		order = " order by " + condition.Order + " "
	}
	if condition.Limit != "" {
		limit = " linit " + condition.Limit + " "
	}
	if condition.Offset != "" {
		offset = " offset " + condition.Offset + " "
	}

	condition.Statement = "select " + fields + " from " + condition.TableName +
		joins + wheres + group + order + limit + offset
	return condition.Statement
}

// -----設置CRUD的所有方法-----end
