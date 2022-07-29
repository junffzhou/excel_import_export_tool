package database

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

//TimeLayout time.Parse 解析数据库时间字段到 GoTime 时的 layout，默认为 "2006-01-02 15:04:05 +0800 CST"
var TimeLayout string

//MysqlTypeToGoType mysql数据库字段类型对应的go数据类型
var MysqlTypeToGoType = map[string]string{
	"enum":               "GoString",
	"set":                "GoString",
	"varchar":            "GoString",
	"char":               "GoString",
	"tinytext":           "GoString",
	"mediumtext":         "GoString",
	"text":               "GoString",
	"longtext":           "GoString",
	"blob":               "GoString",
	"tinyblob":           "GoString",
	"mediumblob":         "GoString",
	"longblob":           "GoString",
	"date":               "GoString",
	"time":               "GoString",
	"binary":             "GoString",
	"varbinary":          "GoString",
	"json":               "GoString",

	"int":                "GoInt",
	"integer":            "GoInt",
	"tinyint":            "GoInt",
	"smallint":           "GoInt",
	"mediumint":          "GoInt",
	"bigint":             "GoInt",
	"int unsigned":       "GoInt",
	"integer unsigned":   "GoInt",
	"tinyint unsigned":   "GoInt",
	"smallint unsigned":  "GoInt",
	"mediumint unsigned": "GoInt",
	"bigint unsigned":    "GoInt",
	"bit":                "GoInt",

	"float":              "GoFloat",
	"double":             "GoFloat",
	"decimal":            "GoFloat",

	"bool":               "GoBool",

	"datetime":           "GoTime",
	"timestamp":          "GoTime",
}

type (
	GoString string //MysqlTypeToGoType中 value=GoString的所有数据库类型
	GoInt int //MysqlTypeToGoType中 value=GoInt的所有数据库类型
	GoFloat float64 //MysqlTypeToGoType中 value=GoFloat的所有数据库类型
	GoBool bool //MysqlTypeToGoType中 value=GoBool的所有数据库类型
	GoTime string //MysqlTypeToGoType中 value=GoTime的所有数据库类型
)


func (s *GoString) Scan(value interface{}) error {

	if value == nil {
		*s = ""
		return nil
	}

	tmp, err := toString(value)
	if err != nil {
		return err
	}

	*s = GoString(tmp)
	return nil
}

func (s *GoInt) Scan(value interface{}) error {

	if value == nil {
		*s = 0
		return nil
	}

	tmp, err := toString(value)
	if err != nil {
		return err
	}

	if tmp == "" {
		*s = 0
		return nil
	}

	i, err := strconv.Atoi(tmp)

	if err == nil {
		*s = GoInt(i)
	}

	return err
}

func (s *GoFloat) Scan(value interface{}) error {

	if value == nil {
		*s = 0
		return nil
	}

	tmp, err := toString(value)
	if err != nil {
		return err
	}

	if tmp == "" {
		*s = 0
		return nil
	}

	i, err := strconv.ParseFloat(tmp, 32)
	if err == nil {
		*s = GoFloat(i)
	}

	return err
}

func (s *GoBool) Scan(value interface{}) error {

	if value == nil {
		*s = false
		return nil
	}

	tmp, err := toString(value)
	if err != nil {
		return err
	}

	if tmp == "true" || tmp == "1" || tmp == "yes" {
		*s = true
		return nil
	}

	*s = false
	return nil
}

func (s *GoTime) Scan(value interface{}) error {

	if value == nil {
		*s = ""
		return nil
	}

	tmp, err := toString(value)
	if err != nil {
		return nil
	}
	if tmp == "" {
		return nil
	}

	if TimeLayout == "" {
		TimeLayout = "2006-01-02 15:04:05 +0800 CST"
	}

	i, err := time.Parse(TimeLayout, tmp)
	if err == nil {
		*s = GoTime(i.Format("2006-01-02 15:04:05"))
	}

	return nil
}

//*GoString 实现 String() 方法, 当格式化输出 *GoString 类型时自动调用该方法
func (s *GoString) String() string{
	res := fmt.Sprintf("%v", *s)
	return res
}

//*GoInt实现 String() 方法
func (s *GoInt) String() string{
	res := fmt.Sprintf("%v", *s)
	return res
}

//*GoFloat 实现 String() 方法
func (s *GoFloat) String() string{
	res := fmt.Sprintf("%v", *s)
	return res
}

//*GoBool 实现 String() 方法
func (s *GoBool) String() string{
	res := fmt.Sprintf("%v", *s)
	return res
}

//*GoTime 实现 String() 方法
func (s *GoTime) String() string{
	res := fmt.Sprintf("%v", *s)
	return res
}

//toString ----数据库查到的原始值转化为对应类型的值
func toString(value interface{}) (res string, err error) {

	defer func() {
		if p := recover(); p != nil {
			err = errors.Errorf("toString error: %v", p)
		}
	}()

	if value == nil {
		return "", nil
	}

	switch val := value.(type) {
	case nil:
		return "", nil
	case *string:
		return fmt.Sprintf("%v", *val), nil
	case *bool:
		return fmt.Sprintf("%v", *val), nil
	case *uint:
		return fmt.Sprintf("%v", *val), nil
	case *uint8:
		return fmt.Sprintf("%v", *val), nil
	case *uint16:
		return fmt.Sprintf("%v", *val), nil
	case *uint32:
		return fmt.Sprintf("%v", *val), nil
	case *uint64:
		return fmt.Sprintf("%v", *val), nil
	case *int:
		return fmt.Sprintf("%v", *val), nil
	case *int8:
		return fmt.Sprintf("%v", *val), nil
	case *int16:
		return fmt.Sprintf("%v", *val), nil
	case *int32:
		return fmt.Sprintf("%v", *val), nil
	case *int64:
		return fmt.Sprintf("%v", *val), nil
	case *float32:
		return fmt.Sprintf("%v", *val), nil
	case *float64:
		return fmt.Sprintf("%v", *val), nil
	case *[]byte:
		return fmt.Sprintf("%v", *val), nil
	case string, bool, uint, uint8, uint16, uint32, uint64,  int, int8, int16, int32, int64, float32, float64:
		return fmt.Sprintf("%v", val), nil
	case []byte:
		return string(val), nil
	case *interface{}:
		return toString(*val)
	case interface{}:
		switch v := val.(type) {
		case string, bool, uint, uint8, uint16, uint32, uint64,  int, int8, int16, int32, int64, float32, float64:
			return fmt.Sprintf("%v", v), nil
		case []byte:
			return string(v), nil
		default:
			return fmt.Sprintf("%v", v), nil
		}
	default:
		return fmt.Sprintf("%v", val), nil
	}
}