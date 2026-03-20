package global

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/utils"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

func ReadJsonToConf(path, as string) {
	absPath, _ := utils.GetAbsolutePath(path)
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

func ReadYamlToConf(path, as string) {
	fmt.Println("ReadYamlToConf ....", path, as)
	absPath, _ := utils.GetAbsolutePath(path)
	yamlFile, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("[Err]读取%s文件失败：%v \n", absPath, err)
		os.Exit(0)
	}

	var tempMap map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &tempMap)
	if err != nil {
		fmt.Printf("解析YAML失败：%v \n", err)
		os.Exit(0)
	}

	customMap := convertToDictType(tempMap)
	fmt.Printf("customMap = %v \n", customMap)

	if _, ok := interpreter.Const.Load(as); ok {
		fmt.Printf("[Wring] 全局常量%s已被定义. \n", as)
		return
	}

	interpreter.Const.Store(as, customMap)

}

func ReadINIToConf(path, as string) {
	fmt.Println("ReadINIToConf ....", path, as)
	absPath, _ := utils.GetAbsolutePath(path)
	cfg, err := ini.Load(absPath)
	if err != nil {
		fmt.Printf("[Err]读取%s文件失败：%v \n", absPath, err)
		os.Exit(0)
	}

	tempMap := make(map[string]interface{})

	globalSection := cfg.Section("")
	for _, key := range globalSection.Keys() {
		keyName := key.Name()
		value := parseValue(key.String()) // 自动解析值类型
		tempMap[keyName] = value
	}

	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		if sectionName == "" { // 跳过全局section（已处理）
			continue
		}

		// 遍历section下的所有键值对
		for _, key := range section.Keys() {
			keyName := fmt.Sprintf("%s.%s", sectionName, key.Name()) // 拼接为 section.key
			value := parseValue(key.String())                        // 自动解析值类型
			tempMap[keyName] = value
		}
	}

	customMap := convertToDictType(tempMap)
	fmt.Printf("customMap = %v \n", customMap)

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

func parseValue(valueStr string) interface{} {
	// 去除首尾空格
	valueStr = strings.TrimSpace(valueStr)

	// 空值返回空字符串
	if valueStr == "" {
		return ""
	}

	// 尝试解析为bool
	lowerVal := strings.ToLower(valueStr)
	if lowerVal == "true" || lowerVal == "false" {
		boolVal, _ := strconv.ParseBool(lowerVal)
		return boolVal
	}

	// 尝试解析为int
	intVal, err := strconv.Atoi(valueStr)
	if err == nil {
		return intVal
	}

	// 尝试解析为float64
	floatVal, err := strconv.ParseFloat(valueStr, 64)
	if err == nil {
		return floatVal
	}

	// 都解析失败则返回原始字符串
	return valueStr
}
