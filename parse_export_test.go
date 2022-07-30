package excel_import_export_tool

import (
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/plugin/export_plugin"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestParseDownReq(t *testing.T) {

	before := time.Now()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s&writeTimeout=60s",
		"root",
		"mysql_pass",
		"localhost",
		"3306",
		"test_mysql_database")

	dir, e := os.Getwd()
	assert.NoError(t, e)
	sysType := runtime.GOOS
	if sysType == "windows" {
		dir = dir + "\\test\\"
	} else if sysType == "linux" {
		dir = dir + "/test/"
	}

	bookName, f, e := singleTable(dsn)
	assert.NoError(t, e)
	e = f.SaveAs(dir + bookName)
	assert.NoError(t, e)

	bookName, f, e = dictChange(dsn)
	assert.NoError(t, e)
	e = f.SaveAs(dir + bookName)
	assert.NoError(t, e)

	bookName, f, e = ontToMany(dsn)
	assert.NoError(t, e)
	e = f.SaveAs(dir + bookName)
	assert.NoError(t, e)

	bookName, f, e = manyToMany(dsn)
	assert.NoError(t, e)
	e = f.SaveAs(dir + bookName)
	assert.NoError(t, e)

	bookName, f, e = completeExample(dsn)
	assert.NoError(t, e)
	e = f.SaveAs(dir + bookName)
	assert.NoError(t, e)

	fmt.Println("=====总耗时=====", time.Now().Sub(before).String())
}

