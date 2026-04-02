# ChromeBot Script 
是基于ATS实现是DSL语言，主要用于ChromeBot执行自动化任务的脚本编写，脚本文件后缀为.cbs

## 运行
chromeBot.exe case.cbs

## 语法

### SDL语法设计
- 注重脚本实现过程专注于流程精准表达目的为目标，语法要追求简单，操作多以命令式语法，省去了函数式编程和面向对象编程
- 由于主要编写执行自动化的脚本，设计极简，只有数值类型、字符串类型、布尔类型、列表类型、字典类型（键值对类型）
- 支持逻辑判断(if else elif), 支持循环(for  while), 支持分支(switch case), 支持运算
- 注释为 # 或 //
- 很多内置函数，编写脚本需参照文档找到需要的方法
- 支持函数链式调用，函数依次调用返回值会自动传递给下个函数
- 列表类型、字典类型（键值对类型）都是以下标取值
- 支持for in, while in 遍历列表或字典
- http关键字，这个关键字执行所有http相关的操作
- chrome的操作关键字，用这些关键字命令式语法编写脚本来操作浏览器
- 每个方法每行指令脚本必须要有反馈，并且反馈还必须提供信息，必须准确合理

### 变量

与大多数脚本语言一样使用关键字 var 来声明变量

```cbs
var a = 1
b = 1
```

#### 数值类型

顾名思义就是整型+浮点类型的数值类型

```cbs
1
1.1
0
0.0001
999999999
```

#### 字符串类型

支持双引号和单引号

```cbs
"asdasdasd"
"escaped \"quote\" \\test \n\t 中文\""
""
'asdasdsad'
```

#### 布尔类型

true 和 false 

#### 列表类型

用中括号括起，元素用逗号分隔，元素可以是数值、字符串、布尔、列表
如果元素为列表，即为多维列表
用下标获取列表的元素，遍历多使用循环
内置函数 len可获取列表长度

```cbs
[1,"a", true]
[[1,2,3],[1,2,3]]
len([1,2,3])
```

#### 字典类型(键值对类型)

用大括号括起，键值对用逗号分隔，键和值用冒号分隔
键必须是字符串，值可以是任意类型（数值、字符串、布尔、列表、字典）
支持点语法或中括号语法访问字典中的值

```cbs
var person = { "name": "张三", "age": 25, "isStudent": false }
var config = { "timeout": 30, "retry": 3, "headers": { "User-Agent": "ChromeBot" } }

# 访问方式
person.name          #  "张三"
person["age"]        #  25
config.headers       #  { "User-Agent": "ChromeBot" }
config["timeout"]    #  30
```

### 运算

+ - * / % ++ --

### 逻辑判断

支持 if、else if、else 结构，条件表达式结果为布尔类型

```cbs
var score = 85
if score >= 90 {
    print("优秀")
} elif score >= 60 {
    print("及格")
} else {
    print("不及格")
}

# 支持嵌套
var a = 10
var b = 20
if a > 5 {
    if b > 15 {
        print("a>5 且 b>15")
    }
}
```

### 循环


for 循环, 遍历列表或字典，或执行指定次数的循环

```cbs
# 遍历列表
var list = [1, 2, 3, 4, 5]
for item in list {
    print(item)
}

# 遍历字典
var dict = { "a": 1, "b": 2, "c": 3 }
for key, value in dict {
    print(key + ": " + value)
}

# 指定次数循环
for i = 0; i < 5; i = i + 1 {
    print("第 " + i + " 次循环")
}
```

while 循环, 满足条件时重复执行代码块

```cbs
var count = 0
while count < 5 {
    print("count: " + count)
    count = count + 1
}

# 无限循环（需谨慎使用）
while true {
    # 执行某些操作
    break  # 可使用 break 退出循环
}
```

### 全局指令与全局常量

- @cron 设置定时执行脚本,语法参考 cron 核心定时参数总览
```
@cron 0 0 0 * * *
```

- @conf_json 设置外部配置文件json,读取json文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用  @conf_json path="" as=conf
```
@conf_json path="./_examples/case_conf.json" as=conf
```

- @conf_yaml 设置外部配置文件yaml,读取yaml文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用  @conf_ini path="" as=conf
```
@conf_yaml path="./_examples/case_conf.yaml" as=conf
```

