package import_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model"
	"github.com/junffzhou/excel_import_export_tool/model/import_model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"regexp"
	"strconv"
	"strings"
)

/*
DataCheckObj ---Data check

RepeatabilityCheckField: 重复性校验字段, Repeatability check, Unique check
LegalityCheckField: 合法性校验字段,Legality check: whether the check rules are met
NonEmptyCheckField: 非空字段,non-null check
NonZeroCheckField: 非零字段, non-zero check
*/
type DataCheckObj struct {
	RepeatabilityCheckField []*RepeatabilityCheck `json:"repeatability_check_field"`
	LegalityCheckField      []*LegalityCheck      `json:"legality_check_field"`
	NonEmptyCheckField      string                `json:"non_empty_check_field"`
	NonZeroCheckField       string                `json:"non_zero_check_field"`
}

/*
RepeatabilityCheck ----Repeatability check detail
Field: 支持多个联合唯一, check fields, support multiple fields jointly unique
Condition: Query conditions during check
*/
type RepeatabilityCheck struct {
	Field     string   `json:"field"`
	Condition []string `json:"condition"`
	Order     string   `json:"order"`
	Limit     float64  `json:"limit"`
	Offset    float64  `json:"offset"`
}

/*
LegalityCheck ----Legality check detail

Field: check fields, support multiple fields with the same check rule
Reg: 合法性校验的规则,支持正则. Check rules, support regular expressions
 */
type LegalityCheck struct {
	Field      string   `json:"field"`
	Reg        string   `json:"reg"`
}


