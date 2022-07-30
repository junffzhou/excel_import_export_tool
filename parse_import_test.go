package excel_import_export_tool

import (
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/plugin/import_plugin"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
	"time"
)


func TestParseUploadReq(t *testing.T) {

	before := time.Now()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s&writeTimeout=60s",
		"root",
		"mysql_pass",
		"localhost",
		"3306",
		"test_mysql_database")


	dir, e := os.Getwd()
	assert.NoError(t, e)
	filePath := dir
	sysType := runtime.GOOS
	if sysType == "windows" {
		filePath = dir + "\\test\\子表信息_import_test.xlsx"
	} else if sysType == "linux" {
		filePath = dir + "/test/子表信息_import_test.xlsx"
	}

	params := fmt.Sprintf(`{
		"table_name": "sub_info",
		"sheet_name":"Sheet1",
		"show_type":"00",
		"primary_key": "id",
		"auto_increment": "01",
		"action_field_index": 2,
		"import_column":{
			"主键ID": "id",
        	"创建人": "creator"
		},
		"additional_condition": {
			"create_add_condition": {
				"extra_column": {
              		"status": "888",
              		"creator": "admin",
					"modifier": "admin"
          		}
      		},
			"update_add_condition": {
          		"extra_column": {
              		"modifier": "admin"
          		},
          		"condition": ["status !=:999"]
			},
			"delete_add_condition": {
          		"extra_column": {
              		"modifier": "admin",
              		"deleted_at": "%s",
              		"status": "999"
          		},
          		"is_soft_delete": "00",
          		"condition": ["status !=:999"]
			}
		}
	}`, time.Now().Format("2006-01-02 15:04:05"))

	client := NewImportClient()
	client.RegisterReadSheetExt(import_plugin.NewReadSheetObj())

	relColumnParam := fmt.Sprintf(`{
    "related_column": [
        {
            "title_to_column": "主表name:name, :id:main_id",
            "from_schema": "test_mysql_database",
            "from_table": "main_info",
            "condition":["status !=: 999"],
            "order": "id",
            "limit": 0,
            "offset": 0
        },
        {
            "title_to_column": "性别:dict_item_name, :dict_item_value:sex",
            "from_schema": "test_mysql_database",
            "from_table": "dict_info",
			"condition":["dict_code:性别", "&status!=:999"],
            "order": "dict_item_value",
            "limit": 0,
            "offset": 0
        },
        {
            "title_to_column": "测试字典_name:dict_item_name, :dict_item_value:dict_field",
            "from_schema": "test_mysql_database",
            "from_table": "dict_info",
            "condition": ["dict_code:字典项测试", "&status!=:999"],
            "order": "dict_item_value",
            "limit": 0,
            "offset": 0
        }
    ]
}`)
	relColumnExt, e := import_plugin.NewRelatedColumnObj([]byte(relColumnParam))
	assert.NoError(t, e)
	client.RegisterHandlerDataExt(relColumnExt)


	fieldCombineParam := fmt.Sprintf(`{
    "field_combination": [
        {
            "field": "string_field_1, string_field_2, string_field_3",
            "new_field": "field_1",
            "field_type": "00",
            "separator": "_"
        },
        {
            "field": "slice_field_1, slice_field_2, slice_field_3, slice_field_4",
            "new_field": "field_2",
            "field_type": "01"
        },
        {
            "field": "map_field_1: key_1, map_filed_2:key_2, map_field_3: key_3",
            "new_field": "field_3",
            "field_type": "02"
        },
        {
            "field": "[\"SliceEleIsMap_field_1: Key_1, SliceEleIsMap_field_2: Key_2\",\"SliceEleIsMap_field_3: Key_3, SliceEleIsMap_field_4: Key_4\",\"SliceEleIsMap_field_5: Key_5, SliceEleIsMap_field_6: Key_6\"]",
            "new_field": "field_4",
            "field_type": "03"
        }
    ]
}`)
	fieldCombineExt, e := import_plugin.NewFieldCombinationObj([]byte(fieldCombineParam))
	assert.NoError(t, e)
	client.RegisterHandlerDataExt(fieldCombineExt)

	fieldSeparatorParam := fmt.Sprintf(`{
    "field_separation":[
        {
            "field":"name_age",
            "new_field":"name, age",
            "separation_regexp":"(.*):(.*)"
        }
	]}`)
	fieldSeparatorExt, e := import_plugin.NewFieldSeparationObj([]byte(fieldSeparatorParam))
	assert.NoError(t, e)
	client.RegisterHandlerDataExt(fieldSeparatorExt)

	dataCheckParam := fmt.Sprintf(`{
    "repeatability_check_field":[
        {
			"field": "sex, name, age",
			"condition":["status !=: 999"]
        }
    ],
	"legality_check_field":[
        {
            "field": "dict_field",
            "reg": "(\\w+)"
        }
    ],
	"non_empty_check_field": "sex",
	"non_zero_check_field": "age"
}`)
	dataCheckExt, e := import_plugin.NewDataCheckObj([]byte(dataCheckParam))
	assert.NoError(t, e)
	client.RegisterHandlerDataExt(dataCheckExt)

	client.RegisterWriteDataBaseExt(import_plugin.NewWriteDatabaseObj())

	e = client.ParseImportReq(dsn, filePath, []byte(params))
	assert.NoError(t, e)

	fmt.Println("=====总耗时=====", time.Now().Sub(before).String())
}
