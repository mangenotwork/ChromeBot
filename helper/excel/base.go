package excel

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strconv"

	_ "golang.org/x/image/bmp"

	"github.com/xuri/excelize/v2"
)

// ShowExcel 读取ecxel并打印excel
func ShowExcel(path, sheetName string) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

// WriteListToExcel 将二维列表写入Excel
// list: 二维字符串列表（[[行1列1, 行1列2], [行2列1, 行2列2]]）
// path: 输出文件路径
// sheetName: 工作表名称
func WriteListToExcel(list [][]string, path, sheetName string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	for rowIdx, row := range list {
		for colIdx, cellValue := range row {
			cell, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				fmt.Println("[Err] Excel CoordinatesToCellName error:", err)
				return err
			}
			f.SetCellValue(sheetName, cell, cellValue)
		}
	}

	return saveAs(f, path)
}

// WriteMapToExcel 将字典列表写入Excel（字典key为表头，value为值）
// data: 字典列表（[{"姓名":"张三","年龄":"25"}, {"姓名":"李四","年龄":"30"}]）
// path: 输出文件路径
// sheetName: 工作表名称
func WriteMapToExcel(data []map[string]string, path, sheetName string) error {
	if len(data) == 0 {
		return fmt.Errorf("数据为空")
	}
	f := excelize.NewFile()
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	headers := make([]string, 0, len(data[0]))
	for k := range data[0] {
		headers = append(headers, k)
	}

	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for rowIdx, item := range data {
		for colIdx, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, item[header])
		}
	}

	return saveAs(f, path)
}

// ReadExcelToList 读取Excel全部数据到二维列表
// path: 文件路径
// sheetName: 工作表名称
// return: 二维列表、错误
func ReadExcelToList(path, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	return f.GetRows(sheetName)
}

// ReadExcelToMap 读取Excel全部数据到字典列表（表头为key）
// path: 文件路径
// sheetName: 工作表名称
// return: 字典列表、错误
func ReadExcelToMap(path, sheetName string) ([]map[string]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("工作表无数据")
	}

	headers := rows[0]
	result := make([]map[string]string, 0, len(rows)-1)

	for _, row := range rows[1:] {
		item := make(map[string]string)
		for colIdx, header := range headers {
			if colIdx < len(row) {
				item[header] = row[colIdx]
			} else {
				item[header] = "" // 空值填充
			}
		}
		result = append(result, item)
	}

	return result, nil
}

// ReadExcelToMapRowHead 读取Excel全部数据到字典列表指定rowHead（表头为key）
// path: 文件路径
// sheetName: 工作表名称
// return: 字典列表、错误
func ReadExcelToMapRowHead(path, sheetName string, rowHead int) ([]map[string]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("工作表无数据")
	}
	if rowHead < 0 {
		rowHead = 0
	}
	if rowHead > len(rows) {
		rowHead = len(rows) - 1
	}

	headers := rows[rowHead-1]
	result := make([]map[string]string, 0, len(rows)-1)

	for i, row := range rows {
		if i == rowHead-1 {
			continue
		}
		item := make(map[string]string)
		for colIdx, header := range headers {
			if colIdx < len(row) {
				item[header] = row[colIdx]
			} else {
				item[header] = "" // 空值填充
			}
		}
		result = append(result, item)
	}

	return result, nil
}

// GetExcelInfo 获取Excel文件基本信息（工作表列表、文件属性等）
// path: 文件路径
// return: 信息字符串、错误
func GetExcelInfo(path string) (string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sheets := f.GetSheetList()

	props, err := f.GetDocProps()
	if err != nil {
		return "", err
	}

	info := "Excel文件信息 : \n"
	info += fmt.Sprintf("  工作表列表：%v\n", sheets)
	info += fmt.Sprintf("  创建者：%s\n", props.Creator)
	info += fmt.Sprintf("  创建时间：%s\n", props.Created)
	info += fmt.Sprintf("  修改时间：%s\n", props.Modified)
	info += fmt.Sprintf("  标题：%s\n", props.Title)
	info += fmt.Sprintf("  主题：%s\n", props.Subject)

	return info, nil
}

