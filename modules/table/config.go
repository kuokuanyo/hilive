package table

// ConfigTable 設置資訊
type ConfigTable struct {
	Driver     string
	CanAdd     bool
	EditAble   bool
	DeleteAble bool
	PrimaryKey PrimaryKey
}

// DefaultConfig 預設Config(struct)
func DefaultConfig() ConfigTable {
	return ConfigTable{
		Driver:    "mysql",
		CanAdd:    true,
		EditAble:  true,
		DeleteAble: true,
		PrimaryKey: PrimaryKey{
			Type: "INT",
			Name: "id",
		},
	}
}

// DefaultConfigTableByDriver 建立預設的ConfigTable(struct)
func DefaultConfigTableByDriver(driver string) ConfigTable {
	return ConfigTable{
		Driver:     driver,
		CanAdd:     true,
		EditAble:   true,
		DeleteAble: true,
		PrimaryKey: PrimaryKey{
			Type: "INT",
			Name: "id",
		},
	}
}
