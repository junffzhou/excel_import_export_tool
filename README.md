
- [一、安装](#安装)
- [二、功能](#head1)
  - [ 一) excel导出支持的功能](#head2)
  - [ 二) excel导入支持的功能](#head3)
- [三、快速开始](#head4)
  - [一) excel导出](#head5)
    - [1 基础使用](#head6)
      - [参数说明](#head7)
    - [2 使用插件](#head8)
      - [1) 数据处理插件](#head9)
        - [1) 字段组合](#head10)
          - [参数说明](#head11)
          - [示例](#head12)
        - [2) 字段分离](#head13)
          - [参数说明](#head14)
          - [示例](#head15)
      - [2) 数据写入excel插件](#head16)
        - [参数说明](#head17)
        - [示例](#head18)
    - [3 示例](#head19)
      - [1 单表](#head20)
      - [2 字典项转换](#head21)
      - [3 一对多关系](#head22)
      - [4 多对多关系](#head23)
      - [5 包含字典项转换,一对多,多对多关系的完整示例](#head24)
  - [二) excel导入](#head25)
    - [1 基础使用](#head26)
      - [参数说明](#excel导入基础参数说明)
      - [示例](#head28)
    - [2 使用插件](#head29)
      - [1)  从excel读取数据插件](#head30)
        - [参数说明](#head31)
        - [示例](#head32)
      - [2)  数据处理插件](#head33)
        - [(1) 数据校验](#head34)
          - [参数说明](#head35)
          - [示例](#head36)
        - [(2) 关联字段的替换](#head37)
          - [参数说明](#head38)
          - [示例](#head39)
        - [(3) 字段组合](#head40)
          - [参数说明](#head41)
          - [示例](#head42)
        - [(4) 字段分离](#head43)
          - [参数说明](#head44)
          - [示例](#head45)
      - [3) 数据写入数据库插件](#head46)
        - [参数说明](#head47)
        - [示例](#head48)
    - [3 示例](#head49)
- [四、自定义插件](#head50)
  - [一) 导出插件](#head51)
    - [1) 数据处理插件](#head52)
    - [2) 数据写入excel插件](#head53)
  - [二) 导入插件](#head54)
    - [1) 读取excel数据插件](#head55)
    - [2) 数据处理插件](#head56)
    - [3) 数据写入数据库插件](#head57)

# 一、安装<a name="安装"></a>
```
go get github.com/junffzhou/excel_import_export_tool
```
#  二、功能<a name="head1"></a>

## 一) excel导出支持的功能<a name="head2"></a>

1. 自定义excel导出时的表头
2. 插件化:  数据处理插件, 数据写入excel表插件
3. 数据处理插件: 支持字段组合, 字段分离
4. 字段组合: 多个字段组合为新的字段
5. 字段分离: json串类型的字段支持平铺展开。支持的json串类型: string, []interface{}, map[string]interface{}, []map[string]interface{}四种类型。 string类型时,根据正则表达式分离为多个新的字段
6. 数据写入excel表插件: 查询的结果,经过数据处理插件的处理后, 数据按行或者按列插入excel表中
7. 支持自定义插件, 通用插件不满足需求时,用户可自定义插件



## 二) excel导入支持的功能<a name="head3"></a>

1. 根据操作列的动作(增/删/改)对数据库进行操作
2. 插件化:  从excel读取数据插件,  数据处理插件,  数据写入数据插件
3. 从excel读取数据插件: 按行或者按列读取excel表内容
4. 数据处理插件: 支持数据校验,  关联字段的替换, 字段组合, 字段分离
5. 数据校验: 对数据校验,  重复性校验(唯一性), 合法性校验(是否满足校验规则), 非空字段校验,非零字段校验
6. 关联字段的替换: 从其他表获取数据来替换当前字段的值, 表之间需要有关联的字段; 如: 导出时,字段的值是关联表的其他字段,数据库中存的是关联表的主键ID, 则需要把其他字段值转换为主键ID的值
7. 字段组合插件: 字段可以组合为 string, []interface{}, map[string]interface, []map[string]interface{}四种类型的json字符串类型的数据, 存入数据库的类型是 json字符串类型
8. 字段分离插件: 字段值根据正则表达式分离为多个新的字段
9. 支持自定义插件, 通用插件不满足需求时,用户可自定义插件



#  三、快速开始<a name="head4"></a>
**注意: 所有使用的数据来自[test_mysql_database.sql定义的数据](./test/test_mysql_database.sql)**
## 一) excel导出<a name="head5"></a>

### 1 基础使用<a name="head6"></a>

#####  参数说明<a name="head7"></a>

| 参数名                  | 类型                    | 是否必需 | 说明                                                         |
| :---------------------- | :---------------------- | :------- | ------------------------------------------------------------ |
| schema_name             | string                  | 否       | 库名                                                         |
| table_name              | string                  | 是       | 导出数据时查询的表                                           |
| table_alias             | string                  | 否       | 表别名                                                       |
| book_name               | string                  | 是       | 导出时excel的表名                                            |
| sheet_name              | string                  | 是       | excel的sheet页名                                             |
| show_type               | string                  | 否       | excel数据展示方式: "00":按行展示,默认; "01":按列展示         |
| primary_key             | string                  | 否       | 导出字段中包含主键列时,主键列的名字                          |
| is_include_action_field | string                  | 否       | 导出时是否包含操作列: "00":包含; "01":不包含,默认            |
| action_field_index      | float64                 | 否       | 包含操作列时,操作列所在excel表头的位置,只有is_include_action_field="00"时,该参数才有意义 |
| export_column           | []ExportColumn          | 是       | 导出字段详情                                                 |
| condition               | []map\[string\][]string | 否       | 执行sql语句时的 where条件                                    |
| order                   | string                  | 否       | 排序                                                         |
| limit                   | float64                 | 否       | 返回数据的数量                                               |
| offset                  | float64                 | 否       | 数据偏移量                                                   |



ExportColumn支持的参数:

| 参数名           | 类型   | 是否必须 | 说明                                                         |
| ---------------- | ------ | -------- | ------------------------------------------------------------ |
| column           | string | 否       | 数据表中的字段,别名,以及导出时的表头                         |
| from_schema      | string | 否       | 库名,from_table和table_name不是同一个库时,必需指定           |
| from_table       | string | 否       | 表名,column来自哪张表,默认来自table_name表                   |
| from_table_alias | string | 否       | from_table表的别名                                           |
| join_type        | string | 否       | column不是来自table_name时, join的方式:"00": left join,默认; "01": right join; "02": inner join |
| related_column   | string | 否       | from_table中关联的字段                                       |
| rel_table_column | string | 否       | join on关联的表和字段                                        |



其中:

1 ExportColumn:

column: 导出的字段,支持多个字段,字段间用英文逗号(",")隔开

每个字段最多包含三段,每段以英文冒号(":")分开

格式: 

```
select_column:column_alias:title_column
```

select_column: 查询时select的字段, 为空表示不选择字段

column_alias: 字段别名,可以为空,不设置别名

title_column: 导出时的表头



**注意:**

1) select_column 理论上不允许重复, 有重复请务必设置column_alias

2) 如果select_column都来自同一张表且不重复, 可以不用设置column_alias

3) 导出select_column来自多张表, 推荐设置column_alias, 因为表的字段可能相同, column_alias理论上不能重复,重复可能会影响 数据处理插件 的正确性

4) title_column理论上不允许重复,重复可能会影响 数据处理插件 的正确性

5)  title_column = "" 且 column_alias != "", title_column = column_alias;

title_column = "" 且 column_alias = "", title_column = select_column



以下情况表示不会选择任何字段

```
column = ""(值只包含空格也会被忽略)
column = ":"
column = "::"
column = ":id"
column = ":id:"
column = ":id:主表ID"
column = "::主表ID"
```



2 rel_table_column: from_table 在join on 时关联的表和字段, 如果关联表使用了别名, 必须要使用别名

join on 的构造为:  

```mysql
JOIN  from_table  (from_table_alias)  ON from_table|from_table_alias.related_column = rel_table_column
```



例如:

```json
{
  "table_name":"sub_info",
  "table_alias":"sub",
  "export_column":[
    {
      "column":"id:sub_id"
    },
    {
      "column":"id:rel_id",
      "join_type":"00",
      "from_schema":"test_mysql_database",
      "from_table":"sub_dept_rel",
      "from_table_alias":"rel",
      "related_column":"sub_id",
      "rel_table_column":"sub.id"
    },
    {
      "column":"id:dept_id:多对多关联的dept_info表的主键ID,dept_name:dept_name:部门名字,status:dept_status:部门状态",
      "join_type":"00",
      "from_schema":"test_mysql_database",
      "from_table":"dept_info",
      "from_table_alias":"dept",
      "related_column":"id",
      "rel_table_column":"rel.dept_id"
    }
  ]
}
```

转为:
注意, 当sub_info, sub_dept_rel, dept_info 都使用了别名, select时也会使用别名作为区分字段来源
```mysql

select 
sub.id as sub_id,
rel.id as rel_id,
dept.id as dept_id, 
dept.dept_name as dept_name, 
dept.status as dept_status 
from sub_info sub 
left join sub_dept_rel rel on rel.sub_id = sub.id
left join dept_info dept on dept.id = rel.dept_id
```


2 condition:  

格式: 

```json
[
  {
    "tableName_1, tableName_2":[
      "过滤条件1",
      "过滤条件2"
    ]
  }
]
```



最外层的数组表示where条件的顺序, 

每个 map的 key: 表名, 支持多张表, 多张表用英文逗号(",")隔开, 如果过滤字段是来自table_name, 表名可以为空;

value: []string类型,

格式:

```
 (过滤字段间的连接符)过滤字段(过滤字段和字段值间的符号):字段值
```



**注意:**

1) 表名来自table_name, table_alias, from_table, from_table_alias 四个字段的值

2) 如果表名为空, 自动使用table_name或table_alias;

3) 如果from_table_alias有值, 则必须使用from_table_alias, 不能使用from_table



过滤字段间的连接符有:

| 符号   | 含义                                 |
| ------ | ------------------------------------ |
| &      | 字段间是 AND,默认                    |
| &#124; | 字段间是 OR                          |
| ()     | ()包含的所有表的过滤字段作为一个组合 |



过滤字段和字段值间的符号有:

| 符号 | 含义                                                  |
| ---- | ----------------------------------------------------- |
| =    | 等于,默认, 当值为null时，转为 is null                 |
| !=   | 不等于, 当值为null时，转为 is not  null               |
| $    | like, 如 "comment $: %desc"表示 comment like "%desc"  |
| ~    | 正则                                                  |
| <>   | 在某个区间,  如 "id <>: [1,4]"表示 id BETWEEN 1 AND 4 |
| \>   | 大于                                                  |
| >=   | 大于等于                                              |
| <    | 小于                                                  |
| <=   | 小于等于                                              |
| {}   | in, 如"id {}:[1,2,3]"表示id in(1,2,3)                 |
| !{}  | not in, 如"id !{}:[1,2,3]"表示id not in(1,2,3)        |



例如:

```json
{
  "table_name":"sub_info",
  "table_alias":"sub",
  "condition":[
    {
      ", dept, main, rel":["&(status=:888)", "&(status!=:999)"]
    },
    {
      "rel":["&(sub_id !=:null", "|id >:0)"]
    },
    {
      "dict_info_1":["&dict_code =:性别"]
    },
    {
      "dict_info_2":["&dict_code =:字典项测试"]
    }
  ]
}
```

转为: 

```mysql
 where (sub.status = 888 and dept.status = 888 and main.status = 888 and rel.status = 888) and (sub.status != 999 and dept.status != 999 and main.status = 999 and rel.status != 999) and (rel.sub_id is not null or rel.id > 0) and dict_info_1 = "性别" and dict_info_2.dict_code = "字典项测试"
```



3 order: 支持多个字段排序,用英文逗号(",")隔开

+:升序排列,默认

-:降序排列

例如:

```json
{"order": "sub.id+, main.id-, main.name+, dept.id, dept.dept_name-"}
```

转为: 

```mysql
order by sub.id ASC, main.id DESC, main.name ASC, dept.id, dept.dept_name DESC
```



### 2 使用插件<a name="head8"></a>

**插件的优先级:**

1 数据处理插件(字段组合, 字段分离) > 数据写入excel插件

2 数据处理插件支持多个, 先注册先执行



#### 1) 数据处理插件<a name="head9"></a>

数据处理插件可以包含多个插件: 字段组合, 字段分离

##### 1) 字段组合<a name="head10"></a>

如果要组合字段中索引位置最小的元素没有被删除, 则组合的新字段会位于要组合字段中索引位置最小的元素后面, 否则, 组合的新字段位于要组合字段中索引位置最小的元素处



######  参数说明<a name="head11"></a>

| 参数名 | 类型              | 是否必需 | 说明                                                         |
| ------ | ----------------- | -------- | ------------------------------------------------------------ |
| field  | map[string]string | 是       | key:组合的字段名,作为表头; valu:要组合的字段,支持多个字段组合,多个字段间用英文逗号隔开(",") |

其中:

```json
{
  "field":{
    "sub_id_name_age":"sub_id:01, sub_name, sub_age:, @sep:-",
    "sub_age_name":"sub_age:01, sub_name:01, @sep::"
  }
}
```

表示: 

"sub_id_name_age":  组合后的新字段名

"sub_id:01, sub_name, sub_age:":  要组合字段为sub_id, sub_name, sub_age, 组合字段来自export_column中的column的第一段(字段名)或者第二段(字段别名,有别名一定使用别名), sub_id组合后不删除,sub_name,sub_age组合后从表头删除。"00":删除,默认; "01":不删除

"@sep::":  多个组合字段之间的分隔符, 表示分隔符为 ":"



######  示例<a name="head12"></a>

```go
client := NewExportClient()
fieldCombinationExt, e := export_ext.NewFieldCombinationObj([]byte(`{
    "field":{
        "sub_id_name_age":"sub_id:01, sub_name, sub_age:, @sep:-",
        "sub_age_name":"sub_age:01, sub_name:01, @sep::"
    }
}`))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldCombinationExt)
```



##### 2) 字段分离<a name="head13"></a>

######  参数说明<a name="head14"></a>

| 参数名           | 类型   | 是否必需 | 说明                                                         |
| ---------------- | ------ | -------- | ------------------------------------------------------------ |
| field_separation | string | 是       | 要分离的字段                                                 |
| field_type       | string | 是       | 要分离的字段的类型, "00":普通字符串类型; "01":[]interface{}类型; "02":map[string]interface{}类型; "03":[]map[string]interface{}类型 |
| is_delete        | string | 否       | 是否删除要分离的字段, "00":表示删除; "01":不删除,默认删除    |
| sep_regexp       | string | 否       | field_type="00", 分离的规则,支持正则                         |
| is_single        | string | 否       | field_type="01",  is_single="00": 组合成一个字段,默认;is_single="01":每个元素对应一个字段 |
| sep              | string | 否       | field_type="01", is_single="00",元素之间的分隔符             |
| new_field        | string | 否       | 字段分离后的新字段名, 作为excel的表头                        |



其中:

1 field_type="00": new_field支持多个字段,多个字段间使用英文逗号(",")隔开。根据sep_regexp规则匹配, 当匹配值的个数大于new_field的字段个数,匹配值中多余的值被删除;否则,new_field中多余的字段在导出时,值为空

2 field_type="01",is_single="00": new_field如果传入多个,只会获取第一个,剩下的会被忽略, 使用sep作为字段间的分隔符

3 field_type="01",is_single="01":  new_field支持多个字段,多个字段间使用英文逗号(",")隔开。当[]interface{}的元素个数大于new_field的字段个数,匹配值中多余的值被删除;否则,new_field中多余的字段在导出时,值为空

4 field_type="02":
new_field的格式为:

```
"key_1: new_field_1,  key_2: new_filed_2, key_3: new_field_3"
```

 每个字段中又包含 key-value对, key: 数据库中map[string]interface{}对应的key,  value:分离后的新字段名。当map[string]interface{}的key的个数大于new_field的字段个数,map[string]interface{}中多余的值被删除;否则,new_field中多余的字段在导出时,值为空

5 field_type="03":
 new_field的格式为:

```
"[\"Key_1: SliceEleIsMap_field_1, Key_2:SliceEleIsMap_field_2\", \"Key_3: SliceEleIsMap_field_3, Key_4: SliceEleIsMap_field_4\",  \"Key_5: SliceEleIsMap_field_5, Key_6: SliceEleIsMap_field_6\"]"
```

数组类型的字符串, 数组的每一个元素对应 []map[string]interface{} 中一个map[string]interface{}, 当 []map[string]interface{} 中map的个数大于new_field的字段个数,[]map[string]interface{} 中多余的值被删除;否则new_field中多余的字段在导出时,值为空
每个元素包含key-value对, key: map[string]interface{}对应的key,  value:分离后的新字段名。当map[string]interface{}的key的个数大于元素中的字段个数,map[string]interface{}中多余的值被删除;否则,元素中多余的字段在导出时,值为空



######  示例<a name="head15"></a>

```go
client := NewExportClient()
fieldSeparationExt, e := export_ext.NewFieldSeparationObj([]byte(`[
    {
        "field_separation":"sub_field_1",
        "field_type":"00",
        "is_delete":"00",
        "new_field":"string_field_1, string_field_2, string_field_3",
        "sep_regexp":"(\\w{2,})_(\\w{2,})_(\\w{2,})_(.*)"
    },
    {
        "field_separation":"sub_field_2",
        "field_type":"01",
        "is_delete":"00",
        "is_single":"00",
        "sep":",",
        "new_field":"slice_field_is_single, slice_field_2, slice_field_3, slice_field_4"
    },
    {
        "field_separation":"sub_field_2",
        "field_type":"01",
        "is_delete":"00",
        "is_single":"01",
        "new_field":"slice_field_1, slice_field_2, slice_field_3, slice_field_4"
    },
    {
        "field_separation":"sub_field_3",
        "field_type":"02",
        "is_delete":"00",
        "new_field":"key_1: map_field_1, key_2: map_filed_2, key_3: map_field_3"
    },
    {
        "field_separation":"sub_field_4",
        "field_type":"03",
        "is_delete":"00",
        "new_field":"[\"Key_1: SliceEleIsMap_field_1, Key_2:SliceEleIsMap_field_2\", \"Key_3: SliceEleIsMap_field_3, Key_4: SliceEleIsMap_field_4\", \"Key_5: SliceEleIsMap_field_5, Key_6: SliceEleIsMap_field_6\"]"
    }
]`))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldSeparationExt)
```



#### 2) 数据写入excel插件<a name="head16"></a>

不设置列宽度和合并单元格时,插件可以不用注册,默认执行

#####  参数说明<a name="head17"></a>

| 参数名    | 类型               | 是否必需 | 说明                                                         |
| --------- | ------------------ | -------- | ------------------------------------------------------------ |
| col_width | map[string]float64 | 否       | 列宽度, key的值支持: A,B: 表示第A,B两列; D:表示D列; E...G:表示E到G列;H...:表示H列到最后一列 |
| merge     | []int              | 否       | 按行展示时:合并的是值相等的列, 同一列的上下两行合并; 按列展示时:合并的是值相等的行,同一行的前后两列合并。当前单元格前面单元格的值(包括当前单元格)与上一条数据的前几列单元格的值全部相等,则合并;否则不合并。值从1开始且必须大于等于1 |

#####  示例<a name="head18"></a>

```go
client := NewExportClient()
insertValueExt, e := export_ext.NewWriteExcelObj([]byte(`{
    "col_width":{
        "A,B":15,
        "C":20,
        "D...E":25,
        "F...":45
    },
    "merge":[1, 2, 3, 4]
}`))
if e != nil {
	return e
}
client.RegisterWriteExcelExt(insertValueExt)
```



### 3 示例<a name="head19"></a>

#### 1 单表<a name="head20"></a>

导出字段只涉及当前表中数据, 不需要其他的处理 ,最简单场景

```go

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s"+
	"&writeTimeout=60s",
	"root",
	"mysql_password",
	"localhost",
	"3306",
	"test_mysql_database")

exportColumn := fmt.Sprintf(`[
    {
        "column":"id::主键ID, dict_code::字典code, dict_item_name::字典项name, dict_item_value::字典项value, comment::描述, creator::创建人"
    }
]`)
condition := fmt.Sprintf(`[
    {
        "":["status=:888", "deleted_at =:null"]
    }
]`)

book_name := fmt.Sprintf(`字典项_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
params := fmt.Sprintf(`{
    "table_name":"dict_info",
    "book_name":"%s",
    "sheet_name":"Sheet1",
    "show_type":"00",
    "primary_key":"id",
    "is_include_action_field":"00",
    "action_field_index":2,
    "export_column":%s,
    "condition":%s,
    "order":"dict_code+, dict_item_value-",
    "limit":100,
    "offset":0
}`, book_name, exportColumn, condition)

client := NewExportClient()

dir, e := os.Getwd()
if e != nil {
	return e
}
sysType := runtime.GOOS
if sysType == "windows" {
    dir += "\\test\\"
} else if sysType == "linux" {
    dir += "/test/"
}

bookName, f, e  := client.ParseExportReq(dsn, []byte(params))
if e != nil {
	return e
}

e = f.SaveAs(dir+bookName)
if e != nil {
	return e
}
```

对应的SQL语句:

```mysql
SELECT id, dict_code, dict_item_name, dict_item_value, comment, creator 
FROM dict_info 
WHERE  status = "888" AND deleted_at is null
ORDER BY dict_code ASC, dict_item_value DESC 
LIMIT 100;
```



#### 2 字典项转换<a name="head21"></a>

导出字段需要转换, 本质上和一对多关系一样

```go

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s"+
	"&writeTimeout=60s", 
	"root", 
	"mysql_password", 
	"localhost", 
	"3306", 
	"test_mysql_database")

exportColumn := fmt.Sprintf(`[
    {
        "column":"id:sub_id:主键ID, name:sub_name:姓名, age:sub_age:年龄"
    },
    {
        "column":"dict_item_name:sub_sex_name:性别",
        "join_type":"00",
        "from_schema":"test_mysql_database",
        "from_table":"dict_info",
        "from_table_alias":"dict_info_1",
        "related_column":"dict_item_value",
        "rel_table_column":"sub.sex"
    },
    {
        "column":"dict_item_name:sub_dict_field_name:测试字典",
        "join_type":"00",
        "from_schema":"test_mysql_database",
        "from_table":"dict_info",
        "from_table_alias":"dict_info_2",
        "related_column":"dict_item_value",
        "rel_table_column":"sub.dict_field"
    },
    {
        "column":"creator:sub_creator:创建人"
    }
]`)
condition := fmt.Sprintf(`[
    {
        ", dict_info_1, dict_info_2":["status=: 888"]
    },
    {
        "dict_info_1":["&dict_code: 性别"]
    },
    {
        "dict_info_2":["&dict_code: 字典项测试"]
    }
]`)
book_name := fmt.Sprintf(`子表信息_字典项转换_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
params := fmt.Sprintf(`{
    "table_name":"sub_info",
    "table_alias":"sub",
    "book_name":"%s",
    "sheet_name":"Sheet1",
    "show_type":"00",
    "primary_key":"id",
    "is_include_action_field":"00",
    "action_field_index":2,
    "export_column":%s,
    "condition":%s,
    "order":"sub_id+, sub_sex_name-, sub_dict_field_name+",
    "limit":100,
    "offset":0
}`, book_name, exportColumn, condition)

client := NewExportClient()

dir, e := os.Getwd()
if e != nil {
	return e
}
sysType := runtime.GOOS
if sysType == "windows" {
	dir += "\\test\\"
} else if sysType == "linux" {
	dir += "/test/"
}

bookName, f, e  := client.ParseExportReq(dsn, []byte(params))
if e != nil {
	return e
}

e = f.SaveAs(dir+bookName)
if e != nil {
	return e
}
```

对应的SQL语句:

```mysql
SELECT sub.id AS sub_id, 
sub.name AS sub_name, 
sub.age AS sub_age, 
dict_info_1.dict_item_name AS sub_sex_name, 
dict_info_2.dict_item_name AS sub_dict_field_name, 
sub.creator AS sub_creator
FROM sub_info sub  
LEFT JOIN  test_mysql_database.dict_info dict_info_1 on dict_info_1.dict_item_value = sub.sex   
LEFT JOIN  test_mysql_database.dict_info dict_info_2 on dict_info_2.dict_item_value = sub.dict_field
WHERE   sub.status = "888" AND  dict_info_1.status = "888" AND  dict_info_2.status = "888" AND dict_info_1.dict_code = "性别" AND dict_info_2.dict_code = "字典项测试" 
ORDER BY sub_id ASC, sub_sex_name DESC, sub_dict_field_name ASC  
LIMIT 100;
```



#### 3 一对多关系<a name="head22"></a>

字段关联的是其他表的主键字段,导出时需要转化为其他字段

```go

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s"+
	"&writeTimeout=60s",
	"root",
	"mysql_password",
	"localhost",
	"3306",
	"test_mysql_database")

exportColumn := fmt.Sprintf(`[
    {"column": "id:sub_id:主键ID"},
    {
        "column": "id:main_id:关联的主表主键ID, name:main_name:主表名字, status:main_status:主表状态",
        "join_type": "00",
        "from_schema": "test_mysql_database",
        "from_table": "main_info",
        "from_table_alias": "main",
        "related_column": "id",
        "rel_table_column": "sub.main_id"
    },
    {"column": "name:sub_name:姓名, age:sub_age:年龄, creator:sub_creator:创建人"}
]`)
condition := fmt.Sprintf(`[
    {", main": ["status=:888"]}
]`)
    
book_name := fmt.Sprintf(`子表信息_一对多关系_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
params := fmt.Sprintf(`{
    "table_name":"sub_info",
    "table_alias": "sub",
    "book_name":"%s",
    "sheet_name":"Sheet1",
    "show_type":"00",
    "primary_key":"id",
    "is_include_action_field":"00",
    "action_field_index":2,
    "export_column": %s,
    "condition": %s,
    "order": "main_id+, main_name-",
    "limit":100,
    "offset":0
}`, book_name, exportColumn, condition)

client := NewExportClient()

dir, e := os.Getwd()
if e != nil {
	return e
}
sysType := runtime.GOOS
if sysType == "windows" {
	dir += "\\test\\"
} else if sysType == "linux" {
	dir += "/test/"
}

bookName, f, e  := client.ParseExportReq(dsn, []byte(params))
if e != nil {
	return e
}

e = f.SaveAs(dir+bookName)
if e != nil {
	return e
}
```



对应的SQL语句:

```mysql
SELECT 
sub.id AS sub_id, 
main.id AS main_id,
main.name AS main_name, 
main.status AS main_status, 
sub.name AS sub_name, 
sub.age AS sub_age, 
sub.creator AS sub_creator 
FROM sub_info sub  
LEFT JOIN  test_mysql_database.main_info main on main.id = sub.main_id  
WHERE   sub.status = "888" AND  main.status = "888" 
ORDER BY main_id ASC, main_name DESC  
LIMIT 100;
```



#### 4 多对多关系<a name="head23"></a>

导出字段来自其他多对多关系的表

```go

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s"+
	"&writeTimeout=60s",
	"root",
	"mysql_password",
	"localhost",
	"3306",
	"test_mysql_database")

exportColumn := fmt.Sprintf(`[
    {
        "column": "id:sub_id:主键ID"
    },
    {
        "column": "",
        "join_type": "00",
        "from_schema": "test_mysql_database",
        "from_table": "sub_dept_rel",
        "from_table_alias": "rel",
        "related_column": "sub_id",
        "rel_table_column": "sub.id"
    },
    {
        "column": "id:dpet_id:多对多关联的dept_info表的主键ID, dept_name:dept_name:部门名字, status:dept_status:部门状态",
        "join_type": "00",
        "from_schema": "test_mysql_database",
        "from_table": "dept_info",
        "from_table_alias": "dept",
        "related_column": "id",
       "rel_table_column": "rel.dept_id"
    },
    {
        "column": "name:sub_name:姓名, age:sub_age:年龄, creator:sub_creator:创建人"
    }
]`)

condition := fmt.Sprintf(`[
    {", rel, dept": ["status =:888"]},
    {"rel": ["(sub_id !=:null", "id <>:[1, 100])"]}
]`)

book_name := fmt.Sprintf(`子表信息_多对多关系_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
params := fmt.Sprintf(`{
    "table_name":"sub_info",
    "table_alias": "sub",
    "book_name":"%s",
    "sheet_name":"Sheet1",
    "show_type":"00",
    "primary_key":"id",
    "is_include_action_field":"00",
    "action_field_index":2,
    "export_column": %s,
    "condition": %s,
    "order": "dept_id+, sub_name-",
    "limit":100,
    "offset":0
}`, book_name, exportColumn, condition)

client := NewExportClient()

dir, e := os.Getwd()
if e != nil {
    return e
}
sysType := runtime.GOOS
if sysType == "windows" {
    dir += "\\test\\"
} else if sysType == "linux" {
    dir += "/test/"
}

bookName, f, e  := client.ParseExportReq(dsn, []byte(params))
if e != nil {
    return e
}

e = f.SaveAs(dir+bookName)
if e != nil {
    return e
}
```



对应的SQL语句:

```mysql
SELECT 
sub.id AS sub_id, 
dept.id AS dpet_id, 
dept.dept_name AS dept_name, 
dept.status AS dept_status, 
sub.name AS sub_name, 
sub.age AS sub_age, 
sub.creator AS sub_creator
FROM sub_info sub  
LEFT JOIN  test_mysql_database.sub_dept_rel rel on rel.sub_id = sub.id   
LEFT JOIN  test_mysql_database.dept_info dept on dept.id = rel.dept_id 
WHERE  sub.status ="888"
AND rel.status = "888"
AND dept.status = "888"
AND (rel.sub_id is not null
AND rel.id BETWEEN 1 AND 100)
ORDER BY dept_id ASC, sub_name DESC  
LIMIT 100;
```



#### 5 包含字典项转换,一对多,多对多关系的完整示例<a name="head24"></a>

```go
	
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s"+
	"&writeTimeout=60s",
	"root",
	"mysql_password",
	"localhost",
	"3306",
	"test_mysql_database")


exportColumn := fmt.Sprintf(`[
    {
        "column":"id:sub_id:主键ID"
    },
    {
        "column":"",
        "join_type":"00",
        "from_schema":"test_mysql_database",
        "from_table":"sub_dept_rel",
        "from_table_alias":"rel",
        "related_column":"sub_id",
        "rel_table_column":"sub.id"
    },
    {
        "column":"dept_name:dept_name:所属部门",
        "join_type":"00",
        "from_schema":"test_mysql_database",
        "from_table":"dept_info",
        "from_table_alias":"dept",
        "related_column":"id",
        "rel_table_column":"rel.dept_id"
    },
    {
        "column":"name:main_name:主表name",
        "join_type":"00",
        "from_schema":"test_mysql_database",
        "from_table":"main_info",
        "from_table_alias":"main",
        "related_column":"id",
        "rel_table_column":"sub.main_id"
    },
    {
        "column":"name:sub_name:姓名, age:sub_age:年龄"
    },
    {
        "column":"dict_item_name:sex_name:性别",
        "join_type":"00",
        "from_schema":"test_mysql_database",
        "from_table":"dict_info",
        "from_table_alias":"dict_info_1",
        "related_column":"dict_item_value",
        "rel_table_column":"sub.sex"
    },
    {
        "column":"dict_item_name:dict_field_name:测试字典_name",
        "join_type":"00",
        "from_schema":"test_mysql_database",
        "from_table":"dict_info",
        "from_table_alias":"dict_info_2",
        "related_column":"dict_item_value",
        "rel_table_column":"sub.dict_field"
    },
    {
        "column":"field_1:sub_field_1:字段1, field_2:sub_field_2:字段2, field_3:sub_field_3:字段3, field_4:sub_field_4:字段4, creator:sub_creator:创建人"
    }
]`)

condition := fmt.Sprintf(`[
    {",dept,main,rel": ["&(status=:888)"]},
    {"rel": ["&(sub_id !=:null", "&id <>:[1,100])"]},
    {"dict_info_1": ["&dict_code =:性别"]},
    {"dict_info_2": ["&dict_code =:字典项测试"]}
]`)


book_name := fmt.Sprintf(`子表信息_completeExample_%v.xlsx`, time.Now().Format(`2006-01-02_15-04-05`))
params := fmt.Sprintf(`{
   	"table_name": "sub_info",
	"table_alias": "sub",
   	"book_name": "子表信息_import_test.xlsx",
	"sheet_name": "Sheet1",
	"show_type": "00",
   	"primary_key": "",
   	"is_include_action_field": "00",
   	"action_field_index": 2,
	"export_column": %s,
	"condition": %s,
	"order": "sub_id+, sub_name-, dept_name-",
	"limit": 100,
   	"offset": 0
}`, book_name, exportColumn, condition)
	
client := NewExportClient()

//注册字段组合插件
fieldCombinationExt, e := export_ext.NewFieldCombinationObj([]byte(`{
    "field": {
            "sub_id_name_age": "sub_id:01, sub_name, sub_age:, @sep:-",
            "sub_name_age": "sub_name, sub_age, @sep::"
    }
}`))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldCombinationExt)

//注册字段分离插件
fieldSeparationExt, e := export_ext.NewFieldSeparationObj([]byte(`[
    {
        "field_separation":"sub_field_1",
        "field_type":"00",
        "is_delete":"00",
        "new_field":"string_field_1, string_field_2, string_field_3",
        "sep_regexp":"(\\w{2,})_(\\w{2,})_(\\w{2,})_(.*)"
    },
    {
        "field_separation":"sub_field_2",
        "field_type":"01",
        "is_delete":"00",
        "is_single":"00",
        "sep":",",
        "new_field":"slice_field_is_single, slice_field_2, slice_field_3, slice_field_4"
    },
    {
        "field_separation":"sub_field_2",
        "field_type":"01",
        "is_delete":"00",
        "is_single":"01",
        "new_field":"slice_field_1, slice_field_2, slice_field_3, slice_field_4"
    },
    {
        "field_separation":"sub_field_3",
        "field_type":"02",
        "is_delete":"00",
        "new_field":"key_1: map_field_1, key_2: map_filed_2, key_3: map_field_3"
    },
    {
        "field_separation":"sub_field_4",
        "field_type":"03",
        "is_delete":"00",
        "new_field":"[\"Key_1: SliceEleIsMap_field_1, Key_2:SliceEleIsMap_field_2\", \"Key_3: SliceEleIsMap_field_3, Key_4: SliceEleIsMap_field_4\", \"Key_5: SliceEleIsMap_field_5, Key_6: SliceEleIsMap_field_6\"]"
    }
]`))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldSeparationExt)

//注册数据写入excel插件
insertValueExt, e := export_ext.NewWriteExcelObj([]byte(`{
    "col_width":{
        "A,B,C":15,
        "D":20,
        "E":25,
        "F...":45
    },
    "merge":[1, 2, 3, 4]
}`))
if e != nil {
	return e
}
client.RegisterWriteExcelExt(insertValueExt)
	
dir, e := os.Getwd()
if e != nil {
	return e
}
sysType := runtime.GOOS
if sysType == "windows" {
    dir += "\\test\\"
} else if sysType == "linux" {
    dir += "/test/"
}
bookName, f, e  := client.ParseExportReq(dsn, []byte(params))
if e != nil {
    return e
}

e = f.SaveAs(dir+bookName)
if e != nil {
    return e
}
	
```


对应的SQL语句:

```mysql
SELECT 
sub.id AS sub_id, 
dept.dept_name AS dept_name, 
main.name AS main_name, 
sub.name AS sub_name, 
sub.age AS sub_age, 
dict_info_1.dict_item_name AS sex_name, 
dict_info_2.dict_item_name AS dict_field_name, 
sub.field_1 AS sub_field_1, 
sub.field_2 AS sub_field_2, 
sub.field_3 AS sub_field_3, 
sub.field_4 AS sub_field_4, 
sub.creator AS sub_creator 
FROM sub_info sub  
LEFT JOIN  test_mysql_database.sub_dept_rel rel on rel.sub_id = sub.id   
LEFT JOIN  test_mysql_database.dept_info dept on dept.id = rel.dept_id  
LEFT JOIN  test_mysql_database.main_info main on main.id = sub.main_id   
LEFT JOIN  test_mysql_database.dict_info dict_info_1 on dict_info_1.dict_item_value = sub.sex  
LEFT JOIN  test_mysql_database.dict_info dict_info_2 on dict_info_2.dict_item_value = sub.dict_field  
WHERE  ( sub.status = "888" 
AND dept.status = "888"
AND main.status = "888" 
AND rel.status = "888" ) 
AND ( rel.sub_id is not null  
AND  rel.id  BETWEEN 1 AND 100 ) 
AND  dict_info_1.dict_code = "性别"  
AND  dict_info_2.dict_code = "字典项测试"   
ORDER BY sub_id ASC, sub_name DESC, dept_name DESC  
LIMIT 100;
```



## 二) excel导入<a name="head25"></a>

### 1 基础使用<a name="head26"></a>

##### 参数说明<a name="excel导入基础参数说明"></a>

| 参数名称                              | 类型                 | 是否必需 | 说明                                                         |
| ------------------------------------- | -------------------- | -------- | ------------------------------------------------------------ |
| schema_name                           | string               | 否       | 库名                                                         |
| table_name<a name="导入操作的表"></a> | string               | 是       | 表名                                                         |
| sheet_name                            | string               | 是       | excel的sheet页名                                             |
| show_type                             | string               | 否       | excel数据展示方式: "00":按行展示,默认; "01":按列展示,按行或按列读取内容 |
| primary_key                           | string               | 否       | excel表头字段中包含主键列时,主键列的名字                     |
| action_field_index                    | float64              | 否       | excel表头包含操作列时,操作列所在excel表头的位置,从1开始      |
| auto_increment                        | string               | 否       | 包含主键字段时,主键是否自增, "00":自增, 默认, "01":非自增    |
| import_column                         | map[string]string    | 是       | excel表头和数据库字段的映射, key: excel表的表头, value: 表头对应数据表中的字段。如果表头字段是从其他表获取的, 可以不用在import_column中声明 |
| additional_condition                  | *AdditionalCondition | 否       | 对每种导入操作(增/删/改)的附加条件                           |



AdditionalCondition支持的参数

| 参数名               | 类型              | 是否必须 | 说明                     |
| -------------------- | ----------------- | -------- | ------------------------ |
| create_add_condition | *BaseAddCondition | 否       | 导入时新增操作的附加条件 |
| update_add_condition | *BaseAddCondition | 否       | 导入时修改操作的附加条件 |
| delete_add_condition | *BaseAddCondition | 否       | 导入时删除操作的附加条件 |

BaseAddCondition支持的参数

| 参数名         | 类型                   | 是否必须 | 说明                                                         |
| -------------- | ---------------------- | -------- | ------------------------------------------------------------ |
| extra_column   | map[string]interface{} | 否       | 额外的字段，insert时作为要插入的字段，update,delete时作为set后面要更新的字段 |
| is_soft_delete | string                 | 否       | 删除操作是否是软删除: "00":软删除,默认; "01":非软删除,直接delete数据。软删除指的是数据库某些字段(可以有多个)是删除标识,执行删除时并不会从数据库删除该条数据,只是将删除标识的值置为删除状态的值,如: status="999" and deleted_at is not null 表示数据是删除状态 |
| condition      | []string               | 否       | 执行update/delete SQL语句时附加的WHERE条件                   |

其中:

1 condition<a name="additional_condition的condition条件"></a>

每个元素(即过滤条件) 的格式:

```
(过滤字段间的连接符)过滤字段(过滤字段和字段值间的符号):字段值
```



过滤字段间的连接符有:

| 符号   | 含义                                 |
| ------ | ------------------------------------ |
| &      | 字段间是 AND,默认                    |
| &#124; | 字段间是 OR                          |
| ()     | ()包含的所有表的过滤字段作为一个组合 |



过滤字段和字段值间的符号有:

| 符号 | 含义                                                  |
| ---- | ----------------------------------------------------- |
| =    | 等于,默认, 当值为null时，转为 is null                 |
| !=   | 不等于, 当值为null时，转为 is not  null               |
| $    | like, 如 "comment $: %desc"表示 comment like "%desc"  |
| ~    | 正则                                                  |
| <>   | 在某个区间,  如 "id <>: [1,4]"表示 id BETWEEN 1 AND 4 |
| \>    | 大于                                                  |
| >=   | 大于等于                                              |
| <    | 小于                                                  |
| <=   | 小于等于                                              |
| {}   | in, 如"id {}:[1,2,3]"表示id in(1,2,3)                 |
| !{}  | not in, 如"id !{}:[1,2,3]"表示id not in(1,2,3)        |



#####  示例<a name="head28"></a>

```go

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s"+
	"&writeTimeout=60s", 
	"root", 
	"mysql_password", 
	"localhost", 
	"3306", 
	"test_mysql_database")

params := fmt.Sprintf(`{
    "table_name":"sub_info",
    "sheet_name":"Sheet1",
    "show_type":"00",
    "primary_key":"id",
    "auto_increment":"00",
    "action_field_index":2,
    "import_column":{
        "sub_id":"id",
        "性别":"sex",
        "测试字典_name":"dict_field",
        "创建人":"creator"
    },
    "additional_condition":{
        "create_add_condition":{
            "extra_column":{
                "status":"888",
                "creator":"admin",
                "modifier":"admin"
            }
        },
        "update_add_condition":{
            "extra_column":{
                "modifier":"admin"
            },
            "condition":[
                "status !=:999"
            ]
        },
        "delete_add_condition":{
            "extra_column":{
                "modifier":"admin",
                "deleted_at":"%s",
                "status":"999"
            },
            "is_soft_delete":"00",
            "condition":[
                "status !=:999"
            ]
        }
    }
}`, time.Now().Format("2006-01-02 15:04:05"))

dir, e := os.Getwd()
if e != nil {
    return e
}
filePath := dir
sysType := runtime.GOOS
if sysType == "windows" {
    filePath = dir + "\\test\\子表信息_import_test.xlsx"
} else if sysType == "linux" {
    filePath = dir + "/test/子表信息_import_test.xlsx"
}

client := NewImportClient()
e := client.ParseImportReq(dsn, filePath, []byte(params))
if e != nil {
	return e
}
```



### 2 使用插件<a name="head29"></a>

**插件的优先级:**

1 从excel读取数据插件 > 数据处理插件(数据校验, 关联字段的替换, 字段组合, 字段分离) > 数据写入数据库插件

2 数据处理插件支持多个, 先注册先执行

#### 1)  从excel读取数据插件<a name="head30"></a>

可以不用注册该插件,该插件默认执行

#####  参数说明<a name="head31"></a>

无参数,使用的参数是在<a href=#excel导入基础参数说明>导入基础参数</a>



#####  示例<a name="head32"></a>

```go
client := NewImportClient()
client.RegisterReadSheetExt(import_ext.NewReadSheetObj())
```



#### 2)  数据处理插件<a name="head33"></a>

数据处理插件可以包含多个插件: 数据校验, 关联字段的替换, 字段组合, 字段分离

##### (1) 数据校验<a name="head34"></a>

######  参数说明<a name="head35"></a>

| 参数名                    | 类型                  | 是否必须 | 说明                                      |
| ------------------------- | --------------------- | -------- | ----------------------------------------- |
| repeatability_check_field | []*RepeatabilityCheck | 否       | 重复性校验,唯一性校验                     |
| legality_check_field      | []*LegalityCheck      | 否       | 合法性校验: 是否满足校验规则              |
| non_empty_check_field     | string                | 否       | 非空校验多个以英文逗号(",")隔开           |
| non_zero_check_field      | string                | 否       | 非零校验,支持多个字段,以英文逗号(",")隔开 |



RepeatabilityCheck支持的参数

| 参数名    | 类型     | 是否必须 | 说明                                                  |
| --------- | -------- | -------- | ----------------------------------------------------- |
| field     | string   | 是       | 要校验的字段,支持多个字段联合唯一,以英文逗号(",")隔开 |
| condition | []string | 否       | 校验时查询的条件,支持多个过滤字段,以英文逗号(",")隔开 |
| order     | string   | 否       | 排序                                                  |
| limit     | float64  | 否       | 返回数据的数量                                        |
| offset    | float64  | 否       | 数据偏移量                                            |



LegalityCheck

| 参数名 | 类型   | 是否必须 | 说明                                                        |
| ------ | ------ | -------- | ----------------------------------------------------------- |
| field  | string | 是       | 要校验的字段, 支持多个字段同一个校验规则以英文逗号(",")隔开 |
| reg    | string | 是       | 校验规则,支持正则表达式                                     |

其中:

1 condition: 和<a href="#additional_condition的condition条件">additional_condition的condition条件</a>一样

2 order<a name="数据校验的order条件"></a>: 支持多个字段排序,用英文逗号(",")隔开

+:升序排列,默认

-:降序排列

```json
{"order": "id+, name-"}

```

转为: 

```mysql
order by id ASC, name DESC
```



######  示例<a name="head36"></a>

```go

client := NewImportClient()
dataCheckParam := fmt.Sprintf(`{
    "repeatability_check_field":[
        {
            "field":"sex, name",
            "condition":[
                "status !=: 999"
            ]
        }
    ],
    "legality_check_field":[
        {
            "field":"dict_field",
            "reg":"(\\w+)"
        }
    ],
    "non_empty_check_field":"sex",
    "non_zero_check_field":"age"
}`)

dataCheckExt, e := import_ext.NewDataCheckObj([]byte(dataCheckParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(dataCheckExt)
```



##### (2) 关联字段的替换<a name="head37"></a>

[基础参数的table_name表](#导入操作的表)和from_table表之间有一个关联字段(可以是主键ID,即一对一,一对多关系; 也可以是其他字段,如字典项的关联)，导出时使用from_table表的其他字段,导入时需要替换成关联字段

######  参数说明<a name="head38"></a>

| 参数名称       | 类型             | 是否必需 | 说明 |
| -------------- | ---------------- | -------- | ---- |
| related_column | []*RelatedColumn | 是       |      |



RelatedColumn支持的参数

| 参数名称        | 类型     | 是否必需 | 说明                                          |
| --------------- | -------- | -------- | --------------------------------------------- |
| title_to_column | string   | 是       | excel表头的和其他表关联的字段                 |
| from_schema     | string   | 否       | 库名, 字段来源表和当前表不在同一库时,必须传入 |
| from_table      | string   | 是       | 表名, 字段来源表                              |
| condition       | []string | 否       | 执行sql语句时的 where条件                     |
| order           | string   | 否       | 排序                                          |
| limit           | float64  | 否       | 返回数据的数量                                |
| offset          | float64  | 否       | 数据偏移量                                    |

其中:

1 title_to_column:  excel表头的和其他表关联的字段, 可以包含多个, 多个以英文逗号(",")隔开

每个字段最多包含三段, 每段以英文冒号(":")分开

格式:

```
excel_title_column:db_column:rel_column
```

多个字段表示: 在导出时excel表头包含多个from_table的字段, 只有当所有的excel_title_column的值(从excel读取)在from_table匹配到(查询数据库)时, 才会把db_column的值赋给rel_column

第一段: excel表头

第二段: excel表头的值在from_table中对应的数据库字段

第三段:  excel表头在[基础参数的table_name表中](#导入操作的表)有关联字段时关联字段的名字, 为空表示[基础参数的table_name表中](#导入操作的表)无关联字段,导入时不包含当前字段

**注意:**

1) excel_title_column不存在时, db_column,rel_column必须存在, 否则忽略该字段

2) excel_title_column存在, db_column不存在时, db_column=excel_title_column 

3) 如果想要关联字段, 必须传入rel_column

4) 以下情况被忽略:

```
title_to_column = ""(值只包含空格也会被忽略)
title_to_column = ":"
title_to_column = "::"
title_to_column = ":id"
title_to_column = "::main_id"
title_to_column = ":id:"

```

5) 以下情况是允许的:

```
1 title_to_column = "主表ID"
2 title_to_column = "主表ID:"
3 title_to_column = "主表ID::"

4 title_to_column = "主表ID:id"
5 title_to_column = "主表ID:id:"  

6 title_to_column = ":id:main_id"
7 title_to_column = "主表ID::main_id"
8 title_to_column = "主表ID:id:main_id"

其中
1,2,3等效: excel表头中主表ID对应from_table的主表ID字段, 基础参数的table_name表不包含主表ID字段
4,5等效: excel表头中主表ID字段对应from_table的id字段, 基础参数的table_name表不包含主表ID字段
6: from_table的id字段在excel表头中不存在, 基础参数的table_name表包含该字段, 且该字段为main_id
7: excel表头中主表ID对应from_table的主表ID字段, 基础参数的table_name表包含该字段, 且该字段为main_id
8: excel表头中主表ID字段对应from_table的id字段, 基础参数的table_name表包含该字段, 且该字段为main_id


```



2 condition: 和<a href="#additional_condition的condition条件">additional_condition的condition条件</a>一样

3 order: 和<a href="#数据校验的order条件">数据校验的order条件</a>一样



######  示例<a name="head39"></a>

```go

/*
	related_column的
	第一个元素是关联的是main_info表的主键,导出时使用了main_info表的name字段值,导入时从main_info表根据	name字段的值获取id字段的值, 在写入数据时 id 对应基础参数的table_name表的main_id字段
	第二,第三个元素是字典项, 导出时使用了dict_info表的dict_item_name字段值,导入时从dict_info表根据	dict_item_name字段的值获取dict_item_value字段的值, 在写入数据时 dict_item_value对应基础参数的	table_name的sex, dict_field字段
*/
client := NewImportClient()
relColumnParam = fmt.Sprintf(`{
    "related_column":[
        {
            "title_to_column":"主表name:name:main_name, :id:main_id",
            "from_schema":"test_mysql_database",
            "from_table":"main_info",
            "condition":[
                "status!=:999"
            ],
            "order":"id",
            "limit":0,
            "offset":0
        },
        {
            "title_to_column":"性别:dict_item_name, :dict_item_value:sex",
            "from_schema":"test_mysql_database",
            "from_table":"dict_info",
            "condition":[
                "dict_code:性别",
                "&status!=:999"
            ],
            "order":"dict_item_value",
            "limit":0,
            "offset":0
        },
        {
            "title_to_column":"测试字典_name:dict_item_name, :dict_item_value:dict_field",
            "from_schema":"test_mysql_database",
            "from_table":"dict_info",
            "condition":[
                "dict_code:字典项测试",
                "&status!=:999"
            ],
            "order":"dict_item_value",
            "limit":0,
            "offset":0
        }
    ]
}`)

relColumnExt, e := import_ext.NewRelatedColumnObj([]byte(relColumnParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(relColumnExt)
```



##### (3) 字段组合<a name="head40"></a>

######  参数说明<a name="head41"></a>

| 参数名称          | 类型                | 是否必需 | 说明         |
| ----------------- | ------------------- | -------- | ------------ |
| field_combination | []*FieldCombination | 是       | 字段组合详情 |

FieldCombination支持的参数

| 参数名称   | 类型   | 是否必需 | 说明                                                         |
| ---------- | ------ | -------- | ------------------------------------------------------------ |
| field      | string | 是       | 要组合的字段                                                 |
| new_field  | string | 是       | 组合后的新字段名, 作为数据库字段存入数据库中                 |
| field_type | string | 否       | 组合的新字段的类型: "00":string类型,默认;"01":[]interface{}类型;"02" :map[string]interface{}类型;"03":[]map[string]interface{}类型 |
| separator  | string |  否        | field_type="00"时,字段间的分隔符                             |



其中:

1 field_type="00": field支持多个字段,多个字段间使用英文逗号(",")隔开, 多个字段组合后字段间的分隔符为separator。当field的字段的字段不在excel表头中,对应字段的值为空

2 field_type="01": field支持多个字段,多个字段间使用英文逗号(",")隔开。当field的字段不在excel表头中,对应字段的值为空

3 field_type="02":  

field的格式:

```
 "excel_title_key_1: db_key_1,  excel_title_key_2:db_key_2, excel_title_key_3: db_key_3"
```

 每个字段中又包含 key-value对,  key: excel表头, value:数据库中map[string]interface{}对应的key。当key不在excel表头,value的值为空



4 field_type="03":  

field的格式:

```
"[\"excel_title_Key_1: SliceEleIsMap_key_1, excel_title_Key_2:SliceEleIsMap_key_2\", \"excel_title_Key_3: SliceEleIsMap_key_3,excel_title_Key_4: SliceEleIsMap_key_4\", \"excel_title_Key_5: SliceEleIsMap_key_5, excel_title_Key_6: SliceEleIsMap_key_6\"]"
```

数组类型的字符串, 数组的每一个元素对应 []map[string]interface{} 中一个map[string]interface{}, 

每个map[string]interface{}中又包含 key-value对,  key: excel表头, value:数据库中map[string]interface{}对应的key。当key不在excel表头,value的值为空



######  示例<a name="head42"></a>

```go

client := NewImportClient()
fieldCombineParam = fmt.Sprintf(`{
    "field_combination":[
        {
            "field":"string_field_1, string_field_2, string_field_3",
            "new_field":"field_1",
            "field_type":"00",
            "separator":"_"
        },
        {
            "field":"slice_field_1, slice_field_2, slice_field_3, slice_field_4",
            "new_field":"field_2",
            "field_type":"01"
        },
        {
            "field":"map_field_1: key_1, map_filed_2:key_2, map_field_3: key_3",
            "new_field":"field_3",
            "field_type":"02"
        },
        {
            "field":"[\"SliceEleIsMap_field_1: Key_1, SliceEleIsMap_field_2: Key_2\",\"SliceEleIsMap_field_3: Key_3, SliceEleIsMap_field_4: Key_4\",\"SliceEleIsMap_field_5: Key_5, SliceEleIsMap_field_6: Key_6\"]",
            "new_field":"field_4",
            "field_type":"03"
        }
    ]
}`)
fieldCombineExt, e := import_ext.NewFieldCombinationObj([]byte(fieldCombineParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldCombineExt)
```



##### (4) 字段分离<a name="head43"></a>

######  参数说明<a name="head44"></a>

| 参数名称         | 类型               | 是否必需 | 说明         |
| ---------------- | ------------------ | -------- | ------------ |
| field_separation | []*FieldSeparation | 是       | 分离字段详情 |

FieldSeparation支持的参数

| 参数名称          | 类型   | 是否必需 | 说明                                                         |
| ----------------- | ------ | -------- | ------------------------------------------------------------ |
| field             | string | 是       | 要分离的字段                                                 |
| new_field         | string | 是       | 分离后的新字段, 支持多个,多个以英文逗号隔开(","), 作为数据库字段存入数据库中 |
| separation_regexp | string | 是       | 分离匹配的规则,支持正则表达式                                |



其中:

separation_regexp: 根据正则表达式获取field的值。如果正则表达式结果的长度 < new_field的个数, 则 new_field剩余元素的值为空,否则, 正则表达式的结果剩余元素被忽略



######  示例<a name="head45"></a>

```go

client := NewImportClient()
fieldSeparatorParam := fmt.Sprintf(`{
    "field_separation":[
        {
            "field":"sub_name_age",
            "new_field":"name, age",
            "separation_regexp":"(.*):(.*)"
        }
    ]
}`)
fieldSeparatorExt, e := import_ext.NewFieldSeparationObj([]byte(fieldSeparatorParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldSeparatorExt)
```

 

#### 3) 数据写入数据库插件<a name="head46"></a>

可以不用注册该插件,该插件默认执行

  

以下情况的数据被忽略:

1) 操作字段的值不是增/删/改

2) 包含操作字段, 不包含主键字段且主键字段不是自增, 任意操作都被忽略

3) 包含操作字段, 不包含主键字段, 主键字段自增, 修改和删除操作被忽略

4) 主键字段值为空, 修改和删除操作被忽略



#####  参数说明<a name="head47"></a>

无参数



#####  示例<a name="head48"></a>

```go
client := NewImportClient()
client.RegisterHandlerDataExt(import_ext.NewWriteDatabaseObj())
```



### 3 示例<a name="head49"></a>

```go

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=15s&readTimeout=15s"+
	"&writeTimeout=60s",
	"root",
	"mysql_pass",
	"localhost",
	"3306",
	"test_mysql_database")

dir, e := os.Getwd()
if e != nil {
	return e
}
filePath := dir
sysType := runtime.GOOS
if sysType == "windows"{
	filePath = dir + "\\test\\子表信息_import_test.xlsx"
} else if sysType == "linux" {
	filePath = dir + "/test/子表信息_import_test.xlsx"
}

params := fmt.Sprintf(`{
    "table_name":"sub_info",
    "sheet_name":"Sheet1",
    "show_type":"00",
    "primary_key":"id",
    "auto_increment":"00",
    "action_field_index":2,
    "import_column":{
        "主键":"id",
        "创建人":"creator"
    },
    "additional_condition":{
        "create_add_condition":{
            "extra_column":{
                "status":"888",
                "creator":"admin",
                "modifier":"admin"
            }
        },
        "update_add_condition":{
            "extra_column":{
                "modifier":"admin"
            },
            "condition":[
                "status !=:999"
            ]
        },
        "delete_add_condition":{
            "extra_column":{
                "modifier":"admin",
                "deleted_at":"%s",
                "status":"999"
            },
            "is_soft_delete":"00",
            "condition":[
                "status !=:999"
            ]
        }
    }
}`, time.Now().Format("2006-01-02 15:04:05"))


client := NewImportClient()

//注册读取excel数据插件
client.RegisterReadSheetExt(import_ext.NewReadSheetObj())

//注册数据处理插件: 关联字段替换
relColumnParam := fmt.Sprintf(`{
    "related_column":[
        {
            "title_to_column":"主表name:name, :id:main_id",
            "from_schema":"test_mysql_database",
            "from_table":"main_info",
            "condition":[
                "status !=: 999"
            ],
            "order":"id",
            "limit":0,
            "offset":0
        },
        {
            "title_to_column":"性别:dict_item_name, :dict_item_value:sex",
            "from_schema":"test_mysql_database",
            "from_table":"dict_info",
            "condition":[
                "dict_code:性别",
                "&status!=:999"
            ],
            "order":"dict_item_value",
            "limit":0,
            "offset":0
        },
        {
            "title_to_column":"测试字典_name:dict_item_name, :dict_item_value:dict_field",
            "from_schema":"test_mysql_database",
            "from_table":"dict_info",
            "condition":[
                "dict_code:字典项测试",
                "&status!=:999"
            ],
            "order":"dict_item_value",
            "limit":0,
            "offset":0
        }
    ]
}`)
relColumnExt, e := import_ext.NewRelatedColumnObj([]byte(relColumnParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(relColumnExt)

//注册数据处理插件: 字段组合
fieldCombineParam := fmt.Sprintf(`{
    "field_combination":[
        {
            "field":"string_field_1, string_field_2, string_field_3",
            "new_field":"field_1",
            "field_type":"00",
            "separator":"_"
        },
        {
            "field":"slice_field_1, slice_field_2, slice_field_3, slice_field_4",
            "new_field":"field_2",
            "field_type":"01"
        },
        {
            "field":"map_field_1: key_1, map_filed_2:key_2, map_field_3: key_3",
            "new_field":"field_3",
            "field_type":"02"
        },
        {
            "field":"[\"SliceEleIsMap_field_1: Key_1, SliceEleIsMap_field_2: Key_2\",\"SliceEleIsMap_field_3: Key_3, SliceEleIsMap_field_4: Key_4\",\"SliceEleIsMap_field_5: Key_5, SliceEleIsMap_field_6: Key_6\"]",
            "new_field":"field_4",
            "field_type":"03"
        }
    ]
}`)
fieldCombineExt, e := import_ext.NewFieldCombinationObj([]byte(fieldCombineParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldCombineExt)

//注册数据处理插件: 字段分离
fieldSeparatorParam := fmt.Sprintf(`{
    "field_separation":[
        {
            "field":"name_age",
            "new_field":"name, age",
            "separation_regexp":"(.*):(.*)"
        }
    ]
}`)
fieldSeparatorExt, e := import_ext.NewFieldSeparationObj([]byte(fieldSeparatorParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(fieldSeparatorExt)


//注册数据处理插件: 数据校验
dataCheckParam := fmt.Sprintf(`{
    "repeatability_check_field":[
        {
            "field":"sex, name, age",
            "condition":[
                "status !=: 999"
            ]
        }
    ],
    "legality_check_field":[
        {
            "field":"dict_field",
            "reg":"(\\w+)"
        }
    ],
    "non_empty_check_field":"sex",
    "non_zero_check_field":"age"
}`)
dataCheckExt, e := import_ext.NewDataCheckObj([]byte(dataCheckParam))
if e != nil {
	return e
}
client.RegisterHandlerDataExt(dataCheckExt)
 

//注册数据写入数据库插件
client.RegisterWriteDataBaseExt(import_ext.NewWriteDatabaseObj())

e := client.ParseImportReq(dsn, filePath, []byte(params))
if e != nil {
	return e
}
```



# 四、自定义插件<a name="head50"></a>

## 一) 导出插件<a name="head51"></a>

### 1) 数据处理插件<a name="head52"></a>

实现了通用的数据处理方法: 字段组合, 字段分离,若不满足需求,用户可以自定义实现该插件,

需要实现以下方法:

```go
//DataHandler ----It provides methods to process database data.
type DataHandler interface {
    HandlerData(dbData []map[string]string, dbClient *database.Client) (columnReplaces []*export_model.ColumnReplace, e error)
}
```



### 2) 数据写入excel插件<a name="head53"></a>

实现通用的将数据写入excel方法, 若不满足需求,用户可以自定义实现该插件

需要实现以下方法:

```go
//WriteExcel ----It provides methods to write database data to the excel.
type WriteExcel interface {
    WriteExcel(sheetName, showType string, exportData *export_model.ExportData) (*excelize.File, error)
}
```



## 二) 导入插件<a name="head54"></a>

### 1) 读取excel数据插件<a name="head55"></a>

实现了通用的读取excel数据,若不满足需求,用户可以自定义实现该插件

需要实现以下方法:

```go
//SheetReader ----It provides methods to read excel content
type SheetReader interface {
    ReadSheet(filePath string, titleColumnToDbColumn map[string]string, baseData *import_model.ImportRequestParam) ([]*import_model.ImportExcelData, error)
}
```



### 2) 数据处理插件<a name="head56"></a>

实现了通用的数据处理: 数据校验, 关联字段的替换, 字段组合, 字段分离,若不满足需求,用户可以自定义实现该插件,

需要实现以下方法:

```go
//DataHandler ----It provides methods to process excel data.
type DataHandler interface {
    HandlerData(titleColumnToDbColumn map[string]string, dbClient *database.Client,  baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) error
}
```

### 3) 数据写入数据库插件<a name="head57"></a>

实现通用的根据操作列将数据写入数据库操作,若不满足需求,用户可以自定义实现该插件

需要实现以下方法:

```go
//DataWriter ----It provides methods to write excel data to the database.
type DataWriter interface {
    DataWriter(titleColumnToDbColumn map[string]string, dbClient *database.Client, baseData *import_model.ImportRequestParam, xlsxData []*import_model.ImportExcelData) error
}
```