// GetSheetInfo 获取指定工作表信息（行数、列数、合并单元格等）
// path: 文件路径
// sheetName: 工作表名称
// return: 信息字符串、错误
func GetSheetInfo(path, sheetName string) (string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", err
	}
	rowCount := len(rows)

	colCount := 0
	for _, row := range rows {
		if len(row) > colCount {
			colCount = len(row)
		}
	}

	mergeCells, err := f.GetMergeCells(sheetName)
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("工作表 [%s] 信息：\n", sheetName)
	info += fmt.Sprintf("  总行数：%d\n", rowCount)
	info += fmt.Sprintf("  总列数：%d\n", colCount)
	info += fmt.Sprintf("  合并单元格数量：%d\n", len(mergeCells))
	for i, mc := range mergeCells {
		info += fmt.Sprintf("    合并单元格%d：%s\n", i+1, mc.GetCellValue())
	}

	return info, nil
}

// GetExcelSheetNames 获取Excel文件中所有工作表名称
// path: 文件路径
// return: 工作表名称数组、错误
func GetExcelSheetNames(path string) ([]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("关闭文件失败: %v\n", err)
		}
	}()

	sheetNames := f.GetSheetList()

	return sheetNames, nil
}

// ReadCell 读取指定单元格数据
// path: 文件路径
// sheetName: 工作表名称
// cell: 单元格名称（如"A1"）
// return: 单元格值、错误
func ReadCell(path, sheetName, cell string) (string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	return f.GetCellValue(sheetName, cell)
}

// WriteCell 写入指定单元格数据
// path: 文件路径
// sheetName: 工作表名称
// cell: 单元格名称（如"A1"）
// value: 写入值（支持string/int/float等）
// return: 错误
func WriteCell(path, sheetName, cell string, value interface{}) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	f.SetCellValue(sheetName, cell, value)
	return f.Save()
}

// ClearCell 删除指定单元格数据（清空值）
// path: 文件路径
// sheetName: 工作表名称
// cell: 单元格名称（如"A1"）
// return: 错误
func ClearCell(path, sheetName, cell string) error {
	return WriteCell(path, sheetName, cell, "")
}

// GetCellValueByPos 按行列位置读取指定单元格的值
// path: 文件路径
// sheetName: 工作表名称
// row: 行号（从1开始）
// col: 列号（从1开始）
// return: 单元格值、错误
func GetCellValueByPos(path, sheetName string, row, col int) (string, error) {
	// 打开Excel文件
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	if row < 1 || col < 1 {
		return "", fmt.Errorf("行列号必须大于0（当前行：%d，列：%d）", row, col)
	}

	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return "", fmt.Errorf("行列号转换失败: %w", err)
	}

	value, err := f.GetCellValue(sheetName, cell)
	if err != nil {
		return "", fmt.Errorf("读取单元格[%s]失败: %w", cell, err)
	}

	return value, nil
}

// SetCellValueByPos 按行列位置写入值到指定单元格
// path: 文件路径
// sheetName: 工作表名称
// row: 行号（从1开始）
// col: 列号（从1开始）
// value: 要写入的值（支持string/int/float/bool等类型）
// return: 错误
func SetCellValueByPos(path, sheetName string, row, col int, value interface{}) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	if row < 1 || col < 1 {
		return fmt.Errorf("行列号必须大于0（当前行：%d，列：%d）", row, col)
	}

	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return fmt.Errorf("行列号转换失败: %w", err)
	}

	if err := f.SetCellValue(sheetName, cell, value); err != nil {
		return fmt.Errorf("写入单元格[%s]失败: %w", cell, err)
	}

	if err := f.Save(); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}

	return nil
}

// ClearCellValueByPos 按行列位置删除（清空）指定单元格的数据
// path: 文件路径
// sheetName: 工作表名称
// row: 行号（从1开始）
// col: 列号（从1开始）
// return: 错误
func ClearCellValueByPos(path, sheetName string, row, col int) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	if row < 1 || col < 1 {
		return fmt.Errorf("行列号必须大于0（当前行：%d，列：%d）", row, col)
	}

	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return fmt.Errorf("行列号转换失败: %w", err)
	}

	if err := f.SetCellValue(sheetName, cell, ""); err != nil {
		return fmt.Errorf("清空单元格[%s]失败: %w", cell, err)
	}

	if err := f.Save(); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}

	return nil
}

// ReadRow 读取指定行数据
// path: 文件路径
// sheetName: 工作表名称
// rowNum: 行号（从1开始）
// return: 行数据列表、错误
func ReadRow(path, sheetName string, rowNum int) ([]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if rowNum < 1 || rowNum > len(rows) {
		return nil, fmt.Errorf("行号超出范围")
	}
	return rows[rowNum-1], nil
}

