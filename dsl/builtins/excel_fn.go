package builtins

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/helper/excel"
	"fmt"

	gt "github.com/mangenotwork/gathertool"
)

var excelFn = map[string]interpreter.Function{
	"ExcelSave":        excelSave,        // ExcelSave(arg, path, 可选参数sheetName) 将变量保存到excel
	"ExcelReadList":    excelReadList,    // ExcelReadList(path, 可选参数sheetName) 读取excel返回二维列表
	"ExcelReadDict":    excelReadDict,    // ExcelReadDict(path, 可选参数sheetName) 读取excel返回字典
	"ExcelShow":        excelShow,        // ExcelShow(path, 可选参数sheetName) 显示excel
	"ExcelInfo":        excelInfo,        // ExcelInfo(path) 获取excel信息
	"ExcelSheetInfo":   excelSheetInfo,   // ExcelSheetInfo(path, sheetName) 获取excel的sheet信息
	"ExcelSheet":       excelSheet,       // ExcelSheet(path) 获取excel的sheet信息
	"ExcelGetByCell":   excelGetByCell,   // ExcelGetByCell(path, cell, 可选参数sheetName) 通过位置标签获取excel数据   cell 标签 A1 B1 C1 ...
	"ExcelGetByPos":    excelGetByPos,    // ExcelGetByPos(path, row, col, 可选参数sheetName) 通过位置获取excel数据
	"ExcelSetByCell":   excelSetByCell,   // ExcelSetByCell(path, cell, value, 可选参数sheetName) 通过位置标签设置excel数据   cell 标签 A1 B1 C1 ...
	"ExcelSetByPos":    excelSetByPos,    // ExcelSetByPos(path, row, col, value, 可选参数sheetName) 通过位置设置excel数据
	"ExcelClearByCell": excelClearByCell, // ExcelClearByCell(path, cell, 可选参数sheetName) 通过位置标签清除excel数据   cell 标签 A1 B1 C1 ...
	"ExcelClearByPos":  excelClearByPos,  // ExcelClearByPos(path, row, col, 可选参数sheetName) 通过位置清除excel数据
	"ExcelReadRow":     excelReadRow,     // ExcelReadRow(path, row, 可选参数sheetName)  读取指定行数据
	"ExcelWriteRow":    excelWriteRow,    // ExcelWriteRow(path, row, list, 可选参数sheetName)  写入指定行数据
	"ExcelDeleteRow":   excelDeleteRow,   // ExcelDeleteRow(path, row, 可选参数sheetName)  删除指定行数据
	"ExcelReadCol":     excelReadCol,     // ExcelReadCol(path, col, 可选参数sheetName)  读取指定列数据
	"ExcelWriteCol":    excelWriteCol,    // ExcelWriteCol(path, col, list, 可选参数sheetName)  写入指定列数据
	"ExcelDeleteCol":   excelDeleteCol,   // ExcelDeleteCol(path, col, 可选参数sheetName)  删除指定列数据
	"ExcelReadCell":    excelReadCell,    // ExcelReadCell(path, cell, 可选参数sheetName)  读取列 cell 标签 A B C ...
	"ExcelWriteCell":   excelWriteCell,   // ExcelWriteCell(path, cell, list, 可选参数sheetName)  写入列 cell 标签 A B C ...
	"ExcelDeleteCell":  excelDeleteCell,  // ExcelDeleteCell(path, cell, 可选参数sheetName)  删除指定列数据  cell 标签 A B C ...
	"ExcelImg":         excelImg,         // ExcelImg(path, cell, imgPath, 可选参数sheetName)   插入图片 cell 标签 A1 B1 C1 ...

	// ExcelCellStyle(path, cell, style, 可选参数sheetName) 设置单元格样式 cell 标签 A1 B1 C1 ...
	// style {
	//  	fontBold: 是否加粗
	// 		fontColor: 字体颜色（十六进制，如"FF0000"）
	// 		bgColor: 背景颜色（十六进制，如"E0E0E0"）
	// 		alignCenter: 是否居中
	// }
	"ExcelCellStyle": excelCellStyle,

	"ExcelMergeCells": excelMergeCells, // ExcelMergeCells(path, startCell, endCell, 可选参数sheetName) 合并单元格 cell 标签 A1 B1 C1 ...
	"ExcelSetFormula": excelSetFormula, // ExcelSetFormula(path, cell, formula, 可选参数sheetName) 给单元格设置公式 标签 A1 B1 C1 ...  formula公式 如"SUM(A1:A3)"

	// todo AddChart 添加图表到单元格

}

