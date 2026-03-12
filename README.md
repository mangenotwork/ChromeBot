谷歌浏览器（Chrome）自动化平台，通过输入指令或脚本(chrome bot script)自动执行操作Chrome，支持采样输出(页面html,监听Chrome信息等等),与Chrome交互由CDP协议实现;
还支持http请求(curl的功能)，系统级别的脚本交互(编写系统脚本)，常用与自动化场景（自动化测试，接口测试，爬虫脚本，系统脚本)和定制模拟操作场景(编写自动化操作机器人)等;


- 最初的想法

![最初的想法](_doc/ef1b8cdd-e942-482c-9b0c-fd2d01a468c4.jpg "最初的想法")

- 最终的目标

![最终的目标](_doc/a63349d6-7d19-4049-adfe-d32eb5f25665.jpg "最终的目标")


## 使用与环境
- 系统: 当前只支持 windows
- 浏览器: 只支持chrome,电脑需要自行安装chrome浏览器
- chromeBot.exe 免安装直接使用
    1. 运行 "chromeBot.exe" 进入REPL模式
    2. 运行 "chromeBot.exe 脚本.cbs" 执行cbs脚本

# ChromeBot Script v0.1
是基于ATS实现是DSL语言，主要用于ChromeBot执行自动化任务的脚本编写，脚本文件后缀为.cbs

## 运行
chromeBot.exe case.cbs

文档 ： [ChromeBot Script 文档](_doc/DSL_v0.1.md)

脚本示例

- [访问百度搜索交互](_examples/case_chrome_2.cbs)
- [访问东方财经网，先进官网再点击排行，获取排名列表并将数据保存到本地文件](_examples/eastmoney_1.cbs)
- [新华网获取要闻聚焦，点击每个要闻访问文章内容并截图保存在本地](_examples/xinhuanet_1.cbs)
- [央视新闻网，获取每个分类标签下的新闻列表打印新闻的标题出来](_examples/newscctv_1.cbs)
- [阳光高考获取全部的学校信息](_examples/gaokaochsi_1.cbs)
- todo...

例子

```cbs
// baidu.cbs  例子1 ：简单访问百度进行查询最后截图保存操作

// 打开浏览器访问百度
chrome init  
chrome req="www.baidu.com"

// 获取当前页面能输入的第一个输入框的xpath，输入 "ChromeBot"
var a = NowTabGetInputFirstXpath() 
print("获取当前页面能输入的第一个输入框的xpath :", a) 
chrome xpath=a input="ChromeBot"

// 获取当前页面内容为“百度一下”可交互的xpath, 执行点击操作
var b = NowTabMatchDemoContentOP("百度一下") 
print("获取当前页面内容为“百度一下”可交互的xpath :", b) 
chrome click=b 

// 截图保存到本地然后关闭浏览器
chrome screenshot="D:\baidu4.png"  
chrome close 
```

# 直接运行会进入ChromeBot REPL终端

## 运行

chromeBot.exe

```
_________
|       |
|  o o  |
|   c   |\
|_______| \_chrome

欢迎使用 ChromeBot v0.0.1
https://github.com/mangenotwork/ChromeBot
输入代码并按回车执行，按Ctrl+Z(Windows)退出
使用 'exit' 或 'quit' 命令退出程序
===================================================================
>>> a
a
>>> BayBay.
```

-h 
```

ChromeBot v0.0.1 ( https://github.com/mangenotwork/ChromeBot )

简介：
  谷歌浏览器（Chrome）自动化平台，通过输入指令或脚本(chrome bot script)自动执行操作Chrome

用法：
  chromebot [选项] [文件名]

选项：
  -h    查看 ChromeBot 帮助信息
  -v    查看 ChromeBot 版本信息

示例：
  chromebot          # 启动交互式 REPL 环境
  chromebot test.cbs # 执行 test.cbs 中的代码
  chromebot -v       # 查看版本信息
  chromebot -h       # 查看帮助信息
```

-v
```
v0.0.1
```