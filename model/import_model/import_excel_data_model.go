package import_model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/database"
	"github.com/junffzhou/excel_import_export_tool/model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"log"
	"strings"
)

/*
ImportExcelData ----excel data

Row: When the excel content is displayed by rows, which row of data in excel.
Col: When the excel content is displayed by columns, which column of data in excel.
PrimaryKeyValue: The value of the primary key field.
Action: The value of the action field.
Data: data from excel, map[string]string: key：excel header field, value：field value
*/
type ImportExcelData struct {
	Row             int
	Col             int
	PrimaryKeyValue string
	Action          string
	Data            map[string]string
}


type ImportExcelDataArray []*ImportExcelData

func (m ImportExcelDataArray) ActionFunc(titleColumnToDbColumn map[string]string, dbClient *database.Client, baseData *ImportRequestParam) (e error) {

	tx := dbClient.Tx
	for _, data := range m {
		e = data.ExecByAction(titleColumnToDbColumn, tx, baseData)
		if e != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil{
				return errors.New(fmt.Sprintf("rollback failed: %v", rollbackErr))
			}
			return e
		}
	}

	e = tx.Commit()
	if e != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil{
			return errors.New(fmt.Sprintf("rollback failed: %v", rollbackErr))
		}
	}

	return
}

/*
ExecByAction ----根据action值执行sql语句
忽略的情况:
1) 操作字段的值不是增/删/改
2) 包含操作字段, 不包含主键字段且主键字段不是自增, 任意操作都被忽略
3) 包含操作字段, 不包含主键字段, 主键字段自增, 修改和删除操作被忽略
4) 主键字段值为空, 修改和删除操作被忽略

execute sql statement according to action value
Ignored cases:
1) The value of the operation field is not 增/删/改
2) Contains the operation field, does not contain the primary key field and the primary key field is not self-incrementing, any operation is ignored
3) Contains the operation field, does not contain the primary key field, the primary key field is auto-incremented, and the modification and deletion operations are ignored
4) The primary key field value is empty, and the modification and deletion operations are ignored
*/
func (m *ImportExcelData) ExecByAction(titleColumnToDbColumn map[string]string, tx *sql.Tx, baseData *ImportRequestParam) (e error) {

	switch m.Action {
	case utils.ActionAdd:
		//不包含主键字段,主键非自增,忽略
		//不包含主键字段,主键自增; 包含主键字段,主键自增; 包含主键字段,主键非自增
		if baseData.PrimaryKey == "" && baseData.AutoIncrement == utils.NotAutoIncrement  {
			return
		}
		e = m.Create(titleColumnToDbColumn, tx, baseData)
		if e != nil {
			return e
		}
	case utils.ActionUpdate:
		//不包含主键字段,或者主键字段值为空,忽略
		if m.PrimaryKeyValue == "" || baseData.PrimaryKey == "" {
			return
		}
		e = m.Update(titleColumnToDbColumn, tx, baseData)
		if e != nil {
			return e
		}
	case utils.ActionDelete:
		//不包含主键字段,或者主键字段值为空,忽略
		if m.PrimaryKeyValue == "" || baseData.PrimaryKey == "" {
			return
		}
		e = m.Delete(tx, baseData)
		if e != nil {
			return e
		}
	default:
		return
	}

	return
}

func (m *ImportExcelData) Create(titleColumnToDbColumn map[string]string, tx *sql.Tx, baseData *ImportRequestParam)  error {

	insertSQL, args := m.getInsertSQL(titleColumnToDbColumn, baseData)
	_, e := tx.Exec(insertSQL, args...)
	if e != nil {
		if m.Row != 0 {
			return errors.New(fmt.Sprintf("row:%v, Create failed:%v", m.Row,  e))
		}

		return errors.New(fmt.Sprintf("col:%v, Create failed:%v", m.Col,  e))
	}

	return nil
}