func singleTable(dsn string) (bookName string, file *excelize.File, e error) {

	exportColumn := fmt.Sprintf(`[
        {"column": "id::主键ID, dict_code::字典code, dict_item_name::字典项name, dict_item_value::字典项value, comment::描述, creator::创建人"}
    ]`)
	condition := fmt.Sprintf(`[
        {"": ["status=:888", "deleted_at =:null"]}
    ]`)

	bookName = fmt.Sprintf(`字典项_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
	params := fmt.Sprintf(`{
    	"table_name": "dict_info",
    	"book_name": "%s",
    	"sheet_name": "Sheet1",
    	"show_type": "00",
    	"primary_key": "id",
    	"is_include_action_field": "00",
    	"action_field_index": 2,
    	"export_column": %s,
		"condition": %s,
		"order": "dict_code+, dict_item_value-",
		"limit": 100,
    	"offset": 0
    }`, bookName, exportColumn, condition)

	return NewExportClient().ParseExportReq(dsn, []byte(params))
}

func dictChange(dsn string) (bookName string, file *excelize.File, e error) {

	exportColumn := fmt.Sprintf(`[
        {"column": "id:sub_id:主键ID, name:sub_name:姓名, age:sub_age:年龄"},
        {
    		"column": "dict_item_name:sub_sex_name:性别",
    		"join_type": "00",
    		"from_schema": "test_mysql_database",
    		"from_table": "dict_info",
    		"from_table_alias": "dict_info_1",
    		"related_column": "dict_item_value",
    		"rel_table_column": "sub.sex"
		},
		{
    		"column": "dict_item_name:sub_dict_field_name:测试字典",
    		"join_type": "00",
    		"from_schema": "test_mysql_database",
    		"from_table": "dict_info",
    		"from_table_alias": "dict_info_2",
    		"related_column": "dict_item_value",
    		"rel_table_column": "sub.dict_field"
		},
        {"column": "creator:sub_creator:创建人"}
    ]`)

	condition := fmt.Sprintf(`[
        {", dict_info_1, dict_info_2": ["status=: 888"]},
        {"dict_info_1": ["&dict_code: 性别"]},
        {"dict_info_2": ["&dict_code: 字典项测试"]}
    ]`)

	bookName = fmt.Sprintf(`子表信息_字典项转换_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
	params := fmt.Sprintf(`{
    	"table_name": "sub_info",
    	"table_alias": "sub",
    	"book_name": "%s",
    	"sheet_name": "Sheet1",
    	"show_type": "00",
    	"primary_key": "id",
    	"is_include_action_field": "00",
    	"action_field_index": 2,
    	"export_column": %s,
		"condition": %s,
		"order": "sub_id+, sub_sex_name-, sub_dict_field_name+",
		"limit": 100,
    	"offset": 0
    }`, bookName, exportColumn, condition)

	return NewExportClient().ParseExportReq(dsn, []byte(params))
}

func ontToMany(dsn string) (bookName string, file *excelize.File, e error) {

	exportColumn := fmt.Sprintf(`[
        {"column": "id:sub_id:主键ID"},
        {
            "column": "id:main_id:关联的主表主键ID, name:main_name:主表名字, status:main_status:主表状态",
            "join_type": "00",
            "from_schema": "test_mysql_database",
            "from_table": "main_info",
            "from_table_alias": "main",
			"related_column": "id",
            "rel_table_column": "sub.main_id"
        },
        {"column": "name:sub_name:姓名, age:sub_age:年龄, creator:sub_creator:创建人"}
    ]`)
	condition := fmt.Sprintf(`[
        {", main": ["status=:888"]}
    ]`)

	bookName = fmt.Sprintf(`子表信息_一对多关系_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
	params := fmt.Sprintf(`{
    	"table_name": "sub_info",
    	"table_alias": "sub",
    	"book_name": "%s",
    	"sheet_name": "Sheet1",
    	"show_type": "00",
    	"primary_key": "id",
    	"is_include_action_field": "00",
    	"action_field_index": 2,
    	"export_column": %s,
		"condition": %s,
		"order": "main_id+, main_name-",
		"limit": 100,
    	"offset": 0
	}`, bookName, exportColumn, condition)

	return NewExportClient().ParseExportReq(dsn, []byte(params))
}

func manyToMany(dsn string) (bookName string, file *excelize.File, e error) {

	exportColumn := fmt.Sprintf(`[
    	{
			"column": "id:sub_id:主键ID"
    	},
    	{
        	"column": "",
        	"join_type": "00",
        	"from_schema": "test_mysql_database",
        	"from_table": "sub_dept_rel",
        	"from_table_alias": "rel",
        	"related_column": "sub_id",
        	"rel_table_column": "sub.id"
    	},
    	{
        	"column": "id:dpet_id:多对多关联的dept_info表的主键ID, dept_name:dept_name:部门名字, status:dept_status:部门状态",
        	"join_type": "00",
        	"from_schema": "test_mysql_database",
        	"from_table": "dept_info",
        	"from_table_alias": "dept",
        	"related_column": "id",
        	"rel_table_column": "rel.dept_id"
    	},
    	{
        	"column": "name:sub_name:姓名, age:sub_age:年龄, creator:sub_creator:创建人"
    	}
	]`)
	condition := fmt.Sprintf(`[
        {", rel, dept": ["status =:888"]},
        {"rel": ["(sub_id !=:null", "id <>:[1, 100])"]}
    ]`)

	bookName = fmt.Sprintf(`子表信息_多对多关系_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
	params := fmt.Sprintf(`{
    	"table_name": "sub_info",
		"table_alias": "sub",
		"book_name": "%s",
    	"sheet_name": "Sheet1",
    	"show_type": "00",
    	"primary_key": "id",
    	"is_include_action_field": "00",
    	"action_field_index": 2,
		"export_column": %s,
		"condition": %s,
    	"order": "dept_id+, sub_name-",
    	"limit": 100,
    	"offset": 0
	}`, bookName, exportColumn, condition)

	return NewExportClient().ParseExportReq(dsn, []byte(params))
}

func completeExample(dsn string) (bookName string, file *excelize.File, e error) {

	exportColumn := fmt.Sprintf(`[
       	{"column": "id:sub_id:主键ID"},
       	{
           	"column": "",
			"join_type": "00",
           	"from_schema": "test_mysql_database",
           	"from_table": "sub_dept_rel",
           	"from_table_alias": "rel",
			"related_column": "sub_id",
           	"rel_table_column": "sub.id"
       	},
		{
           	"column": "dept_name:dept_name:所属部门",
			"join_type": "00",
           	"from_schema": "test_mysql_database",
           	"from_table": "dept_info",
           	"from_table_alias": "dept",
			"related_column": "id",
           	"rel_table_column": "rel.dept_id"
       	},
       	{
           	"column": "name:main_name:主表name",
			"join_type": "00",
           	"from_schema": "test_mysql_database",
           	"from_table": "main_info",
           	"from_table_alias": "main",
			"related_column": "id",
           	"rel_table_column": "sub.main_id"
       	},
		{"column": "name:sub_name:姓名, age:sub_age:年龄"},
		{
   			"column": "dict_item_name:sex_name:性别",
   			"join_type": "00",
   			"from_schema": "test_mysql_database",
   			"from_table": "dict_info",
   			"from_table_alias": "dict_info_1",
   			"related_column": "dict_item_value",
   			"rel_table_column": "sub.sex"
		},
		{
   			"column": "dict_item_name:dict_field_name:测试字典_name",
   			"join_type": "00",
   			"from_schema": "test_mysql_database",
   			"from_table": "dict_info",
   			"from_table_alias": "dict_info_2",
   			"related_column": "dict_item_value",
   			"rel_table_column": "sub.dict_field"
		},
		{"column": "field_1:sub_field_1:字段1, field_2:sub_field_2:字段2, field_3:sub_field_3:字段3, field_4:sub_field_4:字段4, creator:sub_creator:创建人"}
   ]`)
	condition := fmt.Sprintf(`[
       	{",dept,main,rel": ["&(status=:888)"]},
       	{"rel": ["&(sub_id !=:null", "&id <>:[1,100])"]},
		{"dict_info_1": ["&dict_code =:性别"]},
		{"dict_info_2": ["&dict_code =:字典项测试"]}
   ]`)

	bookName = fmt.Sprintf(`子表信息_completeExample_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
	params := fmt.Sprintf(`{
   		"table_name": "sub_info",
		"table_alias": "sub",
   		"book_name": "%s",
		"sheet_name": "Sheet1",
		"show_type": "00",
   		"primary_key": "",
   		"is_include_action_field": "00",
   		"action_field_index": 2,
		"export_column": %s,
		"condition": %s,
		"order": "sub_id+, sub_name-, dept_name-",
		"limit": 100,
   		"offset": 0
	}`, bookName, exportColumn, condition)

	client := NewExportClient()

	//注册字段组合插件
	fieldCombinationExt, e := export_plugin.NewFieldCombinationObj([]byte(`{
		"field": {
       		"sub_id_name_age": "sub_id:01, sub_name, sub_age:, @sep:-",
       		"sub_name_age": "sub_name, sub_age, @sep::"
   		}
	}`))
	if e != nil {
		return
	}
	client.RegisterHandlerDataExt(fieldCombinationExt)

	//注册字段分离插件
	fieldSeparationExt, e := export_plugin.NewFieldSeparationObj([]byte(`[
    	{
        	"field_separation": "sub_field_1",
        	"field_type": "00",
        	"is_delete": "00",
        	"new_field": "string_field_1, string_field_2, string_field_3",
        	"sep_regexp": "(\\w{2,})_(\\w{2,})_(\\w{2,})_(.*)"
		},
    	{
        	"field_separation": "sub_field_2",
        	"field_type": "01",
        	"is_delete": "00",
        	"is_single": "00",
        	"sep": ",",
        	"new_field": "slice_field_is_single, slice_field_2, slice_field_3, slice_field_4"
    	},
    	{
        	"field_separation": "sub_field_2",
        	"field_type": "01",
        	"is_delete": "00",
        	"is_single": "01",
        	"new_field": "slice_field_1, slice_field_2, slice_field_3, slice_field_4"
    	},
    	{
        	"field_separation": "sub_field_3",
        	"field_type": "02",
        	"is_delete": "00",
        	"new_field": "key_1: map_field_1, key_2: map_filed_2, key_3: map_field_3"
    	},
    	{
        	"field_separation": "sub_field_4",
        	"field_type": "03",
        	"is_delete": "00",
        	"new_field": "[\"Key_1: SliceEleIsMap_field_1, Key_2:SliceEleIsMap_field_2\", \"Key_3: SliceEleIsMap_field_3, Key_4: SliceEleIsMap_field_4\", \"Key_5: SliceEleIsMap_field_5, Key_6: SliceEleIsMap_field_6\"]"
    	}
	]`))
	if e != nil {
		return
	}
	client.RegisterHandlerDataExt(fieldSeparationExt)

	//注册数据写入excel插件
	insertValueExt, e := export_plugin.NewWriteExcelObj([]byte(`{
		"col_width":{
			"A,B,C": 15,
			"D": 20,
			"E": 25,
			"F...": 45
  		},
		"merge": [1, 2, 3, 4]
	}`))
	if e != nil {
		return
	}
	client.RegisterWriteExcelExt(insertValueExt)

	return client.ParseExportReq(dsn, []byte(params))
}