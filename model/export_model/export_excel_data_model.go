package export_model

import "github.com/junffzhou/excel_import_export_tool/utils"

/*
ExportData ----导出数据的结构体
Title: Excel export header.
Data: Excel export data.
*/
type ExportData struct {
	Title []string
	Data  [][]string
}


func NewExportData(titleColumn []string) (exportData *ExportData) {

	exportData = new(ExportData)
	for _, titleField := range titleColumn {
		exportData.Title = append(exportData.Title, titleField)
	}

	return
}

//SetData ----设置excel表导出时的数据
func (exportData *ExportData) SetData(dbData []map[string]string, dbColumn []string) {

	for _, dataItem := range dbData {
		data := make([]string, 0)
		for _, dbField := range dbColumn {
			data = append(data, dataItem[dbField])
		}
		exportData.Data = append(exportData.Data, data)
	}

	return
}

//AddActionField ----添加操作列
func (exportData *ExportData) AddActionField(index int) {

	if index == - 1 {
		return
	}
	exportData.Title = utils.AddElementBeforeIndex(index, exportData.Title, []string{utils.ActionField})
	for i, _  := range exportData.Data {
		exportData.Data[i] = utils.AddElementBeforeIndex(index, exportData.Data[i], []string{utils.ActionFieldDefaultValue})
	}

	return
}