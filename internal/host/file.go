package host

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"
)

// 返回 1是文件，2是目录
func CheckPath(pathStr string) (string, int) {
	if pathStr == "" {
		return utils.ScriptDir, 2
	}
	if !filepath.IsAbs(pathStr) {
		pathStr = filepath.Join(utils.ScriptDir, pathStr)
	}
	pathStr = filepath.Clean(pathStr)
	stat, err := os.Stat(pathStr)
	if err != nil {
		return pathStr, 2
	}
	if stat.IsDir() {
		return pathStr, 2
	}
	return pathStr, 1
}

func LS(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "名称\t大小\t修改时间\t类型")
	fmt.Fprintln(w, "----\t----\t--------\t----")

	var totalSize int64

	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		var size int64

		if file.IsDir() {
			size = calcDirSizeParallel(fullPath) // 🔥 并行计算
		} else {
			size = file.Size()
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			file.Name(),
			formatSize(size),
			file.ModTime().Format("2006-01-02 15:04:05"),
			dirOrFile(file.IsDir()),
		)
		totalSize += size
	}

	w.Flush()
	fmt.Println("\n✅ 目录总大小：", formatSize(totalSize))
	return nil
}

var (
	lsSizeChan = make(chan int64, 1000)
	lsWg       sync.WaitGroup
	lsSema     = make(chan struct{}, 50) // 限制并发数，防止卡死
)

// 🔥 并行计算目录大小（超快）
func calcDirSizeParallel(root string) int64 {
	var total int64
	lsWg.Add(1)

	// 并发扫描
	go func() {
		defer lsWg.Done()
		scanDir(root)
	}()

	// 关闭通道并等待
	go func() {
		lsWg.Wait()
		close(lsSizeChan)
	}()

	// 汇总大小
	for s := range lsSizeChan {
		total += s
	}

	// 重置
	lsSizeChan = make(chan int64, 1000)
	return total
}

// 🔥 递归扫描（带并发限流，极快）
func scanDir(path string) {
	lsSema <- struct{}{}        //  acquire
	defer func() { <-lsSema }() // release

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			lsWg.Add(1)
			go func(p string) {
				defer lsWg.Done()
				scanDir(p)
			}(fullPath)
			continue
		}

		// 文件加入大小
		info, err := entry.Info()
		if err == nil {
			lsSizeChan <- info.Size()
		}
	}
}

func dirOrFile(isDir bool) string {
	if isDir {
		return "目录"
	}
	return "文件"
}

func formatSize(s int64) string {
	const unit = 1024
	if s < unit {
		return fmt.Sprintf("%d B", s)
	}
	div, exp := int64(unit), 0
	for n := s / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %sB", float64(s)/float64(div), []string{"K", "M", "G", "T"}[exp])
}

func GetFileInfo(filePath string) (interpreter.DictType, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileInfo := interpreter.DictType{
		"name":    stat.Name(),
		"path":    filePath,
		"size":    stat.Size(),
		"sizeStr": formatSize(stat.Size()),
		"modtime": stat.ModTime().Format("2006-01-02 15:04:05"),
		"isdir":   stat.IsDir(),
		"mode":    stat.Mode().String(),
	}
	return fileInfo, nil
}

func GetDirInfo(dirPath string) (interpreter.DictType, error) {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	fileCount := 0
	dirCount := 0
	entries, err := os.ReadDir(dirPath)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				dirCount++
			} else {
				fileCount++
			}
		}
	}

	dirInfo := interpreter.DictType{
		"name":      stat.Name(),
		"path":      dirPath,
		"modtime":   stat.ModTime().Format("2006-01-02 15:04:05"),
		"isdir":     stat.IsDir(),
		"mode":      stat.Mode().String(),
		"dircount":  dirCount,
		"filecount": fileCount,
		"count":     fileCount + dirCount,
	}
	return dirInfo, nil
}

// SearchResult 搜索结果结构体
type SearchResult struct {
	Path  string `json:"path"`  // 完整路径
	Name  string `json:"name"`  // 名称
	IsDir bool   `json:"isDir"` // 是否是目录
}

// SearchFilesDir
// root: 搜索根目录
// keyword: 搜索关键词
// recursive: 是否递归搜索子目录
// ignoreCase: 是否忽略大小写
func SearchFilesDir(root string, keyword string, recursive, ignoreCase bool) ([]SearchResult, error) {
	var results []SearchResult

	// 处理关键词
	kw := keyword
	if ignoreCase {
		kw = strings.ToLower(keyword)
	}

	// 定义遍历函数
	var walk func(string) error
	walk = func(dir string) error {
		// 读取目录（快速）
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			name := entry.Name()
			fullPath := filepath.Join(dir, name)

			// 匹配名称
			matchName := name
			if ignoreCase {
				matchName = strings.ToLower(name)
			}

			// 关键词匹配
			if strings.Contains(matchName, kw) {
				results = append(results, SearchResult{
					Path:  fullPath,
					Name:  name,
					IsDir: entry.IsDir(),
				})
			}

			// 递归搜索子目录
			if recursive && entry.IsDir() {
				if err := walk(fullPath); err != nil {
					continue // 出错跳过
				}
			}
		}
		return nil
	}

	// 开始搜索
	err := walk(root)
	return results, err
}

// CreateDir 创建目录（支持多级目录，已存在不会报错）
func CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// CreateFile 创建文件
func CreateFile(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func DeleteFile(path string) error {
	return os.RemoveAll(path)
}

func MoveFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// Copy 通用复制函数：自动判断 复制文件 / 复制目录（递归）
func CopyFile(src, dst string) error {
	// 获取源信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 如果是目录，调用目录复制
	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	// 如果是文件，调用文件复制
	return copyFile(src, dst)
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件夹
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// copyDir 递归复制目录
func copyDir(srcDir, dstDir string) error {
	// 读取源目录所有内容
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadFileToString 读取文件内容，直接返回字符串
func ReadFileToString(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	// 字节数组 转 字符串
	return string(content), nil
}

func RenameOrMove(oldPath, newPath string) error {
	err := os.MkdirAll(filepath.Dir(newPath), 0755)
	if err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

// WriteFileOverwrite 全覆盖写入文件
func WriteFileOverwrite(filePath string, content string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(content), 0644)
}

// AppendToFile 追加写入文件
func AppendToFile(filePath string, content string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}
