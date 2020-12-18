package db

import "fmt"

// DatabaseType 資料型態
type DatabaseType string

// Value is a string.
type Value string

var (
	// StringTypeList is a DatabaseType list of string.
	StringTypeList = []DatabaseType{Date, Time, Year, Datetime, Timestamptz, Timestamp, Timetz,
		Varchar, Char, Mediumtext, Longtext, Tinytext,
		Text, JSON, Blob, Tinyblob, Mediumblob, Longblob,
		Interval, Point, Bpchar,
		Line, Lseg, Box, Path, Polygon, Circle, Cidr, Inet, Macaddr, Character, Varyingcharacter,
		Nchar, Nativecharacter, Nvarchar, Clob, Binary, Varbinary, Enum, Set, Geometry, Multilinestring,
		Multipolygon, Linestring, Multipoint, Geometrycollection, Name, UUID, Timestamptz,
		Name, UUID, Inet}

	// BoolTypeList is a DatabaseType list of bool.
	BoolTypeList = []DatabaseType{Bool, Boolean}

	// IntTypeList is a DatabaseType list of integer.
	IntTypeList = []DatabaseType{Int4, Int2, Int8,
		Int,
		Tinyint,
		Mediumint,
		Smallint,
		Smallserial, Serial, Bigserial,
		Integer,
		Bigint}

	// FloatTypeList is a DatabaseType list of float.
	FloatTypeList = []DatabaseType{Float, Float4, Float8, Double, Real, Doubleprecision}

	// UintTypeList is a DatabaseType list of uint.
	UintTypeList = []DatabaseType{Decimal, Bit, Money, Numeric}
)

