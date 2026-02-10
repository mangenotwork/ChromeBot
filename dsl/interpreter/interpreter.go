package interpreter

import (
	"ChromeBot/dsl/ast"
	"ChromeBot/utils"
	"fmt"
	"os"
)

// Value 值接口
type Value interface{}

// DictType 字典类型
type DictType map[Value]Value

// Function 函数定义
type Function func(args []Value) (Value, error)

// Context 执行上下文
type Context struct {
	parent      *Context
	variables   map[string]Value
	functions   map[string]Function
	returnVal   *Value
	hasReturn   bool
	hasBreak    bool
	hasContinue bool
}

// NewContext 创建上下文
func NewContext(parent *Context) *Context {
	return &Context{
		parent:    parent,
		variables: make(map[string]Value),
		functions: make(map[string]Function),
	}
}

// SetVar 设置变量
func (c *Context) SetVar(name string, value Value) {
	c.variables[name] = value
}

// GetVar 获取变量
func (c *Context) GetVar(name string) (Value, bool) {
	val, ok := c.variables[name]
	if !ok && c.parent != nil {
		return c.parent.GetVar(name)
	}
	return val, ok
}

// SetFunc 设置函数
func (c *Context) SetFunc(name string, fn Function) {
	c.functions[name] = fn
}

// GetFunc 获取函数
func (c *Context) GetFunc(name string) (Function, bool) {
	fn, ok := c.functions[name]
	if !ok && c.parent != nil {
		return c.parent.GetFunc(name)
	}
	return fn, ok
}

// Interpreter 解释器
type Interpreter struct {
	global *Context
	errors []error
}

// NewInterpreter 创建解释器
func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		global: NewContext(nil),
		errors: []error{},
	}

	// 注册内置函数
	interp.registerBuiltins()

	return interp
}

// Global 返回全局上下文
func (i *Interpreter) Global() *Context {
	return i.global
}

// Interpret 执行AST
func (i *Interpreter) Interpret(program *ast.Program) (Value, error) {
	utils.Debug("执行AST ....")
	for n, stmt := range program.Statements {
		utils.Debug(n+1, " - Interpret ==> ", stmt)
		_ = i.evaluateStmt(stmt, i.global, n+1)

		if i.global.hasReturn {
			return *i.global.returnVal, nil
		}
	}

	return nil, nil
}

func (i *Interpreter) evaluateStmt(stmt ast.Statement, ctx *Context, hang int) Value {

	utils.Debug("evaluateStmt ==> ", stmt)

	switch s := stmt.(type) {
	case *ast.VarDecl:
		return i.evaluateVarDecl(s, ctx, hang)
	case *ast.AssignStmt:
		return i.evaluateAssignStmt(s, ctx, hang)
	case *ast.ExpressionStmt:
		return i.evaluateExpr(s.Expr, ctx, hang)
	case *ast.BlockStmt:
		return i.evaluateBlockStmt(s, ctx, hang)
	case *ast.IfStmt:
		return i.evaluateIfStmt(s, ctx, hang)
	case *ast.SwitchStmt:
		return i.evaluateSwitchStmt(s, ctx, hang)
	case *ast.WhileStmt:
		return i.evaluateWhileStmt(s, ctx, hang)
	case *ast.ReturnStmt:
		return i.evaluateReturnStmt(s, ctx, hang)
	case *ast.ChromeStmt:
		return i.evaluateChromeStmt(s, ctx, hang)
	case *ast.BreakStmt:
		ctx.hasBreak = true
		return nil
	case *ast.ContinueStmt:
		ctx.hasContinue = true
		return nil
	case *ast.ForStmt:
		return i.evaluateForStmt(s, ctx, hang)
	case *ast.IndexAssignStmt: // 添加下标赋值处理
		return i.evaluateIndexAssignStmt(s, ctx, hang)
	default:
		fmt.Println("[Crash]len:", hang, " | ", fmt.Errorf("不支持的语句类型: %T", stmt))
		os.Exit(0)
	}
	return nil
}

