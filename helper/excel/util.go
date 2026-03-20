package excel

import (
	"ChromeBot/utils"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// 获取目标sheet名称（无传入则取第一个），
// 如果传入的sheet没有则创建 sheet
func getTargetSheetName(f *excelize.File, sheetName string) string {
	if sheetName == "" {
		sheetList := f.GetSheetList()
		if len(sheetList) > 0 {
			sheetName = sheetList[0]
		} else {
			// 无任何sheet时，创建默认Sheet1
			sheetName = "Sheet1"
			f.NewSheet(sheetName)
		}
		return sheetName
	}

	sheetList := f.GetSheetList()
	for _, s := range sheetList {
		if s == sheetName {
			return sheetName
		}
	}

	f.NewSheet(sheetName)
	return sheetName
}

func saveAs(f *excelize.File, path string) error {
	f.SetDocProps(&excelize.DocProperties{
		Subject:     "ChromeBot",
		Title:       "ChromeBot",
		Creator:     "ChromeBot(https://github.com/mangenotwork/ChromeBot)",
		Description: "由 ChromeBot 生成的Excel文件, https://github.com/mangenotwork/ChromeBot",
	})
	path, _ = utils.GetAbsolutePath(path)
	fmt.Println("excel save path = ", path)
	err := f.SaveAs(path)
	if err != nil && err.Error() == "unsupported workbook file format" {
		return fmt.Errorf("文件格式不支持：仅支持 .xlsx/.xlsm，不支持 .xls 或非 Excel 文件")
	}
	return err
}
