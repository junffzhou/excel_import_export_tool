package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"regexp"
	"strings"
)

//GetOrder ==== +:升序,默认; -:降序
func GetOrder(order string) (res []string) {

	for _, item := range strings.Split(order, ",") {
		if item == "" {
			continue
		}
		if strings.HasSuffix(item, "+") {
			res = append(res, item[:len(item)-1] + ` ASC `)
		} else if strings.HasSuffix(item, "-") {
			res = append(res, item[:len(item)-1] + ` DESC `)
		} else {
			res = append(res, item)
		}
	}

	return
}

//GetConditionBySliceMap ----通过[]map[string][]string类型的参数获取condition条件
func GetConditionBySliceMap(hasJoin bool, tableName, alias string, condition []map[string][]string) (cond []string, whereParams []interface{}, e error) {

	cond = make([]string, 0)
	whereParams = make([]interface{}, 0)
	for _, condItem := range condition {
		for tableKey, columnKey := range condItem {
			tableArr := getTableArr(hasJoin, tableName, alias, tableKey)
			for _, col := range columnKey {
				keyArr := strings.Split(col, ":")
				if len(keyArr) <= 1 {
					continue
				}

				conditionItem, whereParamsItem, e := getConditionItem(tableArr, keyArr)
				if e != nil {
					return cond, whereParams, e
				}
				cond = append(cond, conditionItem)
				whereParams = append(whereParams, whereParamsItem...)
			}
		}
	}

	//去掉 where条件的第一个 "AND", "OR" 连接符
	if len(cond) > 0 {
		cond[0] = strings.Replace(cond[0], "AND", "", 1)
		cond[0] = strings.Replace(cond[0], "OR", "", 1)
	}

	return
}

/*
getTableArr----获取表名
需要join, 查询数据的表存在别名时,设置为别名;不存在别名时,设置为原始表名,其他表的表名不变
无需join时,查询数据的表存在别名时设置为别名,不存在别名时设置为空,其他表被忽略
 */
func getTableArr(hasJoin bool, tableName, alias, tableKey string) (tableArr []string){

	tableArr = make([]string, 0)
	if hasJoin {
		/*if alias == "" {
			alias = tableName
		}
		for _, table := range strings.Split(tableKey, ",") {
			if table == "" {
				tableArr = append(tableArr, alias)
			} else {
				tableArr = append(tableArr, table)
			}
		}*/
		tableAlias := alias
		if alias == "" {
			tableAlias = tableName
		}
		for _, table := range strings.Split(tableKey, ",") {
			table = strings.TrimSpace(table)
			if table == "" || table == tableName || table == alias {
				tableArr = append(tableArr, tableAlias)
			} else {
				tableArr = append(tableArr, table)
			}
		}
	} else {
		for _, table := range strings.Split(tableKey, ",") {
			if (table == "" && alias != "") || table == alias {
				tableArr = append(tableArr, alias)
			} else if table == "" || table == tableName {
				tableArr = append(tableArr, "")
			}
		}
	}

	return
}

//GetConditionBySlice ----通过[]string类型的参数 获取condition条件
func GetConditionBySlice(tableName string, condition []string) (cond []string, whereParams []interface{}, e error) {

	cond = make([]string, 0)
	whereParams = make([]interface{}, 0)
	for _, condItem := range condition {
		conditionItem, whereParamsItem, e := getConditionItem([]string{tableName}, strings.Split(condItem, ":"))
		if e != nil {
			return cond, whereParams, e
		}
		cond = append(cond, conditionItem)
		whereParams = append(whereParams, whereParamsItem...)
	}

	//去掉 where条件的第一个 "AND", "OR" 连接符
	if len(cond) > 0 {
		cond[0] = strings.Replace(cond[0], "AND", "", 1)
		cond[0] = strings.Replace(cond[0], "OR", "", 1)
	}

	return
}