func (i *Interpreter) evaluateExpr(expr ast.Expression, ctx *Context, hang int) Value {
	utils.Debug("evaluateExpr ==> ", expr)
	switch e := expr.(type) {

	case *ast.Integer:
		utils.Debug("evaluateExpr ast.Integer ==> ", e.Value)
		return e.Value

	case *ast.Float: // 添加浮点数求值
		utils.Debug("evaluateExpr ast.Float ==> ", e.Value)
		return e.Value

	case *ast.String:
		utils.Debug("evaluateExpr ast.String ==> ", e.Value)
		return e.Value

	case *ast.Boolean:
		utils.Debug("evaluateExpr ast.Boolean ==> ", e.Value)
		return e.Value

	case *ast.Identifier:
		utils.Debug("evaluateExpr ast.Identifier ==> ", e)
		if val, ok := ctx.GetVar(e.Name); ok {
			return val
		}
		fmt.Println("[Crash]len:", hang, " | ", fmt.Errorf("未定义的变量: %s", e.Name))
		os.Exit(0)

	case *ast.BinaryExpr:
		utils.Debug("evaluateExpr ast.BinaryExpr ==> ", e)
		return i.evaluateBinaryExpr(e, ctx, hang)

	case *ast.UnaryExpr:
		utils.Debug("evaluateExpr ast.UnaryExpr ==> ", e)
		return i.evaluateUnaryExpr(e, ctx, hang)

	case *ast.PostfixExpr: // 添加对自增自减的处理
		utils.Debug("evaluateExpr ast.PostfixExpr ==> ", e)
		return i.evaluatePostfixExpr(e, ctx, hang)

	case *ast.CallExpr:
		utils.Debug("evaluateExpr ast.CallExpr ==> ", e)
		return i.evaluateCallExpr(e, ctx, hang)

	case *ast.List: // 添加列表字面量求值
		utils.Debug("evaluateExpr ast.List ==> ", e)
		return i.evaluateList(e, ctx, hang)

	case *ast.IndexExpr: // 添加下标表达式求值
		utils.Debug("evaluateExpr ast.IndexExpr ==> ", e)
		return i.evaluateIndexExpr(e, ctx, hang)

	case *ast.Dict: // 添加字典字面量求值
		utils.Debug("evaluateExpr ast.Dict ==> ", e)
		return i.evaluateDict(e, ctx, hang)

	case *ast.ChainCallExpr: // 添加链式调用求值
		utils.Debug("evaluateExpr ast.ChainCallExpr ==> ", e)
		return i.evaluateChainCall(e, ctx, hang)

	default:
		fmt.Println("[Crash]len:", hang, " | ", fmt.Errorf("不支持的表达式类型: %T", expr))
		os.Exit(0)

	}
	return nil
}

func (i *Interpreter) evaluateVarDecl(decl *ast.VarDecl, ctx *Context, hang int) Value {
	var value Value

	if decl.Expr != nil {
		value = i.evaluateExpr(decl.Expr, ctx, hang)
	} else {
		// 默认值
		switch decl.Type {
		case "int":
			value = int64(0)
		case "string":
			value = ""
		case "bool":
			value = false
		default:
			value = nil
		}
	}
	utils.Debug("evaluateVarDecl ==> ", decl.Name.Name, ":", value)
	ctx.SetVar(decl.Name.Name, value)
	return value
}

func (i *Interpreter) evaluateAssignStmt(stmt *ast.AssignStmt, ctx *Context, hang int) Value {
	value := i.evaluateExpr(stmt.Expr, ctx, hang)
	utils.Debug("evaluateAssignStmt ==> ", stmt.Left.Name, ":", value)
	ctx.SetVar(stmt.Left.Name, value)
	return value
}

func (i *Interpreter) evaluateBinaryExpr(expr *ast.BinaryExpr, ctx *Context, hang int) Value {
	left := i.evaluateExpr(expr.Left, ctx, hang)
	right := i.evaluateExpr(expr.Right, ctx, hang)

	// 类型检查和转换
	switch expr.Op {
	case "+":
		return i.add(left, right)
	case "-":
		return i.sub(left, right)
	case "*":
		return i.mul(left, right)
	case "/":
		return i.div(left, right)
	case "%":
		return i.mod(left, right)
	case "==":
		return i.equal(left, right)
	case "!=":
		return !i.equal(left, right)
	case "<":
		return i.less(left, right)
	case "<=":
		return i.less(left, right) || i.equal(left, right)
	case ">":
		return i.greater(left, right)
	case ">=":
		return i.greater(left, right) || i.equal(left, right)
	case "&&":
		return i.bool(left) && i.bool(right)
	case "||":
		return i.bool(left) || i.bool(right)
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作符: %s", expr.Op))
		return nil
	}
}

