package import_plugin

import (
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
)

type WriteDatabaseObj struct {}


func NewWriteDatabaseObj() DataWriter {
	return new(WriteDatabaseObj)
}

func (m *WriteDatabaseObj) DataWriter(titleColumnToDbColumn map[string]string, dbClient *database.Client, baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) (e error) {

	var importExcelDataArray import_model.ImportExcelDataArray
	importExcelDataArray = append(importExcelDataArray, xlsxData...)

	return importExcelDataArray.ActionFunc(titleColumnToDbColumn, dbClient, baseData)
}