- @conf_ini 设置外部配置文件ini,读取ini文件后将值存储到 as到指定的全局字典常量,以供脚本全局使用   @conf_ini path="" as=conf
```
@conf_ini path="./_examples/case_conf.ini" as=conf
```

- @chrome_check 全局检查是否支持chrome浏览器，提取检查，如果宿主机未安装会提前检查出来 语法 : 直接使用  @chrome_check
```
@chrome_check
```

- @network_check 全局检查网络是否OK, 如果当前宿主机未网络会提前检查出来 语法: 请求地址可以是ip也可以是域名 
```
@network_check "www.baidu.com" 
@network_check "254.254.254.254"
```


### 链式调用

支持将多个函数调用连接在一起，上一个函数的返回值自动作为下一个函数的参数

```cbs
var s = "test";
var s1 = upper(s).repeat(2)
print(s1)
```



### 内置函数

- print(arg1,arg2....) : 在终端输出打印
```cbs
print("hello", " ", "word")
```

- int(arg) : 类型转换 数值字符串转换数值类型
```
var a = "11"
var b = int(a)
print(b)
```

- str(arg) 类型转换 转换为字符串类型
```cbs
var a = 11
var b = str(a)
print(b)
```

- len(arg): 获取传入类型的长度，arg是任意类型，返回长度
```cbs
var list = [1, 2, 3, 4, 5]
print(len(list))
```

- keys(arg)  获取字典的keys
```cbs
var d1 = {"one": 1, "two": 2}
print(keys(d1))
```

- values(arg)  获取字典的values
```cbs
var d1 = {"one": 1, "two": 2}
print(values(d1))
```

- items  获取所有键值对（每个键值对是一个包含两个元素的列表）
```cbs
var d1 = {"one": 1, "two": 2}
print(items(d1))
```

- has 字典或列表是否存在元素, arg第一个是字典或列表， 第二个是要找的元素
```cbs
var d1 = {"one": 1, "two": 2}
print(has(d1, "one"))
print(has(d1, "aa"))
```

- delete 删除字典或列表的指定元素, arg第一个是字典或列表， 第二个是要找的元素
```cbs
var d2 = {"one": 1, "two": 2}
delete(d2, "one")
print(d2)
```

- type_of 获取变量类型
```cbs
var d1 = {"one": 1, "two": 2}
print(type_of(d1))
```

- append(list, item) 给List增加元素
```
var a = []
a = append(a, {"a":1})
a = append(a, {"a":2})
print(a)
```

- exit() 退出程序


### 内置函数 - 数学方法 math

- abs 计算绝对值
```cbs
abs(-1)
```

- max 计算最大值
```cbs
max(1,2)
```

- min 计算最小值
```cbs
min(1,2)
```

### 内置函数 - 字符串方法 str

-  upper 将参数转换为字符串并转为大写
```cbs
print("aa".upper())
```

- repeat 将字符串进行重复, 第二个参数必须是整数
```cbs
print("aa".repeat())
```

- lower 字符串转小写
```cbs
lower("aaaa")
```

- trim 取首字符
```cbs
trim("abc")
```

- split 字符分割
```cbs
split("a/b", "/")
```

- replace 字符串替换
```cbs
replace("hello@word", "@", " ")
```

- replaceN 字符串替换 指定替换几个
```cbs
replaceN("hello@word@", "@", " ", 1)
```

- CleanWhitespace 函数 清理字符串回车，换行符号，还有前后空格
```cbs
CleanWhitespace("\n hello ")
```

- StrDeleteSpace 函数 删除字符串前后的空格
```cbs
StrDeleteSpace("  hello ")
```

- UnicodeDec 函数 字符串进行unicode编码
```cbs
UnicodeDec("hello")
```

- UnescapeUnicode 函数 字符串进行unicode解码
```cbs
UnescapeUnicode("hello")
```

- Base64Encode 函数 字符串进行base64编码
```cbs
Base64Encode("hello")
```

- Base64Decode 函数 字符串进行base64解码
```cbs
Base64Decode("aGVsbG8=")
```

- UrlBase64Encode 函数 url进行base64编码
```cbs
UrlBase64Encode("www.baidu.com?t=123")
```

- UrlBase64Decode 函数 url进行base64解码
```cbs
UrlBase64Decode("d3d3LmJhaWR1LmNvbT90PTEyMw==")
```

- MD5 函数 将字符串进行md5
```cbs
MD5("hello")
```

- MD516 函数 将字符串进行md5，返回16位
```cbs
MD516("hello")
```

