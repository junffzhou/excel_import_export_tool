package export_model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/junffzhou/excel_import_export_tool/model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"strings"
)

/*
ExportRequestParam ----Export request  parameters

SchemaName: Database name.
TableName: Table name.
TableAlias: Table alias.
BookName: Excel name.
SheetName: Sheet name of excel.
ShowType: When the excel content is displayed by rows, it is equal to 00; when excel content is displayed by columns, it is equal to 01. The default value is 00.
PrimaryKey: Primary key column name in database.
IsIncludeActionField: When the excel header contains the operation column, it is equal to "00"; otherwise,it is equal to "01". The default value is 01.
ActionFieldIndex: When the excel content is displayed in rows, the number of columns  where the action column is located,the defaults is the second column;
				  when displaying by column, the number of rows  where the action column is located,the defaulting is the second row.Rows and columns start at 1.
ExportColumn: export column
Condition: where condition when executing sql statement
Order: order
Limit: limit
Offset: offset
*/
type ExportRequestParam struct {
	SchemaName           string                `json:"schema_name"`
	TableName            string                `json:"table_name"`
	TableAlias           string                `json:"table_alias"`
	BookName             string                `json:"book_name"`
	SheetName            string                `json:"sheet_name"`
	ShowType             string                `json:"show_type"`
	PrimaryKey           string                `json:"primary_key"`
	IsIncludeActionField string                `json:"is_include_action_field"`
	ActionFieldIndex     float64               `json:"action_field_index"`
	ExportColumn         []*ExportColumn        `json:"export_column"`
	Condition            []map[string][]string `json:"condition"`
	Order                string                `json:"order"`
	Limit                float64               `json:"limit"`
	Offset               float64               `json:"offset"`
}


/*
ExportColumn ----导出字段详情

Column: 数据库中的字段,别名,以及导出时的表头
FromSchema: 库名, FromTable 和 TableName 不是同一个库时,必需指定
FromTable: 表名, Column来自哪张表, 默认来自 TableName 表
FromTableAlias: FromTable表的别名
JoinType: Column不是来自 TableName 时, join的方式:"00": left join,默认; "01": right join; "02": inner join
RelatedColumn:  FromTable中关联的字段
RelTableColumn: join on 关联的表和字段

JoinType: When Column is not from TableName, the way to join: "00": left join, default; "01": right join; "02": inner join
RelatedColumn: Related fields in FromTable
RelTableColumn: join on associated tables and fields
 */
type ExportColumn struct {
	Column string `json:"column"`
	FromSchema     string `json:"from_schema"`
	FromTable      string `json:"from_table"`
	FromTableAlias string `json:"from_table_alias"`
	JoinType       string `json:"join_type"`
	RelatedColumn  string `json:"related_column"`
	RelTableColumn string `json:"rel_table_column"`
}

//GetExportReq ----Construct request parameter structure, oldField validation
func GetExportReq(req []byte) (res *ExportRequestParam, e error) {

	res = new(ExportRequestParam)
	e = json.Unmarshal(req, &res)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("GetExportReq error: %v", e))
	}

	if res.TableName == "" {
		return nil, errors.New("GetExportReq err: table_name is a required parameter")
	}

	if res.BookName == "" {
		return nil, errors.New("GetExportReq err: book_name is a required parameter")
	}

	if len(res.ExportColumn) == 0 {
		return nil, errors.New("GetExportReq err: export_column is a required parameter")
	}


	res.setShowType()
	res.setInclude()

	if res.PrimaryKey == "" &&  res.IsIncludeActionField == utils.Include && res.ActionFieldIndex <= 0 {
		return nil, errors.New(fmt.Sprintf("GetExportReq err: when primary_key=\"\" and is_include_action_field=%v, " +
			"action_field_index >= 1",  utils.Include))
	}


	return
}

func (m *ExportRequestParam) setShowType() {

	if  m.ShowType == utils.ShowTypeByCol {
		m.ShowType = utils.ShowTypeByCol
		return
	}

	m.ShowType = utils.ShowTypeByRow
	return
}

func (m *ExportRequestParam) setInclude() {

	if m.IsIncludeActionField == utils.Include {
		m.IsIncludeActionField = utils.Include
		return
	}

	m.IsIncludeActionField = utils.NotInclude
	return
}

/*
GetActionFieldIndex ----get action oldField index

不包含操作列,ActionFieldIndex是否有值都会被忽略;
包含操作列,ActionFieldIndex>0,根据ActionFieldIndex值插入;
包含操作列,ActionFieldIndex<=0,指定了主键,默认插入主键之后;
其他情况也会忽略操作列
 */
func (m *ExportRequestParam) GetActionFieldIndex(dbColumn []string) int {

	if m.IsIncludeActionField == utils.NotInclude {
		return -1
	} else {
		if m.ActionFieldIndex > 0 {
			return int(m.ActionFieldIndex) - 1
		} else if m.PrimaryKey != "" {
			index := utils.GetIndex(dbColumn, m.PrimaryKey)
			if index == -1 {
				return -1
			}
			return index + 1
		} else {
			return -1
		}
	}
}

