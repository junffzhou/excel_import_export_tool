package model

import (
	"fmt"
	"log"
	"strings"
)

type Statement struct {
	schema     string
	tableName string
	alias string
	condition  string
	groupBy    string
	having string
	limit      int
	offset     int
	selects   []string
	joins     []string
	orders     []string
	params     []interface{}
}

func NewStatement(schema, tableName, alias string, limit, offset int) *Statement {
	return &Statement{schema: schema, tableName: tableName, alias: alias, limit: limit, offset: offset}
}

func (stmt *Statement) GetTableName() string {
	return stmt.tableName
}

func (stmt *Statement) GetAlias() string {
	return stmt.alias
}

func (stmt *Statement) GetOrder() string {
	return strings.Join(stmt.orders, ", ")
}

func (stmt *Statement) GetParams() []interface{} {
	return stmt.params
}

func (stmt *Statement) SetSchema(schema string) {
	stmt.schema = schema
	return
}

func (stmt *Statement) SetTableName(tableName string) {
	stmt.tableName = tableName
	return
}

func (stmt *Statement) SetAlias(alias string) {
	stmt.alias = alias
	return
}

func (stmt *Statement) SetCondition(condition string) {

	stmt.condition = condition
	return
}

func (stmt *Statement) SetLimit(limit int) {
	stmt.limit = limit
	return
}

func (stmt *Statement) SetOffset(offset int) {
	stmt.offset = offset
	return
}

func (stmt *Statement) SetGroupBy(group ...string) {

	if len(group) == 1 {
		stmt.groupBy = group[0]
	} else if len(group) > 1 {
		stmt.groupBy = "`" + strings.Join(group, "`,`") + "`"
	}

	return
}

func (stmt *Statement) AddSelects(selects ...string) {

	if stmt.selects == nil {
		stmt.selects = make([]string, 0)
	}

	stmt.selects = append(stmt.selects, selects...)
	return
}

func (stmt *Statement) AddCondition(cond ...string) {

	stmt.condition += strings.Join(cond, "")

	return
}

func (stmt *Statement) AddJoins(joins ...string) {

	if stmt.joins == nil {
		stmt.joins = make([]string, 0)
	}

	stmt.joins= append(stmt.joins, joins...)
	return
}

func (stmt *Statement) AddOrders(order ...string) {

	if stmt.orders == nil {
		stmt.orders = make([]string, 0)
	}

	stmt.orders = append(stmt.orders, order...)

	return
}

func (stmt *Statement) AddParams(val ...interface{}) {

	if stmt.params == nil {
		stmt.params = make([]interface{}, 0)
	}

	stmt.params = append(stmt.params, val...)

	return
}

func (stmt *Statement) GetQuerySQL() (sql string, params []interface{}) {

	sqlPre := fmt.Sprint(stmt.tableName)
	if stmt.schema != "" {
		sqlPre = fmt.Sprint(stmt.schema, `.`, sqlPre)
	}
	if stmt.alias != "" {
		sqlPre = fmt.Sprint(sqlPre, ` `, stmt.alias)
	}

	sql = fmt.Sprint(`SELECT `, strings.Join(stmt.selects, `, `), ` FROM `, sqlPre)

	for _, joinItem := range stmt.joins {
		sql = fmt.Sprint(sql, ` `, joinItem, ` `)
	}

	if stmt.condition != "" {
		sql = fmt.Sprint(sql, ` WHERE ` + stmt.condition)
	}
	if stmt.groupBy != "" {
		sql = fmt.Sprint(sql, ` GROUP BY `, stmt.groupBy)
	}
	if stmt.having != "" {
		sql = fmt.Sprint(sql, ` HAVING `, stmt.having)
	}
	if len(stmt.orders) > 0 {
		sql = fmt.Sprint(sql, ` ORDER BY `, stmt.GetOrder())
	}
	if stmt.limit > 0 {
		sql = fmt.Sprint(sql, ` LIMIT `, stmt.limit)
	}
	if stmt.offset > 0 {
		sql = fmt.Sprint(sql, ` OFFSET `, stmt.offset)
	}

	log.Printf("sql statement:%#v\nsql params:%#v\n\n", sql, stmt.params)
	return sql, stmt.params
}
