package excel_import_export_tool

import (
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
	"github.com/junffzhou/excel_import_export_tool/plugin/import_plugin"
)

type ImportClient struct {
	ImportReq        *import_model.ImportRequestParam
	readSheetExt     import_plugin.SheetReader
	handlerDataExts  []import_plugin.DataHandler
	writeDataBaseExt import_plugin.DataWriter
}

func NewImportClient() *ImportClient {
	return new(ImportClient)
}

//RegisterHandlerDataExt ----Register DataHandler Plugin
func (client *ImportClient) RegisterHandlerDataExt(ext import_plugin.DataHandler) {
	client.handlerDataExts = append(client.handlerDataExts, ext)
}

//RegisterReadSheetExt ----Register SheetReader Plugin
func (client *ImportClient) RegisterReadSheetExt(ext import_plugin.SheetReader) {
	client.readSheetExt = ext
}

//RegisterWriteDataBaseExt ----Register DataWriter Plugin
func (client *ImportClient) RegisterWriteDataBaseExt(ext import_plugin.DataWriter) {
	client.writeDataBaseExt = ext
}

func (client *ImportClient) GetHandlerDataExts() []import_plugin.DataHandler {
	return client.handlerDataExts
}

func (client *ImportClient) GetReadSheetExt() import_plugin.SheetReader {
	return client.readSheetExt
}

func (client *ImportClient) GetWriteDataBaseExt() import_plugin.DataWriter {
	return client.writeDataBaseExt
}

//HandlerDataActuator ----DataHandler Plugin Actuator
func (client *ImportClient) HandlerDataActuator(titleColumnToDbColumn map[string]string, dbClient *database.Client,
	baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) (e error) {

	for _, ext := range client.GetHandlerDataExts() {
		e = ext.HandlerData(titleColumnToDbColumn, dbClient, baseData, xlsxData)
		if e != nil {
			return  e
		}
	}

	return nil
}

//ReadSheetActuator ----SheetReader Plugin Actuator
func (client *ImportClient) ReadSheetActuator(filePath string, titleColumnToDbColumn map[string]string,
	baseData *import_model.ImportRequestParam) ([]*import_model.ImportExcelData, error) {

	ext := client.GetReadSheetExt()
	if ext != nil {
		return ext.ReadSheet(filePath, titleColumnToDbColumn, baseData)
	}

	return import_plugin.NewReadSheetObj().ReadSheet(filePath, titleColumnToDbColumn, baseData)
}

//WriteDataBaseActuator ----DataWriter Plugin Actuator
func (client *ImportClient) WriteDataBaseActuator(titleColumnToDbColumn map[string]string,
	dbClient *database.Client, baseData *import_model.ImportRequestParam,
	xlsxData []*import_model.ImportExcelData) (e error) {

	ext := client.GetWriteDataBaseExt()
	if ext != nil {
		return ext.DataWriter(titleColumnToDbColumn, dbClient, baseData, xlsxData)
	}

	return import_plugin.NewWriteDatabaseObj().DataWriter(titleColumnToDbColumn, dbClient, baseData, xlsxData)
}


//ParseImportReq ----Parse excel import request
func (client *ImportClient) ParseImportReq(dataSourceName, filePath string,
	reqbody []byte,) (e error) {

	client.ImportReq, e = import_model.GetImportReq(reqbody)
	if e != nil {
		return errors.New(fmt.Sprintf("ParseExportReq call GetImportReq err:%v", e))
	}

	titleColumnToDbColumn, e := client.GetTitleColumnToDbColumn()
	if e != nil {
		return e
	}

	dbClient, e := database.GetDb(dataSourceName)
	if e != nil {
		return  errors.New(fmt.Sprintf("ParseExportReq call GetDb err:%v", e))
	}

	xlsxData, e := client.ReadSheetActuator(filePath, titleColumnToDbColumn, client.ImportReq)
	if e != nil {
		return errors.New(fmt.Sprintf("ParseExportReq call ReadSheetActuator err:%v", e))
	}


	e = client.HandlerDataActuator(titleColumnToDbColumn, dbClient, client.ImportReq, xlsxData)
	if e != nil {
		return errors.New(fmt.Sprintf("ParseExportReq call HandlerDataActuator err:%v", e))
	}

	return  client.WriteDataBaseActuator(titleColumnToDbColumn, dbClient, client.ImportReq, xlsxData)
}

func (client *ImportClient) GetTitleColumnToDbColumn() (titleColumnToDbColumn map[string]string, e error) {

	if client.ImportReq.ImportColumn == nil {
		return
	}

	titleColumnToDbColumn = make(map[string]string)
	isIncludePrimaryKey := false
	for titleColumn, dbColumn := range client.ImportReq.ImportColumn {
		titleColumnToDbColumn[titleColumn] = dbColumn
		if client.ImportReq.PrimaryKey != "" && dbColumn == client.ImportReq.PrimaryKey {
			isIncludePrimaryKey = true
		}
	}

	if !isIncludePrimaryKey {
		return titleColumnToDbColumn, errors.New(fmt.Sprintf("GetTitleColumnToDbColumn error, " +
			"when primary_key is not empty, import_column must contain %v", client.ImportReq.PrimaryKey))
	}

	return
}