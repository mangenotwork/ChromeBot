package utils

import (
	"testing"
)

func TestProcessCommandLine(t *testing.T) {
	// 测试用例（模拟命令行长命令）
	testInput := `docker run \
  --name my-container \
  -p 8080:80 \
  -v /host/path:/container/path \
  nginx:latest
echo "hello world" \
  && ls -l \
  && pwd
chrome init \
new
`

	// 处理字符串
	result := ProcessCommandLine(testInput)

	// 输出结果
	t.Log("处理前：")
	t.Log(testInput)
	t.Log("\n处理后：")
	t.Log(result)
}

func TestProcessArgs(t *testing.T) {
	args1 := []string{"click=MatchDemoContentOP", "(", "a", "百度一下", ")"}
	result1 := ProcessArgs(args1)
	t.Log("测试用例1结果：", result1, " | len = ", len(result1))

	// 测试用例2：多个括号场景
	args2 := []string{"func1", "(", "1", "2", ")", "func2", "(", "a", "b", "c", ")"}
	result2 := ProcessArgs(args2)
	t.Log("测试用例2结果：", result2, " | len = ", len(result2))

	// 测试用例3：无匹配右括号场景
	args3 := []string{"test", "(", "x", "y"}
	result3 := ProcessArgs(args3)
	t.Log("测试用例3结果：", result3, " | len = ", len(result3))

	// 测试用例4：空数组
	args4 := []string{}
	result4 := ProcessArgs(args4)
	t.Log("测试用例4结果：", result4, " | len = ", len(result4))

	// 测试用例5：一个元素
	args5 := []string{"init"}
	result5 := ProcessArgs(args5)
	t.Log("测试用例5结果：", result5, " | len = ", len(result5))
}
