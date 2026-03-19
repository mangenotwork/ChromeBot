
#### v0.1.2



#### v0.1.1
- [] DSL v0.1, 继续完善语法文档和定版
- [] chrome 脚本操作自动化功能完整，能应对主流业务场景
- [] http+host 脚本操作自动化功能完成，能应对主流业务场景
- 写blog,录制视，4~5月看数据最后总结

#### v0.0.26
- 测试和优化

#### v0.0.25
- chrome  WebAudio ： 此域名允许查看 Web Audio API。
- chrome  WebAuthn ： 该域允许配置虚拟身份验证器来测试 WebAuthn API。

#### v0.0.24
- chrome  Tethering ： 域定义了浏览器端口绑定的方法和事件。
- chrome  Tracing ： 追踪

#### v0.0.23
- chrome  ServiceWorker ： 服务任务
- chrome  Storage ： 存储

#### v0.0.22
- chrome  Runtime ： 运行时域通过远程求值和镜像对象公开 JavaScript 运行时环境。
- chrome  Security ： 安全域

#### v0.0.21
- chrome  Profiler ： 分析器域
- chrome  PWA ： 该域允许与浏览器交互以控制 PWA。

#### v0.0.20
- chrome  PerformanceTimeline ： 按照https://w3c.github.io/performance-timeline/#dom-performanceobserver中的规定，报告性能时间线事件
- chrome  Preload ： 预加载域

#### v0.0.19
- chrome  Page ： 与被检查页面相关的操作和事件属于页面域。
- chrome  Performance ： 性能域

#### v0.0.18
- chrome  Network ： 网络域允许跟踪页面的网络活动。它公开有关 HTTP、文件、数据和其他请求和响应的信息，包括它们的标头、正文、时间等。
- chrome  Overlay ： 叠加域 该域提供与在被检查页面上绘制图形相关的各种功能。

#### v0.0.17
- chrome  Media ： 该域允许对媒体元素进行详细检查。
- chrome  Memory ： 内存相关

#### v0.0.16
- chrome  LayerTree ： 层树
- chrome  Log ： 提供对日志条目的访问权限。

#### v0.0.15
- chrome  Inspector ： 检查域
- chrome  IO ： 对 DevTools 生成的流进行输入/输出操作。

#### v0.0.14
- chrome  IndexedDB : IndexedDB相关的域
- chrome  Input ： 输入域

#### v0.0.13
- chrome  HeadlessExperimental : 此域提供仅在无头模式下支持的实验性命令。
- chrome  HeapProfiler : 堆分析器域

#### v0.0.12
- chrome  Fetch : 允许客户端使用客户端代码替换浏览器网络层的域。
- chrome  FileSystem : 文件系统域

#### v0.0.11
- chrome  Extensions : 定义浏览器扩展的命令和事件。
- chrome  FedCm : 该域允许与 FedCM 对话框进行交互。

#### v0.0.10
- chrome  Emulation : 该域名模拟了页面的不同环境。
- chrome  EventBreakpoints : 事件断点域 允许在 JavaScript 调用的原生代码中发生的操作和事件上设置 JavaScript 断点。

#### v0.0.9
- chrome  DOMSnapshot : 该域便于获取包含 DOM、布局和样式信息的文档快照。
- chrome  DOMStorage  : 查询和修改 DOM 存储。

#### v0.0.8
- 检查当前是否支持该协议域
- chrome  SystemInfo ： SystemInfo 域定义了用于查询底层系统信息的方法和事件
- chrome  Browser ： 浏览器域定义了用于管理浏览器的方法和事件。
- 完善 chrome Target : 目标对象
- chrome  BackgroundService ：  定义后台 Web 平台功能的事件。
- chrome  CacheStorage ： 缓存存储域
- 实现全局 @mysql  全局声明并连接mysql as到指定对象 todo 只设计语法
- 实现全局 @redis  全局声明并连接redis as到指定对象 todo 只设计语法
- http代理
- websocket客户端
- chrome  CSS ： 此域公开 CSS 的读写操作。
- chrome  Debugger ： 调试器域公开了 JavaScript 调试功能。 
- chrome  DOM : 此域公开 DOM 读/写操作。
- chrome  DOMDebugger : DOM调试允许在特定的DOM操作和事件上设置断点。

