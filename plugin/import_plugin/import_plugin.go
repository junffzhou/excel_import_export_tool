package import_plugin

import (
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
)

//SheetReader ----It provides methods to read excel content
type SheetReader interface {
	ReadSheet(filePath string, titleColumnToDbColumn map[string]string, baseData *import_model.ImportRequestParam) ([]*import_model.ImportExcelData, error)
}

//DataHandler ----It provides methods to process excel data.
type DataHandler interface {
	HandlerData(titleColumnToDbColumn map[string]string, dbClient *database.Client,  baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) error
}

//DataWriter ----It provides methods to write excel data to the database.
type DataWriter interface {
	DataWriter(titleColumnToDbColumn map[string]string, dbClient *database.Client, baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) error
}

