package import_plugin

import (
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)


type ReadSheetObj struct {
}

func NewReadSheetObj() SheetReader {
	return new(ReadSheetObj)
}

/*
ReadSheet ----Read excel table content by row or column

Get data from Excel
*/
func (obj *ReadSheetObj) ReadSheet(filePath string, titleColumnToDbColumn map[string]string,
	baseData *import_model.ImportRequestParam) (xlsxData []*import_model.ImportExcelData, e error) {

	f, e := excelize.OpenFile(filePath)
	if e != nil {
		return xlsxData, e
	}
	if baseData.ShowType == utils.ShowTypeByRow {
		xlsxData, e = obj.readSheetByRow(titleColumnToDbColumn, f, baseData)
		if e != nil {
			return xlsxData, e
		}
	} else {
		xlsxData, e = obj.readSheetByCol(titleColumnToDbColumn, f, baseData)
		if e != nil {
			return xlsxData, e
		}
	}
	return xlsxData, nil
}


/*
readSheetByRow ----按行读取数据. excelRead excel table content by row

excelTitleMap: excel表头字段,第一列是表头.The first column is header, get excel header field.
Action: 包含操作字段时,操作字段的值. When contains an action column, gets the value of the action column.

PrimaryKeyValue: 如果包含主键列, 并且当前列等于主键列,则获取主键的值主键字段的值. If the primary key column is included, and the current column is equal to the primary key column, the value of the primary key is obtained
Data: 数据列的值,不包含主键字段.If the primary key column is not included, or the current column is not equal to the primary key column, it is the data column. Remove primary key columns and action columns.
*/
func (obj *ReadSheetObj) readSheetByRow(titleColumnToDbColumn map[string]string, f *excelize.File,
	baseData *import_model.ImportRequestParam) ([]*import_model.ImportExcelData, error) {


	sheetName := baseData.SheetName
	rows, e := f.GetRows(sheetName)
	if e != nil {
		return  nil, e
	}

	actionFieldErrRows := make([]string, 0)
	primaryFieldErrRows := make([]string, 0)
	xlsxData := make([]*import_model.ImportExcelData, 0)
	excelTitleMap := make(map[int]string)
	for i, row := range rows {
		item := new(import_model.ImportExcelData)
		item.Data = make(map[string]string)

		for j, colCell := range row {

			if i == 0 {
				excelTitleMap[j] = colCell
				continue
			}

			if j == (int(baseData.ActionFieldIndex) - 1) {
				if _, ok := utils.ActionMap[colCell]; ok {
					item.Action = colCell
				} else {
					actionFieldErrRows = append(actionFieldErrRows, strconv.Itoa(i+1))
				}
			}


			if i >= 1 && j != (int(baseData.ActionFieldIndex) - 1) {
				excelTitle := excelTitleMap[j]
				if baseData.PrimaryKey == titleColumnToDbColumn[excelTitle] {
					item.PrimaryKeyValue = colCell
				} else {
					item.Data[excelTitle] = colCell
				}
			}
		}

		if i > 0 {
			item.Row = i+1
			if  baseData.PrimaryKey != "" && item.Action == utils.ActionUpdate && item.PrimaryKeyValue== "" {
				primaryFieldErrRows = append(primaryFieldErrRows, strconv.Itoa(i+1))
			}
			xlsxData = append(xlsxData, item)
		}
	}

	if len(actionFieldErrRows) != 0 {
		return nil, errors.New(fmt.Sprintf("readSheetByRow error, sheetName:%v, row:%v, the value of the action column is incorrect, " +
			"the value can only be 增/删/改/无", sheetName, strings.Join(actionFieldErrRows, ",")))
	}


	if len(primaryFieldErrRows) != 0 {
		return nil, errors.New(fmt.Sprintf("readSheetByRow error,sheetName:%v, row:%v, " +
			"when primary_key is not empty and action field is equal to %v, primary_key field value not null",
			sheetName, strings.Join(primaryFieldErrRows, ","),  utils.ActionUpdate))
	}

	return  xlsxData, nil
}


/*
readSheetByCol ----按列读取excel数据. Read excel table content by column

excelTitleMap: excel表头字段,第一行是表头. The first column is header, get excel header field.
Action: 包含操作字段时,操作字段的值.When contains an action column, gets the value of the action column.

PrimaryKeyValue: 如果包含主键列, 并且当前列等于主键列,则获取主键的值主键字段的值. If the primary key column is included, and the current column is equal to the primary key column, the value of the primary key is obtained
Data: 数据列的值,不包含主键字段. If the primary key column is not included, or the current column is not equal to the primary key column, it is the data column. Remove primary key columns and action columns.
 */
func (obj *ReadSheetObj) readSheetByCol(titleColumnToDbColumn map[string]string, f *excelize.File,
	baseData *import_model.ImportRequestParam) ([]*import_model.ImportExcelData, error) {

	sheetName := baseData.SheetName
	cols, e := f.GetCols(sheetName)
	if e != nil {
		return  nil, e
	}

	actionFieldErrColumns := make([]string, 0)
	primaryFieldErrColumns := make([]string, 0)
	xlsxData := make([]*import_model.ImportExcelData, 0)
	excelTitleMap := make(map[int]string)
	for i, col := range cols {
		item := new(import_model.ImportExcelData)
		item.Data = make(map[string]string)

		for j, rowCell := range col {
			if i == 0 {
				excelTitleMap[j] = rowCell
				continue
			}

			if  j == (int(baseData.ActionFieldIndex) - 1) {
				if _, ok := utils.ActionMap[rowCell]; ok {
					item.Action = rowCell
				} else {
					actionFieldErrColumns = append(actionFieldErrColumns, strconv.Itoa(i+1))
				}
			}

			excelTitle := excelTitleMap[j]
			if dbField, ok := titleColumnToDbColumn[excelTitle]; ok {
				if dbField == baseData.PrimaryKey {
					item.PrimaryKeyValue = rowCell
				} else {
					item.Data[dbField] = rowCell
				}
			}
		}

		if i > 0 {
			item.Col = i+1
			if baseData.PrimaryKey != "" && item.Action == utils.ActionUpdate && item.PrimaryKeyValue== "" {
				primaryFieldErrColumns = append(primaryFieldErrColumns, strconv.Itoa(i+1))
			}
			xlsxData = append(xlsxData, item)
		}
	}


	if len(actionFieldErrColumns) != 0 {
		return nil, errors.New(fmt.Sprintf("readSheetByCol error, sheetName:%v, column:%v, the value of the action column is incorrect, " +
			"the value can only be 增/删/改/无", sheetName, strings.Join(actionFieldErrColumns, ",")))
	}

	if len(primaryFieldErrColumns) != 0 {
		return nil, errors.New(fmt.Sprintf("readSheetByCol error, sheetName:%v, column:%v," +
			"when primary_key is not empty and action field is equal to %v, primary_key field value not null",
			sheetName, strings.Join(primaryFieldErrColumns, ","),  utils.ActionUpdate))
	}


	return xlsxData, nil
}