func excelSave(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelSave(arg, path) 需要两个参数")
	}

	path, pathOK := args[1].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelSave(arg, path)  path 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelSave(arg, path, sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	fmt.Printf("%T", args[0])

	dataType := "list"
	dataList := make([][]string, 0)
	dataMap := make([]map[string]string, 0)

	switch args[0].(type) {
	case []interpreter.Value:
		for _, v := range args[0].([]interpreter.Value) {
			switch v := v.(type) {
			case []interpreter.Value:
				dataItem := make([]string, 0)
				for _, vv := range v {
					dataItem = append(dataItem, gt.Any2String(vv))
				}
				dataList = append(dataList, dataItem)
			case interpreter.DictType:
				dataType = "dict"
				dataMapItem := make(map[string]string)
				for k, vv := range v {
					dataMapItem[gt.Any2String(k)] = gt.Any2String(vv)
				}
				dataMap = append(dataMap, dataMapItem)
			default:
				if len(dataList) == 0 {
					dataList = append(dataList, make([]string, 0))
				}
				dataList[0] = append(dataList[0], gt.Any2String(v))
			}
		}

	case interpreter.DictType:
		dataType = "dict"
		dataMapItem := make(map[string]string)
		for k, vv := range args[0].(interpreter.DictType) {
			dataMapItem[gt.Any2String(k)] = gt.Any2String(vv)
		}
		dataMap = append(dataMap, dataMapItem)

	default:
		dataList = append(dataList, make([]string, 0))
		dataList[0] = append(dataList[0], gt.Any2String(args[0]))

	}

	var err error
	if dataType == "list" {
		err = excel.WriteListToExcel(dataList, path, sheetName)
	}
	if dataType == "dict" {
		err = excel.WriteMapToExcel(dataMap, path, sheetName)
	}
	if err != nil {
		fmt.Println("[Err] excel write error: ", err)
	}

	return nil, nil
}