func (m *ExportRequestParam) GetStmtAndColumn() (dbColumn, titleColumn []string, stmt *model.Statement, e error) {

	stmt = model.NewStatement(m.SchemaName, m.TableName, m.TableAlias, int(m.Limit), int(m.Offset))
	stmt.AddOrders(model.GetOrder(strings.TrimSpace(m.Order))...)

	//判断是否需要join, 需要join时,筛选的字段前面需要加上表名: select table1.select_field1, table2.select_field2,... from table1 left join table2  on table.id=table2.id where ...;
	//否则,不加: select select_field1, select_field2,... from table1  where ...
	hasJoin := false
	for _, item := range m.ExportColumn {
		if item.FromTable != "" && item.FromTable != m.TableName && item.FromTable != m.TableAlias {
			item.setJoinType()
			hasJoin = true
			break
		}
	}

	cond, whereParams, e := model.GetConditionBySliceMap(hasJoin, m.TableName, m.TableAlias, m.Condition)
	if e != nil {
		return
	}
	stmt.AddCondition(cond...)
	stmt.AddParams(whereParams...)

	for _, item := range m.ExportColumn {
		dbColumnItem, titleColumnItem, selectsItem, joinsItem := item.getColumn(hasJoin, m.TableName, m.TableAlias)
		stmt.AddSelects(selectsItem...)
		stmt.AddJoins(joinsItem...)
		dbColumn = append(dbColumn, dbColumnItem...)
		titleColumn = append(titleColumn, titleColumnItem...)
	}

	return
}

//setJoinType ---- 设置Join类型: "00": left join,默认; "01":right join; "02":inner join
func (item *ExportColumn) setJoinType() {

	switch item.JoinType {
	case "", utils.LeftJoin:
		item.JoinType = utils.LeftJoin
	case utils.RightJOIN:
		item.JoinType = utils.RightJOIN
	case utils.InnerJoin:
		item.JoinType = utils.InnerJoin
	}

	return
}


func (item *ExportColumn) getColumn(isIncludeJoin bool, tableName, alias string) (dbColumn, titleColumn, selects, joins []string) {

	dbColumn, titleColumn, selects, joins = make([]string, 0), make([]string, 0), make([]string, 0), make([]string, 0)

	//不管是否需要Join,先设置fromTableAlias的值
	//只要别名存在则使用别名;别名不存在,且需要Join时使用原始表名
	//不需要Join且别名不存在,fromTableAlias = ""
	fromTableAlias := ""
	if alias != "" {
		fromTableAlias = alias
	} else if isIncludeJoin {
		fromTableAlias = tableName
	}

	//需要Join时,重新设置fromTableAlias, fromTableAlias = item.FromTable 或者 fromTableAlias = item.FromTableAlias
	if item.FromTable != "" && item.FromTable != tableName && item.FromTableAlias != alias {
		prefix := ""
		prefix, fromTableAlias = getPrefix(item.FromSchema, item.FromTable, item.FromTableAlias)
		joins = append(joins, utils.JoinTypeMap[item.JoinType] + ` ` + prefix + `.` + item.RelatedColumn + ` = ` + item.RelTableColumn)
	}

	dbColumn, titleColumn, selects = getColumnItem(item.Column, fromTableAlias)

	return
}

//getPrefix ----获取Join的前缀: "schema.tableName alias on alias" 或者 "schema.tableName on tableName"
func getPrefix(schema, tableName, alias string) (prefix, tableAlias string) {

	prefix = tableName
	tableAlias = tableName
	if schema != "" {
		prefix = schema + `.` + prefix
	}
	if alias != "" {
		prefix = prefix + ` ` + alias + ` on ` + alias
		tableAlias = alias
	} else {
		prefix += ` on ` + tableName
	}

	return
}

func getColumnItem(baseColumn, fromTableAlias string) (dbColumn, titleColumn, selects []string) {

	dbColumn, titleColumn, selects = make([]string, 0), make([]string, 0), make([]string, 0)
	for _, column := range strings.Split(strings.TrimSpace(baseColumn), ",") {
		columnArr := strings.Split(strings.TrimSpace(column), ":")
		if columnArr[0] == "" {
			continue
		}

		selectColumn := strings.TrimSpace(columnArr[0])
		alias := ""
		dbColumnItem :=  selectColumn
		titleColumnItem :=  selectColumn
		switch len(columnArr) {
		case 2:
			if val := strings.TrimSpace(columnArr[1]); val != "" {
				alias = val
				dbColumnItem = val
				titleColumnItem = val
			} else {
				dbColumnItem = selectColumn
				titleColumnItem = selectColumn
			}
		case 3:
			if val := strings.TrimSpace(columnArr[1]); val != "" {
				alias = val
				dbColumnItem = val
			} else {
				dbColumnItem = selectColumn
			}

			if val := strings.TrimSpace(columnArr[2]); val != "" {
				titleColumnItem = val
			} else {
				if v := strings.TrimSpace(columnArr[1]); v != "" {
					titleColumnItem = v
				} else {
					titleColumnItem = selectColumn
				}
			}
		}

		if alias != "" {
			selectColumn += ` AS ` + alias
		}
		if fromTableAlias != "" {
			selectColumn = fromTableAlias + `.` + selectColumn
		}

		dbColumn = append(dbColumn, dbColumnItem)
		titleColumn =append(titleColumn, titleColumnItem)
		selects = append(selects, selectColumn)
	}

	return dbColumn, titleColumn, selects
}