// WriteRow 写入指定行数据
// path: 文件路径
// sheetName: 工作表名称
// rowNum: 行号（从1开始）
// values: 行数据列表
// return: 错误
func WriteRow(path, sheetName string, rowNum int, values []string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	for colIdx, value := range values {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
		f.SetCellValue(sheetName, cell, value)
	}
	return f.Save()
}

// DeleteRow 删除指定行
// path: 文件路径
// sheetName: 工作表名称
// rowNum: 行号（从1开始）
// return: 错误
func DeleteRow(path, sheetName string, rowNum int) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	err = f.RemoveRow(sheetName, rowNum)
	if err != nil {
		return err
	}
	return f.Save()
}

// ReadColumn 读取指定列数据
// path: 文件路径
// sheetName: 工作表名称
// col: 列数（如1代表A列，2代表B列...）
// return: 列数据列表、错误
func ReadColumn(path, sheetName string, col int) ([]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	// 校验列数合法性
	if col < 1 {
		return nil, excelize.ErrColumnNumber
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(rows))
	for _, row := range rows {
		if col-1 < len(row) {
			result = append(result, row[col-1])
		} else {
			result = append(result, "")
		}
	}
	return result, nil
}

// WriteColumn 写入指定列数据
// path: 文件路径
// sheetName: 工作表名称
// col: 列数（如1代表A列，2代表B列...）
// values: 列数据列表
// return: 错误
func WriteColumn(path, sheetName string, col int, values []string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	// 校验列数合法性
	if col < 1 {
		return excelize.ErrColumnNumber
	}

	// 将列数转换为列名（如1->A，2->B）
	colName, err := excelize.ColumnNumberToName(col)
	if err != nil {
		return err
	}

	for rowIdx, value := range values {
		cell := colName + strconv.Itoa(rowIdx+1)
		if err := f.SetCellValue(sheetName, cell, value); err != nil {
			return err
		}
	}
	return f.Save()
}

// DeleteColumn 删除指定列
// path: 文件路径
// sheetName: 工作表名称
// col: 列数（如1代表A列，2代表B列...）
// return: 错误
func DeleteColumn(path, sheetName string, col int) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	// 校验列数合法性
	if col < 1 {
		return excelize.ErrColumnNumber
	}

	// 将列数转换为列名（如1->A，2->B）
	colName, err := excelize.ColumnNumberToName(col)
	if err != nil {
		return err
	}

	err = f.RemoveCol(sheetName, colName)
	if err != nil {
		return err
	}
	return f.Save()
}

// ReadColumnCell 读取指定列数据
// path: 文件路径
// sheetName: 工作表名称
// colName: 列名（如"A"）
// return: 列数据列表、错误
func ReadColumnCell(path, sheetName, colName string) ([]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	// 转换列名到列索引（A->1, B->2...）
	colIdx, _, err := excelize.CellNameToCoordinates(colName + "1")
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(rows))
	for _, row := range rows {
		if colIdx-1 < len(row) {
			result = append(result, row[colIdx-1])
		} else {
			result = append(result, "")
		}
	}
	return result, nil
}

// WriteColumnCell 写入指定列数据
// path: 文件路径
// sheetName: 工作表名称
// colName: 列名（如"A"）
// values: 列数据列表
// return: 错误
func WriteColumnCell(path, sheetName, colName string, values []string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	for rowIdx, value := range values {
		cell := colName + strconv.Itoa(rowIdx+1)
		f.SetCellValue(sheetName, cell, value)
	}
	return f.Save()
}

// DeleteColumnCell 删除指定列
// path: 文件路径
// sheetName: 工作表名称
// colName: 列名（如"A"）
// return: 错误
func DeleteColumnCell(path, sheetName, colName string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	err = f.RemoveCol(sheetName, colName)
	if err != nil {
		return err
	}
	return f.Save()
}

// InsertImage 插入图片到指定单元格
// path: Excel文件路径
// sheetName: 工作表名称
// cell: 单元格（如"A1"）
// imgPath: 图片文件路径（支持jpg/png等）
// return: 错误
func InsertImage(path, sheetName, cell, imgPath string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	absImgPath, err := filepath.Abs(imgPath)
	if err != nil {
		return err
	}

	err = f.AddPicture(sheetName, cell, absImgPath, nil)
	if err != nil {
		return err
	}
	return f.Save()
}

// todo 读取并下载ecxel里的图片