- GBKToUTF8 函数 将GBK编码的字符串转换为utf-8编码
```cbs
GBKToUTF8("hello")
```

- UTF8ToGBK 函数 将utf-8编码的字符串转换为GBK编码
```cbs
UTF8ToGBK("hello")
```

- TF8ToGB2312 函数 将UTF-8转换为GB2312
```cbs
UTF8ToGB2312("hello")
```

- GB2312ToUTF8 函数 将GB2312转换为UTF-8
```cbs
GB2312ToUTF8("hello")
```

- UTF8ToGB18030 函数 将UTF-8转换为GB18030
```cbs
UTF8ToGB18030("hello")
```

- GB18030ToUTF8 函数 将GB18030转换为UTF-8
```cbs
GB18030ToUTF8("hello")
```

- UTF8ToBIG5 函数 将UTF-8转换为BIG5
```cbs
UTF8ToBIG5("hello")
```

- BIG5ToUTF8 函数 将BIG5转换为UTF-8
```cbs
BIG5ToUTF8("hello")
```

- UTF8ToLatin1 函数 将UTF-8转换为ISO-8859-1（Latin1）
```cbs
UTF8ToLatin1("hello")
```

- Latin1ToUTF8 函数 将ISO-8859-1转换为UTF-8
```cbs
Latin1ToUTF8("hello")
```

- Reg 函数 字符串正则 第一个参数是字符串，第二个参数是正则串
```cbs
Reg("<a>1</a>23<a>4</a>", `(?is:<a.*?</a>)`)
```

- RegHtml 函数 用正则提取html 第一个参数是html字符串，第二个是标签
```cbs
RegHtmlText("<a>1</a>23<a>4</a>", "a")

全部标签: 
"a", "title", "keywords", "description", "tr", "input", "td", "p", "span", "src", "href", "h1", "h2", "h3", "h4", "h5",
"h6", "tbody", "video", "canvas", "code", "img", "ul", "li", "meta", "select", "table", "button", "tableOnly", "div",
"option",
```

- RegHtmlText 函数 用正则提取html只匹配标签内的文本部分 第一个参数是html字符串，第二个是标签名
```cbs
RegHtmlText("<a>1</a>23<a>4</a>", "a")

全部标签: 
"a", "title", "keywords", "description", "tr", "input", "td", "p", "span", "src", "href", "h1", "h2", "h3", "h4", "h5",
"h6", "tbody", "video", "canvas", "code", "img", "ul", "li", "meta", "select", "table", "button", "tableOnly", "div",
"option",
```

