package global

import (
	"ChromeBot/dsl/interpreter"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func ReadJsonToConf(path, as string) {
	absPath, _ := getAbsolutePath(path)
	jsonFile, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("[Err]读取%s文件失败：%v \n", absPath, err)
		os.Exit(0)
	}
	var tempMap map[string]interface{}
	err = json.Unmarshal(jsonFile, &tempMap)
	if err != nil {
		fmt.Printf("解析JOSN失败：%v \n", err)
		os.Exit(0)
	}
	customMap := convertToDictType(tempMap)

	if _, ok := interpreter.Const.Load(as); ok {
		fmt.Printf("[Wring] 全局常量%s已被定义. \n", as)
		return
	}

	interpreter.Const.Store(as, customMap)

}

func convertToDictType(data interface{}) interpreter.Value {
	switch v := data.(type) {
	// 嵌套 Map → 递归转换为 DictType
	case map[string]interface{}:
		customMap := make(interpreter.DictType)
		for k, val := range v {
			// 键：string → Value；值：递归转换
			customMap[interpreter.Value(k)] = convertToDictType(val)
		}
		return customMap

	// 数组 []interface{} → 转换为 []Value
	case []interface{}:
		customSlice := make([]interpreter.Value, 0)
		for _, val := range v {
			// 数组元素递归转换
			customSlice = append(customSlice, convertToDictType(val))
		}
		return customSlice

	case string, int, int64, float64, bool, nil:
		return interpreter.Value(v)

	// 其他未知类型（兜底）
	default:
		fmt.Printf("警告：未知类型 %T，值：%v，已转为空值\n", v, v)
		return interpreter.Value(v)
	}
}

func ReadYamlToConf(path, as string) {
	fmt.Println("ReadYamlToConf ....", path, as)
	absPath, _ := getAbsolutePath(path)
	yamlFile, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("[Err]读取%s文件失败：%v \n", absPath, err)
		os.Exit(0)
	}

	// 2. 定义自定义Map（支持任意嵌套结构）
	var customMap interpreter.DictType

	// 3. 解析YAML内容到Map
	err = yaml.Unmarshal(yamlFile, &customMap)
	if err != nil {
		fmt.Printf("解析YAML失败：%v \n", err)
		os.Exit(0)
	}

	fmt.Printf("customMap = %v \n", customMap)

	if _, ok := interpreter.Const.Load(as); ok {
		fmt.Printf("[Wring] 全局常量%s已被定义. \n", as)
		return
	}

	interpreter.Const.Store(as, customMap)

}

func ReadINIToConf(path, as string) {
	fmt.Println("ReadINIToConf ....", path, as)

}

func cleanPath(path string) string {
	// 移除首尾的引号
	path = strings.Trim(path, `"'`)

	// 如果是Windows路径，移除额外的转义
	if filepath.Separator == '\\' { // Windows
		path = strings.ReplaceAll(path, `\"`, `"`)
		path = strings.ReplaceAll(path, `\\`, `\`)
	}

	return path
}

func getAbsolutePath(path string) (string, error) {
	// 清理路径
	clean := cleanPath(path)

	// 转换为绝对路径
	return filepath.Abs(clean)
}