#### v0.0.7
- [] 系统文件相关交互方法
  [] 1. 查文件或目录
  [] 2. 创建文件或目录
  [] 3. 删除文件或目录
  [] 4. 移动文件或目录
  [] 5. 复制文件或目录
  [] 6. 读写文件
  [] 7. 文件或目录改名

- [] host 方法扩展  // 预计周五完成
  [] 1. disk 磁盘信息 使用率
  [] 2. mem 内存信息 使用率
  [] 3. cpu cpu信息 使用率
  [] 4. ping 
  [] 5. port 已开放的端口
  [] 6. pid 进程信息
  [] 7. zip 压缩解压文件

- [] 示例
  [] 1. 循环ping ip
  [] 2. 定期获取系统 disk mem cpu 达到监控的目的
  [] 3. 查看本机已开放的端口将数据存储到excel
  [] 4. 利用zip压缩文件达到备份的目的   // 预计周六完成

- [] 改测试的bug和优化
- [] bug和优化验收
- [] 更新文档


#### v0.0.6
- [] 支持excel相关操作方法 底层使用 https://github.com/qax-os/excelize 库
  [] 1. 读取excel
  [] 2. 写入excel
  [] 3. 最小粒度的单元格操作
  [] 4. 图片插入
  [] 5. 单元格样式设置（字体、颜色、对齐）
  [] 6. 合并单元格
  [] 7. 添加图表
  [] 8. 设置单元格公式  
  
- [] 支持json与字典互转方法   // 周四完成
  [] 1. json转字典
  [] 2. 字典转json
  [] 3. json查找元素方法
  [] 4. 读写json文件

- [] 改测试的bug和优化
- [] bug和优化验收
- [] 更新文档

#### v0.0.5
- [ok] 设计脚本配置 以 @ 开头, 全局的程序执行前就需要处理，优先解析并保存在全局常量，只有脚本模式下才有
  方案，在脚本解析前，优先判断@开头的行，如果符合以下全局配置则记录到全局配置常量，否则进行报错，脚本过完后将每行的第一个 @改为#，然后执行后面的语法检查
  [ok] 1. @cron 设置定时执行脚本,语法参考 cron 核心定时参数总览
  [ok] 2. @conf_json 设置外部配置文件json,读取json文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用  @conf_json path="" as=conf
  [ok] 3. @conf_yaml 设置外部配置文件yaml,读取yaml文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用  @conf_ini path="" as=conf
  [ok] 4. @conf_ini 设置外部配置文件ini,读取ini文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用   @conf_ini path="" as=conf
  [ok] 5. @chrome_check 全局检查是否支持chrome浏览器，提取检查，如果宿主机未安装会提前检查出来 语法 : 直接使用  @chrome_check
  [ok] 6. @network_check 全局检查网络是否OK, 如果当前宿主机未网络会提前检查出来 语法: 请求地址可以是ip也可以是域名 如: @network_check "www.baidu.com"   @network_check "254.254.254.254" 

- [ok] 实现全局定时任务
- [ok] 实现全局配置 @conf_json, @conf_yaml, @conf_ini
- [ok] 实现全局 @chrome_check
- [ok] 实现全局 @network_check
- [ok] host 宿主机的相关方法 增加 host 关键字
  1. host name
  2. host ip
  3. host info
  4. host to
- [ok] DNS查询方法
- [ok] Whois查询方法
- [ok] 端口扫描方法
- [ok] 网站死链检查
- [ok] 网站证书信息检查
- [ok] 扫描网站的url, 扫描网站的外链
- [] 示例
  1. 定期监控网站证书到期时间
  2. 定期对网站进行死链检查
  3. 定期扫描端口
  4. case_cron.cbs 定期chrome交互网站
- [] 改测试的bug和优化
- [] bug和优化验收
- [] 更新文档 (主要是host,和几个新增的方法)  

