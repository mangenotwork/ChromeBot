package browser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
)

var (
	ChromeLocalRecordFilePath = ""
	fileMutex                 sync.Mutex
)

func init() {
	wd, _ := os.Getwd()
	ChromeLocalRecordFilePath = wd + "\\profiles\\record"
	// 初始化目录（如果不存在则创建）
	dir := filepath.Dir(ChromeLocalRecordFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf(fmt.Sprintf("创建\\profiles\\record目录失败: %v \n", err))
	}
}

// 打开文件
func openLocalRecord() map[string]int {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// 第二步：尝试打开文件（不存在则返回空map）
	file, err := os.OpenFile(ChromeLocalRecordFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("打开记录文件失败: %v\n", err)
		return make(map[string]int)
	}
	defer file.Close()

	// 第四步：读取文件内容
	var data map[string]int
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		// 文件为空/格式错误时返回空map
		if errors.Is(err, os.ErrNotExist) || err.Error() == "EOF" {
			return make(map[string]int)
		}
		fmt.Printf("解析记录文件失败: %v\n", err)
		return make(map[string]int)
	}

	return data
}

// 保存文件，全部覆盖
func saveLocalRecord(data map[string]int) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// 第二步：打开文件（截断写入，不存在则创建）
	file, err := os.OpenFile(ChromeLocalRecordFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("打开文件写入失败: %v\n", err)
		return
	}
	defer file.Close()

	// 第四步：写入JSON数据
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 格式化JSON，便于阅读
	if err := encoder.Encode(&data); err != nil {
		fmt.Printf("写入记录文件失败: %v\n", err)
	}
}

// AddLocalRecord 添加记录
func AddLocalRecord(userPath string, pid int) {
	if userPath == "" || pid <= 0 {
		fmt.Println("无效的记录参数：userPath为空或pid非法")
		return
	}
	data := openLocalRecord()
	data[userPath] = pid
	saveLocalRecord(data)
}

// HasLocalRecord 检查记录
func HasLocalRecord(userPath string) bool {
	if userPath == "" {
		return false
	}
	data := openLocalRecord()
	pid, ok := data[userPath]
	if ok {
		if isRun, _ := isProcessRunning(pid); isRun {
			return true
		} else {
			delete(data, userPath)
			saveLocalRecord(data)
		}
	}
	return false
}

// countDirectSubDirs 统计指定目录下的直接子目录数量
// 参数：
//
//	dirPath: 目标目录路径
//	skipHidden: 是否跳过隐藏目录（Windows/Linux/macOS通用）
//
// 返回值：
//
//	子目录数量 / 错误信息
func countDirectSubDirs(dirPath string, skipHidden bool) (int, error) {
	// 第一步：校验目录是否存在
	info, err := os.Stat(dirPath)
	if err != nil {
		return 0, fmt.Errorf("目录不存在或无访问权限: %w", err)
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("指定路径不是目录: %s", dirPath)
	}

	// 第二步：打开目录并读取所有条目
	dir, err := os.Open(dirPath)
	if err != nil {
		return 0, fmt.Errorf("打开目录失败: %w", err)
	}
	defer dir.Close()

	// 第三步：读取目录下的所有直接条目（不递归）
	entries, err := dir.ReadDir(-1) // -1 表示读取所有条目
	if err != nil {
		return 0, fmt.Errorf("读取目录条目失败: %w", err)
	}

	// 第四步：统计子目录数量
	count := 0
	for _, entry := range entries {
		// 过滤非目录条目
		if !entry.IsDir() {
			continue
		}

		// 可选：跳过隐藏目录
		if skipHidden && isHiddenDir(entry, dirPath) {
			continue
		}

		count++
	}

	return count, nil
}

// isHiddenDir 判断目录是否为隐藏目录（跨平台）
func isHiddenDir(entry os.DirEntry, parentDir string) bool {
	if runtime.GOOS == "windows" {
		// Windows：获取文件属性判断是否隐藏
		fullPath := filepath.Join(parentDir, entry.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			return false
		}
		// 获取Windows文件属性
		winAttr := info.Sys().(*syscall.Win32FileAttributeData)
		return winAttr.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
	} else {
		// Linux/macOS：目录名以.开头即为隐藏
		return len(entry.Name()) > 0 && entry.Name()[0] == '.'
	}
}