func NewDataCheckObj(param []byte) (DataHandler, error) {

	res := new(DataCheckObj)
	e := json.Unmarshal(param, &res)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (m *DataCheckObj) HandlerData(titleFieldToDbField map[string]string, dbClient *database.Client, baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) (e error) {

	e = m.repeatabilityCheck(baseData.TableName, baseData.PrimaryKey, titleFieldToDbField, dbClient, xlsxData)
	if e != nil {
		return e
	}

	e = m.legalityCheck(titleFieldToDbField, xlsxData)
	if e != nil {
		return e
	}

	e = m.nonEmptyCheck(titleFieldToDbField, xlsxData)
	if e != nil {
		return e
	}

	e = m.nonZeroCheck(titleFieldToDbField, xlsxData)
	if e != nil {
		return e
	}

	return
}

/*
repeatabilityCheck ----Repeatability Check
重复性校验会把所有数据校验完,有错误才会被抛出。因为每校验一条数据立即抛出检查结果并不友好
每一个校验都会构造一条查询的SQL语句

createRes: map[string][]string, key: 将要验证的字段合并为key，value:[]string, 数据库中每个验证字段的值的组合, 可以有多个。
updateRes：map[string]map[string]string, 第一个key: 检查的字段组合, 第二个key: 检查的字段数据库中的值组合, value: 主键字段的值。

操作列是增时, excel中检查字段的组合值等于createRes的value时, 表示重复
操作列是改时，excel中检查的的字段组合值等于updateRes的第二个key, 并且主键值等于value时, 表示重复

The repeatability check will check all the data, and errors will be thrown. Because it is not friendly to throw the check result immediately after each piece of data is verified
Each check will construct a query SQL statement

createRes: map[string][]string, key: merge the fields to be verified into key, value:[]string. There can be multiple combinations of values of each verified field in the database.
updateRes: map[string]map[string]string, the first key: the checked field combination, the second key: the checked field combination in the database, and value: the value of the primary key field.
When the operation column is 增, and the combined value of the check field in excel is equal to the value of createRes, it indicates repetition
When the operation column is 改, the field combination value checked in excel is equal to the second key of updateRes, and the primary key value is equal to value, indicating repetition
*/
func (m *DataCheckObj) repeatabilityCheck(tableName, primaryKey string,
	titleFieldToDbField map[string]string, dbClient *database.Client,
	xlsxData []*import_model.ImportExcelData)  (e error) {

	errMessMap := make(map[string]*ErrMess)
	fieldNotExistMap := make(map[string]int)
	fieldNotExist := make([]string, 0)
	for _, item := range m.RepeatabilityCheckField {
		dbFieldArr := make([]string, 0)
		titleFieldArr := make([]string, 0)
		for _, titleField := range strings.Split(item.Field, ",") {
			titleField = strings.TrimSpace(titleField)
			if dbField, ok := titleFieldToDbField[titleField]; ok {
				dbFieldArr = append(dbFieldArr, dbField)
				titleFieldArr = append(titleFieldArr, titleField)
			} else {
				if _, ok := fieldNotExistMap[titleField]; !ok {
					fieldNotExistMap[titleField] = 1
					fieldNotExist = append(fieldNotExist, titleField)
				}
			}
		}

		if num := utils.GetIndex(dbFieldArr, primaryKey); num == -1 {
			dbFieldArr = append(dbFieldArr, primaryKey)
		}

		stmt := model.NewStatement("", tableName, "", 0, 0)
		stmt.AddSelects(dbFieldArr...)
		stmt.AddOrders(model.GetOrder(strings.TrimSpace(item.Order))...)

		cond, whereParams, err := model.GetConditionBySlice("", item.Condition)
		if err != nil {
			return
		}
		stmt.AddCondition(cond...)
		stmt.AddParams(whereParams...)

		sql, params := stmt.GetQuerySQL()
		res, e := dbClient.QueryMaps(sql, params...)
		if e != nil {
			return errors.New(fmt.Sprintf("repeatabilityCheck err:%v", e))
		}

		createRes := make(map[string][]string)
		updateRes := make(map[string]map[string]string)
		key :=  strings.Join(titleFieldArr, "-")
		createRes[key] = make([]string, 0)
		updateRes[key] = make(map[string]string)

		for _, resItem := range res {
			dbVal := ""
			for _, titleField := range titleFieldArr {
				if dbField, ok := titleFieldToDbField[titleField]; ok {
					dbVal += resItem[dbField]
				}
			}
			createRes[key] = append(createRes[key], dbVal)
			updateRes[key][dbVal] = resItem[primaryKey]
		}

		for _, data := range xlsxData {
			xlsxVal := ""
			for _, titleField := range titleFieldArr {
				if dbField, ok := titleFieldToDbField[titleField]; ok {
					xlsxVal += data.Data[dbField]
				}
			}
			switch data.Action {
			case  utils.ActionAdd:
				num := utils.GetIndex(createRes[key], xlsxVal)
				if num != -1 {
					if _, ok := errMessMap[data.Action]; !ok {
						errMessMap[data.Action] = new(ErrMess)
					}
					errMessMap[data.Action].addErrMess(utils.RepeatabilityCheck, data.Row, data.Col, titleFieldArr)
				}
			case utils.ActionUpdate:
				if primaryKeyValue, ok := updateRes[key][xlsxVal]; ok {
					if data.PrimaryKeyValue != primaryKeyValue {
						if _, ok := errMessMap[data.Action]; !ok {
							errMessMap[data.Action] = new(ErrMess)
						}
						errMessMap[data.Action].addErrMess(utils.RepeatabilityCheck, data.Row, data.Col, titleFieldArr)
					}
				}
			default:
				continue
			}
		}
	}

	if len(fieldNotExist) != 0 {
		return errors.New(fmt.Sprintf("repeatabilityCheck Error," +
			"the field to be checked does not exist in the excel header field: %v," +
			"please check and try again", strings.Join(fieldNotExist, ",")))
	}

	return getErrMess(errMessMap)
}

func (m *DataCheckObj) legalityCheck(excelTitleFieldToDbField map[string]string, xlsxData []*import_model.ImportExcelData) (e error) {

	errMessMap := make(map[string]*ErrMess)
	fieldNotExistMap := make(map[string]int)
	fieldNotExist := make([]string, 0)
	for _, item := range m.LegalityCheckField {
		for _, titleField := range strings.Split(item.Field, ",") {
			titleField = strings.TrimSpace(titleField)
			for _, data := range xlsxData {
				switch data.Action {
				case utils.ActionAdd, utils.ActionUpdate:
					if dbField, ok := excelTitleFieldToDbField[titleField]; ok {
						re, e := regexp.Compile(item.Reg)
						if e != nil {
							return e
						}

						if ! re.MatchString(data.Data[dbField]) {

							if _, ok := errMessMap[data.Action]; !ok {
								errMessMap[data.Action] = new(ErrMess)
							}
							errMessMap[data.Action].addErrMess(utils.LegalityCheck,
								data.Row, data.Col, []string{titleField})
						}
					} else {
						if _, ok := fieldNotExistMap[titleField]; !ok {
							fieldNotExistMap[titleField] = 1
							fieldNotExist = append(fieldNotExist, titleField)
						}
					}
				default:
					continue
				}
			}
		}
	}

	if len(fieldNotExist) != 0 {
		return errors.New(fmt.Sprintf("legalityCheck Error," +
			"the field to be checked does not exist in the excel header field: %v," +
			"please check and try again", strings.Join(fieldNotExist, ",")))
	}


	return getErrMess(errMessMap)
}

func (m *DataCheckObj) nonEmptyCheck(excelTitleFieldToDbField map[string]string, xlsxData []*import_model.ImportExcelData) (e error) {

	errMessMap := make(map[string]*ErrMess)
	fieldNotExistMap := make(map[string]int)
	fieldNotExist := make([]string, 0)
	for _, titleField := range strings.Split(m.NonEmptyCheckField, ",") {
		titleField = strings.TrimSpace(titleField)
		for _, data := range xlsxData {
			switch data.Action {
			case utils.ActionAdd, utils.ActionUpdate:
				if dbField, ok := excelTitleFieldToDbField[titleField]; ok {
					if data.Data[dbField] == "" || data.Data[dbField] == "-" {
						if _, ok := errMessMap[data.Action]; !ok {
							errMessMap[data.Action] = new(ErrMess)
						}
						errMessMap[data.Action].addErrMess(utils.NonEmptyCheck, data.Row, data.Col, []string{titleField})
					}
				} else {
					if _, ok := fieldNotExistMap[titleField]; !ok {
						fieldNotExistMap[titleField] = 1
						fieldNotExist = append(fieldNotExist, titleField)
					}
				}
			default:
				continue
			}
		}
	}

	if len(fieldNotExist) != 0 {
		return errors.New(fmt.Sprintf("nonEmptyCheck Error," +
			"the field to be checked does not exist in the excel header field: %v," +
			"please check and try again", strings.Join(fieldNotExist, ",")))
	}


	return getErrMess(errMessMap)
}

func (m *DataCheckObj) nonZeroCheck(excelTitleFieldToDbField map[string]string, xlsxData []*import_model.ImportExcelData) (e error) {

	errMessMap := make(map[string]*ErrMess)
	fieldNotExistMap := make(map[string]int)
	fieldNotExist := make([]string, 0)
	for _, titleField := range strings.Split(m.NonZeroCheckField, ",") {
		titleField = strings.TrimSpace(titleField)
		for _, data := range xlsxData {
			switch data.Action {
			case utils.ActionAdd, utils.ActionUpdate:
				if dbField, ok := excelTitleFieldToDbField[titleField]; ok {
					if data.Data[dbField] == "0" || data.Data[dbField] == "-" {
						if _, ok := errMessMap[data.Action]; !ok {
							errMessMap[data.Action] = new(ErrMess)
						}
						errMessMap[data.Action].addErrMess(utils.NonZeroCheck, data.Row, data.Col, []string{titleField})
					}
				} else {
					if _, ok := fieldNotExistMap[titleField]; !ok {
						fieldNotExistMap[titleField] = 1
						fieldNotExist = append(fieldNotExist, titleField)
					}
				}
			default:
				continue
			}
		}
	}

	if len(fieldNotExist) != 0 {
		return errors.New(fmt.Sprintf("nonZeroCheck Error," +
			"the field to be checked does not exist in the excel header field: %v, " +
			"please check and try again", strings.Join(fieldNotExist, ",")))
	}

	return getErrMess(errMessMap)
}

func getErrMess(errMessMap map[string]*ErrMess) (e error) {

	errMess := make([]string, 0)
	for _, err := range errMessMap {
		errMess = append(errMess,  err.getErrMess()...)
	}

	if len(errMess) != 0 {
		return errors.New(strings.Join(errMess, "\n"))
	}

	return nil
}


/*
ErrMess ----Error message of data check during import 导入时数据校验的错误信息


checkType: 数据校验类型, 包括: Repeatability_Check、Legality_Check、NonEmpty_Check、NonZero_Check
row: 数据在excel中按行显示时数据所在的行数
col: 数据在excel中按列显示时数据所在的列数
titleField: 当checkType等于Repeatability_Check时, excel数据校验失败的表头字段
titleFieldByRow: 当checkType不等于Repeatability_Check时, key：excel数据校验失败的表头字段, value：在excel中按行显示数据时数据所在的行数
titleFieldByCol: 当checkType不等于Repeatability_Check时, key：excel数据校验失败的表头字段, value：数据在excel中按列显示时数据所在的列数


checkType: Type of data check, include:Repeatability_Check, Legality_Check, NonEmpty_Check, NonZero_Check.
row: The number of rows where the data is located when the data is displayed by row in excel.
col: The number of columns where the data is located when the data is displayed by column in excel.
titleField:  When checkType is equal to Repeatability_Check, the header fields that fail the excel data check.
titleFieldByRow: When checkType is not equal to Repeatability_Check, key: the header fields that fail the excel data check, value: The number of rows where the data is located when the data is displayed by row in excel.
titleFieldByCol: When checkType is not equal to Repeatability_Check, key: the header fields that fail the excel data check, value: The number of columns where the data is located when the data is displayed by column in excel.
*/
type ErrMess struct {
	checkType  string
	row        []int
	col        []int
	titleField []string
	titleFieldByRow map[string][]int
	titleFieldByCol map[string][]int
}

func (m *ErrMess) addErrMess(checkType string, row, col int, titleFieldArr []string) {

	m.checkType = checkType

	switch checkType {
	case utils.RepeatabilityCheck:
		m.addRepeatabilityCheckErrMess(row, col, titleFieldArr)
	default:
		m.addErrMessByDefault(row, col, titleFieldArr)
	}
}

func (m *ErrMess) addRepeatabilityCheckErrMess(row, col int, titleFieldArr []string) {

	if row != 0 {
		m.row = append(m.row, row)
	}

	if col != 0 {
		m.col = append(m.col, col)
	}

	m.titleField = append(m.titleField, titleFieldArr...)
	m.titleField = utils.RemoveDuplicationByStringSlice(m.titleField)

	return
}

func (m *ErrMess) addErrMessByDefault(row, col int, titleFieldArr []string) {

	if m.titleFieldByRow == nil {
		m.titleFieldByRow = make(map[string][]int, 0)
	}
	if m.titleFieldByCol == nil {
		m.titleFieldByCol = make(map[string][]int, 0)
	}

	for _, titleField := range titleFieldArr {
		if row != 0 {
			if _, ok := m.titleFieldByRow[titleField]; !ok {
				m.titleFieldByRow[titleField] = make([]int, 0)
			}
			m.titleFieldByRow[titleField] = append(m.titleFieldByRow[titleField], row)
		}

		if col != 0 {
			if _, ok := m.titleFieldByCol[titleField]; !ok {
				m.titleFieldByCol[titleField] = make([]int, 0)
			}
			m.titleFieldByCol[titleField] = append(m.titleFieldByCol[titleField], col)
		}
	}

	return
}

func (m *ErrMess) getErrMess() []string {

	errMess := make([]string, 0)


	switch m.checkType {
	case utils.RepeatabilityCheck:
		errMess = append(errMess, m.getRepeatabilityCheckErrMess()...)
	default:
		errMess = append(errMess, m.getCheckErrMessByDefault()...)
	}


	return errMess
}

func (m *ErrMess) getRepeatabilityCheckErrMess() []string {

	errMess := make([]string, 0)
	if len(m.row) != 0 {
		rows := make([]string, 0)
		for _, row := range m.row {
			rows = append(rows, strconv.Itoa(row))
		}

		errMess = append(errMess, fmt.Sprintf("rows:%v, field:%v, checkType:%v error",
			 strings.Join(rows, ","), strings.Join(m.titleField, ","), m.checkType))
	}

	if len(m.col) != 0 {
		cols := make([]string, 0)
		for _, col := range m.col {
			cols = append(cols, strconv.Itoa(col))
		}
		errMess = append(errMess, fmt.Sprintf("columns:%v, field:%v, checkType:%v error",
			strings.Join(cols, ","),
			strings.Join(m.titleField, ","), m.checkType))
	}

	return errMess
}

func (m *ErrMess) getCheckErrMessByDefault() []string {

	rowsErrMess := make([]string, 0)
	colsErrMess := make([]string, 0)
	for titleField, rows := range m.titleFieldByRow {
		rowsArr := make([]string, 0)
		for _, row := range rows {
			rowsArr = append(rowsArr, strconv.Itoa(row))
		}
		rowsErrMess = append(rowsErrMess, fmt.Sprintf("rows:%v, field:%v",
			strings.Join(rowsArr, ","),
			titleField))
	}

	for titleField, cols := range m.titleFieldByCol {
		colsArr := make([]string, 0)
		for _, col := range cols {
			colsArr = append(colsArr, strconv.Itoa(col))
		}
		colsErrMess = append(colsErrMess, fmt.Sprintf("columns:%v, field:%v",
			strings.Join(colsArr, ","),
			titleField))
	}


	errMess := make([]string, 0)
	if len(rowsErrMess) != 0 {
		errMess = append(errMess, fmt.Sprintf("%v, checkType:%v error", strings.Join(rowsErrMess, ","), m.checkType))
	}

	if len(colsErrMess) != 0 {
		errMess = append(errMess, fmt.Sprintf("%v, checkType:%v error", strings.Join(colsErrMess, ","), m.checkType))
	}


	return errMess
}
