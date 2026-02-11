package parser

import (
	"ChromeBot/dsl/ast"
	"ChromeBot/dsl/lexer"
	"ChromeBot/utils"
	"strings"
)

// Chrome 语法解析 chrome arg1 arg2=123 ...
func (p *Parser) parseChromeStatement() *ast.ChromeStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	utils.Debug("======= parseChromeStatement 开始 =======")
	utils.Debug("当前token: %v", p.curTok)

	// 保存起始位置
	startPos := ast.Position{
		Line:   p.curTok.Line,
		Column: p.curTok.Column,
	}

	// 跳过 chrome 关键字
	p.nextToken()
	utils.Debug("跳过chrome后: %v", p.curTok)

	var args []ast.Expression
	startLine := p.curTok.Line

	// 读取chrome参数
	for p.curTok.Line == startLine && !p.curTokenIs(lexer.TokenEOF) {
		utils.Debug("解析参数，当前token: %v", p.curTok)

		// 构建参数字符串
		argStr := p.readChromeArgs()
		if len(argStr) != 0 {
			for _, arg := range argStr {
				args = append(args, &ast.String{
					StartPos: ast.Position{
						Line:   p.curTok.Line,
						Column: p.curTok.Column,
					},
					Value: arg,
				})
			}

		}

		// 跳过逗号
		if p.curTokenIs(lexer.TokenComma) {
			p.nextToken()
		}
	}

	utils.Debug("parseChromeStatement: 完成，共 %d 个参数", len(args))

	return &ast.ChromeStmt{
		StartPos: startPos,
		Args:     args,
	}
}

func (p *Parser) readChromeArgs() []string {
	var args []string
	startLine := p.curTok.Line
	var currentArg strings.Builder

	// 记录是否在等号表达式中
	inKeyValue := false

	for p.curTok.Line == startLine && !p.curTokenIs(lexer.TokenEOF) {
		token := p.curTok

		// 跳过逗号
		if token.Type == lexer.TokenComma {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
				inKeyValue = false
			}
			p.nextToken()
			continue
		}

		// 遇到等号
		if token.Type == lexer.TokenAssign {
			currentArg.WriteString(token.Literal)
			inKeyValue = true
			p.nextToken()
			continue
		}

		// 普通token
		if currentArg.Len() == 0 {
			// 参数开始
			currentArg.WriteString(token.Literal)
		} else if inKeyValue {
			// 在等号表达式中，直接连接
			currentArg.WriteString(token.Literal)
			inKeyValue = false
		} else {
			// 新参数开始
			args = append(args, currentArg.String())
			currentArg.Reset()
			currentArg.WriteString(token.Literal)
		}

		p.nextToken()
	}

	// 添加最后一个参数
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}

// http 语法解析 http arg1 arg2=123 ...
func (p *Parser) parseHttpStatement() *ast.HttpStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	utils.Debug("======= parseHttpStatement 开始 =======")
	utils.Debug("当前token: %v", p.curTok)

	// 保存起始位置
	startPos := ast.Position{
		Line:   p.curTok.Line,
		Column: p.curTok.Column,
	}

	// 跳过 http 关键字
	p.nextToken()
	utils.Debug("跳过http后: %v", p.curTok)

	var args []ast.Expression
	startLine := p.curTok.Line

	for p.curTok.Line == startLine && !p.curTokenIs(lexer.TokenEOF) {
		utils.Debug("解析参数，当前token: %v", p.curTok)

		// 构建参数字符串
		argStr := p.readChromeArgs()
		if len(argStr) != 0 {
			for _, arg := range argStr {
				args = append(args, &ast.String{
					StartPos: ast.Position{
						Line:   p.curTok.Line,
						Column: p.curTok.Column,
					},
					Value: arg,
				})
			}

		}

		// 跳过逗号
		if p.curTokenIs(lexer.TokenComma) {
			p.nextToken()
		}
	}

	utils.Debug("parseHttpStatement: 完成，共 %d 个参数", len(args))

	return &ast.HttpStmt{
		StartPos: startPos,
		Args:     args,
	}
}
