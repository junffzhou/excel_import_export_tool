package export_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/export_model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"regexp"
	"strings"
)

//FieldSeparationObj ----field Separation
type FieldSeparationObj struct {
	FieldSeparation []*FieldSeparation
}

/*
FieldSeparation ----field Separation Detail

FieldSeparation: Fields to be separated.
FieldType: The type of the field to be detached, "00": string type; "01":[]interface{} type; "02":map[string]interface{} type; "03":[]map[string]interface{} type.
IsDelete: Whether to delete the field to be separated, "00": delete; "01": not delete, delete by default.
SepRegexp: FieldType="00", separate rules, support regular.
IsSingle: FieldType="01", if IsSingle="00", combine into one field; IsSingle="01", each element corresponds to one field; default "00".
Sep: FieldType="01" and IsSingle="00", separator between elements.
NewField: New field after separation.
*/
type FieldSeparation struct {
	FieldSeparation string `json:"field_separation"`
	FieldType       string `json:"field_type"`
	IsDelete        string `json:"is_delete"`
	SepRegexp       string `json:"sep_regexp"`
	IsSingle        string `json:"is_single"`
	Sep             string `json:"sep"`
	NewField        string `json:"new_field"`
	handlerData     map[string]HandlerData
}

type HandlerData func(dbData []map[string]string) (columnReplace *export_model.ColumnReplace, e error)

func NewFieldSeparationObj(param []byte) (DataHandler, error) {

	data := make([]*FieldSeparation, 0)
	e := json.Unmarshal(param, &data)
	if e != nil {
		return nil, e
	}

	for _, item := range data {
		if item.FieldType == "" {
			item.FieldType = utils.FieldTypeIsString
		}
		item.handlerData = make(map[string]HandlerData)
		switch item.FieldType {
		case utils.FieldTypeIsString:
			item.handlerData[item.FieldType] = item.handlerByString
		case utils.FieldTypeIsSlice:
			item.handlerData[item.FieldType] = item.handlerBySlice
		case utils.FieldTypeIsMap:
			item.handlerData[item.FieldType] = item.handlerByMap
		case utils.FieldTypeIsSliceByEleIsMap:
			item.handlerData[item.FieldType] = item.handlerBySliceEleIsMap
		default:
			return nil, errors.New(fmt.Sprintf("NewFieldSeparationObj error, The value of field_type is optional:%v, %v, %v, %v, %v",
				"null", utils.FieldTypeIsString, utils.FieldTypeIsMap, utils.FieldTypeIsSlice, utils.FieldTypeIsSliceByEleIsMap))
		}
	}

	res := new(FieldSeparationObj)
	res.FieldSeparation= data

	return res, nil
}

func (m *FieldSeparationObj) HandlerData(dbData []map[string]string, dbClient *database.Client) (columnReplaces []*export_model.ColumnReplace, e error) {

	columnReplaces = make([]*export_model.ColumnReplace, 0)
	for _, item := range m.FieldSeparation {
		var columnReplace *export_model.ColumnReplace
		columnReplace, e = item.handlerData[item.FieldType](dbData)
		if e != nil {
			return
		}

		columnReplaces = append(columnReplaces, columnReplace)
	}

	return
}

func (item *FieldSeparation) getNewField() (newField []string) {

	newField = make([]string, 0)
	for _, field := range strings.Split(item.NewField, ",") {
		if strings.TrimSpace(field) == "" {
			continue
		}
		newField = append(newField, strings.TrimSpace(field))
	}

	return
}

func (item *FieldSeparation) getNewFieldByMap() (newField []map[string]string) {

	newField = make([]map[string]string, 0)
	for _, field := range strings.Split(item.NewField, ",") {
		val := strings.Split(field, ":")
		if len(val) == 1 {
			continue
		}
		oldKey := strings.TrimSpace(val[0])
		newKey := strings.TrimSpace(val[1])
		if  oldKey == "" || newKey == "" {
			continue
		}
		newField = append(newField, map[string]string{oldKey: newKey})
	}

	return
}