// SetCellStyle 设置单元格样式（字体、颜色、对齐、背景）
// path: Excel文件路径
// sheetName: 工作表名称
// cell: 单元格（如"A1"）
// fontBold: 是否加粗
// fontColor: 字体颜色（十六进制，如"FF0000"）
// bgColor: 背景颜色（十六进制，如"E0E0E0"）
// alignCenter: 是否居中
// return: 错误
func SetCellStyle(path, sheetName, cell string, fontBold bool, fontColor, bgColor string, alignCenter bool) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	style := &excelize.Style{
		Font: &excelize.Font{
			Bold:  fontBold,
			Color: fontColor,
			Size:  12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{bgColor},
			Pattern: 1,
		},
	}
	if alignCenter {
		style.Alignment = &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		}
	}

	styleID, err := f.NewStyle(style)
	if err != nil {
		return err
	}
	err = f.SetCellStyle(sheetName, cell, cell, styleID)
	if err != nil {
		return err
	}
	return f.Save()
}

// MergeCells 合并单元格
// path: Excel文件路径
// sheetName: 工作表名称
// startCell: 起始单元格（如"A1"）
// endCell: 结束单元格（如"C3"）
// return: 错误
func MergeCells(path, sheetName, startCell, endCell string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	err = f.MergeCell(sheetName, startCell, endCell)
	if err != nil {
		return err
	}
	return f.Save()
}

// AddChart 添加柱状图到指定单元格
// path: Excel文件路径
// sheetName: 工作表名称
// cell: 图表插入位置（如"D2"）
// title: 图表标题
// categories: 分类数据范围（如"Sheet1!$A$2:$A$4"）
// values: 数值数据范围（如"Sheet1!$B$2:$B$4"）
// return: 错误
func AddChart(path, sheetName, cell, title, categories, values string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	chart := &excelize.Chart{
		Type: "col",
		Series: []excelize.ChartSeries{
			{
				Name:       fmt.Sprintf("%s!$B$1", sheetName),
				Categories: categories,
				Values:     values,
			},
		},
		Title: excelize.ChartTitle{
			Name: title,
		},
	}

	err = f.AddChart(sheetName, cell, chart)
	if err != nil {
		return err
	}
	return f.Save()
}

// SetCellFormula 设置单元格公式
// path: Excel文件路径
// sheetName: 工作表名称
// cell: 单元格（如"A4"）
// formula: 公式（如"SUM(A1:A3)"）
// return: 错误
func SetCellFormula(path, sheetName, cell, formula string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName = getTargetSheetName(f, sheetName)

	// 设置公式
	err = f.SetCellFormula(sheetName, cell, formula)
	if err != nil {
		return err
	}
	// 计算公式（可选）
	_, err = f.CalcCellValue(sheetName, cell)
	if err != nil {
		return err
	}
	return f.Save()
}

// todo SqlResultToExcel 模拟将SQL查询结果写入Excel（需结合实际数据库驱动）
// sqlResult: SQL查询结果（二维列表，第一行为表头）
// path: Excel输出路径
// sheetName: 工作表名称
// return: 错误
func SqlResultToExcel(sqlResult [][]string, path, sheetName string) error {
	// 实际使用时需替换为真实的SQL查询逻辑（如database/sql库）
	if len(sqlResult) == 0 {
		return fmt.Errorf("SQL查询结果为空")
	}
	return WriteListToExcel(sqlResult, path, sheetName)
}

// todo ExcelToSql 模拟将Excel数据转换为INSERT SQL语句
// path: Excel文件路径
// sheetName: 工作表名称
// tableName: 数据库表名
// return: SQL语句列表、错误
func ExcelToSql(path, sheetName, tableName string) ([]string, error) {
	// 读取Excel为字典列表
	data, err := ReadExcelToMap(path, sheetName)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("Excel无数据")
	}

	// 生成INSERT语句
	sqlList := make([]string, 0, len(data))
	for _, item := range data {
		// 拼接字段和值
		fields := make([]string, 0, len(item))
		values := make([]string, 0, len(item))
		for k, v := range item {
			fields = append(fields, k)
			values = append(values, fmt.Sprintf("'%s'", v)) // 简单转义，实际需处理SQL注入
		}
		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
			tableName,
			fmt.Sprintf("%s", fields),
			fmt.Sprintf("%s", values),
		)
		sqlList = append(sqlList, sql)
	}
	return sqlList, nil
}