func getConditionItem(tableArr, keyArr []string) (cond string, whereParams []interface{}, e error) {

	connector, column, operator, e := matchOperator( strings.TrimSpace(keyArr[0]))
	if e != nil {
		return cond, whereParams, e
	}

	value := strings.Replace(strings.TrimSpace(keyArr[1]), ")", "", -1)
	ending := " "
	if strings.Contains(keyArr[1], ")") {
		ending = ") "
	}

	item := ""
	for _, table := range tableArr {
		preFix := column
		if table != "" {
			preFix = `.` + preFix
		}

		switch operator {
		case "IN", "NOT IN":
			suffix, paramArr, e := handIN(value)
			if e != nil {
				return cond, whereParams, e
			}
			item =  preFix + ` ` + operator + suffix
			whereParams = append(whereParams, paramArr...)
		case "<>":
			item =  preFix + ` ` + ` BETWEEN ? AND ? `
			paramArr, e := stringTransformSliceInterface(value)
			if e != nil {
				return cond, whereParams, e
			}
			whereParams = append(whereParams, paramArr...)
		case "=", "!=":
			if value == utils.NullValue {
				item = preFix + ` ` + utils.NullMap[operator] + ` `
			} else {
				item = preFix + ` ` + operator + ` ? `
				whereParams = append(whereParams, value)
			}
		default:
			item = preFix + ` ` + operator + ` ? `
			whereParams = append(whereParams, value)
		}
	}
	cond = connector + ` ` + strings.Join(tableArr,  item + strings.Split(connector, " ")[0] + ` `) + item + ending

	return
}

func matchOperator(key string) (connector, column, operator string, e error) {

	/*re, e := regexp.Compile(`(\w+.*\w+)`)
	if e != nil {
		return
	}

	column = re.FindString(key)
	res := re.Split(key, -1)

	switch strings.TrimSpace(res[0]) {
	case "", "&":
		connector = "AND"
	case "|":
		connector = "OR"
	case "(", "&(":
		connector = "AND ("
	case "|(":
		connector = "OR ("
	}

	switch strings.TrimSpace(res[1]) {
	case "", "=":
		operator = "="
	case "!=":
		operator = "!="
	case "$":
		operator = "LIKE"
	case "~":
		operator = "REGEXP"
	case "<":
		operator = "<"
	case "<=":
		operator = "<="
	case ">":
		operator = ">"
	case ">=":
		operator = ">="
	case "<>":
		operator = "<>"
	case "{}":
		operator = "IN"
	case "!{}":
		operator = "NOT IN"
	case "!":
		operator = "!"
	default:
		operator = ""
	}*/

	connRe, e := regexp.Compile(`^([|&])*(\(*)`)
	if e != nil {
		return connector, column, operator, e
	}

	connRes := connRe.FindStringSubmatch(key)
	if len(connRes) == 0 {
		return
	}

	connector = "AND "
	if connRes[1] == "|" {
		connector = "OR " + connRes[2]
	} else {
		connector +=  connRes[2]
	}

	operatorRe, e := regexp.Compile(`(=|!=|\$|~|<|<=|>|>=|<>|{}|!{}|!)$`)
	if e != nil {
		return connector, column, operator, e
	}
	opRes := operatorRe.FindStringSubmatch(connRe.ReplaceAllString(key, ""))
	operator = "="
	if len(opRes) != 0  {
		switch opRes[1] {
		case "=":
			operator = "="
		case "!=":
			operator = "!="
		case "$":
			operator = "LIKE"
		case "~":
			operator = "REGEXP"
		case "<":
			operator = "<"
		case "<=":
			operator = "<="
		case ">":
			operator = ">"
		case ">=":
			operator = ">="
		case "<>":
			operator = "<>"
		case "{}":
			operator = "IN"
		case "!{}":
			operator = "NOT IN"
		case "!":
			operator = "!"
		}
	}

	column = strings.TrimSpace(operatorRe.ReplaceAllString(connRe.ReplaceAllString(key, ""), ""))

	return
}

//handIN ----处理 IN 或者 NOT IN 语句
func handIN(value string) (suffix string, paramArr []interface{}, e error) {

	paramArr, e = stringTransformSliceInterface(value)
	if e != nil {
		return
	}

	suffix = `(?`
	for i := 1; i < len(paramArr); i++ {
		suffix = suffix + `, ?`
	}
	suffix += `) `

	return
}

/*
stringTransformSliceInterface

key: `\"[\"aaa\", \"bbb\", \"ccc\", ...]\"`
*/
func stringTransformSliceInterface(key string) (res []interface{}, e error) {

	res = make([]interface{}, 0)
	e = json.Unmarshal([]byte(key), &res)
	if e != nil {
		return res, errors.New(fmt.Sprintf("stringTransformSliceInterface error:%v", e))
	}

	return
}