func (item *FieldSeparation) getNewFieldBySliceEleIsMap() (newField [][]map[string]string, e error) {

	newField = make([][]map[string]string, 0)

	baseData := make([]string, 0)
	e = json.Unmarshal([]byte(item.NewField), &baseData)
	if e != nil {
		return
	}
	for _, dataItem := range baseData {
		newFieldItem := make([]map[string]string, 0)
		for _, field := range strings.Split(dataItem, ",") {
			val := strings.Split(field, ":")
			if len(val) == 1 {
				continue
			}
			oldKey := strings.TrimSpace(val[0])
			newKey := strings.TrimSpace(val[1])
			if  oldKey == "" || newKey == "" {
				continue
			}
			newFieldItem = append(newFieldItem, map[string]string{oldKey: newKey})
		}

		newField = append(newField, newFieldItem)
	}

	return
}

func (item *FieldSeparation) handlerByString(dbData []map[string]string) (columnReplace *export_model.ColumnReplace, e error) {

	newField := item.getNewField()
	e = item.setDbDataByString(newField, dbData)
	if e != nil {
		return
	}

	columnReplace = item.setColumnReplaceByString(newField)

	return
}
func (item *FieldSeparation) handlerBySlice(dbData []map[string]string) (columnReplace *export_model.ColumnReplace, e error) {

	newField := item.getNewField()

	e = item.setDbDataBySlice(newField, dbData)
	if e != nil {
		return
	}
	columnReplace = item.setColumnReplaceBySlice(newField)

	return
}
func (item *FieldSeparation) handlerByMap(dbData []map[string]string) (columnReplace *export_model.ColumnReplace, e error) {


	newFeild := item.getNewFieldByMap()
	e = item.setDbDataByMap(newFeild, dbData)
	if e != nil {
		return
	}
	columnReplace = item.setColumnReplaceByMap(newFeild)

	return
}
func (item *FieldSeparation) handlerBySliceEleIsMap(dbData []map[string]string) (columnReplace *export_model.ColumnReplace, e error) {

	newField, e := item.getNewFieldBySliceEleIsMap()
	if e != nil {
		return columnReplace, e
	}

	e = item.setDbDataBySliceEleIsMap(newField, dbData)
	if e != nil {
		return
	}
	columnReplace = item.setColumnReplaceBySliceEleIsMap(newField)

	return
}

/*
setDbDataByString ----field_type="00", set ExportData

If len(res) < len(newField), The values of the remaining elements of newField are ignored in the exportData,
but the header will exist, and the final result of the remaining elements of newField is the empty string;
otherwise,extra elements in newData are ignored
*/
func (item *FieldSeparation) setDbDataByString(newField []string, dbData []map[string]string) error {

	re, e := regexp.Compile(item.SepRegexp)
	if e != nil {
		return errors.New(fmt.Sprintf("setDbDataByString error:%v", e))
	}

	for _, dataItem := range dbData {
		res := re.FindStringSubmatch(dataItem[item.FieldSeparation])
		if len(res) == 0 {
			continue
		}
		res = res[1:]
		for i, field := range newField {
			if i > len(res) - 1 {
				break
			}
			dataItem[field] = res[i]
		}

	}

	return nil
}
func (item *FieldSeparation) setColumnReplaceByString(newField []string) (columnReplace *export_model.ColumnReplace) {


	columnReplace = export_model.NewColumnReplace()
	columnReplace.AddOldField(item.FieldSeparation, item.IsDelete)
	columnReplace.AddNewField(newField...)
	return
}