func (i *Interpreter) evaluateUnaryExpr(expr *ast.UnaryExpr, ctx *Context, hang int) Value {
	right := i.evaluateExpr(expr.Expr, ctx, hang)

	switch expr.Op {
	case "-":
		switch v := right.(type) {
		case int64:
			return -v
		case float64:
			return -v
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: -%T", right))
			return nil
		}
	case "!":
		return !i.bool(right)
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作符: %s", expr.Op))
		return nil
	}
}

func (i *Interpreter) evaluateCallExpr(expr *ast.CallExpr, ctx *Context, hang int) Value {
	fn, ok := ctx.GetFunc(expr.Function.Name)
	if !ok {
		i.errors = append(i.errors, fmt.Errorf("未定义的函数: %s", expr.Function.Name))
		return nil
	}

	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		args[idx] = i.evaluateExpr(arg, ctx, hang)
	}

	result, err := fn(args)
	if err != nil {
		i.errors = append(i.errors, fmt.Errorf("函数调用错误 %s: %v", expr.Function.Name, err))
		return nil
	}

	return result
}

func (i *Interpreter) evaluateBlockStmt(block *ast.BlockStmt, ctx *Context, hang int) Value {
	// 创建一个新的作用域
	newCtx := NewContext(ctx)
	utils.Debug("evaluateBlockStmt ==> ", block.Stmts)

	for _, stmt := range block.Stmts {
		utils.Debug("evaluateBlockStmt is stmt item ==> ", stmt)

		switch stmt.(type) {

		case *ast.BreakStmt:
			ctx.hasBreak = true
			return nil

		case *ast.ContinueStmt:
			ctx.hasContinue = true
			return nil

		default:
			_ = i.evaluateStmt(stmt, newCtx, hang)

			if newCtx.hasReturn || newCtx.hasBreak || newCtx.hasContinue {
				ctx.hasReturn = true
				ctx.returnVal = newCtx.returnVal
				return *ctx.returnVal
			}
		}

	}

	return nil
}

func (i *Interpreter) evaluateIfStmt(stmt *ast.IfStmt, ctx *Context, hang int) Value {
	condition := i.evaluateExpr(stmt.Condition, ctx, hang)
	utils.Debug("evaluateBlockStmt ==> ", stmt)
	if i.bool(condition) {
		return i.evaluateBlockStmt(stmt.Then, ctx, hang)
	} else if stmt.Else != nil {
		switch e := stmt.Else.(type) {
		case *ast.BlockStmt:
			return i.evaluateBlockStmt(e, ctx, hang)
		case *ast.IfStmt:
			return i.evaluateIfStmt(e, ctx, hang)
		}
	}

	return nil
}

func (i *Interpreter) evaluateWhileStmt(stmt *ast.WhileStmt, ctx *Context, hang int) Value {
	utils.Debug("evaluateWhileStmt ==> ", stmt)
	for {
		// 检查循环条件
		condition := i.evaluateExpr(stmt.Condition, ctx, hang)
		if !i.bool(condition) {
			break
		}

		// 执行循环体
		loopCtx := NewContext(ctx)

		// 执行循环体
		for _, stmtItem := range stmt.Body.Stmts {
			_ = i.evaluateStmt(stmtItem, loopCtx, hang)

			// 检查是否需要提前退出
			if loopCtx.hasReturn || loopCtx.hasBreak || loopCtx.hasContinue {
				//log.Println("检查到要提前退出")
				break
			}
		}

		// 处理控制流
		if loopCtx.hasReturn {
			ctx.hasReturn = true
			ctx.returnVal = loopCtx.returnVal
			return *ctx.returnVal
		}

		if loopCtx.hasBreak {
			break
		}

		// 如果父作用域中本来没有这个变量，但现在有了，也要设置
		for k, v := range loopCtx.variables {
			ctx.SetVar(k, v) // 直接设置，覆盖原有的值
		}

		if loopCtx.hasContinue {
			continue
		}

	}

	return nil
}

func (i *Interpreter) evaluateForStmt(stmt *ast.ForStmt, ctx *Context, hang int) Value {

	if stmt.Init != nil {
		_ = i.evaluateStmt(stmt.Init, ctx, hang)
		if ctx.hasReturn || ctx.hasBreak || ctx.hasContinue {
			return nil
		}
	}

	for {
		// 检查循环条件
		if stmt.Cond != nil {
			condition := i.evaluateExpr(stmt.Cond, ctx, hang)
			if !i.bool(condition) {
				break
			}
		}

		// 执行循环体
		for _, stmtItem := range stmt.Body.Stmts {
			_ = i.evaluateStmt(stmtItem, ctx, hang)

			// 检查控制流
			if ctx.hasReturn {
				return *ctx.returnVal
			}
			if ctx.hasBreak {
				ctx.hasBreak = false
				return nil
			}
			if ctx.hasContinue {
				ctx.hasContinue = false
				break
			}
		}

		// 如果是 continue，直接执行后置语句
		if ctx.hasContinue {
			ctx.hasContinue = false
		}

		// 执行后置语句
		if stmt.Post != nil {
			_ = i.evaluateStmt(stmt.Post, ctx, hang)
			if ctx.hasReturn || ctx.hasBreak || ctx.hasContinue {
				return nil
			}
		}

	}

	return nil
}

