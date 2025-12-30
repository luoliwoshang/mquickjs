# CLAUDE.md - MQuickJS 开发指南

本文档记录了 MQuickJS 的构建、嵌入式部署和 Go 绑定生成的探索过程。

## 项目概述

MQuickJS 是一个轻量级的 JavaScript 引擎，专为嵌入式系统设计：
- **ROM**: ~100KB
- **RAM**: ~10KB 起
- **无动态内存分配**: 不依赖 malloc/free

## 源文件结构

核心源文件（构建静态库需要的文件）：

| 文件 | 作用 |
|------|------|
| `mquickjs.c` | JS 引擎核心（解析器、解释器、GC） |
| `dtoa.c` | 浮点数转字符串（double to ascii） |
| `libm.c` | 软件浮点数学库（sin, cos, sqrt 等） |
| `cutils.c` | 通用工具函数 |

头文件：
- `mquickjs.h` - 公共 API（用户只需引入此文件）
- `mquickjs_priv.h` - 内部实现
- `libm.h` - 数学库头文件
- `softfp_template.h` - 软浮点模板

## libc 依赖

MQuickJS 对 libc 的依赖非常小：
- `memcpy`, `memset`, `memmove`, `memcmp`
- `strlen`, `strcmp`, `strncmp`
- `setjmp`, `longjmp`（异常处理）

**不依赖**：`malloc`, `free`, `printf`, `FILE*` 等

## 构建静态库

### 主机平台（macOS/Linux）

```bash
# 编译目标文件
gcc -Os -c -o mquickjs.o mquickjs.c
gcc -Os -c -o dtoa.o dtoa.c
gcc -Os -c -o libm.o libm.c
gcc -Os -c -o cutils.o cutils.c

# 创建静态库
ar rcs libmquickjs.a mquickjs.o dtoa.o libm.o cutils.o
```

### ESP32（Xtensa 架构）

使用 `build_esp32.sh` 脚本：

```bash
./build_esp32.sh
```

关键编译参数：
- `-mlongcalls`: Xtensa 特有，启用间接函数调用（CALLX0），允许调用超过 ±1MB 范围的函数
- `-Os`: 优化大小
- `-ffunction-sections -fdata-sections`: 配合链接器 `--gc-sections` 移除未使用代码

## ESP-IDF 集成

### 作为组件集成

将源文件放入 `components/mquickjs/` 目录：

```
components/mquickjs/
├── CMakeLists.txt
├── mquickjs.c
├── mquickjs.h
├── mquickjs_priv.h
├── dtoa.c
├── libm.c
├── libm.h
├── softfp_template.h
├── cutils.c
└── mquickjs_atom_esp32.h      # 32位 atom 表
└── mqjs_stdlib_esp32.h        # 32位标准库
```

`CMakeLists.txt` 内容：

```cmake
idf_component_register(
    SRCS "mquickjs.c" "dtoa.c" "libm.c" "cutils.c"
    INCLUDE_DIRS "."
)

target_compile_options(${COMPONENT_LIB} PRIVATE
    -Os
    -fno-math-errno
    -fno-trapping-math
    -Wno-error=format
    -Wno-error=type-limits
)
```

### 32位/64位 Atom 表

**重要**: ESP32 是 32 位平台，必须使用 32 位的 atom 表和标准库头文件：

```bash
# 在主机上生成 32 位头文件
./mqjs_stdlib -m32 -a > mquickjs_atom_esp32.h
./mqjs_stdlib -m32 > mqjs_stdlib_esp32.h
```

如果使用 64 位头文件，会出现：
- `var` 关键字无法识别
- 字符串拼接崩溃
- 各种奇怪的运行时错误

### 示例代码

参见 `examples/esp32/main/main.c`：

