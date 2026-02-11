package ast

import (
	"fmt"
	"strings"
)

// 扩展提供交互的关键字定义

// ChromeStmt chrome 关键字，操作chrome
type ChromeStmt struct {
	StartPos Position
	Args     []Expression
}

func (c *ChromeStmt) Pos() Position { return c.StartPos }
func (c *ChromeStmt) String() string {
	args := make([]string, len(c.Args))
	for i, arg := range c.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("chrome %s ", strings.Join(args, " "))
}
func (c *ChromeStmt) stmtNode() {}

// HttpStmt http 关键字,http相关操作
type HttpStmt struct {
	StartPos Position
	Args     []Expression
}

func (c *HttpStmt) Pos() Position { return c.StartPos }
func (c *HttpStmt) String() string {
	args := make([]string, len(c.Args))
	for i, arg := range c.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("http %s ", strings.Join(args, " "))
}
func (c *HttpStmt) stmtNode() {}
