package export_plugin

import (
	"encoding/json"
	"errors"
	"github.com/junffzhou/excel_import_export_tool/model/export_model"
	"github.com/junffzhou/excel_import_export_tool/utils"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

/*
WriteExcelObj ----Write Excel
ColWidth: map[string]float64, The value of key supports:
		A,B,C: represents two columns A and B and C;
		D: represents D column;
		E...G: represents E to G column;
		H...: represents H column to the last column.
Merge: When rows are displayed, the same columns are merged; When columns are displayed, the same rows are merged
*/
type WriteExcelObj struct {
	ColWidth  map[string]float64 `json:"col_width"`
	Merge []int `json:"merge"`
}


func NewWriteExcelObj(param []byte) (WriteExcel, error) {

	res := new(WriteExcelObj)
	e := json.Unmarshal(param, &res)
	if e != nil {
		return nil, e
	}

	newMerge := make([]int, 0)
	for _, i := range res.Merge {
		if i <= 0 {
			continue
		}
		newMerge = append(newMerge, i-1)
	}
	res.Merge = newMerge

	return res, nil
}

/*
WriteExcel ----Write Excel
showType="00"：the excel content is displayed by rows,default.
showType="01"：The excel content is displayed by columns.
*/
func (m *WriteExcelObj) WriteExcel(sheetName, showType string, exportData *export_model.ExportData) (file *excelize.File, e error) {

	file = excelize.NewFile()
	if sheetName != utils.DefaultSheetName {
		file.NewSheet(sheetName)
		file.DeleteSheet(utils.DefaultSheetName)
	}

	if len(exportData.Title) == 0 {
		return nil, errors.New("title can not be null")
	}


	if showType == utils.ShowTypeByRow {
		file, e = m.writeExcelByRow(utils.ShowTypeByRow, sheetName, file, exportData)
		if e != nil {
			return
		}
	} else  {
		file, e = m.writeExcelByColumn(utils.ShowTypeByCol, sheetName, file, exportData)
		if e != nil {
			return
		}
	}

	return
}

//writeExcelByRow ----Insert data by row, if there is no data, there will also be a header
func (m *WriteExcelObj) writeExcelByRow(showType, sheetName string,
	file *excelize.File, exportData *export_model.ExportData) (*excelize.File, error) {

	titleInitAxis := "A1"
	col := columnNumberToName(len(exportData.Title))
	titleEndAxis := col + "1"
	e := file.SetSheetRow(sheetName, titleInitAxis, &exportData.Title)
	if e != nil {
		return nil, e
	}

	for i, data := range exportData.Data {
		if len(data) == 0 {
			continue
		}

		e = file.SetSheetRow(sheetName, "A" + strconv.Itoa(i+2), &data)
		if e != nil {
			return nil, e
		}

		file, e = m.mergeCell(i, showType, sheetName, exportData, file)
		if e != nil {
			return nil, e
		}
	}

	dataInitAxis := ""
	dataEndAxis := ""
	if len(exportData.Data) != 0 {
		dataInitAxis = "A2"
		dataEndAxis = col + strconv.Itoa(len(exportData.Data)+1)
	}

	file, e = setCellStyle(file, sheetName, titleInitAxis, titleEndAxis, dataInitAxis, dataEndAxis)
	if e != nil {
		return nil, e
	}

	file, e = setColWidth(sheetName, col, file, m.ColWidth)
	if e != nil {
		return nil, e
	}

	return file, nil
}

//writeExcelByColumn ----Insert data by column, if there is no data, there will also be a header
func (m *WriteExcelObj) writeExcelByColumn(showType, sheetName string,
	file *excelize.File, exportData *export_model.ExportData)  (*excelize.File, error) {


	var e error
	for i, title := range exportData.Title {
		e = file.SetCellValue(sheetName, `A` + strconv.Itoa(i+1), title)
		if e != nil {
			return nil, e
		}
	}

	dataEndAxis := ""
	for i, data := range exportData.Data {
		dataEndAxis = columnNumberToName(i+2)
		for j, cellValue := range data {
			e = file.SetCellValue(sheetName, dataEndAxis + strconv.Itoa(j+1), cellValue)
			if e != nil {
				return nil, e
			}
		}

		file, e = m.mergeCell(i, showType, sheetName, exportData, file)
		if e != nil {
			return nil, e
		}
	}

	file, e = setCellStyle(file, sheetName, `A1`, `A` + strconv.Itoa(len(exportData.Title)), `B1`, dataEndAxis + strconv.Itoa(len(exportData.Title)))
	if e != nil {
		return nil, e
	}

	file, e = setColWidth(sheetName, columnNumberToName(len(exportData.Title)), file, m.ColWidth)
	if e != nil {
		return nil, e
	}

	return file, nil
}

/*
mergeCell ----单元格合并
行展示时,合并的是相同的列; 列展示时,合并的是相同的行
第一条数据不合并,忽略
从第二条数据开始: 若要合并的单元格是第一行或列(k=0), 则只需要判断当前单元格的值(exData.Data[dataIndex][k])与前一条数据相同单元格的值(exData.Data[dataIndex-1][k])是否相等,若相等,则合并
			   若不是第一行或列(k>=1),则需要判断当前单元格前面单元格(包括当前单元格)的值(exData.Data[dataIndex][0:k+1])与上一条数据的前几列单元格的值(exData.Data[dataIndex-1][0:k+1])是否全部相等,若全部相等,则合并

When rows are displayed, the same columns are merged; When columns are displayed, the same rows are merged
The first data is not merged, ignored
Start with number two: If the cell to be merged is the first row or column(k=0), it only needs to determine whether the value of the current cell is equal
					   to the value of the same cell in the previous data;
					   if it were not for the first row or column(k>=1), we need to determine whether the values of the cells in front of
                       the current cell and the values of the cells in the previous columns of the previous data are all phase And so on,
                       if all equal, then merge.
*/
func (m *WriteExcelObj) mergeCell(dataIndex int, showType, sheetName string,
	exportData *export_model.ExportData, file *excelize.File) (*excelize.File, error) {

	//第一条数据不合并,忽略
	if dataIndex == 0 {
		return file, nil
	}

	for _, k := range m.Merge {
		ok := true
		//要合并的单元格是第一行或列(k=0), 则只需要判断当前单元格的值(exData.Data[dataIndex][k])与前一条数据相同单元格的值是否exportData.Data[dataIndex-1][k]相等
		if k == 0 {
			if exportData.Data[dataIndex][k] != exportData.Data[dataIndex-1][k] {
				ok = false
			}
		} else {
			//不是第一行或列(k>=1),则需要判断当前单元格前面单元格(包括当前单元格)的值(exData.Data[dataIndex][0:k+1])与上一条数据的前几列单元格的值exportData.Data[dataIndex-1][0:k+1]是否全部相等,不包括当前单元格时: j := k-1
			for j := k; j >= 0 && ok; j-- {
				if exportData.Data[dataIndex][j] == exportData.Data[dataIndex-1][j] {
					ok = true
				} else {
					ok = false
				}
			}
		}

		if ok {
			if showType == utils.ShowTypeByRow {
				hCell := columnNumberToName(k+1) + strconv.Itoa(dataIndex+1)
				vCell := columnNumberToName(k+1) + strconv.Itoa(dataIndex+2)
				e := file.MergeCell(sheetName, hCell, vCell)
				if e != nil {
					return nil, e
				}
			} else {
				hCell := columnNumberToName(dataIndex+1) + strconv.Itoa(k+1)
				vCell := columnNumberToName(dataIndex+2) + strconv.Itoa(k+1)
				e := file.MergeCell(sheetName, hCell, vCell)
				if e != nil {
					return nil, e
				}
			}
		}
	}

	return file, nil
}

func setCellStyle(file *excelize.File, sheetName, titleInitAxis,
	titleEndAxis, dataInitAxis, dataEndAxis string) (*excelize.File, error) {


	if titleInitAxis != "" || titleEndAxis != "" {
		/*
			表头样式设置
			alignment:单元格对齐方式
			font: 单元格字体
			border: 单元格边框
			fill: 单元格填充
		*/
		style := new(excelize.Style)
		style.Alignment = &excelize.Alignment{
			Horizontal: "left",
			Vertical: "center",
			WrapText: true,
		}
		style.Protection = &excelize.Protection{
			Hidden: false,
			Locked: false,
		}
		style.Font = &excelize.Font{
			Bold: true,
			Italic: false,
			Family: "微软雅黑",
			Size: 10,
			Color: "000000",
		}
		style.Border = []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		}
		style.Fill = excelize.Fill{
			Type: "pattern",
			Color: []string{"#FFFF33"},
			Pattern: 1,
		}
		titleStyle, e := file.NewStyle(style)
		if e != nil {
			return nil, e
		}

		//设置表头样式
		e = file.SetCellStyle(sheetName, titleInitAxis, titleEndAxis, titleStyle)
		if e != nil {
			return nil, e
		}
	}

	if dataInitAxis != "" || dataEndAxis != "" {
		//数据行的样式
		style := new(excelize.Style)
		style.Alignment = &excelize.Alignment{
			Horizontal: "left",
			Vertical: "center",
			WrapText: true,
		}
		style.Protection = &excelize.Protection{
			Hidden: false,
			Locked: false,
		}
		style.Font = &excelize.Font{
			Bold: false,
			Italic: false,
			Family: "微软雅黑",
			Size: 10,
			Color: "000000",
		}
		style.Border = []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		}
		dataStyle, e := file.NewStyle(style)
		if e != nil {
			return nil, e
		}

		//设置数据行样式
		e = file.SetCellStyle(sheetName, dataInitAxis, dataEndAxis, dataStyle)
		if e != nil {
			return nil, e
		}
	}

	return file, nil
}

