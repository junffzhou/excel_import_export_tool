package import_model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/utils"
)

/*
ImportRequestParam ----Import Request parameters

ShowType: excel数据展示方式: "00":按行展示,默认; "01":按列展示,按行或按列读取内容
AutoIncrement: 包含主键字段时,主键是否自增, "00":自增, 默认, "01":非自增
ImportColumn: excel表头和数据库字段的映射, key: excel表的表头, value: 表头对应数据表中的字段。如果表头字段是从其他表获取的, 可以不用在ImportColumn中声明
AdditionalCondition: 对每种导入操作(增/删/改)的附加条件

SchemaName: Database name.
TableName: table name, Required parameter.
SheetName: sheet name, Required parameter.
ShowType: When the excel content is displayed by rows, it is equal to 00; when excel content is displayed by columns, it is equal to 01. The default value is 00.
PrimaryKey: Primary key field name in database
AutoIncrement: When the primary key is incremented, it is equal to "00"; otherwise, it is equal to "01". The default value is 00.
ActionFieldIndex: When the excel header contains an operation column, the position of the operation column in the excel header, starts from 1.
ImportColumn: Required parameter. The mapping between the excel header and the database field,key: the header of the excel table,value: the header corresponds to the field in the data table.
              If the header field is obtained from another table,it is not necessary to declare it in ImportColumn
AdditionalCondition: Additional conditions for each import operation(增/删/改).
*/
type ImportRequestParam struct {
	SchemaName          string               `json:"schema_name"`
	TableName           string               `json:"table_name"`
	SheetName           string               `json:"sheet_name"`
	ShowType            string               `json:"show_type"`
	PrimaryKey          string               `json:"primary_key"`
	AutoIncrement       string               `json:"auto_increment"`
	ActionFieldIndex    float64              `json:"action_field_index"`
	ImportColumn        map[string]string    `json:"import_column"`
	AdditionalCondition *AdditionalCondition `json:"additional_condition"`
}

/*
AdditionalCondition ----Additional conditions
CreateAddCondition: 导入时新增操作的附加条件, Additional conditions for create actions.
UpdateAddCondition: 导入时修改操作的附加条件, Additional Conditions for update actions.
DeleteAddCondition: 导入时删除操作的附加条件, Additional conditions for delete actions.
 */
type AdditionalCondition struct {
	CreateAddCondition *BaseAddCondition `json:"create_add_condition"`
	UpdateAddCondition *BaseAddCondition `json:"update_add_condition"`
	DeleteAddCondition *BaseAddCondition `json:"delete_add_condition"`
}

/*
BaseAddCondition ----Additional conditions detail
ExtraColumn: 额外的字段，insert时作为要插入的字段，update,delete时作为set后面要更新的字段
IsSoftDelete: 删除操作是否是软删除: "00":软删除,默认; "01":非软删除,直接delete数据
Condition: 执行update/delete SQL语句时附加的WHERE条件, 支持多个过滤字段,以英文逗号(",")隔开

ExtraColumn: Additional fields, map[string]interface{}: when create: insert into table(key, ...) values(value, ...), when update or delete: update table set key=value ...
IsSoftDelete: If it is a soft delete, the value is equal to 00, the default; otherwise, the value is equal to 01.
Condition: The where condition attached when executing the update/delete SQL statement supports multiple filter fields separated by English commas (",").
*/
type BaseAddCondition struct {
	ExtraColumn  map[string]interface{} `json:"extra_column"`
	IsSoftDelete string                 `json:"is_soft_delete"`
	Condition    []string               `json:"condition"`
}

//GetImportReq ----Construct request parameter structure, field validation
func GetImportReq(req []byte) (res *ImportRequestParam, e error) {

	res = new(ImportRequestParam)
	e = json.Unmarshal(req, &res)
	if e != nil {
		return res,  errors.New(fmt.Sprintf("GetImportReq error:%v", e))
	}

	if res.TableName == "" {
		return nil, errors.New("GetImportReq err: table_name is a required parameter")
	}

	if res.SheetName == "" {
		return nil, errors.New("GetImportReq err: sheet_name is a required parameter")
	}

	if res.ImportColumn == nil {
		return res, errors.New(fmt.Sprintf("GetImportReq err: import_column is a required parameter"))
	}

	res.setShowType()
	res.setAutoIncrement()
	res.setIsSoftDelete()

	return
}

func (m *ImportRequestParam) setShowType() {

	if  m.ShowType == utils.ShowTypeByCol {
		m.ShowType = utils.ShowTypeByCol
		return
	}

	m.ShowType = utils.ShowTypeByRow
	return
}

func (m *ImportRequestParam) setAutoIncrement() {

	if m.AutoIncrement == utils.NotAutoIncrement {
		m.AutoIncrement = utils.NotAutoIncrement
		return
	}

	m.AutoIncrement = utils.AutoIncrement
	return
}

func (m *ImportRequestParam) setIsSoftDelete() {

	if m.AdditionalCondition != nil {
		if m.AdditionalCondition.DeleteAddCondition != nil {
			if m.AdditionalCondition.DeleteAddCondition.IsSoftDelete == utils.NotSoftDelete {
				m.AdditionalCondition.DeleteAddCondition.IsSoftDelete = utils.NotSoftDelete
				return
			}

			m.AdditionalCondition.DeleteAddCondition.IsSoftDelete = utils.IsSoftDelete
			return
		}
	}

	return
}