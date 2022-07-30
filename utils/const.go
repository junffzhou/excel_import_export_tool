package utils

const (
	Include       = "00" //包含
	NotInclude    = "01" //不包含

	IsSoftDelete = "00" //软删除,默认
	NotSoftDelete = "01" //非软删除

	ShowTypeByRow = "00" //导出时按行展示, 默认
	ShowTypeByCol = "01" //导出时按列展示

	AutoIncrement = "00" //主键字段自增, 默认
	NotAutoIncrement = "01" //主键字段非自增

	Delete    = "00" //删除要替换的字段, 默认
	NotDelete = "01" //不删除要替换的字段

	FieldTypeIsString            = "00" //普通字符串类型, 默认
	FieldTypeIsSlice             = "01" //[]interface{}类型
	FieldTypeIsMap               = "02" //map[string]interface{}类型
	FieldTypeIsSliceByEleIsMap   = "03" //[]map[string]interface{}类型

	IsSingle = "00"
	NotIsSingle = "01"

	DefaultSheetName = "Sheet1" //默认的sheet页name
	ActionField = "操作(增/删/改)" //操作列
	ActionFieldDefaultValue = "无" //操作列默认值

	ActionAdd = "增"
	ActionDelete = "删"
	ActionUpdate = "改"
	ActionNo = "无"


	RepeatabilityCheck  = "Repeatability_Check"//重复性校验
	LegalityCheck       = "Legality_Check" //合法性校验
	NonEmptyCheck       = "NonEmpty_Check" //非空检验
	NonZeroCheck        = "NonZero_Check" //非零校验

	NullValue = "null"
	IsNotNull = `is not null`
	IsNull = `is null`

	//joinType
	LeftJoin = "00"
	RightJOIN = "01"
	InnerJoin = "02"

	Equal   = "="  // 等于
	UnEqual = "!=" // 不等于
)

var NullMap = map[string]string{
	Equal:   IsNull,
	UnEqual: IsNotNull,
}

var JoinTypeMap = map[string]string{
	LeftJoin: " LEFT JOIN ",
	RightJOIN: " RIGHT JOIN ",
	InnerJoin: " INNER JOIN ",
}

//ActionMap 操作类型map
var ActionMap = map[string]struct{}{
	ActionAdd:    {},
	ActionUpdate: {},
	ActionDelete: {},
	ActionNo:     {},
}
