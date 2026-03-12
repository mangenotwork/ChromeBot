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