#### v0.0.4
- [ok] 系统级别的交互确认弹窗
- [ok] <bug转需求> 本地系统记录chrome userpath进程，如果存在默认隔离环境的chrome进程提供交互选择关闭还是新建隔离环境
- [ok] chrome init 新增new指令初始化隔离环境
- [ok] <bug转需求> chrome连不上的时候重试
- [ok] Xpath格式检查相关的函数
- [ok] 语法: 支持 \ 作为代码的换行连接符，与命令行的换行一样
- [ok] 语法: 命令式值能绑定上函数并执行函数使用函数的返回值 如: chrome click=NowTabMatchDemoContentOP("百度一下")
- [ok] chrome 非法参数应该提示
- [ok] 支持输出chrome信息
- [ok] 如果本地未安装chrome提示安装
- [ok] <bug转需求> tab失焦后进行提示, 提供交互确认，如果是继续则默认选择第一个tab,如果是结束就结束脚本
- [ok] 测试并新增示例
  1. [ok] doubao_1.cbs 交互豆包问豆包问题
  2. [ok] zzttop_1.cbs 交互网站手机号随机生成，采集将生成的手机号
  3. [ok] toutiao_1.cbs 头条采集今日要闻访问后截图保存到本地
  4. [ok] douyin_1.cbs 抖音交互下滑视频
  
- [ok] 改测试的bug和优化
- [ok] bug和优化验收 
- [ok] 更新文档 

#### v0.0.3
- [ok] 谷歌浏览器打开方法以及语句 增加 chrome 关键字
  1. 打开窗口以及页签
  2. 访问网址
  3. 新开一个页签
  4. 查看当前页签
  5. 关闭页签
  6. 切换页签
 
- [ok] 输出html以及语句
  1. 输出页面的html
  2. 解析HTML demo树
  
- [ok] 操作点击方法以及语句
  1. 点击xpath
  2. 输入xpath

- [ok] 检查操作，检查页面是否存在指定xpath
- [ok] 滚动操作，滚动页面
- [ok] 截图操作，浏览器截图操作

- [ok] 新增xpath相关的函数 
  1. 打印当前页面的demo数(含xpath)
  2. 匹配页面标签属性值获取对应xpath  
  3. 匹配页面内容获取对应的xpath

- [ok] 等待页面完全响应完
- [ok] 整理日志打印
- [ok] 多增加实例的测试用例以及测试 
  1. [ok] eastmoney_1.cbs  东方财经网，先进官网再访问排名，获取排名列表并保存到本地文件
  2. [ok] xinhuanet_1.cbs  新华网获取要闻聚焦，点击每个要闻访问文章内容并截图保存在本地
  3. [ok] newscctv_1.cbs  央视新闻网，获取每个分类标签下的新闻列表打印新闻的标题出来
  4. [ok] gaokaochsi_1.cbs 阳光高考获取全部的学校信息
  
- [ok] 改测试的bug和优化
- [ok] bug和优化验收
- [ok] 更新文档

#### v0.0.2
- [ok] http请求方法 增加 http 关键字
  1. http参数解析
  2. http请求的功能
  3. 请求后返回值的存储
  4. 关闭gt的日志
- [ok] 字符处理方法
- [ok] 正则方法
- [ok] 时间相关方法
- [ok] http增加save方法，将请求返回的body保存到本地
- [ok] 测试内置函数，语法层面的
- [ok] 除了基础测试的代码，多增加实例的测试用例
- [ok] 改测试的bug和优化
- [ok] 更新文档

#### v0.0.1
- [ok] 项目目录结构
- [ok] 开始函数与REPL
- [ok] DSL v0.1, 基础语法
- [ok] DSL v0.1, 测试语法，更多的示例
- [ok] DSL v0.1, 语法文档
- [ok] 改测试的bug和优化
- [ok] bug和优化验收
- [ok] 更新文档  支持了 elif, switch, for in, while in, 自增自减

## 里程碑
1. 第一阶段实现 DSL 实现脚本编写操作浏览器自动化任务
2. 第二阶段大量产出测试用例给AI进行训练
3. 训练定向模型，最终做到 自热语言 -> DSL脚本 -> 自动化操作浏览器

## 需求池
- 脚本命令式应该支持 output参数来指定数据输出到哪里（支持文件，数据库，远端服务等等）

----
 
浏览器自动化相关资料

