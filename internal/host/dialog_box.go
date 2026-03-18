package host

import (
	"fmt"
	"sort"
	"sync"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")

	// Windows API
	createWindowEx   = user32.NewProc("CreateWindowExW")
	registerClass    = user32.NewProc("RegisterClassW")
	getMessage       = user32.NewProc("GetMessageW")
	translateMessage = user32.NewProc("TranslateMessage")
	dispatchMessage  = user32.NewProc("DispatchMessageW")
	defWindowProc    = user32.NewProc("DefWindowProcW")
	showWindow       = user32.NewProc("ShowWindow")
	updateWindow     = user32.NewProc("UpdateWindow")
	destroyWindow    = user32.NewProc("DestroyWindow")
	setWindowText    = user32.NewProc("SetWindowTextW")
	getClientRect    = user32.NewProc("GetClientRect")
	moveWindow       = user32.NewProc("MoveWindow")
	getModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	sendMessage      = user32.NewProc("SendMessageW")
	loadCursor       = user32.NewProc("LoadCursorW")
	loadIcon         = user32.NewProc("LoadIconW")
	getStockObject   = gdi32.NewProc("GetStockObject")
	getClassInfo     = user32.NewProc("GetClassInfoW")
	setWindowLong    = user32.NewProc("SetWindowLongW")
	getWindowLong    = user32.NewProc("GetWindowLongW")

	setTextColor = gdi32.NewProc("SetTextColor")
	setBkColor   = gdi32.NewProc("SetBkColor")
	setBkMode    = gdi32.NewProc("SetBkMode")
)

// Windows 常量
const (
	// 窗口样式
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_POPUP            = 0x80000000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000
	WS_CHILD            = 0x40000000
	WS_VISIBLE          = 0x10000000
	WS_TABSTOP          = 0x00010000

	// 扩展样式
	WS_EX_DLGMODALFRAME = 0x00000001

	// 按钮样式
	BS_PUSHBUTTON    = 0x00000000
	BS_DEFPUSHBUTTON = 0x00000001

	// 静态文本样式
	SS_LEFT     = 0x00000000
	SS_NOPREFIX = 0x00000080
	SS_NOTIFY   = 0x00000100 // 发送单击消息给父窗口

	// 消息
	WM_CREATE   = 0x0001
	WM_DESTROY  = 0x0002
	WM_CLOSE    = 0x0010
	WM_COMMAND  = 0x0111
	WM_PAINT    = 0x000F
	WM_NCCREATE = 0x0081

	// 游标
	IDC_ARROW = 32512

	// 图标
	IDI_APPLICATION = 32512

	// 显示窗口
	SW_SHOW = 5

	// 默认对象
	WHITE_BRUSH      = 0
	DEFAULT_GUI_FONT = 17

	// 字体
	DEFAULT_CHARSET   = 1
	WM_SETFONT        = 0x0030
	WM_DPICHANGED     = 0x02E0
	WM_CTLCOLORSTATIC = 0x0138
	WM_CTLCOLORBTN    = 0x0135
	TRANSPARENT       = 1
	OPAQUE            = 2
	SS_WORDELLIPSIS   = 0x0000000C // 文本太长时显示省略号
)

// Windows 系统颜色索引常量
const (
	COLOR_SCROLLBAR           = 0  // 滚动条灰色区域
	COLOR_BACKGROUND          = 1  // Windows桌面背景
	COLOR_ACTIVECAPTION       = 2  // 活动窗口标题栏
	COLOR_INACTIVECAPTION     = 3  // 非活动窗口标题栏
	COLOR_MENU                = 4  // 菜单背景
	COLOR_WINDOW              = 5  // 窗口背景
	COLOR_WINDOWFRAME         = 6  // 窗口边框
	COLOR_MENUTEXT            = 7  // 菜单文本
	COLOR_WINDOWTEXT          = 8  // 窗口文本
	COLOR_CAPTIONTEXT         = 9  // 标题栏文本
	COLOR_ACTIVEBORDER        = 10 // 活动窗口边框
	COLOR_INACTIVEBORDER      = 11 // 非活动窗口边框
	COLOR_APPWORKSPACE        = 12 // 多文档界面(MDI)背景
	COLOR_HIGHLIGHT           = 13 // 选中项背景
	COLOR_HIGHLIGHTTEXT       = 14 // 选中项文本
	COLOR_BTNFACE             = 15 // 按钮表面
	COLOR_BTNSHADOW           = 16 // 按钮阴影
	COLOR_GRAYTEXT            = 17 // 灰色（禁用）文本
	COLOR_BTNTEXT             = 18 // 按钮文本
	COLOR_INACTIVECAPTIONTEXT = 19 // 非活动标题栏文本
	COLOR_BTNHIGHLIGHT        = 20 // 按钮高亮（3D高亮）

	// Windows 95/NT4 及以上
	COLOR_3DDKSHADOW = 21 // 3D 暗阴影
	COLOR_3DLIGHT    = 22 // 3D 亮色
	COLOR_INFOTEXT   = 23 // 工具提示文本
	COLOR_INFOBK     = 24 // 工具提示背景

	// Windows XP 及以上
	COLOR_HOTLIGHT                = 26 // 热点项颜色
	COLOR_GRADIENTACTIVECAPTION   = 27 // 活动窗口标题栏渐变颜色
	COLOR_GRADIENTINACTIVECAPTION = 28 // 非活动窗口标题栏渐变颜色

	// Windows Vista/7 及以上
	COLOR_MENUHILIGHT = 29 // 菜单高亮
	COLOR_MENUBAR     = 30 // 菜单栏
)

