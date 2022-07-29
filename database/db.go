package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	//引入mysql
	_ "github.com/go-sql-driver/mysql"
)

// Client ----mysql client
type Client struct {
	dsn string
	Db *sql.DB
	Tx *sql.Tx
}

//GetDb ----connection mysql
func GetDb(dsn string) (dbClient *Client, e error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil,  errors.New(fmt.Sprintf("GetDb,connection mysql failed: %v", err))
	}

	tx, e := db.Begin()
	if e != nil {
		return nil, errors.New(fmt.Sprintf("GetDb, Begin starts a transaction failed: %v", e))
	}

	dbClient = &Client{
		dsn: dsn,
		Db:   db,
		Tx:    tx,
	}

	return dbClient, nil
}

//BaseQuery ----base query
func (c *Client) BaseQuery(query string, args ...interface{}) (rows *sql.Rows, e error) {

	rows, e = c.Db.Query(query, args...)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("sql:%v, args:%v failed, BaseQuery err:%v", query, args, e))
	}

	return rows, nil
}

//QueryMaps ----query data, []map[string]string, key：database field, value：the value of database field.
func (c *Client) QueryMaps(query string, args ...interface{}) (res []map[string]string, e error) {

	rows, e := c.BaseQuery(query, args...)
	if e != nil {
		return nil, e
	}

	res = make([]map[string]string, 0)
	for rows.Next() {
		item, e := nextRows(rows)
		if e != nil {
			return nil, e
		}

		res = append(res, item)
	}


	return res, nil
}

//nextRows ----nextRows
func nextRows(rows *sql.Rows) (res map[string]string, e error){

	columns, e := rows.Columns()
	if e != nil {
		return nil, errors.New(fmt.Sprintf("nextRows, get Columns err: %v", e))
	}

	columnTypes, e := rows.ColumnTypes()
	if e != nil {
		return nil, errors.New(fmt.Sprintf("nextRows, get ColumnTypes err: %v", e))
	}

	var row []interface{}
	columnLen := len(columns)

	/*
	根据 mysql 的字段类型设置对应的 go数据类型，且必须准确设置，否则转换内容会出错
	此时并没有真正的把查询到的数据库数据赋值,只是把数据设置成对应的类型
	 */
	for i :=0 ; i < columnLen; i++ {
		columnType := columnTypes[i]
		dbType := strings.ToLower(columnType.DatabaseTypeName())
		goType, ok := MysqlTypeToGoType[dbType]
		if !ok {
			goType = "GoString"
		}

		switch goType {
		case "GoBool":
			var val GoBool
			row = append(row, &val)
		case "GoInt":
			var val GoInt
			row = append(row, &val)
		case "GoFloat":
			var val GoFloat
			row = append(row, &val)
		case "GoString":
			var val GoString
			row = append(row, &val)
		case "GoTime":
			var val GoTime
			row = append(row, &val)
		default:
			var val interface{}
			row = append(row, &val)
		}
	}

	//查询到的数据库数据赋值，根据具体类型调用Scan方法
	e = rows.Scan(row...)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("nextRows, Scan err: %v", e))
	}

	//值类型都转化为string类型，excel导出时使用string格式的更容易操作
	res = make(map[string]string)
	if row != nil {
		for k, column := range columns {
			v := row[k]
			if v == nil {
				continue
			}

			vType := reflect.TypeOf(v).String()
			switch vType {
			case "*database.GoTime":
				res[column] = fmt.Sprintf("%v", v.(*GoTime))
			case "*database.GoInt":
				res[column] = fmt.Sprintf("%v", v.(*GoInt))
			case "*database.GoString":
				res[column] = fmt.Sprintf("%v",v.(*GoString))
			case "*database.GoFloat":
				res[column] = fmt.Sprintf("%v",v.(*GoFloat))
			case "*database.GoBool":
				res[column] = fmt.Sprintf("%v",v.(*GoBool))
			}
		}
	}

	return res, nil
}