/*
setColWidth ----set col width
colWidth: map[string]float64, The value of key supports: A,B: Represents two columns A and B;
		 D: Represents D column;
		E...G: Represents E to G column;
		H...: Represents H column to the last column(endCol column).
 */
func setColWidth(sheetName, endCol string, file *excelize.File, colWidth map[string]float64) (*excelize.File, error) {

	for col, width := range colWidth {
		if strings.Contains(col, ",") {
			for _, column := range strings.Split(col, ",") {
				e := file.SetColWidth(sheetName, strings.TrimSpace(column), strings.TrimSpace(column), width)
				if e != nil {
					return nil, e
				}
			}
		} else if strings.Contains(col, "...") {
			columns := strings.Split(col, "...")
			startCol := strings.TrimSpace(columns[0])
			if columns[1] != "" {
				endCol = strings.TrimSpace(columns[1])
			}
			e := file.SetColWidth(sheetName, startCol, endCol, width)
			if e != nil {
				return nil, e
			}
		} else {
			e := file.SetColWidth(sheetName, strings.TrimSpace(col), strings.TrimSpace(col), width)
			if e != nil {
				return nil, e
			}
		}
	}

	return file, nil
}

//columnNumberToName  num：0..26/27..702/703..18278 ----> col：A..Z/AA..ZZ/AAA...ZZZ
func columnNumberToName(num int) string {

	var col string
	for num > 0 {
		col = string(rune((num-1)%26+65)) + col
		num = (num - 1) / 26
	}

	return col
}