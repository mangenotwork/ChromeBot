# ChromeBot Script 
是基于ATS实现是DSL语言，主要用于ChromeBot执行自动化任务的脚本编写，脚本文件后缀为.cbs

## 运行
chromeBot.exe case.cbs

## 语法

### SDL语法设计
1. 注重脚本实现过程专注于流程精准表达目的为目标，语法要追求简单，操作多以命令式语法，省去了函数式编程和面向对象编程
2. 由于主要编写执行自动化的脚本，设计极简，只有数值类型、字符串类型、布尔类型、列表类型、字典类型（键值对类型）
3. 支持逻辑判断(if else), 支持循环 (for  while), 由于已经有if了所以省去了switch case
4. chrome的操作关键字，用这些关键字命令式语法编写脚本来操作浏览器
5. 很多内置函数，编写脚本需参照文档找到需要的方法
6. 注释为 # 或 //
7. 支持函数链式调用，函数依次调用返回值会自动传递给下个函数


### 变量

与大多数脚本语言一样使用关键字 var 来声明变量

```cbs
var a = 1
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

### 逻辑判断

支持 if、else if、else 结构，条件表达式结果为布尔类型

```cbs
var score = 85
if score >= 90 {
    print("优秀")
} else if score >= 60 {
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

### 内置函数

- len(arg): 获取传入类型的长度，arg是任意类型，返回长度 
```cbs
var list = [1, 2, 3, 4, 5]
print(len(list))
```

todo....

### 链式调用

支持将多个函数调用连接在一起，上一个函数的返回值自动作为下一个函数的参数

```cbs
var s = "test";
var s1 = upper(s).repeat(2)
print(s1)
```

### Chrome关键字

todo....

