#include <stdio.h>
#include <string.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "esp_timer.h"
#include "mquickjs.h"

// JS 引擎使用的内存 (32KB)
static uint8_t js_mem[32768];

// === 标准库需要的 C 函数实现 ===

static JSValue js_print(JSContext *ctx, JSValue *this_val, int argc, JSValue *argv)
{
    int i;
    JSValue v;

    for (i = 0; i < argc; i++) {
        if (i != 0)
            putchar(' ');
        v = argv[i];
        if (JS_IsString(ctx, v)) {
            JSCStringBuf buf;
            const char *str;
            size_t len;
            str = JS_ToCStringLen(ctx, &len, v, &buf);
            fwrite(str, 1, len, stdout);
        } else {
            JS_PrintValueF(ctx, argv[i], JS_DUMP_LONG);
        }
    }
    putchar('\n');
    fflush(stdout);
    return JS_UNDEFINED;
}

static JSValue js_gc(JSContext *ctx, JSValue *this_val, int argc, JSValue *argv)
{
    JS_GC(ctx);
    return JS_UNDEFINED;
}

static JSValue js_date_now(JSContext *ctx, JSValue *this_val, int argc, JSValue *argv)
{
    int64_t ms = esp_timer_get_time() / 1000;  // 微秒转毫秒
    return JS_NewInt64(ctx, ms);
}

static JSValue js_performance_now(JSContext *ctx, JSValue *this_val, int argc, JSValue *argv)
{
    double ms = (double)esp_timer_get_time() / 1000.0;
    return JS_NewFloat64(ctx, ms);
}

// load 在嵌入式环境不支持文件系统
static JSValue js_load(JSContext *ctx, JSValue *this_val, int argc, JSValue *argv)
{
    return JS_ThrowTypeError(ctx, "load() not supported on ESP32");
}

// setTimeout/clearTimeout 简化实现（不支持）
static JSValue js_setTimeout(JSContext *ctx, JSValue *this_val, int argc, JSValue *argv)
{
    return JS_ThrowTypeError(ctx, "setTimeout() not supported");
}

static JSValue js_clearTimeout(JSContext *ctx, JSValue *this_val, int argc, JSValue *argv)
{
    return JS_ThrowTypeError(ctx, "clearTimeout() not supported");
}

// 包含标准库（必须在上面函数定义之后）
#include "mqjs_stdlib_esp32.h"

// 简单的 print 函数输出
static void js_write_func(void *opaque, const void *buf, size_t buf_len)
{
    fwrite(buf, 1, buf_len, stdout);
    fflush(stdout);
}

void app_main(void)
{
    printf("=== MQuickJS on ESP32 ===\n\n");

    // 创建 JS 上下文
    JSContext *ctx = JS_NewContext(js_mem, sizeof(js_mem), &js_stdlib);
    if (!ctx) {
        printf("Failed to create JS context\n");
        return;
    }

    // 设置输出函数
    JS_SetLogFunc(ctx, js_write_func);

    JSValue ret;

    // 测试 1: 纯字符串
    const char *code1 = "print('Hello from MQuickJS on ESP32!')";
    printf("Test 1: %s\n", code1);
    ret = JS_Eval(ctx, code1, strlen(code1), "test1.js", 0);
    if (JS_IsException(ret)) {
        printf("Exception: ");
        JS_PrintValueF(ctx, JS_GetException(ctx), JS_DUMP_LONG);
        printf("\n");
    }
    printf("\n");

    // 测试 2: 简单数字
    const char *code2 = "print(123)";
    printf("Test 2: %s\n", code2);
    ret = JS_Eval(ctx, code2, strlen(code2), "test2.js", 0);
    if (JS_IsException(ret)) {
        printf("Exception: ");
        JS_PrintValueF(ctx, JS_GetException(ctx), JS_DUMP_LONG);
        printf("\n");
    }
    printf("\n");

    // 测试 3: 数学运算
    const char *code3 = "print(10 + 20)";
    printf("Test 3: %s\n", code3);
    ret = JS_Eval(ctx, code3, strlen(code3), "test3.js", 0);
    if (JS_IsException(ret)) {
        printf("Exception: ");
        JS_PrintValueF(ctx, JS_GetException(ctx), JS_DUMP_LONG);
        printf("\n");
    }
    printf("\n");

    // 测试 4: var 声明
    const char *code4 = "var x = 1";
    printf("Test 4: %s\n", code4);
    ret = JS_Eval(ctx, code4, strlen(code4), "test4.js", 0);
    if (JS_IsException(ret)) {
        printf("Exception: ");
        JS_PrintValueF(ctx, JS_GetException(ctx), JS_DUMP_LONG);
        printf("\n");
    }
    printf("\n");

    // 测试 5: 字符串拼接
    const char *code5 = "print('a' + 'b')";
    printf("Test 5: %s\n", code5);
    ret = JS_Eval(ctx, code5, strlen(code5), "test5.js", 0);
    if (JS_IsException(ret)) {
        printf("Exception: ");
        JS_PrintValueF(ctx, JS_GetException(ctx), JS_DUMP_LONG);
        printf("\n");
    }
    printf("\n");

    // 释放上下文
    JS_FreeContext(ctx);

    printf("=== Done ===\n");

    // 保持运行
    while (1) {
        vTaskDelay(pdMS_TO_TICKS(1000));
    }
}