// Windows 索引常量
const (
	GWL_USERDATA   = uintptr(0xFFFFFFEB)
	GWL_WNDPROC    = -4
	GWL_HINSTANCE  = -6
	GWL_HWNDPARENT = -8
	GWL_ID         = -12
	GWL_STYLE      = -16
	GWL_EXSTYLE    = -20
)

// 结构体
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type MSG struct {
	HWnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct {
		X int32
		Y int32
	}
}

type WNDCLASS struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
}

type CREATESTRUCT struct {
	LpCreateParams uintptr
	HInstance      syscall.Handle
	HMenu          syscall.Handle
	HwndParent     syscall.Handle
	Cy             int32
	Cx             int32
	Y              int32
	X              int32
	Style          int32
	LpszName       *uint16
	LpszClass      *uint16
	ExStyle        uint32
	DwExStyle      uint32
}

// 对话框数据
type DialogData struct {
	Title   string
	Message string
	Buttons []struct {
		ID   int
		Text string
	}
	SelectedID int
	ResultChan chan int
	hWnd       syscall.Handle
	hInstance  syscall.Handle
	hFont      syscall.Handle // 添加字体句柄
	buttonHWND []syscall.Handle
	staticHWND syscall.Handle
}

// 全局变量
var (
	dialogMutex       sync.Mutex
	dialogs           = make(map[syscall.Handle]*DialogData)
	isClassRegistered bool
	className         = "CustomDialogClass"
	hInstance         syscall.Handle
)

// 初始化 - 获取实例句柄
func init() {
	hInst, _, _ := getModuleHandle.Call(0)
	hInstance = syscall.Handle(hInst)
}

// 字符串转UTF16指针
func stringToUTF16Ptr(s string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return ptr
}

// 注册窗口类
func registerWindowClass() error {
	dialogMutex.Lock()
	defer dialogMutex.Unlock()

	if isClassRegistered {
		return nil
	}

	// 检查类是否已经注册
	var wc WNDCLASS
	ret, _, _ := getClassInfo.Call(
		uintptr(hInstance),
		uintptr(unsafe.Pointer(stringToUTF16Ptr(className))),
		uintptr(unsafe.Pointer(&wc)),
	)

	if ret != 0 {
		isClassRegistered = true
		return nil
	}

	hIcon, _, _ := loadIcon.Call(0, uintptr(IDI_APPLICATION))
	hCursor, _, _ := loadCursor.Call(0, uintptr(IDC_ARROW))

	wc = WNDCLASS{
		Style:         0x0003, // CS_HREDRAW | CS_VREDRAW
		LpfnWndProc:   syscall.NewCallback(wndProc),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     hInstance,
		HIcon:         syscall.Handle(hIcon),
		HCursor:       syscall.Handle(hCursor),
		HbrBackground: syscall.Handle(COLOR_WINDOW),
		LpszMenuName:  nil,
		LpszClassName: stringToUTF16Ptr(className),
	}

	atom, _, err := registerClass.Call(uintptr(unsafe.Pointer(&wc)))
	if atom == 0 && err != nil && err.Error() != "The operation completed successfully." {
		return fmt.Errorf("注册窗口类失败: %v", err)
	}

	isClassRegistered = true
	return nil
}

