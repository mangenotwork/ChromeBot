package utils

import (
	"testing"
)

func TestValidateXPathPureNative(t *testing.T) {
	// 测试用例覆盖常见场景
	testCases := []string{
		"//div[@class='content']",            // 有效
		"//a[@href='/test' and @name='btn']", // 有效
		"//div[contains(text(), 'test')]",    // 无效：包含无效空括号模式: ()
		"//div[",                             // 无效：括号未闭合
		"//div[@class=content]",              // 语法有效（引号缺失属于逻辑错误，非语法错误）
		"//div[@]",                           // 无效：@后无属性名
		"//div[]",                            // 无效：空括号
		"//div[@123='test']",                 // 无效：属性名首字符非法
		"//div[@class='test]",                // 无效：未闭合单引号
		"",                                   // 无效：空表达式
	}

	// 执行测试并输出结果
	for _, expr := range testCases {
		valid, err := ValidateXPathPureNative(expr)
		if valid {
			t.Logf("✅ XPath [%s] 语法有效\n", expr)
		} else {
			t.Logf("❌ XPath [%s] 语法无效: %v\n", expr, err)
		}
	}
}
