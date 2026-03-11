package host

import (
	"ChromeBot/utils"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

/*

设计:

1. SystemConfirmBox  确认弹框，固定的 y/n 返回true/false
2. DialogBox 自定义弹框
3. 终止脚本弹窗，显示报错信息
4. 提示弹框，主要是告知提示，点击继续

*/

const (
	MB_YESNO             = 0x00000004 // 弹窗类型：确认+取消按钮（系统级确认框核心）
	MB_ICONQUESTION      = 0x00000020 // 弹窗图标：问号（增强确认感）
	IDYES                = 6          // 返回值：用户点击了"是"
	IDNO                 = 7          // 返回值：用户点击了"否"
	IDCANCEL             = 2          // 取消按钮
	IDABORT              = 3          // 中止按钮
	IDRETRY              = 4          // 重试按钮
	IDIGNORE             = 5          // 忽略按钮
	IDTRYAGAIN           = 10         // 重试按钮（MB_CANCELTRYCONTINUE）
	IDCONTINUE           = 11         // 继续按钮（MB_CANCELTRYCONTINUE）
	HWND_DESKTOP         = 0          // 弹窗归属：桌面窗口（确保是系统级别，不依赖当前进程窗口）
	MB_OK                = 0x00000000 // 只有一个"确定"按钮
	MB_YESNOCANCEL       = 0x00000003 // 是+否+取消按钮
	MB_RETRYCANCEL       = 0x00000005 // 重试+取消按钮
	MB_ABORTRETRYIGNORE  = 0x00000002 // 中止+重试+忽略按钮
	MB_CANCELTRYCONTINUE = 0x00000006 // 取消+重试+继续按钮

	// 其他选项
	MB_DEFBUTTON1 = 0x00000000 // 第一个按钮为默认
	MB_DEFBUTTON2 = 0x00000100 // 第二个按钮为默认
	MB_DEFBUTTON3 = 0x00000200 // 第三个按钮为默认
	MB_DEFBUTTON4 = 0x00000300 // 第四个按钮为默认

	// 图标类型
	MB_ICONEXCLAMATION = 0x00000030 // 感叹号图标
	MB_ICONWARNING     = 0x00000030 // 警告图标（同感叹号）
	MB_ICONINFORMATION = 0x00000040 // 信息图标（i）
	MB_ICONASTERISK    = 0x00000040 // 信息图标（同i）
	MB_ICONSTOP        = 0x00000010 // 停止图标
	MB_ICONERROR       = 0x00000010 // 错误图标（同停止）
	MB_ICONHAND        = 0x00000010 // 错误图标（同停止）
)

var (
	messageBoxW = user32.NewProc("MessageBoxW")
)

// SystemConfirmBox 系统级确认框封装函数
// title: 弹窗标题
// message: 弹窗内容
// return: true(用户点是)/false(用户点否)，error(调用API失败)
func SystemConfirmBox(title, message string) (bool, error) {
	// 将 Go 字符串转为 Windows 要求的 UTF-16 编码
	titleUTF16, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return false, err
	}
	messageUTF16, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		return false, err
	}

	// 调用 Windows API: MessageBoxW(hwnd, text, caption, type)
	ret, _, err := messageBoxW.Call(
		uintptr(HWND_DESKTOP),                 // 父窗口：桌面（系统级别）
		uintptr(unsafe.Pointer(messageUTF16)), // 弹窗内容
		uintptr(unsafe.Pointer(titleUTF16)),   // 弹窗标题
		uintptr(MB_YESNO|MB_ICONQUESTION),     // 弹窗样式：确认+取消 + 问号图标
	)

	// 处理返回值
	switch ret {
	case IDYES:
		fmt.Println("点击了是")
		return true, nil
	case IDNO:
		fmt.Println("点击了否")
		return false, nil
	default:
		return false, err // API 调用失败（如权限问题）
	}
}

func SystemExitBox() error {
	// 将 Go 字符串转为 Windows 要求的 UTF-16 编码
	titleUTF16, err := syscall.UTF16PtrFromString("是否终止ChromeBot?")
	if err != nil {
		return err
	}
	messageUTF16, err := syscall.UTF16PtrFromString("点击是会终止当前ChromeBot进程，否则继续执行。")
	if err != nil {
		return err
	}

	// 调用 Windows API: MessageBoxW(hwnd, text, caption, type)
	ret, _, err := messageBoxW.Call(
		uintptr(HWND_DESKTOP),                 // 父窗口：桌面（系统级别）
		uintptr(unsafe.Pointer(messageUTF16)), // 弹窗内容
		uintptr(unsafe.Pointer(titleUTF16)),   // 弹窗标题
		uintptr(MB_YESNO|MB_ICONQUESTION),     // 弹窗样式：确认+取消 + 问号图标
	)

	if ret == IDYES {
		fmt.Println("点击了是")
		if utils.RunMode == "REPL" {
			utils.SigChan <- syscall.SIGTERM
		} else {
			fmt.Println("终止脚本")
			os.Exit(0)
		}
	}

	return nil
}