## 浏览器自动化操作
点击操作	左键单击 / 双击、右键单击、坐标点击、模拟点击（避免被反爬检测）	触发按钮 / 链接 / 复选框等交互（如登录按钮、勾选同意协议）
输入操作	输入文本、清空输入框、模拟键盘输入（回车 / 退格 / 快捷键）	表单填写（用户名 / 密码 / 搜索框）、触发搜索（输入后按回车）
选择操作	下拉框选择（按索引 / 值 / 文本）、勾选 / 取消勾选复选框 / 单选框	表单中的下拉选择、同意条款勾选
拖拽操作	元素拖拽（如滑块验证）、文件拖拽上传	滑块验证码、文件上传、拖拽排序
元素状态验证	检查元素是否可见 / 可点击 / 存在 / 被选中、获取元素属性 / 文本 / 尺寸 / 位置	验证操作是否生效（如输入框是否有值、按钮是否可点击）
元素滚动	滚动到元素可见位置、页面滚动（向上 / 向下 / 到顶部 / 到底部）	操作不在可视区域的元素（如页面底部的提交按钮）
元素截图	单独截取元素区域、高亮元素后截图	验证元素显示是否正确、留存测试证据


#### 元素定位

定位方式	语法示例（Selenium/Playwright）	适用场景
ID 定位	find_element(By.ID, "username") / page.locator("#username")	元素有唯一 ID（最稳定）
名称（Name）定位	find_element(By.NAME, "password") / page.locator("[name='password']")	表单元素（输入框 / 按钮）常用 name 属性
类名（Class）定位	find_element(By.CLASS_NAME, "btn-submit") / page.locator(".btn-submit")	按样式类定位，注意类名可能重复
标签名定位	find_element(By.TAG_NAME, "input") / page.locator("input")	批量定位同类型元素（如所有输入框）
XPath 定位	find_element(By.XPATH, "//input[@id='username']") / page.locator("//input[@id='username']")	万能定位，支持复杂路径 / 逻辑（如父节点 / 兄弟节点 / 文本匹配）
CSS 选择器定位	find_element(By.CSS_SELECTOR, "#username") / page.locator("#username")	高效定位，支持层级 / 属性 / 伪类（如 input[type='text']:first-child）
链接文本定位	find_element(By.LINK_TEXT, "登录") / page.locator("text=登录")	定位超链接（<a>标签）
部分链接文本	find_element(By.PARTIAL_LINK_TEXT, "登")	链接文本过长时模糊匹配
自定义定位	基于文本 / 正则 / 相对位置（Playwright 的locator.filter()）	复杂场景（如 “包含‘提交’且 class 为 btn 的按钮”）

#### 页面导航与加载

操作类型	具体操作	用途
页面跳转	访问指定 URL（get(url)/goto(url)）、后退 / 前进 / 刷新页面	核心导航，如打开目标网站、返回上一页
加载控制	等待页面加载完成（等待 DOM / 网络 / 渲染完成）、中断页面加载	避免操作过早执行（如页面未加载完就点击按钮）
页面信息获取	获取页面 URL / 标题 / 源代码 / 截图 / HTML、获取页面性能数据（加载时间）	验证页面是否正确跳转、保存页面快照、分析性能
页面状态控制	最大化 / 最小化 / 调整窗口大小、全屏模式、设置页面缩放比例	适配不同分辨率、模拟移动端窗口（如 375x667）


#### 会话 / 驱动管理

操作类型	具体操作	用途
驱动初始化	启动浏览器（Chrome/Firefox/Edge）、配置启动参数（无头模式 / 窗口大小 / 代理）	初始化自动化环境，如 puppeteer.launch({headless: true})、webdriver.Chrome()
会话控制	新建标签页 / 窗口、关闭标签页 / 窗口、切换标签页 / 窗口、获取当前窗口句柄	多标签 / 多窗口场景（如同时操作多个页面）
浏览器配置	设置 Cookie / 本地存储 / 会话存储、清除缓存 / 历史记录、设置用户代理（UA）	模拟不同用户环境、绕过反爬、持久化登录状态
退出驱动	关闭浏览器、释放驱动资源	自动化结束后清理环境，避免进程残留


#### 键盘与鼠标模拟

操作类型	具体操作	用途
键盘模拟	按下 / 释放单个按键（KeyDown/KeyUp）、组合键（Ctrl+C/Ctrl+V）、输入特殊字符	快捷键操作（复制 / 粘贴）、模拟用户打字节奏、触发键盘事件
鼠标模拟	鼠标移动（到元素 / 坐标）、鼠标按下 / 释放、滚轮滚动（向上 / 向下）	模拟用户鼠标移动、悬浮提示框、滚轮翻页
悬浮操作	鼠标悬浮到元素上（触发 hover 事件）	显示悬浮菜单、验证悬浮提示文本


#### 弹窗处理

