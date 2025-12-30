# ESP32 Clang 编译 MQuickJS 构建指南

本文档说明如何使用 esp-clang 直接编译 MQuickJS，不依赖 ESP-IDF 构建系统。

## 需要编译的文件

| 文件 | 作用 |
|------|------|
| `mquickjs.c` | JS 引擎核心（解析器、解释器、GC） |
| `dtoa.c` | 浮点数转字符串 |
| `libm.c` | 软件浮点数学库 |
| `cutils.c` | 通用工具函数 |

## 编译器路径

```bash
# esp-clang 安装后的路径
CLANG=/path/to/.espressif/tools/esp-clang/esp-18.1.2_20240912/esp-clang/bin/clang
AR=/path/to/.espressif/tools/esp-clang/esp-18.1.2_20240912/esp-clang/bin/llvm-ar
```

---

## 必需参数

### 目标架构（必需）

```bash
--target=xtensa-esp-elf    # 目标三元组
-mcpu=esp32                # CPU 类型
```

### 代码生成（必需）

```bash
-fno-jump-tables           # 禁用跳转表（ESP32 IRAM/Flash 架构要求）
-ffunction-sections        # 每个函数独立 section（配合链接器 gc-sections）
-fdata-sections            # 每个数据独立 section
```

#### 为什么必须禁用跳转表？

```
┌─────────────────────────────────────────┐
│              ESP32 芯片                  │
│  ┌─────────┐     ┌─────────────────┐   │
│  │  IRAM   │     │   Flash Cache   │   │
│  │ ~160KB  │     │  (可能被禁用)    │   │
│  │ 始终可用 │     └────────┬────────┘   │
│  └─────────┘              │            │
└───────────────────────────┼────────────┘
                            │ SPI
                    ┌───────┴───────┐
                    │  外部 Flash   │
                    └───────────────┘
```

- 跳转表数据放在 Flash
- IRAM 代码在 Flash Cache 禁用时无法访问跳转表
- **必须禁用以避免运行时崩溃**

---

## 推荐参数

### 优化

```bash
-Os                        # 优化代码大小（嵌入式推荐）
-fno-math-errno            # 数学函数不设置 errno（减小代码）
-fno-trapping-math         # 假设浮点不触发异常（优化）
```

### C 标准

```bash
-std=gnu17                 # GNU C17 标准
```

### 警告抑制（可选，消除编译警告）

```bash
-Wno-format                # MQuickJS 有一些 format 警告
-Wno-type-limits           # 类型范围警告
```

---

## 完整编译脚本

```bash
#!/bin/bash

# === 配置 ===
CLANG_PATH="/path/to/.espressif/tools/esp-clang/esp-18.1.2_20240912/esp-clang/bin"
CLANG="$CLANG_PATH/clang"
AR="$CLANG_PATH/llvm-ar"

# === 必需参数 ===
TARGET_FLAGS="--target=xtensa-esp-elf -mcpu=esp32"
CODEGEN_FLAGS="-fno-jump-tables -ffunction-sections -fdata-sections"

# === 推荐参数 ===
OPT_FLAGS="-Os -fno-math-errno -fno-trapping-math"
STD_FLAGS="-std=gnu17"
WARN_FLAGS="-Wno-format -Wno-type-limits"

# === 合并所有参数 ===
CFLAGS="$TARGET_FLAGS $CODEGEN_FLAGS $OPT_FLAGS $STD_FLAGS $WARN_FLAGS"

# === 编译 ===
echo "Compiling mquickjs.c..."
$CLANG $CFLAGS -c -o mquickjs.o mquickjs.c || exit 1

echo "Compiling dtoa.c..."
$CLANG $CFLAGS -c -o dtoa.o dtoa.c || exit 1

echo "Compiling libm.c..."
$CLANG $CFLAGS -c -o libm.o libm.c || exit 1

echo "Compiling cutils.c..."
$CLANG $CFLAGS -c -o cutils.o cutils.c || exit 1

# === 创建静态库 ===
echo "Creating static library..."
$AR rcs libmquickjs_esp32.a mquickjs.o dtoa.o libm.o cutils.o || exit 1

echo ""
echo "=== Build successful! ==="
ls -lh libmquickjs_esp32.a *.o
```

---

## 参数速查表

| 参数 | 必需 | 说明 |
|------|:----:|------|
| `--target=xtensa-esp-elf` | ✅ | 目标架构 |
| `-mcpu=esp32` | ✅ | CPU 类型 |
| `-fno-jump-tables` | ✅ | 禁用跳转表（IRAM 兼容） |
| `-ffunction-sections` | ✅ | 函数独立 section |
| `-fdata-sections` | ✅ | 数据独立 section |
| `-Os` | 推荐 | 优化大小 |
| `-fno-math-errno` | 推荐 | 数学优化 |
| `-fno-trapping-math` | 推荐 | 浮点优化 |
| `-std=gnu17` | 推荐 | C 标准 |
| `-Wno-format` | 可选 | 消除警告 |
| `-Wno-type-limits` | 可选 | 消除警告 |

---

## 最小编译命令

如果只想快速编译，使用这个最小命令：

```bash
clang --target=xtensa-esp-elf -mcpu=esp32 -fno-jump-tables -Os -c mquickjs.c -o mquickjs.o
clang --target=xtensa-esp-elf -mcpu=esp32 -fno-jump-tables -Os -c dtoa.c -o dtoa.o
clang --target=xtensa-esp-elf -mcpu=esp32 -fno-jump-tables -Os -c libm.c -o libm.o
clang --target=xtensa-esp-elf -mcpu=esp32 -fno-jump-tables -Os -c cutils.c -o cutils.o
llvm-ar rcs libmquickjs_esp32.a mquickjs.o dtoa.o libm.o cutils.o
```

---

## 参考资料

- ESP-IDF commit `ee2f8b1a62`: `-fno-jump-tables` 的原因
- [GCC Xtensa Options](https://gcc.gnu.org/onlinedocs/gcc/Xtensa-Options.html)
