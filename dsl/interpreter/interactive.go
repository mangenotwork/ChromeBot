package interpreter

import (
	"ChromeBot/dsl/ast"
	"ChromeBot/utils"
	"fmt"
	"strings"
)

// 解析 chrome 关键字语法
func (i *Interpreter) evaluateChromeStmt(expr *ast.ChromeStmt, ctx *Context, hang int) Value {
	utils.Debug("evaluateChromeStmt args = ", expr.Args)
	fn, ok := ctx.GetFunc("chrome")
	if !ok {
		i.errors = append(i.errors, fmt.Errorf("未定义Chrome"))
		return nil
	}
	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		args[idx] = i.evaluateExpr(arg, ctx, hang)
	}

	result, err := fn(args)
	if err != nil {
		i.errors = append(i.errors, fmt.Errorf("Chrome调用错误: %v", err))
		return nil
	}

	return result
}

// 解析 http 关键字语法
func (i *Interpreter) evaluateHttpStmt(expr *ast.HttpStmt, ctx *Context, hang int) Value {
	utils.Debug("evaluateHttpStmt args = ", expr.Args)
	fn, ok := ctx.GetFunc("http")
	if !ok {
		i.ErrorShow(hang, "未定义http")
		return nil
	}

	if len(expr.Args) < 1 {
		i.ErrorShow(hang, "http后面没有参数")
		return nil
	}

	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		args[idx] = i.evaluateExpr(arg, ctx, hang)
	}

	oneArg := strings.ToLower(args[0].(string))
	if oneArg != "get" && oneArg != "post" && oneArg != "put" && oneArg != "delete" &&
		oneArg != "options" && oneArg != "head" && oneArg != "patch" {
		i.ErrorShow(hang, "http第一个参数必须是 get, post, put, delete")
		return nil
	}

	isUrl := false

	argMap := make(map[string]Value)

	argMap["method"] = oneArg

	for idx, arg := range args {
		if idx == 0 {
			continue
		}
		utils.Debug("evaluateHttpStmt args[", idx, "] = ", arg)
		utils.Debug(fmt.Sprintf("%T", arg))
		switch e := arg.(type) {
		case string:
			sl := strings.SplitN(e, "=", 2)
			if len(sl) > 0 && (sl[0] == "url" || sl[0] == "body" || sl[0] == "header" || sl[0] == "ctype" ||
				sl[0] == "cookie" || sl[0] == "timeout" || sl[0] == "proxy" || sl[0] == "stress" || sl[0] == "save") {
				if sl[0] == "url" {
					isUrl = true
				}
				if len(sl) < 2 {
					i.ErrorShow(hang, fmt.Sprintf("http操作%s参数没设置值", sl[0]))
					return nil
				}
				newVal, has := ctx.GetVar(sl[1])
				if has {
					argMap[sl[0]] = newVal
				} else {
					argMap[sl[0]] = sl[1]
				}
			}

			if sl[0] == "to" {
				if len(sl) < 2 {
					i.ErrorShow(hang, fmt.Sprintf("http操作%s参数没设置值", sl[0]))
					return nil
				}
				argMap[sl[0]] = sl[1]
			}
		}
	}

	if !isUrl {
		i.ErrorShow(hang, "http操作必须指定url参数")
		return nil
	}

	utils.Debug("argMap = ", argMap)

	inArg := []Value{
		argMap,
	}

	result, err := fn(inArg)
	if err != nil {
		i.errors = append(i.errors, fmt.Errorf("http调用错误: %v", err))
		return nil
	}

	return result
}
