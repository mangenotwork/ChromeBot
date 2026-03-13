package browser

import (
	"debug/pe"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// 手动声明 Windows API 函数（解决 syscall 符号未定义问题）
var (
	// 链接 Windows version.dll
	versionDll = syscall.NewLazyDLL("version.dll")

	// GetFileVersionInfoSizeW - 获取版本信息大小（宽字符版本）
	procGetFileVersionInfoSizeW = versionDll.NewProc("GetFileVersionInfoSizeW")
	// GetFileVersionInfoW - 获取版本信息（宽字符版本）
	procGetFileVersionInfoW = versionDll.NewProc("GetFileVersionInfoW")
	// VerQueryValueW - 查询版本信息（宽字符版本）
	procVerQueryValueW = versionDll.NewProc("VerQueryValueW")
)

// VS_FIXEDFILEINFO 对应 Windows API 中的结构体
type VS_FIXEDFILEINFO struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

// GetFileVersionInfoSize 封装 Windows API：获取文件版本信息大小
func GetFileVersionInfoSize(filePath string) (uint32, error) {
	widePath, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return 0, err
	}

	var dummy uint32
	// 调用 GetFileVersionInfoSizeW
	size, _, err := procGetFileVersionInfoSizeW.Call(
		uintptr(unsafe.Pointer(widePath)),
		uintptr(unsafe.Pointer(&dummy)),
	)
	if size == 0 {
		return 0, err
	}
	return uint32(size), nil
}

// GetFileVersionInfo 封装 Windows API：获取文件版本信息字节数组
func GetFileVersionInfo(filePath string, size uint32) ([]byte, error) {
	widePath, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, size)
	// 调用 GetFileVersionInfoW
	ret, _, err := procGetFileVersionInfoW.Call(
		uintptr(unsafe.Pointer(widePath)),
		0,
		uintptr(size),
		uintptr(unsafe.Pointer(&buf[0])),
	)
	if ret == 0 {
		return nil, err
	}
	return buf, nil
}

// VerQueryValue 封装 Windows API：查询版本信息中的指定字段
func VerQueryValue(buf []byte, query string) (uintptr, uint32, error) {
	wideQuery, err := syscall.UTF16PtrFromString(query)
	if err != nil {
		return 0, 0, err
	}

	var valBuf uintptr
	var valLen uint32
	// 调用 VerQueryValueW
	ret, _, err := procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(wideQuery)),
		uintptr(unsafe.Pointer(&valBuf)),
		uintptr(unsafe.Pointer(&valLen)),
	)
	if ret == 0 {
		return 0, 0, err
	}
	return valBuf, valLen, nil
}

// GetChromeInfo 从 chrome.exe 获取浏览器核心信息
func GetChromeInfo(exePath string) (map[string]string, error) {
	// 1. 验证文件存在
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("chrome.exe 不存在: %s", exePath)
	}

	// 2. 解析PE文件（可选，仅验证是合法PE文件）
	file, err := os.Open(exePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	if _, err := pe.NewFile(file); err != nil {
		return nil, fmt.Errorf("不是合法的PE文件: %w", err)
	}

	// 3. 获取版本信息大小
	size, err := GetFileVersionInfoSize(exePath)
	if err != nil {
		return nil, fmt.Errorf("获取版本信息大小失败: %w", err)
	}

	// 4. 获取版本信息字节数组
	versionBuf, err := GetFileVersionInfo(exePath, size)
	if err != nil {
		return nil, fmt.Errorf("获取版本信息失败: %w", err)
	}

	// 5. 解析版本信息
	info := parseVersionInfo(versionBuf)
	return info, nil
}

// parseVersionInfo 解析版本信息字节数组，提取关键字段
func parseVersionInfo(buf []byte) map[string]string {
	info := make(map[string]string)

	// 语言和代码页（中文-中国），英文可改为 0x040904B0
	langCode := "080404B0"
	// 要提取的核心字段
	fields := map[string]string{
		"ProductVersion":  "产品版本",
		"FileVersion":     "文件版本",
		"ProductName":     "产品名称",
		"CompanyName":     "公司名称",
		"FileDescription": "文件描述",
		"LegalCopyright":  "版权信息",
	}

	// 解析字符串版本信息
	for field, desc := range fields {
		query := fmt.Sprintf(`\StringFileInfo\%s\%s`, langCode, field)
		valBuf, valLen, err := VerQueryValue(buf, query)
		if err != nil || valLen == 0 {
			continue
		}
		// 转换为字符串
		val := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(valBuf))[:valLen-1])
		info[desc] = val
	}

	// 解析数字格式版本信息
	valBuf, valLen, err := VerQueryValue(buf, `\`)
	if err == nil && valLen > 0 {
		fixedInfo := (*VS_FIXEDFILEINFO)(unsafe.Pointer(valBuf))
		// 文件版本（数字格式）
		fileVer := fmt.Sprintf("%d.%d.%d.%d",
			(fixedInfo.FileVersionMS>>16)&0xFFFF,
			fixedInfo.FileVersionMS&0xFFFF,
			(fixedInfo.FileVersionLS>>16)&0xFFFF,
			fixedInfo.FileVersionLS&0xFFFF,
		)
		info["数字文件版本"] = fileVer

		// 产品版本（数字格式）
		productVer := fmt.Sprintf("%d.%d.%d.%d",
			(fixedInfo.ProductVersionMS>>16)&0xFFFF,
			fixedInfo.ProductVersionMS&0xFFFF,
			(fixedInfo.ProductVersionLS>>16)&0xFFFF,
			fixedInfo.ProductVersionLS&0xFFFF,
		)
		info["数字产品版本"] = productVer
	}

	return info
}
