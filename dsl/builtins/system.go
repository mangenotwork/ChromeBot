package builtins

import (
	"ChromeBot/dsl/interpreter"
	"ChromeBot/internal/host"
	"fmt"
)

// 系统相关的方法
var systemFn = map[string]interpreter.Function{
	"sysConfirmBox": sysConfirmBox, // sysConfirmBox(title, msg) 确认弹框； 点击返回 true, false

	// sysDialogBox(title, msg, buttons, window, height) buttons 要求是字典整型值为key按钮名为value如 {1:"按钮名字"， 2:"按钮名字2"}；
	// 如果 window, height = 0那么会默认宽高； 点击返回 buttons 对应的key值; 注意返回0是关闭
	"sysDialogBox": sysDialogBox,

	"sysExitBox": sysExitBox, // sysExitBox() 是否终止当前ChromeBot进程的确认框, 点击是会终止进程

	"sysInfoTip":    sysInfoTip,    // sysInfoTip(msg) 信息提示框
	"sysWarningTip": sysWarningTip, // sysWarningTip(msg) 警告提示框
	"sysErrorTip":   sysErrorTip,   // sysErrorTip(msg) 错误提示框
	"sysSuccessTip": sysSuccessTip, // sysSuccessTip(msg) 成功提示框

	"sysBoxTest": sysBoxTest, // 测试三个按钮的提示框
}

func sysConfirmBox(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("sysConfirmBox(title, msg) 需要两个参数")
	}
	title, titleOK := args[0].(string)
	if !titleOK {
		return nil, fmt.Errorf("sysConfirmBox(title, msg) 第一个参数要求是字符串")
	}
	msg, msgOK := args[1].(string)
	if !msgOK {
		return nil, fmt.Errorf("sysConfirmBox(title, msg) 第二个参数要求是字符串 ")
	}
	boxRes, err := host.SystemConfirmBox(title, msg)
	if err != nil {
		fmt.Println("系统弹出框错误，err : ", err.Error())
		return nil, err
	}

	return boxRes, nil
}

func sysDialogBox(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("sysDialogBox(title, msg, buttons, window, height) 需要五个参数")
	}

	title, titleOK := args[0].(string)
	if !titleOK {
		return nil, fmt.Errorf("sysDialogBox(title, msg, buttons, window, height) title要求是字符串")
	}
	msg, msgOK := args[1].(string)
	if !msgOK {
		return nil, fmt.Errorf("ssysDialogBox(title, msg, buttons, window, height) msg要求是字符串 ")
	}
	buttons, buttonsOK := args[2].(interpreter.DictType)
	if !buttonsOK {
		return nil, fmt.Errorf("ssysDialogBox(title, msg, buttons, window, height) buttons要求是字典 ")
	}
	window, windowOK := args[3].(int64)
	if !windowOK {
		return nil, fmt.Errorf("ssysDialogBox(title, msg, buttons, window, height) window要求是数值 ")
	}
	height, heightOK := args[4].(int64)
	if !heightOK {
		return nil, fmt.Errorf("ssysDialogBox(title, msg, buttons, window, height) height要求是数值 ")
	}

	buttonsData := make(map[int]string)
	for k, v := range buttons {
		kVal, kValOK := k.(int64)
		if !kValOK {
			return nil, fmt.Errorf("ssysDialogBox(title, msg, buttons, window, height) buttons要求是字典, key是数值类型 ")
		}
		vStr, vStrOK := v.(string)
		if !vStrOK {
			return nil, fmt.Errorf("ssysDialogBox(title, msg, buttons, window, height) buttons要求是字典, val是字符串类型 ")
		}
		buttonsData[int(kVal)] = vStr
	}

	result, err := host.ShowCustomDialog(title, msg, buttonsData, int32(window), int32(height))
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("点击了对话框 <%s> 按钮 : %s 值: %d   \n", title, buttonsData[result], result)
	}

	return result, nil
}

func sysExitBox(args []interpreter.Value) (interpreter.Value, error) {
	err := host.SystemExitBox()
	if err != nil {
		fmt.Println("sysExitBox err : ", err.Error())
	}
	return nil, nil
}

func sysInfoTip(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sysInfoTip(msg) 需要一个参数")
	}
	msg, msgOK := args[0].(string)
	if !msgOK {
		return nil, fmt.Errorf("sysInfoTip(msg) 参数要求是字符串 ")
	}
	_, _ = host.InfoTipBox(msg)
	return nil, nil
}

func sysWarningTip(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sysWarningTip(msg) 需要一个参数")
	}
	msg, msgOK := args[0].(string)
	if !msgOK {
		return nil, fmt.Errorf("sysWarningTip(msg) 参数要求是字符串 ")
	}
	_, _ = host.WarningTipBox(msg)
	return nil, nil
}

func sysErrorTip(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sysErrorTip(msg) 需要一个参数")
	}
	msg, msgOK := args[0].(string)
	if !msgOK {
		return nil, fmt.Errorf("sysErrorTip(msg) 参数要求是字符串 ")
	}
	_, _ = host.ErrorTipBox(msg)
	return nil, nil
}

func sysSuccessTip(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sysSuccessTip(msg) 需要一个参数")
	}
	msg, msgOK := args[0].(string)
	if !msgOK {
		return nil, fmt.Errorf("sysSuccessTip(msg) 参数要求是字符串 ")
	}
	_, _ = host.SuccessTipBox(msg)
	return nil, nil
}

func sysBoxTest(args []interpreter.Value) (interpreter.Value, error) {

	fmt.Println("\n 复杂的三按钮对话框:")
	result, err := host.NewThreeButtonDialog("备份选项", "请选择备份方式:").
		SetButtons("完全备份", "增量备份", "差异备份").
		SetBoxType(host.MB_ABORTRETRYIGNORE).
		SetIcon(host.MB_ICONQUESTION).
		SetDefaultButton(1).
		Show()

	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("用户选择了备份方式: 按钮%d\n", result)
	}
	return nil, nil
}