var (
	// 添加字体相关API
	createFontIndirectW = gdi32.NewProc("CreateFontIndirectW")
	deleteObject        = gdi32.NewProc("DeleteObject")
	getDeviceCaps       = gdi32.NewProc("GetDeviceCaps")

	// 添加DPI感知相关API
	setProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	getDpiForWindow               = user32.NewProc("GetDpiForWindow")
	getSystemMetricsForDpi        = user32.NewProc("GetSystemMetricsForDpi")
)

// 字体质量常量
const (
	DEFAULT_QUALITY           = 0
	DRAFT_QUALITY             = 1
	PROOF_QUALITY             = 2
	NONANTIALIASED_QUALITY    = 3
	ANTIALIASED_QUALITY       = 4
	CLEARTYPE_QUALITY         = 5
	CLEARTYPE_NATURAL_QUALITY = 6
)

// DPI 感知常量
const (
	DPI_AWARENESS_CONTEXT_UNAWARE              = ^uintptr(5) // -1
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = ^uintptr(4) // -2
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = ^uintptr(3) // -3
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = ^uintptr(2) // -4
)

// LOGFONT 结构
type LOGFONT struct {
	Height         int32
	Width          int32
	Escapement     int32
	Orientation    int32
	Weight         int32
	Italic         byte
	Underline      byte
	StrikeOut      byte
	CharSet        byte
	OutPrecision   byte
	ClipPrecision  byte
	Quality        byte
	PitchAndFamily byte
	FaceName       [32]uint16
}

// 在init函数中启用DPI感知
func init() {
	// 启用DPI感知（Windows 10 1607+）
	setProcessDpiAwarenessContext.Call(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)

	hInst, _, _ := getModuleHandle.Call(0)
	hInstance = syscall.Handle(hInst)
}

// 创建高清字体
func createHighQualityFont(hWnd syscall.Handle, baseSize int) syscall.Handle {
	// 获取DPI
	dpi := uint32(96)
	if getDpiForWindow.Find() == nil {
		dpiRet, _, _ := getDpiForWindow.Call(uintptr(hWnd))
		dpi = uint32(dpiRet)
	}

	// 根据DPI缩放字体大小
	fontSize := int32(float64(baseSize) * float64(dpi) / 96.0)

	lf := LOGFONT{
		Height:  -fontSize, // 负值表示字符高度
		Width:   0,         // 0表示最佳宽度
		Weight:  400,       // FW_NORMAL
		CharSet: DEFAULT_CHARSET,
		Quality: CLEARTYPE_QUALITY, // ClearType抗锯齿
	}

	// 使用高清字体
	fontName := "Microsoft YaHei UI" // 微软雅黑UI，支持ClearType

	// 复制字体名称
	namePtr, _ := syscall.UTF16PtrFromString(fontName)
	for i := 0; i < 32 && (*namePtr) != 0; i++ {
		lf.FaceName[i] = *namePtr
		namePtr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(namePtr)) + 2))
	}

	hFont, _, _ := createFontIndirectW.Call(uintptr(unsafe.Pointer(&lf)))
	return syscall.Handle(hFont)
}

