package export_plugin

import (
	"encoding/json"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/export_model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"strings"
)

//FieldCombinationObj ----field Combination
type FieldCombinationObj struct {
	Field map[string]string `json:"field"`
}



type StanderFieldCombination struct {
	fieldDetail []*FieldDetail
	sep         string
}


type FieldDetail struct {
	field    string
	isDelete string
}

func NewFieldCombinationObj(param []byte) (DataHandler, error) {

	res := new(FieldCombinationObj)
	e := json.Unmarshal(param, &res)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (m *FieldCombinationObj) stander() (standerData map[string]*StanderFieldCombination) {

	standerData = make(map[string]*StanderFieldCombination)
	for nweField, value := range m.Field {
		standerData[nweField] = new(StanderFieldCombination)
		sep := ""

		for _, val := range strings.Split(value, ",") {
			va := strings.Split(val, ":")
			if v := strings.TrimSpace(va[0]); v != "@sep" {
				detail := new(FieldDetail)
				detail.field = v
				switch len(va) {
				case 1:
					detail.isDelete = utils.Delete
				case 2:
					if strings.TrimSpace(va[1]) == "" {
						detail.isDelete = utils.Delete
					} else {
						detail.isDelete = strings.TrimSpace(va[1])
					}
				}
				standerData[nweField].fieldDetail = append(standerData[nweField].fieldDetail, detail)
			} else {
				if len(va) == 2 {
					sep = strings.TrimSpace(va[1])
				}
				if len(va) == 3 {
					sep = ":"
				}
			}

		}

		standerData[nweField].sep = sep
	}

	return
}

func (m *FieldCombinationObj) HandlerData(dbData []map[string]string, dbClient *database.Client) (columnReplaces []*export_model.ColumnReplace, e error) {

	standerData := m.stander()
	m.setDbData(dbData, standerData)
	columnReplaces = m.setColumnReplace(standerData)

	return
}

func (m *FieldCombinationObj) setDbData(dbData []map[string]string, standerData map[string]*StanderFieldCombination)  {

	for newField, item := range standerData {
		for _, dataItem := range dbData {
			newVal := make([]string, 0)
			for _, detail := range item.fieldDetail {
				if val, ok := dataItem[detail.field]; ok {
					newVal = append(newVal, val)
				}
			}
			dataItem[newField] = strings.Join(newVal, item.sep)
		}
	}

	return
}

func (m *FieldCombinationObj) setColumnReplace(standerData map[string]*StanderFieldCombination) (columnReplaces []*export_model.ColumnReplace) {

	columnReplaces = make([]*export_model.ColumnReplace, 0)
	for newField, item := range standerData {
		resItem := export_model.NewColumnReplace()
		resItem.AddNewField(newField)
		for _, detail := range item.fieldDetail {
			resItem.AddOldField(detail.field, detail.isDelete)
		}
		columnReplaces = append(columnReplaces, resItem)
	}


	return columnReplaces
}