- RegFn 函数 内置了很多用正则提取的常用场景方法 第一个参数是字符串，第二个是方法名
```cbs
RegFn("aaa127.0.0.1aa", "RegIPv4")

全部方法名以及对应的正则: 
"RegTime":         `(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`,
"RegLink":         `(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`,
"RegEmail":        `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`,
"RegIPv4":         `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`,
"RegIPv6":         `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`,
"RegIP":           `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`,
"RegMD5Hex":       `[0-9a-fA-F]{32}`,
"RegSHA1Hex":      `[0-9a-fA-F]{40}`,
"RegSHA256Hex":    `[0-9a-fA-F]{64}`,
"RegGUID":         `[0-9a-fA-F]{8}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{12}`,
"RegMACAddress":   `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`,
"RegEmail2":       `^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$`,
"RegUUID3":        `^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`,
"RegUUID4":        `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
"RegUUID5":        `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
"RegUUID":         `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
"RegInt":          `^(?:[-+]?(?:0|[1-9][0-9]*))$`,
"RegFloat":        `^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$`,
"RegRGBColor":     `^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$`,
"RegFullWidth":    `[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]`,
"RegHalfWidth":    `[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]`,
"RegBase64":       `^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`,
"RegLatitude":     `^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$`,
"RegLongitude":    `^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$`,
"RegDNSName":      `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`,
"RegFullURL":      `^(?:ftp|tcp|udp|wss?|https?):\/\/[\w\.\/#=?&]+$`,
"RegURLSchema":    `((ftp|tcp|udp|wss?|https?):\/\/)`,
"RegURLUsername":  `(\S+(:\S*)?@)`,
"RegURLPath":      `((\/|\?|#)[^\s]*)`,
"RegURLPort":      `(:(\d{1,5}))`,
"RegURLIP":        `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`,
"RegURLSubdomain": `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`,
"RegWinPath":      `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`,
"RegUnixPath":     `^(/[^/\x00]*)+/?$`,
```

- RegDel 函数 常见的删除方法支持html删除指定标签内容 第一个参数是字符串，第二个是方法名或标签名
```cbs
RegDel("<a>1</a>23<a>4</a>", "a")

全部标签: 
"html","number","a","title","tr","input","td","p","span","src","href","video","canvas","code","img","ul","li","meta",
"select","table","button","h1","h2","h3","h4","h5","h6","tbody",
```

- RegHas 函数 使用正则判断是否存在某内容 第一个参数是字符串，第二个是方法名
```cbs
RegHas("123", "IsNumber")

全部方法名以及正则:
"IsNumber":        `^[0-9]*$`,
"IsNumber2Len":    `[0-9]{%d}`,
"IsNumber2Heard":  `^(%d)[0-9]*$`,
"IsFloat":         `^(-?\d+\.\d+)?$`,
"IsFloat2Len":     `^(-?\d+\.\d{%d})?$`,
"IsEngAll":        `^[A-Za-z]*$`,
"IsEngLen":        `^[A-Za-z]{%d}$`,
"IsEngNumber":     `^[A-Za-z0-9]*$`,
"IsLeastNumber":   `[0-9]{%d,}?`,
"IsLeastCapital":  `[A-Z]{%d,}?`,
"IsLeastLower":    `[a-z]{%d,}?`,
"IsLeastSpecial":  `[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\&\$\.\*\+\?\{\}\(\)\[\]\|\!\_\@\#\%\-\=]{%d,}?`,
"HaveNumber":      `[0-9]+`,
"HaveSpecial":     `[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\&\$\.\*\+\?\{\}\(\)\[\]\|\!\_\@\#\%\-\=]+`,
"IsEmail":         `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`,
"IsDomain":        `[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(/.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+/.?`,
"IsURL":           `//([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$`,
"IsPhone":         `^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`,
"IsLandline":      `^(\(?\d{3,4}-)?\d{7,8}$`,
"AccountRational": `^[a-zA-Z][a-zA-Z0-9_]{4,15}$`,
"IsXMLFile":       `^.*\.[xX][mM][lL]$`,
"IsUUID3":         `^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`,
"IsUUID4":         `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
"IsUUID5":         `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
"IsRGB":           `^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$`,
"IsFullWidth":     `[^\x{0020}-\x{007E}\x{FF61}-\x{FF9F}\x{FFA0}-\x{FFDC}\x{FFE8}-\x{FFEE}0-9a-zA-Z]`,
"IsHalfWidth":     `[\x{0020}-\x{007E}\x{FF61}-\x{FF9F}\x{FFA0}-\x{FFDC}\x{FFE8}-\x{FFEE}0-9a-zA-Z]`,
"IsBase64":        `^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`,
"IsLatitude":      `^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$`,
"IsLongitude":     `^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$`,
"IsDNSName":       `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`,
"IsIPv4":          `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`,
"IsWindowsPath":   `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`,
"IsUnixPath":      `^(/[^/\x00]*)+/?$`,
```

### 内置函数 - 时间相关的方法 time

- now 获取当前时间的时间戳
```cbs
now()
```

- sleep 休眠 单位ms
```cbs
sleep(2)
```

-  Timestamp 时间戳
```cbs
Timestamp()
```

- TimestampMilli 时间戳 milliseconds
```cbs
TimestampMilli()
```

- date 获取日期
```cbs
date()
```

- TimestampToDate 时间戳转日期  一个参数（时间戳）
```cbs
TimestampToDate(Timestamp())
```

- TimestampToDateAT 指定时间格式 第一个参数(时间戳) 第二个参数时间格式 YYYYMMDD YYYY-MM-DD YYYYMMDDHHmmss YYYY-MM-DD HH:mm:ss MMdd HHmmss
```cbs
TimestampToDateAT(Timestamp(), "YYYYMMDDHHmmss")
```

- BeginDayUnix 获取当天0点的时间戳
```cbs
BeginDayUnix()
```

- EndDayUnix 获取当天24点的时间戳
```cbs
EndDayUnix()
```

- MinuteAgo 获取多少分钟前的时间戳  一个参数
```cbs
MinuteAgo(4)
```

- HourAgo 获取多少小时前的时间戳  一个参数
```cbs
DayAgo(4)
```

- DayDiffAtUnix 两个时间戳的插值  两个参数都是时间戳
```cbs
DayDiffAtUnix(Timestamp(), MinuteAgo(4))
```

- DayDiff 两个时间字符串的日期差, 返回的是天 两个参数都是时间字符串，格式是 YYYY-MM-DD HH:mm:ss
```cbs
DayDiff(TimestampToDate(Timestamp()), TimestampToDate(MinuteAgo(4)))
```

- NowToEnd 计算当前时间到这天结束还有多久,单位秒
```cbs
NowToEnd()
```

- IsToday 判断时间戳是否是今天，返回今天的时分秒  一个参数（时间戳）
```cbs
IsToday(Timestamp())
```

- Timestamp2Week 传入的时间戳是周几  一个参数（时间戳）
```cbs
Timestamp2Week(Timestamp())
```

- Timestamp2WeekXinQi 传入的时间戳是星期几  一个参数（时间戳）
```cbs
Timestamp2WeekXinQi(Timestamp())
```

### http关键字

http的请求，命令式语法，支持所有类型的请求，能将返回接口保存到变量，也能保存到本地文件；
额外支持并发请求，用于压力测试场景

参数说明：

- method ：请求方式 get post put delete options head patch
- url : 请求的url,要求类型是str
- body ： 请求的body,要求类型者是str或是List和字典（根据ctype解析为from-data，json这些）
- header ： 请求的header,要求类型是字典或者是json str
- ctype ： 请求的 是 header key 为 Content-Type, 要求类型是str
- cookie ：请求的cookie 是 header key 为 Cookie, 要求类型者是str(k=v;)或是List和字典（会解析为 k=v;） list是 ["k1=v1", "k2=v2"...]
- timeout ：设置请求的超时时间单位为毫秒, 要求类型是数值
- proxy ：设置请求的代理，目前只支持 http/https代理, 要求类型是str
- stress ：压力请求，并发请求设置的数量，要求类型是数值
- save : 指定将响应内容存储，要求类型是str,本地文件路径
- to : 将请求的返回存入到指定变量-如果变量未声明这里会自动声明变量
- save : 将请求的返回存入到指定文件

下面是相关例子

```cbs
// 简单的get请求
http get url="www.baidu.com"

// 简单的post请求
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" body="{\"id\": \"2b65775d-d68b-485a-af17-99f13ceb167a\"}"

// 请求参数变量
var b1 = {
    "id": "2b65775d-d68b-485a-af17-99f13ceb167a"
}
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" body=b1

// 请求参数变量未List
var b2 = ["aaaa", "bbbb"]
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" body=b2

// 带请求头
var h1="{\"id\":1}"
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" header=h1

// 请求头参数
var h2= {
    "id":1
}
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" header=h2

// 带很多参数的请求
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" body="{\"id\": \"2b65775d-d68b-485a-af17-99f13ceb167a\"}" header="{\"id\":1}" ctype="application/json" cookie="language=zh-CN" timeout=5 proxy="127.0.0.1:9080"

// 将请求接口存储到变量
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" body="{\"id\": \"2b65775d-d68b-485a-af17-99f13ceb167a\"}" to=rse

// 将请求参数存储到本地文件
http post url="https://api.ecosmos.cc/webapi/industrial/company2/share" body="{\"id\": \"2b65775d-d68b-485a-af17-99f13ceb167a\"}" save="D:\share.txt"

// 下载图片
http get url="https://resource.ecosmos.vip/AD/ad_h5.png?t=1772181855" save="D:\ad_h5.png"
```

### Chrome关键字

自动操作浏览器指令，命令式语法；环境会进行隔离，不影响当前用户已用的chrome; 注意：一个ChromeBot进程对应一个chrome子进程, 一行命令只支持一个操作； 
如果想多开运行多个ChromeBot进程执行脚本需要再启动时候添加new参数来进行隔离； 支持弹出窗进行交互；

参数说明：

- init : 初始化打开浏览器，如果已经打开后续语句再出现init会忽略
- close : 关闭浏览器
- size : 设置浏览器窗口大小与init参数一起用,值为: 宽*高 （900*600） <值类型是字符串>
- proxy : 设置浏览器代理与init参数一起用 <值类型是字符串>
- userpath : 设置浏览器在本机的隔离目录与init参数一起用,对应浏览器的--user-data-dir，建议隔离 <值类型是字符串>
- new : 设置浏览器新建一个隔离环境与init参数一起用；与userPath同时在时，优先使用userPath
- tab : 页签, 值有get:获取；set:指定哪个标签切换到指定的页签; new：新建一个页签；1<number>:第一个页签；select：返回当前选中的页签; 注意: 如果是没有选中页签下文操作默认当前浏览器的页签进行操作 <值类型是指定的字符串>
- req :  请求网址， 值为网址 <值类型是字符串>
- click : 点击操作，值为xpath <值类型是字符串>
- xpath : 当前选中的xpath, 输入的时候用
- input : 输入操作，输入内容  <值类型是字符串>
- check : 检查操作，检查页面是否存在指定xpath  <值类型是字符串>
- wait : 默认会执行等待页面加载完成，这个参数给定操作时候设置等待的时间  <值类型是数值类型>
- pause : 默认会执行等待页面加载完成，这个参数给定操作时候设置等待的时间  <值类型是数值类型>
- scroll : 滚动操作，滚动页面  正数往下，负数往上 <值类型是数值类型>  注意: 该滚动存在局限性只针对根节点进行滚动，嵌套容器要想精确请使用 scrollxpath
- scrollpixel : scroll by pixel 滚动操作,滚动到指定坐标， 值为(x,y)如(2000, 500)   注意: 该滚动存在局限性只针对根节点进行滚动, 嵌套容器要想精确请使用 scrollxpath
- scrollxpath : 滚动操作,滚动到指定xpath <值类型是字符串>
- screenshot : 截图操作，浏览器截图操作  值为保存位置  <值类型是字符串>
- to : 将当前操作的页面html返回存入到指定变量-如果变量未声明这里会自动声明变量  <值类型是字符串>
- save : 将将当前操作的页面html存入到指定文件  <值类型是字符串>
- info : 获取chrome 的信息
- as : 将指令的结果赋值给变量

下面是相关例子
```cbs
// 例子1 ：简单访问百度进行查询最后截图保存操作
chrome init  // 打开浏览器
chrome req="www.baidu.com" // 访问 百度
chrome xpath=NowTabGetInputFirstXpath() input="ChromeBot" // 获取当前页面能输入的第一个输入框的xpath
chrome click=NowTabMatchDemoContentOP("百度一下") // 获取当前页面内容为“百度一下”可交互的xpath
chrome screenshot="D:\baidu5.png"  // 截图保存到本地
chrome close  // 关闭浏览器

// 例子2 ： 央视新闻网，获取每个分类标签下的新闻列表打印新闻的标题出来
// https://news.cctv.com/
chrome init  // 打开浏览器
chrome req="https://news.cctv.com/"
var label = ["新闻","国内","国际","经济","社会","法治","图片","文娱","科技","生活","军事","快看"]
for item in label {
    print("寻找 --> ", item)
    var xpath = NowTabMatchDemoContentOP(item)
    chrome click=xpath // 点击
    var newsList = NowTabGetPointIDHTML("ul", "newslist")
    var title = RegHtml(newsList, "h3")
    for titleItem in  title {
        print("新闻title : ", titleItem)
    }
    chrome pause=2
}
chrome close

// 例子3 ： 抖音交互下滑视频
// https://www.douyin.com/?recommend=1
chrome init
chrome req="https://www.douyin.com/?recommend=1 "
sleep(2000)

// 1. 关闭登录
chrome click=`//*[@id="douyin_login_comp_flat_panel"]/div/div[1]/div[2]`
sleep(1000)

// 2. 点击推荐  NowTabMatchDemoContentOP("推荐")
chrome click=NowTabMatchDemoContentOP("推荐")
sleep(1000)

// 3. 看5秒点击下一个
chrome click=`//*[@id="douyin-right-container"]/div[3]/div[1]/div/div/div[2]`
sleep(5000)

// 4. 截图，然后结束
chrome screenshot="D:\\douyin_1.png"
chrome close


// 例子4 ： 交互豆包问豆包问题
// https://www.doubao.com/chat/
chrome init
chrome req="https://www.doubao.com/chat/"
sleep(2000) // 加载完页面等两秒
chrome xpath=`//textarea[@data-testid='chat_input_input']` input="你好豆包，介绍一下你自己"
chrome click=`//*[@id='flow-end-msg-send']`
// 等待回复，检查是否回复完，最多44检查
var isWait = 0
for var wait= 0; wait < 45; wait++ {
    if isWait > 1 {
        print("已经回复完") // 回复完了截图
        chrome screenshot="D:\doubao_1.png"
        break
    }
    chrome check=`//div[contains(@class, 'send-btn-wrapper') and (contains(@class, '!hidden'))]` as=has
    if has == false {
        isWait++
    }
}
if isWait==0 {
    print("回复太慢，还在回复吗?请检查")
}
chrome close
```

### Chrome 自动化场景下的相关方法

- ShowDemoTree 显示当前demo树
```
ShowDemoTree()
```

- MatchDemoContent 获取匹配到标签内容的xpath
```
MatchDemoContent(html, "首页")
```

- MatchDemoContentOP 获取匹配到标签内容的xpath, 能用于操作的xpath
```
MatchDemoContentOP(html, "首页")
```

- NowTabMatchDemoContentOP 获取当前操作的页面匹配到标签内容的xpath, 能用于操作的xpath
```
NowTabMatchDemoContentOP("首页")
```

- NowTabGetInputFirstXpath 获取当前操作的页面匹配到能输入的标签的xpath，返回匹配到的第一个
```
NowTabGetInputFirstXpath()
```

- NowTabGetPointHTML 获取指定位置的HTML， 用标签， 标签属性， 属性值来定位
```
NowTabGetPointHTML(label, attr, val) // label:标签  attr:标签属性  val:属性值
```

- NowTabGetPointIDHTML 获取指定位置的HTML， 用标签， 标签属性为id， 属性值来定位
```
NowTabGetPointIDHTML(label, val) // label:标签 id的val:属性值
```

- NowTabGetPointClassHTML 获取指定位置的HTML， 用标签， 标签属性为class， 属性值来定位
```
NowTabGetPointClassHTML(label, val) // label:标签 class的val:属性值
```

- HTMLGetPoint(html, label, attr, val)  获取指定位置的HTML， 用标签， 标签属性， 属性值来定位
  
- HTMLGetPointID(html, label, val) 获取指定位置的HTML， 用标签， 标签属性为id， 属性值来定位
  
- HTMLGetPointClass(html, label, val) 获取指定位置的HTML， 用标签， 标签属性为class， 属性值来定位
  
- HtmlToTableSaveExcel(html, path, 可选参数sheetName) 提取html内的表格数据保存为Excel

### host 系统方法关键字

host 系统方法关键字,主要编写系统脚本所用到，如文件管理，系统信息查询，系统相关的等等方法，与系统bat文件脚本一样; 
注意一个命令只执行一个参数。

参数说明：

info : 获取系统的信息
name : 获取系统的名称
ip : 获取系统的ip
to : 将当前操作返回的值存入到指定变量-如果变量未声明这里会自动声明变量  <值类型是字符串>
disk : todo 系统的磁盘信息
ls : 列出文件或目录
file : 操作系统文件
  - s=<search word> root=<path> : 搜索文件或目录
  - c=<path> : 创建文件或目录
  - d=<path> : 删除文件或目录
  - m=<path> goto=<path> : 移动文件或目录
  - cp=<path> goto=<path> : 复制文件或目录
  - r=<path> to=<arg> : 读文件
  - renm=<path> goto=<path> : 文件或目录改名, 路径不同则移动
  - info=<path> : 文件或目录信息
  - w=<path> from=<arg> : 将文件内容写入文件
  - a=<path> from=<arg> : 将文件内容追加写入文件

ping : ping命令
port : 查看本机开放端口
zip src=<path> dst=<path>  : zip压缩
unzip src=<path> dst=<path> : unzip解压

```
host info
```

### 网站站点相关方法

- WebSiteScanBadLink(domain, depth) 网站死链检查, depth是遍历网站的深度
```
WebSiteScanBadLink("www.baidu.com", 1)
```    

- WebCertificateInfo(domain) 网站证书信息
```
WebCertificateInfo("www.baidu.com")
```

- WebScanUrl(domain, depth) NewHostScanUrl 创建扫描站点
```
WebScanUrl("www.baidu.com", 1)
```

- WebScanExtLinks(domain, depth) 创建站点链接采集，只支持get请求
```
- WebScanExtLinks("www.baidu.com", 1) 
```

- WebPageSpeedCheck(domain, depth) 创建站点所有url测速，只支持get请求
```
WebPageSpeedCheck("www.baidu.com", 1)
```


### 网络相关的方法

- NsLookUp(host) DNS查询方法
```
NsLookUp("www.baidu.com")
```

- Whois(host) Whois查询方法
```
Whois("www.baidu.com")
```

- SearchPort(ip) 端口扫描方法
```
SearchPort("192.168.1.1")
```

### Json相关方法
- jsonDict(str) 将json字符串转换成字典
```
jsonDict("{'name':'zhangsan'}")
```

- json(arg) 将字典转换成json字符串
```
json([1,2,3,4])
```

- jsonFind(str, find) 查找json字符串  find是查询节点 如： {a:[{b:1},{b:2}]}  find=/a/[0]  =>   {b:1}   find=a/[0]/b  =>  1
```
jsonFind("{'a':[{'b':1},{'b':2}]}", "a/[0]/b")
```  

- jsonIS(str) 判断是否是json字符串
```
jsonIS("{'name':'zhangsan'}")
```

- jsonSave(arg, path) 将变量转存到本地文件，数据内容为json(格式化输出)
```
jsonSave([1,2,3,4], "D:\\json.txt")
```


### Excel相关方法
- ExcelSave(path, arg, 可选参数sheetName) 将变量保存到excel
  
- ExcelReadList(path, 可选参数sheetName) 读取excel返回二维列表

- ExcelReadDict(path, 可选参数sheetName) 读取excel返回字典

- ExcelShow(path, 可选参数sheetName) 显示excel

- ExcelInfo(path) 获取excel信息

- ExcelSheetInfo(path, sheetName) 获取excel的sheet信息

- ExcelSheet(path) 获取excel的sheet信息

- ExcelGetByCell(path, cell, 可选参数sheetName) 通过位置标签获取excel数据   cell 标签 A1 B1 C1 ...

- ExcelGetByPos(path, row, col, 可选参数sheetName) 通过位置获取excel数据

- ExcelSetByCell(path, cell, value, 可选参数sheetName) 通过位置标签设置excel数据   cell 标签 A1 B1 C1 ...

- ExcelSetByPos(path, row, col, value, 可选参数sheetName) 通过位置设置excel数据

- ExcelClearByCell(path, cell, 可选参数sheetName) 通过位置标签清除excel数据   cell 标签 A1 B1 C1 ...

- ExcelClearByPos(path, row, col, 可选参数sheetName) 通过位置清除excel数据

- ExcelReadRow(path, row, 可选参数sheetName)  读取指定行数据

- ExcelWriteRow(path, row, list, 可选参数sheetName)  写入指定行数据

- ExcelDeleteRow(path, row, 可选参数sheetName)  删除指定行数据

- ExcelReadCol(path, col, 可选参数sheetName)  读取指定列数据

- ExcelWriteCol(path, col, list, 可选参数sheetName)  写入指定列数据

- ExcelDeleteCol(path, col, 可选参数sheetName)  删除指定列数据

- ExcelReadCell(path, cell, 可选参数sheetName)  读取列 cell 标签 A B C ...

- ExcelWriteCell(path, cell, list, 可选参数sheetName)  写入列 cell 标签 A B C ...

- ExcelDeleteCell(path, cell, 可选参数sheetName)  删除指定列数据  cell 标签 A B C ...

- ExcelImg(path, cell, imgPath, 可选参数sheetName)   插入图片 cell 标签 A1 B1 C1 ...

- ExcelCellStyle(path, cell, style, 可选参数sheetName) 设置单元格样式 cell 标签 A1 B1 C1 ...
```
style {
  	fontBold: 是否加粗
 	fontColor: 字体颜色（十六进制，如"FF0000"）
 	bgColor: 背景颜色（十六进制，如"E0E0E0"）
 	alignCenter: 是否居中
}

// 修改单元格样式
ExcelCellStyle("D:\\test4.xlsx", "A1", {"fontBold":true, "fontColor":"FF1234", "bgColor":"E0E0E0", "alignCenter":true})

```

- ExcelMergeCells(path, startCell, endCell, 可选参数sheetName) 合并单元格 cell 标签 A1 B1 C1 ...

- ExcelSetFormula(path, cell, formula, 可选参数sheetName) 给单元格设置公式 标签 A1 B1 C1 ...  formula公式 如"SUM(A1:A3)"

- ExcelToJson(path, rowHead, 可选参数sheetName) rowHead:第几行作为key 如果是0key默认为 标签 A1 B1 C1 ...

- ExcelFromJson(path, json, 可选参数sheetName) json:json字符串数据

### todo....