// 优化的窗口过程
func wndProc(hWnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	// 通过映射获取对话框实例
	dialogMutex.Lock()
	dialogData, exists := dialogs[hWnd]
	dialogMutex.Unlock()

	if !exists {
		if msg == WM_NCCREATE {
			// nolint:unsafeptr // Windows API: CREATESTRUCT与lParam内存布局匹配，转换安全
			createStruct := (*CREATESTRUCT)(unsafe.Pointer(lParam))
			// nolint:unsafeptr // Windows API: DialogData是LpCreateParams的自定义数据结构
			dialogData = (*DialogData)(unsafe.Pointer(createStruct.LpCreateParams))

			dialogMutex.Lock()
			dialogs[hWnd] = dialogData
			dialogMutex.Unlock()

			dialogData.hWnd = hWnd
		} else {
			ret, _, _ := defWindowProc.Call(uintptr(hWnd), uintptr(msg), wParam, lParam)
			return ret
		}
	}

	switch msg {
	case WM_CREATE:
		if dialogData != nil {
			dialogData.hWnd = hWnd

			// 创建高清字体
			dialogData.hFont = createHighQualityFont(hWnd, 14) // 基础12pt

			var clientRect RECT
			getClientRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&clientRect)))

			width := clientRect.Right - clientRect.Left
			height := clientRect.Bottom - clientRect.Top

			// 创建消息文本
			staticX := int32(20)
			staticY := int32(20)
			staticWidth := width - 40

			// 根据文本长度计算高度
			textHeight := int32(100)
			lines := len(dialogData.Message)/40 + 1
			textHeight = int32(lines * 20)
			if textHeight < 60 {
				textHeight = 60
			}

			hStatic, _, _ := createWindowEx.Call(
				0,
				uintptr(unsafe.Pointer(stringToUTF16Ptr("STATIC"))),
				uintptr(unsafe.Pointer(stringToUTF16Ptr(dialogData.Message))),
				uintptr(WS_CHILD|WS_VISIBLE|SS_LEFT|SS_NOPREFIX|SS_NOTIFY|SS_WORDELLIPSIS),
				uintptr(staticX),
				uintptr(staticY),
				uintptr(staticWidth),
				uintptr(textHeight),
				uintptr(hWnd),
				1000,
				uintptr(dialogData.hInstance),
				0,
			)

			if hStatic != 0 {
				dialogData.staticHWND = syscall.Handle(hStatic)
				// 设置高清字体
				if dialogData.hFont != 0 {
					sendMessage.Call(
						uintptr(hStatic),
						WM_SETFONT,
						uintptr(dialogData.hFont),
						1, // 立即重绘
					)
				}
			}

			buttonCount := len(dialogData.Buttons)
			if buttonCount == 0 {
				dialogData.Buttons = append(dialogData.Buttons, struct {
					ID   int
					Text string
				}{ID: 1, Text: "确定"})
				buttonCount = 1
			}

			// 根据DPI调整按钮大小
			buttonWidth := int32(90)
			buttonHeight := int32(32)
			buttonSpacing := int32(12)

			//// 如果按钮多，调整宽度
			//if buttonCount > 3 {
			//	buttonWidth = 160
			//}

			totalWidth := int32(buttonCount)*buttonWidth + int32(buttonCount-1)*buttonSpacing
			startX := (width - totalWidth) / 2
			buttonY := height - buttonHeight - 25 // 增加底部间距

			dialogData.buttonHWND = make([]syscall.Handle, buttonCount)

			for i, btn := range dialogData.Buttons {
				buttonX := startX + int32(i)*(buttonWidth+buttonSpacing)

				style := WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON | WS_TABSTOP
				if i == 0 {
					style |= BS_DEFPUSHBUTTON
				}

				hButton, _, _ := createWindowEx.Call(
					0,
					uintptr(unsafe.Pointer(stringToUTF16Ptr("BUTTON"))),
					uintptr(unsafe.Pointer(stringToUTF16Ptr(btn.Text))),
					uintptr(style),
					uintptr(buttonX),
					uintptr(buttonY),
					uintptr(buttonWidth),
					uintptr(buttonHeight),
					uintptr(hWnd),
					uintptr(1001+i),
					uintptr(dialogData.hInstance),
					0,
				)

				if hButton != 0 {
					dialogData.buttonHWND[i] = syscall.Handle(hButton)
					// 设置按钮字体
					if dialogData.hFont != 0 {
						sendMessage.Call(
							hButton,
							WM_SETFONT,
							uintptr(dialogData.hFont),
							1,
						)
					}
				}
			}
		}
		return 0

	case WM_CTLCOLORSTATIC:
		// 静态文本颜色处理
		hdc := syscall.Handle(wParam)

		// 设置文本颜色
		setTextColor.Call(
			uintptr(hdc),
			uintptr(0x00000000), // 黑色文本
		)

		//// 设置背景颜色
		//setBkColor.Call(
		//	uintptr(hdc),
		//	uintptr(0x00FFFFFF), // 白色背景
		//)

		// 返回白色画刷
		hBrush, _, _ := getStockObject.Call(uintptr(WHITE_BRUSH))
		return hBrush

	case WM_CTLCOLORBTN:
		// 按钮颜色处理
		hdc := syscall.Handle(wParam)

		// 设置按钮文本颜色
		setTextColor.Call(
			uintptr(hdc),
			uintptr(0x00000000), // 黑色文本
		)

		// 返回按钮背景画刷
		hBrush, _, _ := getStockObject.Call(uintptr(COLOR_BTNFACE))
		return hBrush

	case WM_COMMAND:
		if dialogData != nil {
			id := int(wParam & 0xFFFF)

			if id >= 1001 && id <= 1001+len(dialogData.Buttons) {
				buttonIndex := id - 1001
				if buttonIndex < len(dialogData.Buttons) {
					dialogData.SelectedID = dialogData.Buttons[buttonIndex].ID
					select {
					case dialogData.ResultChan <- dialogData.SelectedID:
					default:
					}
					destroyWindow.Call(uintptr(hWnd))
					return 0
				}
			}
		}

	case WM_CLOSE:
		if dialogData != nil {
			dialogData.SelectedID = 0
			select {
			case dialogData.ResultChan <- 0:
			default:
			}
		}
		destroyWindow.Call(uintptr(hWnd))
		return 0

	case WM_DESTROY:
		// 清理字体资源
		if dialogData != nil && dialogData.hFont != 0 {
			deleteObject.Call(uintptr(dialogData.hFont))
		}

		// 清理对话框实例
		dialogMutex.Lock()
		delete(dialogs, hWnd)
		dialogMutex.Unlock()
		return 0

	case WM_DPICHANGED:
		// DPI变化时重新创建字体
		if dialogData != nil && dialogData.hFont != 0 {
			// 删除旧字体
			deleteObject.Call(uintptr(dialogData.hFont))

			// 创建新字体
			newDpi := uint32(wParam >> 16)
			newSize := int32(float64(16) * float64(newDpi) / 96.0)
			dialogData.hFont = createHighQualityFont(hWnd, int(newSize))

			// 重新设置所有控件字体
			if dialogData.staticHWND != 0 {
				sendMessage.Call(
					uintptr(dialogData.staticHWND),
					WM_SETFONT,
					uintptr(dialogData.hFont),
					1,
				)
			}

			for _, hButton := range dialogData.buttonHWND {
				if hButton != 0 {
					sendMessage.Call(
						uintptr(hButton),
						WM_SETFONT,
						uintptr(dialogData.hFont),
						1,
					)
				}
			}

			return 0
		}
	}

	ret, _, _ := defWindowProc.Call(uintptr(hWnd), uintptr(msg), wParam, lParam)
	return ret
}

