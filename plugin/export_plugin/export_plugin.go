package export_plugin

import (
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/export_model"
	"github.com/xuri/excelize/v2"
)

//DataHandler ----It provides methods to process database data.
type DataHandler interface {
	HandlerData(dbData []map[string]string, dbClient *database.Client) (columnReplaces []*export_model.ColumnReplace, e error)
}

//WriteExcel ----It provides methods to write database data to the excel.
type WriteExcel interface {
	WriteExcel(sheetName, showType string, exportData *export_model.ExportData) (*excelize.File, error)
}