func (m *ImportExcelData) Update(titleColumnToDbColumn map[string]string, tx *sql.Tx, baseData *ImportRequestParam) error {

	updateSQL, args, e := m.getUpdateSQL(titleColumnToDbColumn, baseData)
	if e != nil {
		return e
	}

	_, e = tx.Exec(updateSQL, args...)
	if e != nil {
		if m.Row != 0 {
			return errors.New(fmt.Sprintf("row:%v, Update failed:%v", m.Row, e))
		}

		return errors.New(fmt.Sprintf("col:%v, Update failed:%v", m.Col, e))
	}


	return nil
}

func (m *ImportExcelData) Delete(tx *sql.Tx, baseData *ImportRequestParam,) error {

	deleteSQL, args, e := m.getDeleteSQL(baseData)
	if e != nil {
		return e
	}
	_, e = tx.Exec(deleteSQL, args...)
	if e != nil {
		if m.Row != 0 {
			return errors.New(fmt.Sprintf("row:%v, Delete failed:%v", m.Row, e))
		}
		return errors.New(fmt.Sprintf("col:%v, Delete failed:%v", m.Col, e))
	}

	return nil
}

/*
getInsertSQL ----build insert sql statements
INSERT INTO table_name (column1, column2,...) VALUES (value1, value2,....)
*/
func (m *ImportExcelData) getInsertSQL(titleColumnToDbColumn map[string]string, baseData *ImportRequestParam) (sql string, args []interface{}) {

	args = make([]interface{}, 0)
	dbColumnArr := make([]string, 0)
	suffixArr := make([]string, 0)
	//主键字段不是自增,PrimaryKey一定不能为空: 值不为空,使用当前值作为主键值; 值为空,设置主键字段的值. 插入时包括主键字段
	//主键字段自增: 值不为空时,使用当前值作为主键的值; 值为空,插入时不包括主键字段
	if baseData.AutoIncrement  == utils.NotAutoIncrement {
		if m.PrimaryKeyValue == "" {
			m.PrimaryKeyValue = utils.GetUUid()
		}
		args = append(args, m.PrimaryKeyValue)
		dbColumnArr = append(dbColumnArr, baseData.PrimaryKey)
		suffixArr = append(suffixArr, "?")
	} else {
		if m.PrimaryKeyValue  != "" {
			args = append(args, m.PrimaryKeyValue)
			dbColumnArr = append(dbColumnArr, baseData.PrimaryKey)
			suffixArr = append(suffixArr, "?")
		}
	}

	for titleColumn, value := range m.Data {
		if dbColumn, ok := titleColumnToDbColumn[titleColumn]; ok && dbColumn != baseData.PrimaryKey {
			args = append(args, value)
			dbColumnArr = append(dbColumnArr, dbColumn)
			suffixArr = append(suffixArr, "?")
		}
	}

	var (
		addParams    []interface{}
		addColumnArr []string
		addSuffixArr []string
	)
	if baseData.AdditionalCondition != nil && baseData.AdditionalCondition.CreateAddCondition != nil {
		var whereParams []interface{}
		for dbColumn, value := range baseData.AdditionalCondition.CreateAddCondition.ExtraColumn {
			if num := utils.GetIndex(dbColumnArr, dbColumn); num == -1 && dbColumn != baseData.PrimaryKey {
				addParams = append(addParams, value)
				addColumnArr = append(addColumnArr, dbColumn)
				addSuffixArr = append(addSuffixArr, "?")
			}
		}
		addParams = append(addParams, whereParams...)
	}


	sql = `INSERT INTO ` + baseData.TableName + ` (` + strings.Join(append(dbColumnArr, addColumnArr...), ", ") +
		`) VALUES(` + strings.Join(append(suffixArr, addSuffixArr...), ", ") + `)`
	args = append(args, addParams...)

	log.Printf("getInsertSQL, INSERT INTO的SQL：%v, 参数值args：%v", sql, args)
	return sql, args
}