// 显示自定义对话框
func ShowCustomDialog(title, message string, buttons map[int]string, window, height int32) (int, error) {
	// 注册窗口类
	if err := registerWindowClass(); err != nil {
		return 0, err
	}

	// 准备按钮数据
	var buttonList []struct {
		ID   int
		Text string
	}

	keys := make([]int, 0, len(buttons))
	for k := range buttons {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		buttonList = append(buttonList, struct {
			ID   int
			Text string
		}{ID: k, Text: buttons[k]})
	}

	if len(buttonList) == 0 {
		buttonList = append(buttonList, struct {
			ID   int
			Text string
		}{ID: 1, Text: "确定"})
	}

	// 创建对话框数据
	dialogData := &DialogData{
		Title:      title,
		Message:    message,
		Buttons:    buttonList,
		SelectedID: 0,
		ResultChan: make(chan int, 1),
		hInstance:  hInstance,
	}

	// 计算窗口大小
	windowWidth := int32(len(buttons) * 120)
	if window != 0 {
		windowWidth = window
	}
	windowHeight := int32((len(message)*10)/(int(windowWidth)))*100 + 200
	if windowHeight > 700 {
		windowHeight = 700
	}
	if height != 0 {
		windowHeight = height
	}

	// 使用简单居中位置
	x := int32(100)
	y := int32(100)

	// 创建窗口
	hwnd, _, _ := createWindowEx.Call(
		uintptr(WS_EX_DLGMODALFRAME),
		uintptr(unsafe.Pointer(stringToUTF16Ptr(className))),
		uintptr(unsafe.Pointer(stringToUTF16Ptr(title))),
		uintptr(WS_CAPTION|WS_SYSMENU|WS_POPUP),
		uintptr(x),
		uintptr(y),
		uintptr(windowWidth),
		uintptr(windowHeight),
		0,
		0,
		uintptr(hInstance),
		uintptr(unsafe.Pointer(dialogData)), // 传递对话框数据
	)

	if hwnd == 0 {
		return 0, fmt.Errorf("创建窗口失败")
	}

	dialogData.hWnd = syscall.Handle(hwnd)

	// 显示窗口
	showWindow.Call(hwnd, SW_SHOW)
	updateWindow.Call(hwnd)

	// 消息循环
	var msg MSG
	for {
		ret, _, _ := getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)

		if ret == 0 {
			break
		}

		translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))

		// 非阻塞检查结果
		select {
		case result := <-dialogData.ResultChan:
			return result, nil
		default:
		}
	}

	return dialogData.SelectedID, nil
}