func (i *Interpreter) evaluateReturnStmt(stmt *ast.ReturnStmt, ctx *Context, hang int) Value {
	var value Value
	if stmt.Expr != nil {
		value = i.evaluateExpr(stmt.Expr, ctx, hang)
	}
	utils.Debug("evaluateReturnStmt ==> ", stmt, value)
	ctx.hasReturn = true
	ctx.returnVal = &value
	return value
}

func (i *Interpreter) evaluateList(list *ast.List, ctx *Context, hang int) Value {
	elements := make([]Value, len(list.Elements))
	for idx, element := range list.Elements {
		elements[idx] = i.evaluateExpr(element, ctx, hang)
	}
	return elements
}

func (i *Interpreter) evaluateIndexExpr(expr *ast.IndexExpr, ctx *Context, hang int) Value {
	// 求值左边的表达式（应该是列表或字典）
	left := i.evaluateExpr(expr.Left, ctx, hang)

	// 求值下标
	index := i.evaluateExpr(expr.Index, ctx, hang)

	// 检查左边是列表还是字典
	switch container := left.(type) {
	case []Value: // 列表
		// 检查下标是否是整数
		idx, ok := index.(int64)
		if !ok {
			i.errors = append(i.errors, fmt.Errorf("列表下标必须是整数，得到: %T", index))
			return nil
		}

		// 检查下标是否越界
		if idx < 0 || idx >= int64(len(container)) {
			i.errors = append(i.errors, fmt.Errorf("列表下标越界: 长度=%d, 下标=%d", len(container), idx))
			return nil
		}

		return container[idx]

	case DictType: // 字典
		// 检查键是否是可哈希的类型
		if !i.isHashable(index) {
			i.errors = append(i.errors, fmt.Errorf("字典键必须是可哈希的类型，得到: %T", index))
			return nil
		}

		// 查找键对应的值
		value, exists := container[index]
		if !exists {
			i.errors = append(i.errors, fmt.Errorf("字典中不存在键: %v", index))
			return nil
		}

		return value

	default:
		i.errors = append(i.errors, fmt.Errorf("下标操作只支持列表或字典，得到: %T", left))
		return nil
	}
}

func (i *Interpreter) evaluateDict(dict *ast.Dict, ctx *Context, hang int) Value {
	result := make(DictType)

	for keyExpr, valueExpr := range dict.Pairs {
		// 求值键
		key := i.evaluateExpr(keyExpr, ctx, hang)

		// 检查键的类型（在Go中，只有可比较的类型才能作为map的键）
		// 我们只支持基本类型作为键
		if !i.isHashable(key) {
			i.errors = append(i.errors, fmt.Errorf("字典键必须是可哈希的类型，得到: %T", key))
			return nil
		}

		// 求值值
		value := i.evaluateExpr(valueExpr, ctx, hang)

		// 添加到字典
		result[key] = value
	}

	return result
}

// 检查值是否可以作为字典的键
func (i *Interpreter) isHashable(value Value) bool {
	switch value.(type) {
	case int64, float64, string, bool:
		return true
	default:
		return false
	}
}

func (i *Interpreter) evaluateIndexAssignStmt(stmt *ast.IndexAssignStmt, ctx *Context, hang int) Value {
	// 求值右边的表达式
	value := i.evaluateExpr(stmt.Expr, ctx, hang)

	// 获取目标容器和键/下标
	target := stmt.Target.Left
	key := stmt.Target.Index

	// 求值目标容器
	container := i.evaluateExpr(target, ctx, hang)

	// 求值键/下标
	index := i.evaluateExpr(key, ctx, hang)

	// 检查容器类型并赋值
	switch c := container.(type) {
	case []Value: // 列表
		// 检查下标是否是整数
		idx, ok := index.(int64)
		if !ok {
			i.errors = append(i.errors, fmt.Errorf("列表下标必须是整数，得到: %T", index))
			return nil
		}

		// 检查下标是否越界
		if idx < 0 || idx >= int64(len(c)) {
			i.errors = append(i.errors, fmt.Errorf("列表下标越界: 长度=%d, 下标=%d", len(c), idx))
			return nil
		}

		// 赋值
		c[idx] = value

	case DictType: // 字典
		// 检查键是否是可哈希的类型
		if !i.isHashable(index) {
			i.errors = append(i.errors, fmt.Errorf("字典键必须是可哈希的类型，得到: %T", index))
			return nil
		}

		// 赋值（添加或修改）
		c[index] = value

	default:
		i.errors = append(i.errors, fmt.Errorf("下标赋值只支持列表或字典，得到: %T", container))
		return nil
	}

	return value
}

