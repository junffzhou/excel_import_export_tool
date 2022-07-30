package import_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"strings"
)

//FieldCombinationObj ----field combination
type FieldCombinationObj struct {
	FieldCombination []*FieldCombination `json:"field_combination"`
}

/*
FieldCombination ----field combination detail

Field: Fields to be combined, supports multiple field combinations.
NewField: The combined new field name is stored in the database as a database field.
FieldType:  the type of the combined new field. "00":string,default; "01":map[string]interface{}; "02":[]interface{}; "03":[]map[string]interface{}.
Separator: Separator between field combinations, when FieldType="00", required parameter.
*/
type FieldCombination struct {
	Field string `json:"field"`
	NewField    string   `json:"new_field"`
	FieldType string `json:"field_type"`
	Separator   string   `json:"separator"`
}

func NewFieldCombinationObj(param []byte) (DataHandler, error) {

	res := new(FieldCombinationObj)
	e := json.Unmarshal(param, &res)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func  (m *FieldCombinationObj) HandlerData(titleColumnToDbColumn map[string]string, dbClient *database.Client,
	baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) (e error) {

	for _, item := range m.FieldCombination {
		switch item.FieldType {
		case utils.FieldTypeIsString:
			e = item.combinationByString(titleColumnToDbColumn, xlsxData)
			if e != nil {
				return e
			}
		case utils.FieldTypeIsSlice:
			e = item.combinationBySlice(titleColumnToDbColumn, xlsxData)
			if e != nil {
				return e
			}
		case utils.FieldTypeIsMap:
			e = item.combinationByMap(titleColumnToDbColumn, xlsxData)
			if e != nil {
				return e
			}
		case utils.FieldTypeIsSliceByEleIsMap:
			e = item.combinationBySliceEleIsMap(titleColumnToDbColumn, xlsxData)
			if e != nil {
				return e
			}
		}
	}

	return
}

//combinationByString ----FieldType = "00"
func (item *FieldCombination) combinationByString(titleColumnToDbColumn map[string]string, xlsxData []*import_model.ImportExcelData) (e error) {

	newField := strings.TrimSpace(item.NewField)
	titleColumnToDbColumn[newField ] = newField

	for _, data := range xlsxData {
		newVal := make([]string, 0)
		for _, field := range strings.Split(item.Field, ",") {
			newVal = append(newVal, data.Data[strings.TrimSpace(field)])
		}
		data.Data[newField ] = strings.Join(newVal, item.Separator)
	}

	return
}

//combinationBySlice ----FieldType = "02"
func (item *FieldCombination) combinationBySlice(titleColumnToDbColumn map[string]string, xlsxData []*import_model.ImportExcelData) (e error) {

	newField := strings.TrimSpace(item.NewField)
	titleColumnToDbColumn[newField ] = newField

	for _, data := range xlsxData {
		newVal := make([]interface{}, 0)
		for _, field := range strings.Split(item.Field, ",") {
			newVal = append(newVal, data.Data[strings.TrimSpace(field)])
		}

		res, err := json.Marshal(newVal)
		if err != nil {
			return errors.New(fmt.Sprintf("combinationBySlice err:%v", err))
		}
		data.Data[newField] = string(res)
	}

	return
}

//combinationByMap----FieldType = "02"
func (item *FieldCombination) combinationByMap(titleColumnToDbColumn map[string]string, xlsxData []*import_model.ImportExcelData) (e error) {

	newField := strings.TrimSpace(item.NewField)
	titleColumnToDbColumn[newField] = newField

	for _, data := range xlsxData {
		newVal := make(map[string]interface{}, 0)
		for _, field := range strings.Split(item.Field, ",") {
			fields := strings.Split(field, ":")
			titleKey := strings.TrimSpace(fields[0])
			dbKey := titleKey
			if len(fields) >= 1 {
				dbKey = strings.TrimSpace(fields[1])
			}
			newVal[dbKey] = data.Data[titleKey]
		}

		res, err := json.Marshal(newVal)
		if err != nil {
			return errors.New(fmt.Sprintf("combinationByMap err:%v", err))
		}
		data.Data[item.NewField] = string(res)
	}

	return
}

//combinationBySliceEleIsMap----FieldType = "03"
func ( item *FieldCombination) combinationBySliceEleIsMap(titleColumnToDbColumn map[string]string, xlsxData []*import_model.ImportExcelData) (e error) {

	newField := strings.TrimSpace(item.NewField)
	titleColumnToDbColumn[newField] = newField

	baseData := make([]string, 0)
	e = json.Unmarshal([]byte(item.Field), &baseData)
	if e != nil {
		return
	}

	for _, data := range xlsxData {
		newVal := make([]map[string]interface{}, 0)
		for _, ele := range baseData {
			valTmp := make(map[string]interface{})
			for _, field := range strings.Split(ele, ",") {
				fields := strings.Split(field, ":")
				titleKey := strings.TrimSpace(fields[0])
				dbKey := titleKey
				if len(fields) >= 1 {
					dbKey = strings.TrimSpace(fields[1])
				}
				valTmp[dbKey] = data.Data[titleKey]
			}
			newVal = append(newVal, valTmp)
		}

		res, err := json.Marshal(newVal)
		if err != nil {
			return errors.New(fmt.Sprintf("combinationBySliceEleIsMap err:%v", err))
		}
		data.Data[newField] = string(res)
	}

	return
}