func excelReadList(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ExcelReadList(path, 可选参数sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelReadList(path, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 2 {
		sheetName, sheetNameOK = args[1].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelReadList(path, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	data, err := excel.ReadExcelToList(path, sheetName)
	if err != nil {
		fmt.Println("[Err]读取Excel文件失败:", err)
	}

	return data, nil
}

func excelReadDict(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ExcelReadDict(path, 可选参数sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelReadDict(path, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 2 {
		sheetName, sheetNameOK = args[1].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelReadDict(path, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	data, err := excel.ReadExcelToMap(path, sheetName)

	if err != nil {
		fmt.Println("[Err]读取Excel文件失败:", err)
	}

	return data, nil
}

func excelShow(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ExcelShow(path, 可选参数sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelShow(path, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 2 {
		sheetName, sheetNameOK = args[1].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelShow(path, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	excel.ShowExcel(path, sheetName)
	return nil, nil
}

func excelInfo(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ExcelInfo(path) 需要一个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelInfo(path) path 参数要求是字符串 ")
	}

	info, err := excel.GetExcelInfo(path)
	if err != nil {
		fmt.Println("[Err] 获取Excel信息失败：", err.Error())
		return "", err
	}

	fmt.Println("[Info] Excel信息:", info)

	return info, nil
}

func excelSheetInfo(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ExcelSheetInfo(path, sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelSheetInfo(path, sheetName)  path 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 2 {
		sheetName, sheetNameOK = args[1].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelSheetInfo(path, sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	info, err := excel.GetSheetInfo(path, sheetName)
	if err != nil {
		fmt.Println("[Err] 获取Excel信息失败：", err.Error())
		return "", err
	}

	fmt.Println("[Info] Excel Sheet信息:", info)

	return info, nil
}

func excelSheet(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ExcelSheet(path) 需要一个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelSheet(path) path 参数要求是字符串 ")
	}

	info, err := excel.GetExcelSheetNames(path)
	if err != nil {
		fmt.Println("[Err] 获取Excel信息失败：", err.Error())
		return "", err
	}

	fmt.Println("[Info] Excel信息:", info)

	return info, nil
}

func excelGetByCell(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelGetByCell(path, cell, 可选参数sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelGetByCell(path, cell, 可选参数sheetName) path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelGetByCell(path, cell, 可选参数sheetName) cell 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelGetByCell(path, cell, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	data, err := excel.ReadCell(path, sheetName, cell)
	if err != nil {
		fmt.Println("获取单元格数据失败 err = ", err.Error())
	}
	return data, nil
}

func excelGetByPos(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelGetByPos(path, row, col, 可选参数sheetName) 需要三个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelGetByPos(path, row, col, 可选参数sheetName) path 参数要求是字符串 ")
	}

	row, rowOK := args[1].(int64)
	if !rowOK {
		return nil, fmt.Errorf("ExcelGetByPos(path, row, col, 可选参数sheetName) row 参数要求是整数 ")
	}

	col, colOK := args[2].(int64)
	if !colOK {
		return nil, fmt.Errorf("ExcelGetByPos(path, row, col, 可选参数sheetName) col 参数要求是整数 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelGetByPos(path, row, col, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	data, err := excel.GetCellValueByPos(path, sheetName, int(row), int(col))
	if err != nil {
		fmt.Println("获取单元格数据失败 err = ", err.Error())
	}
	return data, nil
}

func excelSetByCell(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelSetByCell(path, cell, value, 可选参数sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelSetByCell(path, cell, value, 可选参数sheetName) path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelSetByCell(path, cell, value, 可选参数sheetName) cell 参数要求是字符串 ")
	}

	value := args[2]

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelSetByCell(path, cell, value, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.WriteCell(path, sheetName, cell, value)
	if err != nil {
		fmt.Println("[Err]数据写入失败 err = ", err.Error())
	}
	fmt.Println("数据写入成功")

	return nil, nil
}

func excelSetByPos(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("ExcelSetByPos(path, row, col, value, 可选参数sheetName) 需要四个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelSetByPos(path, row, col, value, 可选参数sheetName) path 参数要求是字符串 ")
	}

	row, rowOK := args[1].(int64)
	if !rowOK {
		return nil, fmt.Errorf("ExcelSetByPos(path, row, col, value, 可选参数sheetName) row 参数要求是整数 ")
	}

	col, colOK := args[2].(int64)
	if !colOK {
		return nil, fmt.Errorf("ExcelSetByPos(path, row, col, value, 可选参数sheetName) col 参数要求是整数 ")
	}

	value := args[3]

	sheetName := ""
	sheetNameOK := false
	if len(args) == 5 {
		sheetName, sheetNameOK = args[4].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelSetByPos(path, row, col, value, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.SetCellValueByPos(path, sheetName, int(row), int(col), value)
	if err != nil {
		fmt.Println("[Err]数据写入失败 err = ", err.Error())
	}
	fmt.Println("数据写入成功")

	return nil, nil
}

func excelClearByCell(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelClearByCell(path, cell, 可选参数sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelClearByCell(path, cell, 可选参数sheetName) path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelClearByCell(path, cell, 可选参数sheetName) cell 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelClearByCell(path, cell, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.ClearCell(path, sheetName, cell)
	if err != nil {
		fmt.Println("单元格清除失败 err = ", err.Error())
	}
	fmt.Println("单元格清除成功")
	return nil, nil
}

func excelClearByPos(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelClearByPos(path, row, col, 可选参数sheetName)  需要三个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelClearByPos(path, row, col, 可选参数sheetName) path 参数要求是字符串 ")
	}

	row, rowOK := args[1].(int64)
	if !rowOK {
		return nil, fmt.Errorf("ExcelClearByPos(path, row, col, 可选参数sheetName) row 参数要求是整数 ")
	}

	col, colOK := args[2].(int64)
	if !colOK {
		return nil, fmt.Errorf("ExcelClearByPos(path, row, col, 可选参数sheetName) col 参数要求是整数 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelClearByPos(path, row, col, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.ClearCellValueByPos(path, sheetName, int(row), int(col))
	if err != nil {
		fmt.Println("单元格清除失败 err = ", err.Error())
	}
	fmt.Println("单元格清除成功")
	return nil, nil
}

func excelReadRow(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelReadRow(path, row, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelReadRow(path, row, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	row, rowOK := args[1].(int64)
	if !rowOK {
		return nil, fmt.Errorf("ExcelReadRow(path, row, 可选参数sheetName)  row 参数要求是整数 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelReadRow(path, row, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	data, err := excel.ReadRow(path, sheetName, int(row))
	if err != nil {
		fmt.Println("[Err]读取行失败 err = ", err.Error())
	}

	return data, err
}

func excelWriteRow(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelWriteRow(path, row, list, 可选参数sheetName) 需要三个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelWriteRow(path, row, list, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	row, rowOK := args[1].(int64)
	if !rowOK {
		return nil, fmt.Errorf("ExcelWriteRow(path, row, list, 可选参数sheetName)  row 参数要求是整数 ")
	}

	list := make([]string, 0)

	switch args[2].(type) {
	case []interpreter.Value:
		for _, v := range args[2].([]interpreter.Value) {
			list = append(list, gt.Any2String(v))
		}
	default:
		return nil, fmt.Errorf("ExcelWriteRow(path, row, list, 可选参数sheetName)  list 参数要求是List ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelWriteRow(path, row, list, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.WriteRow(path, sheetName, int(row), list)
	if err != nil {
		fmt.Println("[Err]写入行失败 err = ", err.Error())
	}

	fmt.Println("数据写入成功")

	return nil, err
}

func excelDeleteRow(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelDeleteRow(path, row, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelDeleteRow(path, row, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	row, rowOK := args[1].(int64)
	if !rowOK {
		return nil, fmt.Errorf("ExcelDeleteRow(path, row, 可选参数sheetName)  row 参数要求是整数 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelDeleteRow(path, row, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.DeleteRow(path, sheetName, int(row))
	if err != nil {
		fmt.Println("[Err]删除行失败 err = ", err.Error())
	}

	fmt.Println("删除行成功")

	return nil, err
}

func excelReadCol(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelReadCol(path, col, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelReadCol(path, col, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	col, colOK := args[1].(int64)
	if !colOK {
		return nil, fmt.Errorf("ExcelReadCol(path, col, 可选参数sheetName)  col 参数要求是整数 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelReadCol(path, col, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	data, err := excel.ReadColumn(path, sheetName, int(col))
	if err != nil {
		fmt.Println("[Err]读取列失败 err = ", err.Error())
	}

	return data, err
}

func excelWriteCol(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelWriteCol(path, col, list, 可选参数sheetName) 需要三个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelWriteCol(path, col, list, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	col, colOK := args[1].(int64)
	if !colOK {
		return nil, fmt.Errorf("ExcelWriteCol(path, col, list, 可选参数sheetName)  col 参数要求是整数 ")
	}

	list := make([]string, 0)

	switch args[2].(type) {
	case []interpreter.Value:
		for _, v := range args[2].([]interpreter.Value) {
			list = append(list, gt.Any2String(v))
		}
	default:
		return nil, fmt.Errorf("ExcelWriteCol(path, col, list, 可选参数sheetName)  list 参数要求是List ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelWriteCol(path, col, list, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.WriteColumn(path, sheetName, int(col), list)
	if err != nil {
		fmt.Println("[Err]写入列失败 err = ", err.Error())
	}

	fmt.Println("数据写入成功")

	return nil, err
}

func excelDeleteCol(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelDeleteCol(path, col, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelDeleteCol(path, col, 可选参数sheetName) path 参数要求是字符串 ")
	}

	col, colOK := args[1].(int64)
	if !colOK {
		return nil, fmt.Errorf("ExcelDeleteCol(path, col, 可选参数sheetName) col 参数要求是整数 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelDeleteCol(path, col, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.DeleteColumn(path, sheetName, int(col))
	if err != nil {
		fmt.Println("[Err]删除行失败 err = ", err.Error())
	}

	fmt.Println("删除行成功")

	return nil, err
}

func excelReadCell(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelReadCell(path, cell, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelReadCell(path, cell, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelReadCell(path, cell, 可选参数sheetName)  cell 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelReadCell(path, cell, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	data, err := excel.ReadColumnCell(path, sheetName, cell)
	if err != nil {
		fmt.Println("[Err]读取列失败 err = ", err.Error())
	}

	return data, err
}

func excelWriteCell(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelWriteCell(path, cell, list, 可选参数sheetName) 需要三个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelWriteCell(path, cell, list, 可选参数sheetName)  path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelWriteCell(path, cell, list, 可选参数sheetName)  cell 参数要求是字符串 ")
	}

	list := make([]string, 0)

	switch args[2].(type) {
	case []interpreter.Value:
		for _, v := range args[2].([]interpreter.Value) {
			list = append(list, gt.Any2String(v))
		}
	default:
		return nil, fmt.Errorf("ExcelWriteCell(path, cell, list, 可选参数sheetName) list 参数要求是List ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelWriteCell(path, cell, list, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.WriteColumnCell(path, sheetName, cell, list)
	if err != nil {
		fmt.Println("[Err]写入列失败 err = ", err.Error())
	}

	fmt.Println("数据写入成功")

	return nil, err
}

func excelDeleteCell(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ExcelDeleteCell(path, cell, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelDeleteCell(path, cell, 可选参数sheetName) path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelDeleteCell(path, cell, 可选参数sheetName) cell 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 3 {
		sheetName, sheetNameOK = args[2].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelDeleteCell(path, cell, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.DeleteColumnCell(path, sheetName, cell)
	if err != nil {
		fmt.Println("[Err]删除行失败 err = ", err.Error())
	}

	fmt.Println("删除行成功")

	return nil, err
}

func excelImg(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) cell 参数要求是字符串 ")
	}

	imgPath, imgPathOK := args[2].(string)
	if !imgPathOK {
		return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) imgPath 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.InsertImage(path, sheetName, cell, imgPath)
	if err != nil {
		fmt.Println("[Err]单元格插入图片失败 err = ", err.Error())
		return nil, err
	}
	fmt.Println("插入图片成功")
	return nil, nil
}

func excelCellStyle(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelCellStyle(path, cell, style, 可选参数sheetName) 需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelCellStyle(path, cell, style, 可选参数sheetName) path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelCellStyle(path, cell, style, 可选参数sheetName) cell 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	style := args[2]
	styleDict := style.(interpreter.DictType)
	fontBold := styleDict["fontBold"].(bool)
	fontColor := gt.Any2String(styleDict["fontColor"])
	bgColor := gt.Any2String(styleDict["bgColor"])
	alignCenter := styleDict["alignCenter"].(bool)

	err := excel.SetCellStyle(path, sheetName, cell, fontBold, fontColor, bgColor, alignCenter)
	if err != nil {
		fmt.Println("[Err]设置样式失败 err = ", err.Error())
		return nil, err
	}
	fmt.Println("设置样式成功")
	return nil, nil
}

func excelMergeCells(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelMergeCells(path, startCell, endCell, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelMergeCells(path, startCell, endCell, 可选参数sheetName) path 参数要求是字符串 ")
	}

	startCell, startCellOK := args[1].(string)
	if !startCellOK {
		return nil, fmt.Errorf("ExcelMergeCells(path, startCell, endCell, 可选参数sheetName) startCell 参数要求是字符串 ")
	}

	endCell, endCellOK := args[2].(string)
	if !endCellOK {
		return nil, fmt.Errorf("ExcelMergeCells(path, startCell, endCell, 可选参数sheetName) endCell 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.MergeCells(path, sheetName, startCell, endCell)
	if err != nil {
		fmt.Println("[Err]合并单元格失败 err = ", err.Error())
		return nil, err
	}
	fmt.Println("合并单元格成功")
	return nil, nil
}

func excelSetFormula(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ExcelSetFormula(path, cell, formula, 可选参数sheetName)  需要两个参数")
	}

	path, pathOK := args[0].(string)
	if !pathOK {
		return nil, fmt.Errorf("ExcelSetFormula(path, cell, formula, 可选参数sheetName) path 参数要求是字符串 ")
	}

	cell, cellOK := args[1].(string)
	if !cellOK {
		return nil, fmt.Errorf("ExcelSetFormula(path, cell, formula, 可选参数sheetName) cell 参数要求是字符串 ")
	}

	formula, formulaOK := args[2].(string)
	if !formulaOK {
		return nil, fmt.Errorf("ExcelSetFormula(path, cell, formula, 可选参数sheetName) formula 参数要求是字符串 ")
	}

	sheetName := ""
	sheetNameOK := false
	if len(args) == 4 {
		sheetName, sheetNameOK = args[3].(string)
		if !sheetNameOK {
			return nil, fmt.Errorf("ExcelImg(path, cell, imgPath, 可选参数sheetName) 可选参数 sheetName 参数要求是字符串 ")
		}
	}

	err := excel.SetCellFormula(path, sheetName, cell, formula)
	if err != nil {
		fmt.Println("[Err]合并单元格失败 err = ", err.Error())
		return nil, err
	}
	fmt.Println("合并单元格成功")
	return nil, nil
}
