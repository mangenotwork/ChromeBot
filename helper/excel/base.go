package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ReadExcel(path string) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
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

// todo 将List写入excel

// todo 将字典写入excel

// todo 读取excel全部存储到List

// todo 读取excel全部存储到字典

// todo 获取excel的信息

// todo 获取Sheet的信息

// todo 读取指定单元格

// todo 写入指定单元格

// todo 删除指定单元格的数据

// todo 读取指定行的数据

// todo 读取指定列的数据

// todo 写入指定行

// todo 写入指定列

// todo 删除指定行

// todo 删除指定列

// todo [] 4. 图片插入

// todo [] 5. 单元格样式设置（字体、颜色、对齐）

// todo [] 6. 合并单元格

// todo [] 7. 添加图表

// todo [] 8. 设置单元格公式

// todo 写sql操作excel
