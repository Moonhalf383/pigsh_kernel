# Pigsh 使用文档

Pigsh 是一个运行在 QQ 群中的虚拟文件系统与极简脚本语言。所有交互通过 `pigsh` 命令完成。

## 文件系统

Pigsh 提供一个平铺的虚拟文件系统，没有目录概念。文件分为两种类型，由扩展名自动判定：

- **文本文件** — 所有非 `.pigsh` 后缀的文件（如 `note.txt`、`data`）
- **脚本文件** — `.pigsh` 后缀的文件（如 `hello.pigsh`、`count.pigsh`）

文件内容在机器人重启后持久化保存。

## Shell 命令

### pigsh / pigsh help

显示帮助信息。

```
pigsh
pigsh help
```

### pigsh whoami

显示当前用户的 QQ 昵称。

```
pigsh whoami
```

### pigsh ls

列出所有文件，按名称排序。

```
pigsh ls
```

输出示例：

```
count.pigsh
hello.txt
```

### pigsh write

创建或覆写一个文件。内容中如有空格需要用引号包裹。

```
pigsh write <文件名> <内容>
```

示例：

```
pigsh write hello.txt Hello World
pigsh write greet.pigsh print hello
pigsh write note.txt "this has spaces"
```

### pigsh append

向已有文件追加内容。文件不存在时返回错误。

```
pigsh append <文件名> <内容>
```

示例：

```
pigsh append hello.txt New line here
```

### pigsh cat

查看文件内容。

```
pigsh cat <文件名>
```

示例：

```
pigsh cat hello.txt
```

### pigsh rm

删除文件。

```
pigsh rm <文件名>
```

示例：

```
pigsh rm hello.txt
```

### pigsh run

执行一个 `.pigsh` 脚本文件。如果脚本中包含 `input` 指令导致挂起，需要通过 `pigsh input` 提供输入。

```
pigsh run <文件名>
```

示例：

```
pigsh run hello.pigsh
```

### pigsh input

当脚本执行到 `input` 指令挂起后，通过此命令提供输入值。

```
pigsh input <值>
```

示例：

```
pigsh input 42
pigsh input hello
```

## 脚本语言

Pigsh 脚本是一种极简的汇编风格语言。每条语句由 3-4 个词组成，第一个词始终是关键字。`#` 开头的行是注释。

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

将变量或字面量的值输出。

```
var msg "hello"
print msg
print 42
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

## 完整示例

### 示例 1：Hello World

```
pigsh write hello.pigsh print hello
pigsh run hello.pigsh
```

输出：`hello`

### 示例 2：斐波那契数列

写入脚本：

```
pigsh write fib.pigsh "var a 0
var b 1
var n 10
var i 0
label loop
print a
var tmp 0
mov tmp b
add b a b
mov a tmp
inc i
blt i n loop
halt"
```

执行：

```
pigsh run fib.pigsh
```

输出：

```
0
1
1
2
3
5
8
13
21
34
```

### 示例 3：带子程序的脚本

```
pigsh write greet.pigsh "call say_hi
call say_bye
halt
label say_hi
print hello
ret
label say_bye
print goodbye
ret"
pigsh run greet.pigsh
```

输出：

```
hello
goodbye
```

### 示例 4：累加器

```
pigsh write sum.pigsh "var total 0
var i 1
var n 10
label loop
add total total i
inc i
ble i n loop
print total
halt"
pigsh run sum.pigsh
```

输出：`55`

### 示例 5：使用 input 交互

```
pigsh write ask.pigsh "input name who_are_you
print name"
pigsh run ask.pigsh
```

脚本挂起，显示提示 `who_are_you`。用户回复：

```
pigsh input pig
```

输出：`pig`

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

## 注意事项

- 脚本中每条语句的各部分之间用空格分隔，不支持在一行内使用带空格的字符串。需要多行内容时，每行写一条指令，用 `append` 逐行追加。
- `input` 指令的提示文本不能包含空格。
- 脚本没有执行步数限制，请避免编写无限循环。
- 文件系统没有大小限制。
- 脚本的输入等待状态保存在内存中，机器人重启后丢失。