func (i *Interpreter) evaluateChainCall(chain *ast.ChainCallExpr, ctx *Context, hang int) Value {
	utils.Debug("evaluateChainCall ==> ", chain)

	// 存储前一个调用的结果
	var lastResult Value = nil

	// 遍历所有调用
	for callIndex, call := range chain.Calls {
		utils.Debug("执行链式调用第 %d 个调用: %s", callIndex+1, call.Function.Name)

		// 处理第一个调用
		if callIndex == 0 {
			// 检查是否是特殊标记 "_value"
			if call.Function.Name == "_value" {
				// 这是值包装，直接取值
				if len(call.Args) == 1 {
					lastResult = i.evaluateExpr(call.Args[0], ctx, hang)
					utils.Debug("第一个是值: %v", lastResult)
				} else {
					i.errors = append(i.errors, fmt.Errorf("值包装应该有1个参数，得到 %d 个", len(call.Args)))
					return nil
				}

			} else {
				// 普通函数调用
				fn, ok := ctx.GetFunc(call.Function.Name)
				if !ok {
					// 可能是变量
					if val, ok := ctx.GetVar(call.Function.Name); ok {
						lastResult = val
						utils.Debug("第一个是变量 %s: %v", call.Function.Name, lastResult)
					} else {
						i.errors = append(i.errors, fmt.Errorf("未定义的函数或变量: %s", call.Function.Name))
						return nil
					}

				} else {
					// 准备参数
					args := make([]Value, len(call.Args))
					for argIdx, arg := range call.Args {
						args[argIdx] = i.evaluateExpr(arg, ctx, hang)
					}

					// 执行函数
					result, err := fn(args)
					if err != nil {
						i.errors = append(i.errors, fmt.Errorf("链式调用错误 %s: %v", call.Function.Name, err))
						return nil
					}

					lastResult = result
					utils.Debug("函数 %s 返回: %v", call.Function.Name, result)
				}
			}
		} else {
			// 后续调用必须是函数
			fn, ok := ctx.GetFunc(call.Function.Name)
			if !ok {
				i.errors = append(i.errors, fmt.Errorf("未定义的函数: %s", call.Function.Name))
				return nil
			}

			// 准备参数
			args := make([]Value, len(call.Args))
			for argIdx, arg := range call.Args {
				args[argIdx] = i.evaluateExpr(arg, ctx, hang)
			}

			// 将前一个结果作为第一个参数
			newArgs := make([]Value, len(args)+1)
			newArgs[0] = lastResult
			copy(newArgs[1:], args)

			utils.Debug("调用 %s 参数: %v", call.Function.Name, newArgs)

			// 执行函数
			result, err := fn(newArgs)
			if err != nil {
				i.errors = append(i.errors, fmt.Errorf("链式调用错误 %s: %v", call.Function.Name, err))
				return nil
			}

			lastResult = result
			utils.Debug("函数 %s 返回: %v", call.Function.Name, result)
		}
	}

	return lastResult
}

func (i *Interpreter) evaluateSwitchStmt(stmt *ast.SwitchStmt, ctx *Context, hang int) Value {
	// 计算 switch 表达式的值
	switchValue := i.evaluateExpr(stmt.Expr, ctx, hang)
	utils.Debug("evaluateSwitchStmt ==> ", stmt)

	// 遍历所有 case
	for _, caseClause := range stmt.Cases {
		// 检查是否有匹配的 case
		for _, caseValueExpr := range caseClause.Values {
			caseValue := i.evaluateExpr(caseValueExpr, ctx, hang)

			// 比较值是否相等
			if i.equal(switchValue, caseValue) {

				// 为 case 块创建新的上下文
				caseCtx := NewContext(ctx)
				result := i.evaluateBlockStmt(caseClause.Body, caseCtx, hang)

				// 检查是否执行了 break
				if caseCtx.hasBreak {
					return result
				}

				// 如果没有 break，返回结果
				return result
			}
		}
	}

	// 执行 default
	if stmt.Default != nil {
		defaultCtx := NewContext(ctx)
		return i.evaluateBlockStmt(stmt.Default, defaultCtx, hang)
	}

	return nil
}

