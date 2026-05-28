# Pigsh Kernel

一个玩笑语言的解释器内核。

## 构建&使用方法

```sh
go build -o pigsh ./cmd/pigsh 
```

```sh
./pigsh xxx.pigsh # 解释运行代码
```

## 基本语法

### 数据类型

| 类型 | 示例 |
|------|------|
| 整数 | `10` `-5` `0` |
| 浮点数 | `3.14` `-0.5` |
| 字符串 | `"hello"` `'world'` |
| 布尔值 | `true` `false` |

### 指令集

#### 变量声明 — var

声明并赋值一个变量。如果变量已存在则覆盖。

```
var x 10
var name "hello"
var flag true
var pi 3.14
```

#### 算术运算 — add / sub / mul / div / mod

三操作数格式：`指令 目标 操作数1 操作数2`，结果存入目标变量。

```
add result a b      # result = a + b
sub diff big small  # diff = big - small
mul prod x y        # prod = x * y
div ratio total n   # ratio = total / n（浮点除法）
mod rem value d     # rem = value % d
```

操作数可以是变量名或字面量：

```
var x 0
add x x 5       # x = 0 + 5 = 5
add x x 10      # x = 5 + 10 = 15
mul x x 2       # x = 15 * 2 = 30
```

字符串拼接使用 `add`：

```
var greeting ""
add greeting "hello" "world"   # greeting = "helloworld"
```

#### 条件分支 — beq / bne / blt / bgt / ble / bge

比较两个值，满足条件时跳转到目标行号或标签。

```
beq a b target    # a == b 时跳转
bne a b target    # a != b 时跳转
blt a b target    # a < b 时跳转
bgt a b target    # a > b 时跳转
ble a b target    # a <= b 时跳转
bge a b target    # a >= b 时跳转
```

示例 — 循环 10 次：

```
var i 0
var n 10
label loop
print i
inc i
blt i n loop
halt
```

#### 逻辑运算 — and / or / xor / not

```
and r a b    # r = a and b
or  r a b    # r = a or b
xor r a b    # r = a xor b
not r a      # r = not a
```

#### 数据移动 — mov

将值从一个变量或字面量复制到另一个变量。

```
mov x 42      # x = 42
mov y x       # y = x
```

#### 无条件跳转 — jump

跳转到指定行号或标签。

```
jump 5
jump end
```

#### 标签 — label

在当前位置定义一个标签，供跳转指令使用。运行时等同于 nop。

```
label start
label loop
label end
```

#### 打印 — print

引号包裹的字符串直接输出，无引号则查找变量，未找到则输出空行。

```
print "Hello, Pigsh!"   # 输出: Hello, Pigsh!
var x 42
print x                  # 输出: 42
print y                  # y 未定义，输出空行
```

#### 输入 — input

暂停执行，等待用户通过 `pigsh input` 提供输入。可选指定提示文本。

```
input x               # 等待输入，存入 x
input name "your name" # 带提示文本
```

执行到 `input` 时脚本挂起，用户在群内发送 `pigsh input <值>` 后继续执行。

#### 栈操作 — push / pop

```
push x    # 将 x 的值压入栈
pop y     # 弹出栈顶值存入 y
```

栈为空时执行 `pop` 会报错。

#### 子程序调用 — call / ret

`call` 将返回地址压入调用栈并跳转到目标位置。`ret` 弹出返回地址并跳回。

```
call greet
print done
halt

label greet
print hello
ret
```

支持嵌套调用：

```
call a
halt

label a
print in_a
call b
print back_a
ret

label b
print in_b
ret
```

#### 自增/自减 — inc / dec

对变量加 1 或减 1。

```
var x 10
inc x    # x = 11
dec x    # x = 10
```

#### 取反 — neg

对数值变量取反。

```
var x 10
neg x    # x = -10
```

#### 停机 — halt

立即终止脚本执行。

```
halt
```

#### 空操作 — nop

什么也不做，跳到下一行。

```
nop
```

## 错误信息

| 错误 | 含义 |
|------|------|
| `File not found: xxx` | 文件不存在 |
| `Not a script: xxx` | 试图运行非 `.pigsh` 文件 |
| `Unknown command: xxx` | 未知的子命令 |
| `Usage: pigsh xxx ...` | 命令参数不足 |
| `Error at line N: undefined variable: x` | 使用了未声明的变量 |
| `Error at line N: undefined label: x` | 跳转到不存在的标签 |
| `Error at line N: unknown instruction: x` | 无法识别的关键字 |
| `Error at line N: division by zero` | 除以零 |
| `Error at line N: stack underflow` | 空栈执行 pop |
| `Error at line N: call stack underflow` | 空调用栈执行 ret |