func MessageBox(title, message string, flags uint) (int, error) {
	titleUTF16, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}

	messageUTF16, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		return 0, err
	}

	ret, _, err := messageBoxW.Call(
		uintptr(HWND_DESKTOP),
		uintptr(unsafe.Pointer(messageUTF16)),
		uintptr(unsafe.Pointer(titleUTF16)),
		uintptr(flags),
	)

	return int(ret), err
}

func TipBox(title, message string, iconType ...uint) (bool, error) {
	// 默认使用信息图标
	icon := iconType[0]

	// 显示只有一个"确定"按钮的提示框
	_, err := MessageBox(title, message, MB_OK|uint(icon))
	if err != nil {
		return false, err
	}

	// 用户点击了确定/继续按钮
	return true, nil
}

// 预定义的各种提示框函数

// InfoTipBox 信息提示框
func InfoTipBox(message string) (bool, error) {
	return TipBox("信息提示", message, MB_ICONINFORMATION)
}

// WarningTipBox 警告提示框
func WarningTipBox(message string) (bool, error) {
	return TipBox("警告提示", message, MB_ICONWARNING)
}

// ErrorTipBox 错误提示框
func ErrorTipBox(message string) (bool, error) {
	return TipBox("错误提示", message, MB_ICONERROR)
}

// SuccessTipBox 成功提示框
func SuccessTipBox(message string) (bool, error) {
	return TipBox("成功提示", message, MB_ICONINFORMATION)
}

// ThreeButtonResult 三按钮对话框的返回结果
type ThreeButtonResult int

const (
	ButtonFirst  ThreeButtonResult = 1 // 第一个按钮
	ButtonSecond ThreeButtonResult = 2 // 第二个按钮
	ButtonThird  ThreeButtonResult = 3 // 第三个按钮
	ButtonCancel ThreeButtonResult = 0 // 取消/关闭
)

// ThreeButtonDialog 完整的三个按钮对话框配置
type ThreeButtonDialog struct {
	Title         string
	Message       string
	Button1       string
	Button2       string
	Button3       string
	BoxType       uint // MB_YESNOCANCEL, MB_ABORTRETRYIGNORE, MB_CANCELTRYCONTINUE
	IconType      uint
	DefaultButton int // 1, 2, 3
}

// NewThreeButtonDialog 创建三按钮对话框配置
func NewThreeButtonDialog(title, message string) *ThreeButtonDialog {
	return &ThreeButtonDialog{
		Title:         title,
		Message:       message,
		Button1:       "按钮1",
		Button2:       "按钮2",
		Button3:       "按钮3",
		BoxType:       MB_YESNOCANCEL,
		IconType:      MB_ICONINFORMATION,
		DefaultButton: 1,
	}
}

// SetButtons 设置按钮文本
func (d *ThreeButtonDialog) SetButtons(btn1, btn2, btn3 string) *ThreeButtonDialog {
	d.Button1 = btn1
	d.Button2 = btn2
	d.Button3 = btn3
	return d
}

// SetBoxType 设置对话框类型
func (d *ThreeButtonDialog) SetBoxType(boxType uint) *ThreeButtonDialog {
	d.BoxType = boxType
	return d
}

// SetIcon 设置图标
func (d *ThreeButtonDialog) SetIcon(iconType uint) *ThreeButtonDialog {
	d.IconType = iconType
	return d
}

// SetDefaultButton 设置默认按钮
func (d *ThreeButtonDialog) SetDefaultButton(button int) *ThreeButtonDialog {
	if button >= 1 && button <= 3 {
		d.DefaultButton = button
	}
	return d
}

// Show 显示对话框
func (d *ThreeButtonDialog) Show() (ThreeButtonResult, error) {
	// 添加按钮描述
	fullMessage := d.Message + "\n\n请选择:\n" +
		"1. " + d.Button1 + "\n" +
		"2. " + d.Button2 + "\n" +
		"3. " + d.Button3

	// 设置默认按钮
	var defaultButton uint
	switch d.DefaultButton {
	case 1:
		defaultButton = 0x00000000 // MB_DEFBUTTON1
	case 2:
		defaultButton = 0x00000100 // MB_DEFBUTTON2
	case 3:
		defaultButton = 0x00000200 // MB_DEFBUTTON3
	}

	flags := d.BoxType | d.IconType | defaultButton

	ret, err := MessageBox(d.Title, fullMessage, flags)
	if err != nil {
		return ButtonCancel, err
	}

	// 根据对话框类型映射结果
	switch d.BoxType {
	case MB_YESNOCANCEL:
		switch ret {
		case IDYES:
			return ButtonFirst, nil
		case IDNO:
			return ButtonSecond, nil
		case IDCANCEL:
			return ButtonThird, nil
		}

	case MB_ABORTRETRYIGNORE:
		switch ret {
		case IDABORT:
			return ButtonFirst, nil
		case IDRETRY:
			return ButtonSecond, nil
		case IDIGNORE:
			return ButtonThird, nil
		}

	case MB_CANCELTRYCONTINUE:
		switch ret {
		case IDCANCEL:
			return ButtonFirst, nil
		case IDTRYAGAIN:
			return ButtonSecond, nil
		case IDCONTINUE:
			return ButtonThird, nil
		}
	}

	return ButtonCancel, nil
}
