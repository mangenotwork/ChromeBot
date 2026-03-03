package browser

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func FindChrome() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("only Windows is supported")
	}

	// 方法1: 尝试常见路径
	if path, err := findChromeByCommonPaths(); err == nil {
		return path, nil
	}

	// 方法3: 尝试PowerShell WMI查询
	if path, err := findChromeByWMIReliable(); err == nil {
		return path, nil
	}

	// 方法4: 尝试安全的WMIC查询
	if path, err := findChromeByWMICSafe(); err == nil {
		return path, nil
	}

	// 方法5: 尝试原始WMI查询
	if path, err := findChromeByWMI(); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Chrome not found on this system")
}

// 优化的WMI查找方法
func findChromeByWMI() (string, error) {
	// 尝试多种WMI查询
	queries := []struct {
		wmicArgs []string
		property string
	}{
		// 方法1: 查询安装的程序
		{[]string{"product", "where", "name like '%Chrome%'", "get", "InstallLocation"}, "InstallLocation"},
		// 方法2: 查询快捷方式
		{[]string{"process", "where", "name='chrome.exe'", "get", "executablepath"}, "ExecutablePath"},
		// 方法3: 查询所有文件
		{[]string{"datafile", "where", "filename='chrome' and extension='exe'", "get", "name"}, "Name"},
		// 方法4: 查询Win32_Product
		{[]string{"path", "Win32_Product", "where", "Name like '%Chrome%'", "get", "InstallLocation"}, "InstallLocation"},
		// 方法5: 查询已安装的软件
		{[]string{"path", "Win32_InstalledWin32Program", "where", "Name like '%Chrome%'", "get", "InstallLocation"}, "InstallLocation"},
	}

	for _, query := range queries {
		path, err := executeWMIQuery(query.wmicArgs, query.property)
		if err == nil && path != "" {
			// 验证路径
			if validateChromePath(path) {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("Chrome not found via WMI")
}

// 执行WMI查询
func executeWMIQuery(args []string, property string) (string, error) {
	// 构建命令
	cmdArgs := append([]string{"/c", "wmic"}, args...)
	cmd := exec.Command("cmd", cmdArgs...)

	// 设置正确的编码
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=UTF-8")

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// 解码输出（处理可能的UTF-16编码）
	outputStr := decodeWMIOutput(output)

	// 解析输出
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.EqualFold(line, property) {
			continue
		}

		// 尝试提取路径
		if path := extractPathFromLine(line); path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("property not found in output")
}

// 解码WMI输出（处理UTF-16和ANSI编码）
func decodeWMIOutput(data []byte) string {
	// 检查是否是UTF-16LE
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		// UTF-16LE BOM
		utf16Data := make([]uint16, (len(data)-2)/2)
		for i := 0; i < len(utf16Data); i++ {
			utf16Data[i] = uint16(data[2+i*2]) | uint16(data[2+i*2+1])<<8
		}
		return string(utf16.Decode(utf16Data))
	}

	// 尝试UTF-8
	if utf8.Valid(data) {
		return string(data)
	}

	// 尝试使用系统默认编码（Windows通常是GBK/GB2312）
	// 这里使用一个简单的回退策略
	return string(data)
}

// 从行中提取路径
func extractPathFromLine(line string) string {
	line = strings.TrimSpace(line)

	// 移除可能的引号
	line = strings.Trim(line, `"`)

	// 检查是否已经是完整路径
	if strings.HasSuffix(strings.ToLower(line), "chrome.exe") {
		return line
	}

	// 如果只是目录，添加chrome.exe
	if line != "" {
		// 尝试几种可能的组合
		possiblePaths := []string{
			filepath.Join(line, "chrome.exe"),
			filepath.Join(line, "Application", "chrome.exe"),
			filepath.Join(line, "Google", "Chrome", "Application", "chrome.exe"),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// 验证Chrome路径
func validateChromePath(path string) bool {
	if path == "" {
		return false
	}

	// 检查文件是否存在
	if _, err := os.Stat(path); err != nil {
		return false
	}

	// 验证文件名
	baseName := strings.ToLower(filepath.Base(path))
	if baseName != "chrome.exe" && baseName != "msedge.exe" {
		return false
	}

	return true
}

// 更可靠的WMI查询方法
func findChromeByWMIReliable() (string, error) {
	// 使用PowerShell进行更可靠的查询
	psScript := `
$paths = @()
# 方法1: 查询已安装程序
try {
    $chrome = Get-WmiObject -Class Win32_Product | Where-Object { $_.Name -like "*Chrome*" }
    if ($chrome) {
        $paths += $chrome.InstallLocation
    }
} catch {}

# 方法2: 查询注册表
try {
    $regPath = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe"
    if (Test-Path $regPath) {
        $paths += (Get-ItemProperty -Path $regPath).'(default)'
    }
} catch {}

# 方法3: 查询开始菜单快捷方式
try {
    $shortcut = Get-ChildItem -Path "$env:ProgramData\Microsoft\Windows\Start Menu\Programs" -Recurse -Filter "*chrome*.lnk" | Select-Object -First 1
    if ($shortcut) {
        $shell = New-Object -ComObject WScript.Shell
        $paths += $shell.CreateShortcut($shortcut.FullName).TargetPath
    }
} catch {}

# 方法4: 查询文件系统
try {
    $commonPaths = @(
        "$env:ProgramFiles\Google\Chrome\Application\chrome.exe",
        "${env:ProgramFiles(x86)}\Google\Chrome\Application\chrome.exe",
        "$env:LOCALAPPDATA\Google\Chrome\Application\chrome.exe"
    )
    foreach ($p in $commonPaths) {
        if (Test-Path $p) {
            $paths += $p
        }
    }
} catch {}

# 返回第一个有效的路径
foreach ($p in $paths) {
    if ($p -and (Test-Path $p)) {
        $p
        break
    }
}
`

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("PowerShell error: %v, stderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if output != "" && validateChromePath(output) {
		return output, nil
	}

	return "", fmt.Errorf("Chrome not found via PowerShell")
}

// 使用WMIC的更安全方法
func findChromeByWMICSafe() (string, error) {
	// 使用更精确的查询条件
	queries := []string{
		// 精确匹配Google Chrome
		`wmic product where "name='Google Chrome'" get InstallLocation`,
		// 查询快捷方式
		`wmic path Win32_ShortcutFile where "name like '%chrome%.lnk'" get Target`,
		// 查询进程（如果Chrome正在运行）
		`wmic process where "name='chrome.exe'" get executablepath`,
	}

	for _, query := range queries {
		// 执行查询
		cmd := exec.Command("cmd", "/c", query)
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		// 解析输出
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.Contains(strings.ToLower(line), "installlocation") ||
				strings.Contains(strings.ToLower(line), "executablepath") ||
				strings.Contains(strings.ToLower(line), "target") {
				continue
			}

			// 清理路径
			path := cleanWMIPath(line)
			if path != "" && validateChromePath(path) {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("Chrome not found via safe WMI queries")
}

// 清理WMI返回的路径
func cleanWMIPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "\"")
	path = strings.Trim(path, "'")
	path = strings.TrimSpace(path)

	// 移除可能的Windows路径前缀
	path = strings.ReplaceAll(path, "\\??\\", "")
	path = strings.ReplaceAll(path, "\\\\?\\", "")

	return path
}

// 查找常见路径
func findChromeByCommonPaths() (string, error) {
	paths := []string{
		filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("ProgramW6432"), "Google", "Chrome", "Application", "chrome.exe"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("not found in common paths")
}
