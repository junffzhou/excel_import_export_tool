package import_plugin

import (
	"encoding/json"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
	"strings"
)

type RelFieldObj struct {
	RelatedColumn []*RelatedColumn `json:"related_column"`
}

type RelatedColumn struct {
	TitleToColumn string   `json:"title_to_column"`
	FromSchema    string   `json:"from_schema"`
	FromTable     string   `json:"from_table"`
	Condition     []string `json:"condition"`
	Order         string   `json:"order"`
	Limit         float64  `json:"limit"`
	Offset        float64  `json:"offset"`
}


type standerData struct {
	oldColumn map[string]string
	newColumn map[string]string
}


func NewRelatedColumnObj(param []byte) (DataHandler, error) {

	res := new(RelFieldObj)
	e := json.Unmarshal(param, &res)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (obj *RelFieldObj) HandlerData(titleColumnToDbColumn map[string]string, dbClient *database.Client,
	baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) (e error) {


	for _, item := range obj.RelatedColumn {
		stmt, dataRep := item.getStmtAndStanderData(titleColumnToDbColumn)
		sql, params := stmt.GetQuerySQL()
		res, err := dbClient.QueryMaps(sql, params...)
		if err != nil {
			return err
		}

		for _, dataItem := range xlsxData {
			for _, resItem := range res {
				equal := true
				for key, val := range dataRep.oldColumn {
					if dataItem.Data[key] != resItem[val] {
						equal = false
						break
					} else {
						equal = true
					}
				}

				if equal {
					for newKey, key := range dataRep.newColumn {
						dataItem.Data[newKey] = resItem[key]
					}
				}
			}
		}
	}

	return
}


func (item *RelatedColumn) getStmtAndStanderData(titleColumnToDbColumn map[string]string) (stmtItem *model.Statement, stdData *standerData) {

	stmtItem = model.NewStatement(item.FromSchema, item.FromTable, "", int(item.Limit), int(item.Offset))

	var selects []string
	selects, stdData = item.getColumn(titleColumnToDbColumn)
	stmtItem.AddSelects(selects...)
	stmtItem.AddOrders(model.GetOrder(strings.TrimSpace(item.Order))...)

	cond, whereParams, err := model.GetConditionBySlice("", item.Condition)
	if err != nil {
		return stmtItem, stdData
	}

	stmtItem.AddCondition(cond...)
	stmtItem.AddParams(whereParams...)

	return
}



func (item *RelatedColumn) getColumn(titleColumnToDbColumn map[string]string) (selects []string, stdData *standerData) {

	selects = make([]string, 0)
	stdData = new(standerData)
	stdData.oldColumn = make(map[string]string)
	stdData.newColumn = make(map[string]string)

	/*
	column有三段: excel表头:选择字段:关联字段
	1 excel表头不存在, 选择字段, 关联字段必须存在
	2 excel表头存在,选择字段不存在: 选择字段=excel表头
	3 如果想要关联字段, 必须传入第三段
	 */
	for _, column := range strings.Split(item.TitleToColumn, ",") {

		columnArr := strings.Split(strings.TrimSpace(column), ":")
		var selectColumn, val1, val2, val3 string

		val1 = strings.TrimSpace(columnArr[0])
		if len(columnArr) >= 2 {
			val2 = strings.TrimSpace(columnArr[1])
		}
		if len(columnArr) >= 3 {
			val3 = strings.TrimSpace(columnArr[2])
		}


		if val1 == "" {
			if val2 == "" || val3 == "" {
				continue
			} else {
				selectColumn = val2
			}
		} else {
			if val2 == "" {
				selectColumn = val1
			} else {
				selectColumn = val2
			}
		}

		if val3 != "" {
			stdData.newColumn[val3] = selectColumn
			titleColumnToDbColumn[val3] = val3
		}

		if val1 != "" {
			delete(titleColumnToDbColumn, val1)
			stdData.oldColumn[val1] = selectColumn
		}

		selects = append(selects, selectColumn)
	}

	return
}