弹窗类型	具体操作	用途
警告弹窗（Alert）	接受弹窗（OK）、获取弹窗文本	处理系统警告（如 “操作成功” 提示）
确认弹窗（Confirm）	接受 / 取消弹窗、获取弹窗文本	处理确认操作（如 “是否删除”）
提示弹窗（Prompt）	输入文本后确认、取消弹窗、获取弹窗提示文本	处理需要输入的弹窗（如 “请输入验证码”）
浏览器弹窗	处理下载弹窗、文件上传弹窗、权限请求弹窗（摄像头 / 麦克风 / 通知）	下载文件、上传文件、允许 / 拒绝权限请求
自定义弹窗	定位并关闭自定义弹窗（如广告弹窗、登录弹窗）	关闭干扰操作的弹窗（如页面广告）


#### 网络请求与响应

操作类型	具体操作	用途
网络监听	监听所有网络请求 / 响应、过滤指定 URL / 方法（GET/POST）	抓包分析、获取接口数据、验证接口返回
网络拦截	拦截指定请求、修改请求参数 / 响应内容、模拟接口返回（Mock）	测试异常场景（如接口返回 500）、绕过接口限制、加速测试（Mock 数据）
网络等待	等待指定请求完成、等待所有请求完成	确保接口数据加载完成后再操作（如列表数据加载）
Cookie 操作	添加 / 获取 / 删除 Cookie、设置 Cookie 有效期 / 域名 / 路径	持久化登录状态、模拟不同用户登录
本地存储操作	操作 LocalStorage/SessionStorage（添加 / 获取 / 删除 / 清空）	模拟前端存储数据、验证存储逻辑


#### JavaScript 执行

操作类型	具体操作	用途
执行 JS 代码	执行任意 JavaScript 代码、获取 JS 执行结果	操作前端无法通过常规 API 实现的功能（如修改元素样式、触发自定义事件）
注入 JS 脚本	向页面注入外部脚本（如 jQuery）、修改页面全局变量	增强页面交互能力、绕过前端限制
异步 JS 执行	执行异步 JS 代码（如等待 Promise 完成）、获取异步执行结果	处理前端异步逻辑（如接口请求完成后获取数据）


#### 文件操作

操作类型	具体操作	用途
文件上传	定位文件上传输入框、传入本地文件路径	表单中的文件上传（如头像上传、文档上传）
文件下载	设置下载路径、等待文件下载完成、验证下载文件（存在 / 大小 / 内容）	下载报表、文件后验证完整性
截图 / 录屏	页面全屏截图、指定区域截图、录制页面操作视频	测试结果留存、问题复现、自动化报告附件


#### 高级控制

操作类型	具体操作	用途
框架 /iframe 切换	切换到指定 iframe、切回主文档、嵌套 iframe 切换	操作 iframe 内的元素（如支付页面、内嵌表单）
多窗口 / 标签页	新建窗口 / 标签页、切换窗口 / 标签页、关闭窗口 / 标签页、获取所有窗口句柄	多页面并行操作（如同时登录多个账号）
超时设置	设置全局超时、元素定位超时、页面加载超时、脚本执行超时	避免自动化卡死（如元素长时间未出现时超时退出）
异常处理	捕获操作异常、重试失败操作、忽略指定异常	提高自动化稳定性（如网络波动导致的点击失败重试）
代理设置	配置 HTTP/HTTPS/SOCKS 代理、切换代理、验证代理有效性	模拟不同地区访问、绕过 IP 限制
认证处理	处理 HTTP 基本认证、OAuth 认证、验证码识别（对接打码平台）	登录需要认证的网站、处理验证码（如滑块 / 图片验证码）


## 其他文案

- ChromeBot 的作用： 专注与自动化操作浏览器脚本的编写，未减少心智脚本语言cbs语法设计简单为首要目标，然后满足所有的脚本语法，支持命令式语法，只要是有编程语言基础的人能快速上手使用，降低学习成本。
- ChromeBot, 最有帮助的是什么? 最有价值的是什么？最有趣的是什么？最值得推广的是什么？
- ChromeBot核心用户群体: 测试，爬虫，自动化任务，爱好学习者，炫技者，日常电脑办公
- ChromeBot 核心是编写自动化脚本来完成自动化任务，首要支持是Chrome浏览器的自动化，其次是支持http相关方法(如:curl的功能)，系统的脚本功能，读写多存储应用(excel,mysql,redis....)

