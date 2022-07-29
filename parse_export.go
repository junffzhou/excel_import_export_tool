package excel_import_export_tool

import (
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/export_model"
	"github.com/junffzhou/excel_import_export_tool/plugin/export_plugin"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"github.com/xuri/excelize/v2"
	"sort"
)

type ExportClient struct {
	ExportReq       *export_model.ExportRequestParam
	handlerDataExts []export_plugin.DataHandler
	writeExcelExt   export_plugin.WriteExcel
}

func NewExportClient() *ExportClient{
	return new(ExportClient)
}

//RegisterHandlerDataExt ----Register DataHandler Plugin
func (client *ExportClient) RegisterHandlerDataExt(ext export_plugin.DataHandler) {
	client.handlerDataExts = append(client.handlerDataExts, ext)
}

//RegisterWriteExcelExt ----Register WriteExcel Plugin
func (client *ExportClient) RegisterWriteExcelExt(ext export_plugin.WriteExcel) {
	client.writeExcelExt = ext
}

func (client *ExportClient) GetHandlerDataExts() []export_plugin.DataHandler {
	return client.handlerDataExts
}

func (client *ExportClient) GetWriteExcelExt() export_plugin.WriteExcel {
	return client.writeExcelExt
}

//HandlerDataActuator ----DataHandler Plugin Actuator
func (client *ExportClient) HandlerDataActuator(dbData []map[string]string, dbClient *database.Client) (columnReplaces []*export_model.ColumnReplace, e error) {

	for _, ext := range client.GetHandlerDataExts() {
		columnReplace, e := ext.HandlerData(dbData, dbClient)
		if e != nil {
			return columnReplaces, e
		}
		columnReplaces = append(columnReplaces, columnReplace...)
	}

	return columnReplaces, e
}

//WriteExcelExtActuator ----WriteExcel Plugin Actuator
func (client *ExportClient) WriteExcelExtActuator(sheetName, showType string, exportData *export_model.ExportData) (*excelize.File, error) {

	ext := client.GetWriteExcelExt()
	if ext != nil {
		return ext.WriteExcel(sheetName, showType, exportData)
	}

	defaultExt, e := export_plugin.NewWriteExcelObj([]byte(`{}`))
	if e != nil {
		return nil, e
	}

	return defaultExt.WriteExcel(sheetName, showType, exportData)
}

//ParseExportReq ----Parse excel export request
func (client *ExportClient) ParseExportReq(dataSourceName string, reqbody []byte) (string, *excelize.File, error) {


	dbClient, e := database.GetDb(dataSourceName)
	if e != nil {
		return "", nil, errors.New(fmt.Sprintf("ParseExportReq call GetDb err:%v", e))
	}

	client.ExportReq, e = export_model.GetExportReq(reqbody)
	if e != nil {
		return "", nil, errors.New(fmt.Sprintf("ParseExportReq call GetExportReq err:%v", e))
	}

	dbColumn, titleColumn, stmt, e := client.ExportReq.GetStmtAndColumn()
	if e != nil {
		return "", nil, errors.New(fmt.Sprintf("ParseExportReq call GetStmtAndColumn err:%v", e))
	}

	query, args := stmt.GetQuerySQL()
	dbData, e := dbClient.QueryMaps(query, args...)
	if e != nil {
		return "", nil, errors.New(fmt.Sprintf("GetExportXlsxData call QueryMaps err:%v", e))
	}

	columnReplaces, e := client.HandlerDataActuator(dbData, dbClient)
	if e != nil {
		return "", nil, e
	}


	dbColumn, titleColumn = client.ResetDbColumnAndTitleColumn(dbColumn, titleColumn, columnReplaces)
	exportData := export_model.NewExportData(titleColumn)
	exportData.SetData(dbData, dbColumn)
	exportData.AddActionField(client.ExportReq.GetActionFieldIndex(dbColumn))

	file, e := client.WriteExcelExtActuator(client.ExportReq.SheetName, client.ExportReq.ShowType, exportData)
	if e != nil {
		return "", nil, e
	}


	return client.ExportReq.BookName, file, nil
}

/*
ResetDbColumnAndTitleColumn ----reset dbColumn, titleColumn

dbColumnToTitleColumn: map[string]string, key:dbColumn的元素,value:titleColumn相同索引位置的元素
indexArr: oldField在 dbColumn中的索引,去重且取最小索引
titleColumnMap: make(map[int][]string), key: oldField在 dbColumn中的索引,value:新的表头字段
deleteIndex: dbColumn中不在dbColumnToTitleColumn中的元素的索引,即dbColumn中被删除的元素
 */
func (client *ExportClient) ResetDbColumnAndTitleColumn(dbColumn, titleColumn []string, columnReplaces []*export_model.ColumnReplace) ([]string, []string) {

	dbColumnToTitleColumn := make(map[string]string)
	for index, column := range dbColumn {
		if index > len(titleColumn) {
			continue
		}
		dbColumnToTitleColumn[column] = titleColumn[index]
	}

	indexArr := make([]int, 0)
	uniqIndexMap := make(map[int]int)
	titleColumnMap := make(map[int][]string)
	for _, m := range columnReplaces {
		indexItem := make([]int, 0)
		for _, item := range m.GetOldField() {
			index := utils.GetIndex(dbColumn, item.GetOldField())
			if index == -1 {
				continue
			}
			indexItem = append(indexItem, index)
			//删除 isDelete = utils.Delete 的字段
			if item.GetIsDelete() == utils.Delete  {
				delete(dbColumnToTitleColumn, item.GetOldField())
			}
		}
		if len(indexItem) == 0 {
			continue
		}
		sort.Ints(indexItem)
		index := indexItem[0]
		if _, ok := uniqIndexMap[index]; !ok {
			indexArr = append(indexArr, index)
			uniqIndexMap[index] = 1
		}

		titleColumnMap[index] = append(titleColumnMap[index], m.GetNewField()...)
		//添加新的表头字段
		for _, field := range m.GetNewField() {
			dbColumnToTitleColumn[field] = field
		}
	}


	//索引=index 后插入 titleColumnMap[index] 元素
	sort.Sort(sort.Reverse(sort.IntSlice(indexArr)))
	for _, index := range indexArr {
		dbColumn = utils.AddElementAfterIndex(index, dbColumn, titleColumnMap[index])
		titleColumn = utils.AddElementAfterIndex(index, titleColumn, titleColumnMap[index])
	}

	//获取已删除字段的索引,从大到小排序
	deleteIndex := make([]int, 0)
	for index := len(dbColumn) -1; index >= 0; index-- {
		if _, ok := dbColumnToTitleColumn[dbColumn[index]]; !ok {
			deleteIndex = append(deleteIndex, index)
		}
	}

	//删除 索引=index 处的元素
	for _, index := range deleteIndex {
		dbColumn = utils.DeleteElement(index, dbColumn)
		titleColumn = utils.DeleteElement(index, titleColumn)
	}

	return dbColumn, titleColumn
}