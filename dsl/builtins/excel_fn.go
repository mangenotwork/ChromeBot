package builtins

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/helper/excel"
	"fmt"

	gt "github.com/mangenotwork/gathertool"
)

var excelFn = map[string]interpreter.Function{
	"ExcelSave":      excelSave,      // ExcelSave(arg, path, 可选参数sheetName) 将变量保存到excel
	"ExcelReadList":  excelReadList,  // ExcelReadList(path, 可选参数sheetName) 读取excel返回二维列表
	"ExcelReadDict":  excelReadDict,  // ExcelReadDict(path, 可选参数sheetName) 读取excel返回字典
	"ExcelShow":      excelShow,      // ExcelShow(path, 可选参数sheetName) 显示excel
	"ExcelInfo":      excelInfo,      // ExcelInfo(path) 获取excel信息
	"ExcelSheetInfo": excelSheetInfo, // ExcelSheetInfo(path, sheetName) 获取excel的sheet信息
	"ExcelSheet":     excelSheet,     // ExcelSheet(path) 获取excel的sheet信息
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