```c
#include "mquickjs.h"

static uint8_t js_mem[32768];  // 32KB 内存池

void app_main(void) {
    JSContext *ctx = JS_NewContext(js_mem, sizeof(js_mem), &js_stdlib);

    const char *code = "print('Hello from MQuickJS!')";
    JSValue ret = JS_Eval(ctx, code, strlen(code), "test.js", 0);

    if (JS_IsException(ret)) {
        JS_PrintValueF(ctx, JS_GetException(ctx), JS_DUMP_LONG);
    }

    JS_FreeContext(ctx);
}
```

## ESP-IDF 构建系统

### 构建流程

```
idf_component_register()
        ↓
    CMake 配置
        ↓
    生成 build/build.ninja
        ↓
    Ninja 执行
        ↓
    xtensa-esp32-elf-gcc 编译 .c → .o
        ↓
    xtensa-esp32-elf-ar 打包 .o → .a
        ↓
    xtensa-esp32-elf-ld 链接 → .elf
```

### 查看构建参数

```bash
# 查看编译命令
cat build/compile_commands.json | jq '.[] | select(.file | contains("mquickjs"))'

# 查看 ninja 构建规则
grep "mquickjs" build/build.ninja
```

### -mlongcalls 说明

这是 Xtensa 架构特有的 GCC 选项：

- **默认**: 使用 `CALL0` 指令，只能调用 ±1MB 范围内的函数
- **启用 -mlongcalls**: 使用 `CALLX0` 间接调用，可以调用任意地址

ESP-IDF 默认启用此选项，因为 Flash 映射地址可能超出 1MB 范围。

官方文档：https://gcc.gnu.org/onlinedocs/gcc/Xtensa-Options.html

## Go 绑定 (llcppg)

### 配置文件

`mquickjs/llcppg.cfg`：

```json
{
    "name": "mquickjs",
    "cflags": "-I.",
    "libs": "-lm",
    "include": [
        "mquickjs.h"
    ],
    "staticLib": true
}
```

### 生成绑定

```bash
cd mquickjs
llcppg .
```

生成的 Go 绑定位于 `mquickjs/mquickjs.go`。

### 主要 API 映射

| C 函数 | Go 方法 |
|--------|---------|
| `JS_NewContext()` | `JSNewContext()` |
| `JS_FreeContext()` | `(*JSContext).JSFreeContext()` |
| `JS_Eval()` | `(*JSContext).JSEval()` |
| `JS_NewString()` | `(*JSContext).JSNewString()` |
| `JS_NewInt32()` | `(*JSContext).JSNewInt32()` |
| `JS_NewFloat64()` | `(*JSContext).JSNewFloat64()` |
| `JS_GetException()` | `(*JSContext).JSGetException()` |
| `JS_IsException()` | `JS_IsException()` (inline) |

## 常见问题

### Q: ESP32 上 `var` 关键字不识别？

A: 使用了 64 位 atom 表。需要用 `-m32` 重新生成：
```bash
./mqjs_stdlib -m32 -a > mquickjs_atom_esp32.h
```

### Q: 链接时报 `undefined reference to __wrap_longjmp`？

A: ESP-IDF 会 wrap longjmp。解决方案是从源码编译而不是使用预编译的 .a 文件。

### Q: 如何减小代码体积？

A:
1. 使用 `-Os` 优化
2. 启用 `-ffunction-sections -fdata-sections`
3. 链接时使用 `-Wl,--gc-sections`
4. 考虑 strip 调试符号

## 目录结构

```
mquickjs/
├── mquickjs.c              # 核心引擎
├── mquickjs.h              # 公共 API
├── dtoa.c                  # 浮点转换
├── libm.c                  # 数学库
├── cutils.c                # 工具函数
├── build_esp32.sh          # ESP32 构建脚本
├── mquickjs/               # Go 绑定
│   ├── llcppg.cfg
│   └── mquickjs.go
└── examples/
    └── esp32/              # ESP-IDF 示例项目
        ├── CMakeLists.txt
        ├── main/
        │   └── main.c
        └── components/
            └── mquickjs/   # MQuickJS 组件
```