const (
	// =================================
	// integer
	// =================================

	// Int int
	Int DatabaseType = "INT"
	// Tinyint Tinyint
	Tinyint DatabaseType = "TINYINT"
	// Mediumint Mediumint
	Mediumint DatabaseType = "MEDIUMINT"
	// Smallint Smallint
	Smallint DatabaseType = "SMALLINT"
	// Bigint Bigint
	Bigint DatabaseType = "BIGINT"
	// Bit Bit
	Bit DatabaseType = "BIT"
	// Int8 Int8
	Int8 DatabaseType = "INT8"
	// Int4 Int4
	Int4 DatabaseType = "INT4"
	// Int2 Int2
	Int2 DatabaseType = "INT2"

	// Integer Integer
	Integer DatabaseType = "INTEGER"
	// Numeric Numeric
	Numeric DatabaseType = "NUMERIC"
	// Smallserial Smallserial
	Smallserial DatabaseType = "SMALLSERIAL"
	// Serial Serial
	Serial DatabaseType = "SERIAL"
	// Bigserial Bigserial
	Bigserial DatabaseType = "BIGSERIAL"
	// Money Money
	Money DatabaseType = "MONEY"

	// =================================
	// float
	// =================================

	// Real Real
	Real DatabaseType = "REAL"
	// Float Float
	Float DatabaseType = "FLOAT"
	// Float4 Float4
	Float4 DatabaseType = "FLOAT4"
	// Float8 Float8
	Float8 DatabaseType = "FLOAT8"
	// Double Double
	Double DatabaseType = "DOUBLE"
	// Decimal Decimal
	Decimal DatabaseType = "DECIMAL"
	// Doubleprecision Doubleprecision
	Doubleprecision DatabaseType = "DOUBLEPRECISION"

	// =================================
	// string
	// =================================

	// Date Date
	Date DatabaseType = "DATE"
	// Time Time
	Time DatabaseType = "TIME"
	// Year Year
	Year DatabaseType = "YEAR"
	// Datetime Datetime
	Datetime DatabaseType = "DATETIME"
	// Timestamp Timestamp
	Timestamp DatabaseType = "TIMESTAMP"

	// Text Text
	Text DatabaseType = "TEXT"
	// Longtext Longtext
	Longtext DatabaseType = "LONGTEXT"
	// Mediumtext Mediumtext
	Mediumtext DatabaseType = "MEDIUMTEXT"
	// Tinytext Tinytext
	Tinytext DatabaseType = "TINYTEXT"

	// Varchar Varchar
	Varchar DatabaseType = "VARCHAR"
	// Char Char
	Char DatabaseType = "CHAR"
	// Bpchar Bpchar
	Bpchar DatabaseType = "BPCHAR"
	// JSON JSON
	JSON DatabaseType = "JSON"

	// Blob Blob
	Blob DatabaseType = "BLOB"
	// Tinyblob Tinyblob
	Tinyblob DatabaseType = "TINYBLOB"
	// Mediumblob Mediumblob
	Mediumblob DatabaseType = "MEDIUMBLOB"
	// Longblob Longblob
	Longblob DatabaseType = "LONGBLOB"

	// Interval Interval
	Interval DatabaseType = "INTERVAL"
	// Boolean Boolean
	Boolean DatabaseType = "BOOLEAN"
	// Bool Bool
	Bool DatabaseType = "BOOL"

	// Point Point
	Point DatabaseType = "POINT"
	// Line Line
	Line DatabaseType = "LINE"
	// Lseg Lseg
	Lseg DatabaseType = "LSEG"
	// Box Box
	Box DatabaseType = "BOX"
	// Path Path
	Path DatabaseType = "PATH"
	// Polygon Polygon
	Polygon DatabaseType = "POLYGON"
	// Circle Circle
	Circle DatabaseType = "CIRCLE"

	// Cidr Cidr
	Cidr DatabaseType = "CIDR"
	// Inet Inet
	Inet DatabaseType = "INET"
	// Macaddr Macaddr
	Macaddr DatabaseType = "MACADDR"

	// Character Character
	Character DatabaseType = "CHARACTER"
	// Varyingcharacter Varyingcharacter
	Varyingcharacter DatabaseType = "VARYINGCHARACTER"
	// Nchar Nchar
	Nchar DatabaseType = "NCHAR"
	// Nativecharacter Nativecharacter
	Nativecharacter DatabaseType = "NATIVECHARACTER"
	// Nvarchar Nvarchar
	Nvarchar DatabaseType = "NVARCHAR"
	// Clob Clob
	Clob DatabaseType = "CLOB"

	// Binary Binary
	Binary DatabaseType = "BINARY"
	// Varbinary Varbinary
	Varbinary DatabaseType = "VARBINARY"
	// Enum Enum
	Enum DatabaseType = "ENUM"
	// Set Set
	Set DatabaseType = "SET"

	// Geometry Geometry
	Geometry DatabaseType = "GEOMETRY"

	// Multilinestring Multilinestring
	Multilinestring DatabaseType = "MULTILINESTRING"
	// Multipolygon Multipolygon
	Multipolygon DatabaseType = "MULTIPOLYGON"
	// Linestring Linestring
	Linestring DatabaseType = "LINESTRING"
	// Multipoint Multipoint
	Multipoint DatabaseType = "MULTIPOINT"
	// Geometrycollection Geometrycollection
	Geometrycollection DatabaseType = "GEOMETRYCOLLECTION"

	// Name Name
	Name DatabaseType = "NAME"
	// UUID UUID
	UUID DatabaseType = "UUID"

	// Timestamptz Timestamptz
	Timestamptz DatabaseType = "TIMESTAMPTZ"
	// Timetz Timetz
	Timetz DatabaseType = "TIMETZ"
)

// DT string to DatavaseType
func DT(s string) DatabaseType {
	return DatabaseType(s)
}

// Contains 判斷是否包含
func Contains(v DatabaseType, a []DatabaseType) bool {
	for _, i := range a {
		if v == i {
			return true
		}
	}
	return false
}

// GetValueFromDatabaseType 藉由欄位類型取得data值
func GetValueFromDatabaseType(typ DatabaseType, value interface{}) Value {
	switch {
	case Contains(typ, StringTypeList):
		if v, ok := value.(string); ok {
			return Value(v)
		}
		return ""
	case Contains(typ, BoolTypeList):
		if v, ok := value.(bool); ok {
			if v {
				return "true"
			}
			return "false"
		}
		if v, ok := value.(int64); ok {
			if v == 0 {
				return "false"
			}
			return "true"
		}
		return "false"
	case Contains(typ, IntTypeList):
		if v, ok := value.(int64); ok {
			return Value(fmt.Sprintf("%d", v))
		}
		return "0"
	case Contains(typ, FloatTypeList):
		if v, ok := value.(float64); ok {
			return Value(fmt.Sprintf("%f", v))
		}
		return "0"
	case Contains(typ, UintTypeList):
		if v, ok := value.([]uint8); ok {
			return Value(string(v))
		}
		return "0"
	}
	panic("錯誤databasetype?" + string(typ))
}