/*
getUpdateSQL ----build update sql statements
UPDATE table_name SET column1 = newValue WHERE column = value
*/
func (m *ImportExcelData) getUpdateSQL(titleColumnToDbColumn map[string]string, baseData *ImportRequestParam) (sql string, args []interface{}, e error) {

	dbColumnArr := make([]string, 0)
	args = make([]interface{}, 0)
	for titleColumn, value := range m.Data {
		if dbColumn, ok := titleColumnToDbColumn[titleColumn]; ok && dbColumn != baseData.PrimaryKey {
			dbColumnArr = append(dbColumnArr, dbColumn + ` = ?`)
			args = append(args, value)
		}
	}


	var (
		addParams    []interface{}
		addColumnArr []string
		cond []string
		addWhere   string
	)
	if baseData.AdditionalCondition != nil && baseData.AdditionalCondition.UpdateAddCondition  != nil {
		 whereParams := make([]interface{}, 0)
		for dbColumn, value := range  baseData.AdditionalCondition.UpdateAddCondition.ExtraColumn {
			if num := utils.GetIndex(dbColumnArr, dbColumn); num == -1 && dbColumn != baseData.PrimaryKey {
				addColumnArr = append(addColumnArr, dbColumn + " = ?")
				addParams = append(addParams, value)
			}
		}

		cond, whereParams, e = model.GetConditionBySlice("", baseData.AdditionalCondition.UpdateAddCondition.Condition)
		if e != nil {
			return
		}
		addParams = append(addParams, whereParams...)
	}


	if len(cond) != 0 {
		addWhere = ` WHERE ` + strings.Join(cond, "") + ` AND ` + baseData.PrimaryKey + ` = ?`
	} else {
		addWhere = ` WHERE ` + baseData.PrimaryKey + ` = ?`
	}


	sql = `UPDATE ` + baseData.TableName + ` SET ` +  strings.Join(append(dbColumnArr, addColumnArr...), `, `) + addWhere
	args = append(args, append(addParams, m.PrimaryKeyValue)...)

	log.Printf("getUpdateSQL, UPDATE的SQL：%v, 参数值args：%v", sql, args)
	return
}

/*
getDeleteSQL ---- build delete sql statements
If baseData.AdditionalCondition.DeleteAddCondition.IsSoftDelete is equal to 00, build update sql statements: UPDATE table_name SET column1 = newValue WHERE column = value;
otherwise, build delete sql statements: DELETE table_name  WHERE column = value
*/
func (m *ImportExcelData) getDeleteSQL(baseData *ImportRequestParam) (sql string, args []interface{}, e error) {

	var (
		addParams    []interface{}
		addColumnArr []string
		addWhere     string
		cond []string
		whereParams  []interface{}
		isSoftDelete string
	)

	if baseData.AdditionalCondition != nil && baseData.AdditionalCondition.DeleteAddCondition != nil {
		isSoftDelete = baseData.AdditionalCondition.DeleteAddCondition.IsSoftDelete
		for dbColumn, value := range baseData.AdditionalCondition.DeleteAddCondition.ExtraColumn {
			if dbColumn != baseData.PrimaryKey {
				addColumnArr = append(addColumnArr, dbColumn + " = ?")
				addParams = append(addParams, value)
			}
		}
		cond, whereParams, e = model.GetConditionBySlice("", baseData.AdditionalCondition.UpdateAddCondition.Condition)
		if e != nil {
			return
		}
		addParams = append(addParams, whereParams...)
	}

	if len(cond) != 0 {
		addWhere =  ` WHERE ` + strings.Join(cond, "") + ` AND ` + baseData.PrimaryKey + ` = ?`
	} else {
		addWhere = ` WHERE ` + baseData.PrimaryKey + ` = ?`
	}

	args = make([]interface{}, 0)
	if isSoftDelete == utils.NotSoftDelete {
		sql = `DELETE ` + baseData.TableName + addWhere
		args = append(args,  append(whereParams, m.PrimaryKeyValue)...)
	} else {
		sql = `UPDATE ` + baseData.TableName + ` SET ` +  strings.Join(addColumnArr, `, `) + addWhere
		args = append(args, append(addParams, m.PrimaryKeyValue)...)
	}

	log.Printf("getDeleteSQL, SQL：%v, 参数值args：%v", sql, args)
	return
}

