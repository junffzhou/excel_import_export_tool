package import_plugin

import (
	"encoding/json"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
	"regexp"
	"strings"
)

//FieldSeparationObj ----field separation
type FieldSeparationObj struct {
	FieldSeparation []*FieldSeparation `json:"field_separation"`
}

/*
FieldSeparation ----field separation detail

Field: 要分离的字段
NewField: 分离后的新字段, 支持多个,多个以英文逗号隔开(","), 作为数据库字段存入数据库中
		  如果正则表达式结果的长度 < NewField的个数, 则 NewField剩余元素的值为空,否则, 正则表达式的结果剩余元素被忽略
SeparationRegexp: 分离匹配的规则,支持正则表达式

Field: Separate fields.
NewField: The new field after separation, get the value according to the regular expression.
		  If the length of the regular expression result < len(NewField), the value of the remaining elements of NewField is empty,
		  Otherwise, the remaining elements of the result of the regular expression are ignored
SeparationRegexp: Separate rules, support for regular expressions
*/
type FieldSeparation struct {
	Field string `json:"field"`
	NewField    string   `json:"new_field"`
	SeparationRegexp  string   `json:"separation_regexp"`
}

func NewFieldSeparationObj(param []byte) (DataHandler, error) {

	res := new(FieldSeparationObj)
	e := json.Unmarshal(param, &res)
	if e != nil {
		return nil, e
	}

	return res, nil
}


func (m *FieldSeparationObj) HandlerData(titleColumnToDbColumn map[string]string, dbClient *database.Client,
	baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) (e error) {

	for _, item := range m.FieldSeparation {
		re, err := regexp.Compile(item.SeparationRegexp)
		if err != nil {
			return err
		}

		for j, data := range xlsxData {
			newData := re.FindStringSubmatch(data.Data[item.Field])
			if len(newData) == 0 {
				continue
			}
			newData = newData[1:]
			for i, field := range strings.Split(item.NewField, ",") {
				reField := strings.TrimSpace(field)
				if i > len(newData) - 1 {
					data.Data[reField] = ""
					continue
				}
				data.Data[reField] = newData[i]

				if j == 0 {
					titleColumnToDbColumn[reField] = reField
				}
			}
		}
	}

	return
}