// 添加evaluatePostfixExpr函数
func (i *Interpreter) evaluatePostfixExpr(expr *ast.PostfixExpr, ctx *Context, hang int) Value {
	utils.Debug("evaluatePostfixExpr: 开始，表达式=%v", expr)

	// 先获取原始值
	var originalValue Value
	var target interface{} // 用于存储修改的目标

	switch left := expr.Left.(type) {
	case *ast.Identifier:
		// 变量自增自减
		val, ok := ctx.GetVar(left.Name)
		if !ok {
			i.errors = append(i.errors, fmt.Errorf("未定义的变量: %s", left.Name))
			return nil
		}
		originalValue = val
		target = left // 存储标识符以便修改

	case *ast.IndexExpr:
		// 列表或字典元素自增自减
		// 先获取容器和索引
		container := i.evaluateExpr(left.Left, ctx, hang)
		index := i.evaluateExpr(left.Index, ctx, hang)

		// 根据容器类型获取原始值
		switch c := container.(type) {
		case []Value: // 列表
			idx, ok := index.(int64)
			if !ok {
				i.errors = append(i.errors, fmt.Errorf("列表下标必须是整数"))
				return nil
			}
			if idx < 0 || idx >= int64(len(c)) {
				i.errors = append(i.errors, fmt.Errorf("列表下标越界"))
				return nil
			}
			originalValue = c[idx]
			target = &indexTarget{
				container: container,
				index:     idx,
				isList:    true,
			}

		case DictType: // 字典
			if !i.isHashable(index) {
				i.errors = append(i.errors, fmt.Errorf("字典键必须是可哈希类型"))
				return nil
			}
			val, exists := c[index]
			if !exists {
				i.errors = append(i.errors, fmt.Errorf("字典中不存在键: %v", index))
				return nil
			}
			originalValue = val
			target = &indexTarget{
				container: container,
				index:     index,
				isDict:    true,
			}

		default:
			i.errors = append(i.errors, fmt.Errorf("下标操作只支持列表或字典"))
			return nil
		}

	case *ast.Integer, *ast.Float:
		// 数字字面量自增自减
		fmt.Printf("自增自减操作不能用于数字字面量")

	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的左值类型: %T", left))
		return nil
	}

	utils.Debug("原始值: %v, 操作符: %s", originalValue, expr.Op)

	// 根据操作符计算新值
	var newValue Value
	switch expr.Op {
	case "++":
		newValue = i.increment(originalValue)
	case "--":
		newValue = i.decrement(originalValue)
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作符: %s", expr.Op))
		return nil
	}

	// 如果目标存在，更新值
	if target != nil {
		switch t := target.(type) {
		case *ast.Identifier:
			// 更新变量
			ctx.SetVar(t.Name, newValue)
			utils.Debug("更新变量 %s: %v -> %v", t.Name, originalValue, newValue)

		case *indexTarget:
			// 更新列表或字典元素
			if t.isList {
				idx := t.index.(int64)
				container := t.container.([]Value)
				container[idx] = newValue
				utils.Debug("更新列表元素[%d]: %v -> %v", idx, originalValue, newValue)
			} else if t.isDict {
				container := t.container.(DictType)
				container[t.index] = newValue
				utils.Debug("更新字典元素[%v]: %v -> %v", t.index, originalValue, newValue)
			}
		}
	}

	// 后置操作返回原始值
	return originalValue
}

// 辅助结构，用于存储下标目标
type indexTarget struct {
	container interface{}
	index     interface{}
	isList    bool
	isDict    bool
}

// 添加increment函数
func (i *Interpreter) increment(value Value) Value {
	switch v := value.(type) {
	case int64:
		return v + 1
	case float64:
		return v + 1.0
	default:
		i.errors = append(i.errors, fmt.Errorf("自增操作不支持的类型: %T", value))
		return value
	}
}

// 添加decrement函数
func (i *Interpreter) decrement(value Value) Value {
	switch v := value.(type) {
	case int64:
		return v - 1
	case float64:
		return v - 1.0
	default:
		i.errors = append(i.errors, fmt.Errorf("自减操作不支持的类型: %T", value))
		return value
	}
}