/*
setDbDataBySlice ----field_type="01", set ExportData

If is_single="00", after separated, it is combined into a field.
otherwise,split into multiple fields, if len(res) < len(newField),The values of the remaining elements of newField are ignored in the exportData,
									  but the header will exist, and the final result of the remaining elements of newField is the empty string;
									  if len(res)  > len(newField) , extra elements in res are ignored.
*/
func (item *FieldSeparation) setDbDataBySlice(newField []string, dbData []map[string]string) (e error) {

	if len(newField) == 0 {
		return
	}
	for _, dataItem := range dbData {
		res := make([]interface{}, 0)
		e = json.Unmarshal([]byte(dataItem[item.FieldSeparation]), &res)
		if e != nil {
			return errors.New(fmt.Sprintf("setDbDataBySlice error:%v", e))
		}
		if len(res) == 0 {
			continue
		}

		switch item.IsSingle {
		case utils.IsSingle:
			val := make([]string, 0)
			for _, v := range res {
				val = append(val, fmt.Sprintf("%v", v))
			}
			dataItem[newField[0]] = strings.Join(val, item.Sep)
		case utils.NotIsSingle:
			for i, field := range newField {
				if i > len(res) - 1 {
					break
				}
				dataItem[field] = fmt.Sprintf("%v", res[i])
			}
		}
	}

	return
}
func (item *FieldSeparation) setColumnReplaceBySlice(newField []string) (columnReplace *export_model.ColumnReplace) {

	columnReplace = export_model.NewColumnReplace()
	columnReplace.AddOldField(item.FieldSeparation, item.IsDelete)
	if item.IsSingle == utils.IsSingle {
		columnReplace.AddNewField(newField[0])
	} else {
		columnReplace.AddNewField(newField...)
	}

	return
}

/*
setDbDataByMap ----field_type="02", set ExportData

If the key in newField is not in res,it will be ignored; but the header will exist,
and the final result of the key that does not exist in newField will be an empty string
*/
func (item *FieldSeparation) setDbDataByMap(newField []map[string]string, dbData []map[string]string) (e error) {

	for _, dataItem := range dbData {
		res := make(map[string]interface{}, 0)
		e = json.Unmarshal([]byte(dataItem[item.FieldSeparation]), &res)
		if e != nil {
			return errors.New(fmt.Sprintf("setDbDataByMap error:%v", e))
		}

		for _, mapItem := range newField {
			for oldKey, newKey := range mapItem {
				if val, ok := res[oldKey]; ok {
					dataItem[newKey] = fmt.Sprintf("%v", val)
				}
			}
		}
	}

	return
}
func (item *FieldSeparation) setColumnReplaceByMap(newField []map[string]string) (columnReplace *export_model.ColumnReplace) {

	columnReplace = export_model.NewColumnReplace()
	columnReplace.AddOldField(item.FieldSeparation, item.IsDelete)
	for _, mapItem := range newField {
		for _, newKey := range mapItem {
			columnReplace.AddNewField(newKey)
		}
	}

	return
}

/*
setDbDataBySliceEleIsMap ----field_type="03", set ExportData

If len(res) < len(newField), exportData will ignore the remaining elements of newField, but the header will exist,
and the final result of the remaining elements of newField will be an empty string;
otherwise, extra elements in res are ignored.
*/
func (item *FieldSeparation) setDbDataBySliceEleIsMap(newField [][]map[string]string, dbDta []map[string]string) (e error) {

	for _, dataItem := range dbDta {
		res := make([]map[string]interface{}, 0)
		e = json.Unmarshal([]byte(dataItem[item.FieldSeparation]), &res)
		if e != nil {
			return errors.New(fmt.Sprintf("setDbDataBySliceEleIsMap error:%v", e))
		}

		for i, eleItem := range newField {
			if i > len(res) - 1 {
				break
			}
			for _, mapItem := range eleItem {
				for oldKey, newKey := range mapItem {
					if val, ok := res[i][oldKey]; ok {
						dataItem[newKey] = fmt.Sprintf("%v", val)
					}
				}
			}
		}
	}

	return
}
func (item *FieldSeparation) setColumnReplaceBySliceEleIsMap(newField [][]map[string]string) (columnReplace *export_model.ColumnReplace) {


	columnReplace = export_model.NewColumnReplace()
	columnReplace.AddOldField(item.FieldSeparation, item.IsDelete)
	for _, eleItem := range newField {
		for _, mapItem := range eleItem {
			for _, newKey := range mapItem {
				columnReplace.AddNewField(newKey)
			}
		}
	}

	return
}