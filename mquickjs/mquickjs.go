package mquickjs

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const JSW = 4
const JS_TAG_SPECIAL_BITS = 5
const JS_EX_NORMAL = 0
const JS_EX_CALL = 1
const JS_EVAL_REGEXP_FLAGS_SHIFT = 8
const JS_BYTECODE_MAGIC = 0xacfb

type JSContext struct {
	Unused [8]uint8
}
type JSWord c.Uint32T
type JSValue c.Uint32T

const (
	JS_TAG_INT           c.Int = 0
	JS_TAG_PTR           c.Int = 1
	JS_TAG_SPECIAL       c.Int = 3
	JS_TAG_BOOL          c.Int = 3
	JS_TAG_NULL          c.Int = 7
	JS_TAG_UNDEFINED     c.Int = 11
	JS_TAG_EXCEPTION     c.Int = 15
	JS_TAG_SHORT_FUNC    c.Int = 19
	JS_TAG_UNINITIALIZED c.Int = 23
	JS_TAG_STRING_CHAR   c.Int = 27
	JS_TAG_CATCH_OFFSET  c.Int = 31
)

type JSObjectClassEnum c.Int

const (
	JS_CLASS_OBJECT          JSObjectClassEnum = 0
	JS_CLASS_ARRAY           JSObjectClassEnum = 1
	JS_CLASS_C_FUNCTION      JSObjectClassEnum = 2
	JS_CLASS_CLOSURE         JSObjectClassEnum = 3
	JS_CLASS_NUMBER          JSObjectClassEnum = 4
	JS_CLASS_BOOLEAN         JSObjectClassEnum = 5
	JS_CLASS_STRING          JSObjectClassEnum = 6
	JS_CLASS_DATE            JSObjectClassEnum = 7
	JS_CLASS_REGEXP          JSObjectClassEnum = 8
	JS_CLASS_ERROR           JSObjectClassEnum = 9
	JS_CLASS_EVAL_ERROR      JSObjectClassEnum = 10
	JS_CLASS_RANGE_ERROR     JSObjectClassEnum = 11
	JS_CLASS_REFERENCE_ERROR JSObjectClassEnum = 12
	JS_CLASS_SYNTAX_ERROR    JSObjectClassEnum = 13
	JS_CLASS_TYPE_ERROR      JSObjectClassEnum = 14
	JS_CLASS_URI_ERROR       JSObjectClassEnum = 15
	JS_CLASS_INTERNAL_ERROR  JSObjectClassEnum = 16
	JS_CLASS_ARRAY_BUFFER    JSObjectClassEnum = 17
	JS_CLASS_TYPED_ARRAY     JSObjectClassEnum = 18
	JS_CLASS_UINT8C_ARRAY    JSObjectClassEnum = 19
	JS_CLASS_INT8_ARRAY      JSObjectClassEnum = 20
	JS_CLASS_UINT8_ARRAY     JSObjectClassEnum = 21
	JS_CLASS_INT16_ARRAY     JSObjectClassEnum = 22
	JS_CLASS_UINT16_ARRAY    JSObjectClassEnum = 23
	JS_CLASS_INT32_ARRAY     JSObjectClassEnum = 24
	JS_CLASS_UINT32_ARRAY    JSObjectClassEnum = 25
	JS_CLASS_FLOAT32_ARRAY   JSObjectClassEnum = 26
	JS_CLASS_FLOAT64_ARRAY   JSObjectClassEnum = 27
	JS_CLASS_USER            JSObjectClassEnum = 28
)

type JSCFunctionEnum c.Int

const (
	JS_CFUNCTION_bound JSCFunctionEnum = 0
	JS_CFUNCTION_USER  JSCFunctionEnum = 1
)

/* temporary buffer to hold C strings */

type JSCStringBuf struct {
	Buf [5]c.Uint8T
}

type JSGCRef struct {
	Val  JSValue
	Prev *JSGCRef
}

/* stack of JSGCRef */
// llgo:link (*JSContext).JSPushGCRef C.JS_PushGCRef
func (recv_ *JSContext) JSPushGCRef(ref *JSGCRef) *JSValue {
	return nil
}

// llgo:link (*JSContext).JSPopGCRef C.JS_PopGCRef
func (recv_ *JSContext) JSPopGCRef(ref *JSGCRef) JSValue {
	return 0
}

/* list of JSGCRef (they can be removed in any order, slower) */
// llgo:link (*JSContext).JSAddGCRef C.JS_AddGCRef
func (recv_ *JSContext) JSAddGCRef(ref *JSGCRef) *JSValue {
	return nil
}

// llgo:link (*JSContext).JSDeleteGCRef C.JS_DeleteGCRef
func (recv_ *JSContext) JSDeleteGCRef(ref *JSGCRef) {
}

// llgo:link (*JSContext).JSNewFloat64 C.JS_NewFloat64
func (recv_ *JSContext) JSNewFloat64(d c.Double) JSValue {
	return 0
}

// llgo:link (*JSContext).JSNewInt32 C.JS_NewInt32
func (recv_ *JSContext) JSNewInt32(val c.Int32T) JSValue {
	return 0
}

// llgo:link (*JSContext).JSNewUint32 C.JS_NewUint32
func (recv_ *JSContext) JSNewUint32(val c.Uint32T) JSValue {
	return 0
}

// llgo:link (*JSContext).JSNewInt64 C.JS_NewInt64
func (recv_ *JSContext) JSNewInt64(val c.Int64T) JSValue {
	return 0
}

// llgo:link (*JSContext).JSIsNumber C.JS_IsNumber
func (recv_ *JSContext) JSIsNumber(val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSIsString C.JS_IsString
func (recv_ *JSContext) JSIsString(val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSIsError C.JS_IsError
func (recv_ *JSContext) JSIsError(val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSIsFunction C.JS_IsFunction
func (recv_ *JSContext) JSIsFunction(val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSGetClassID C.JS_GetClassID
func (recv_ *JSContext) JSGetClassID(val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSSetOpaque C.JS_SetOpaque
func (recv_ *JSContext) JSSetOpaque(val JSValue, opaque c.Pointer) {
}

// llgo:link (*JSContext).JSGetOpaque C.JS_GetOpaque
func (recv_ *JSContext) JSGetOpaque(val JSValue) c.Pointer {
	return nil
}

// llgo:type C
type JSCFunction func(*JSContext, *JSValue, c.Int, *JSValue) JSValue

// llgo:type C
type JSCFinalizer func(*JSContext, c.Pointer)
type JSCFunctionDefEnum c.Int

const (
	JS_CFUNC_generic           JSCFunctionDefEnum = 0
	JS_CFUNC_generic_magic     JSCFunctionDefEnum = 1
	JS_CFUNC_constructor       JSCFunctionDefEnum = 2
	JS_CFUNC_constructor_magic JSCFunctionDefEnum = 3
	JS_CFUNC_generic_params    JSCFunctionDefEnum = 4
	JS_CFUNC_f_f               JSCFunctionDefEnum = 5
)

type JSCFunctionType struct {
	Generic *JSCFunction
}

type JSCFunctionDef struct {
	Func     JSCFunctionType
	Name     JSValue
	DefType  c.Uint8T
	ArgCount c.Uint8T
	Magic    c.Int16T
}

type JSSTDLibraryDef struct {
	StdlibTable        *JSWord
	CFunctionTable     *JSCFunctionDef
	CFinalizerTable    *JSCFinalizer
	StdlibTableLen     c.Uint32T
	StdlibTableAlign   c.Uint32T
	SortedAtomsOffset  c.Uint32T
	GlobalObjectOffset c.Uint32T
	ClassCount         c.Uint32T
}

// llgo:type C
type JSWriteFunc func(c.Pointer, c.Pointer, c.Int)

// llgo:type C
type JSInterruptHandler func(*JSContext, c.Pointer) c.Int

//go:linkname JSNewContext C.JS_NewContext
func JSNewContext(mem_start c.Pointer, mem_size c.Int, stdlib_def *JSSTDLibraryDef) *JSContext

/* if prepare_compilation is true, the context will be used to compile
   to a binary file. JS_NewContext2() is not expected to be used in
   the embedded version */
//go:linkname JSNewContext2 C.JS_NewContext2
func JSNewContext2(mem_start c.Pointer, mem_size c.Int, stdlib_def *JSSTDLibraryDef, prepare_compilation c.Int) *JSContext

// llgo:link (*JSContext).JSFreeContext C.JS_FreeContext
func (recv_ *JSContext) JSFreeContext() {
}

// llgo:link (*JSContext).JSSetContextOpaque C.JS_SetContextOpaque
func (recv_ *JSContext) JSSetContextOpaque(opaque c.Pointer) {
}

// llgo:link (*JSContext).JSSetInterruptHandler C.JS_SetInterruptHandler
func (recv_ *JSContext) JSSetInterruptHandler(interrupt_handler JSInterruptHandler) {
}

// llgo:link (*JSContext).JSSetRandomSeed C.JS_SetRandomSeed
func (recv_ *JSContext) JSSetRandomSeed(seed c.Uint64T) {
}

// llgo:link (*JSContext).JSGetGlobalObject C.JS_GetGlobalObject
func (recv_ *JSContext) JSGetGlobalObject() JSValue {
	return 0
}

// llgo:link (*JSContext).JSThrow C.JS_Throw
func (recv_ *JSContext) JSThrow(obj JSValue) JSValue {
	return 0
}

//go:linkname JSThrowError C.JS_ThrowError
func JSThrowError(ctx *JSContext, error_num JSObjectClassEnum, fmt *c.Char, __llgo_va_list ...interface{}) JSValue

// llgo:link (*JSContext).JSThrowOutOfMemory C.JS_ThrowOutOfMemory
func (recv_ *JSContext) JSThrowOutOfMemory() JSValue {
	return 0
}

// llgo:link (*JSContext).JSGetPropertyStr C.JS_GetPropertyStr
func (recv_ *JSContext) JSGetPropertyStr(this_obj JSValue, str *c.Char) JSValue {
	return 0
}

// llgo:link (*JSContext).JSGetPropertyUint32 C.JS_GetPropertyUint32
func (recv_ *JSContext) JSGetPropertyUint32(obj JSValue, idx c.Uint32T) JSValue {
	return 0
}

// llgo:link (*JSContext).JSSetPropertyStr C.JS_SetPropertyStr
func (recv_ *JSContext) JSSetPropertyStr(this_obj JSValue, str *c.Char, val JSValue) JSValue {
	return 0
}

// llgo:link (*JSContext).JSSetPropertyUint32 C.JS_SetPropertyUint32
func (recv_ *JSContext) JSSetPropertyUint32(this_obj JSValue, idx c.Uint32T, val JSValue) JSValue {
	return 0
}

// llgo:link (*JSContext).JSNewObjectClassUser C.JS_NewObjectClassUser
func (recv_ *JSContext) JSNewObjectClassUser(class_id c.Int) JSValue {
	return 0
}

// llgo:link (*JSContext).JSNewObject C.JS_NewObject
func (recv_ *JSContext) JSNewObject() JSValue {
	return 0
}

// llgo:link (*JSContext).JSNewArray C.JS_NewArray
func (recv_ *JSContext) JSNewArray(initial_len c.Int) JSValue {
	return 0
}

/* create a C function with an object parameter (closure) */
// llgo:link (*JSContext).JSNewCFunctionParams C.JS_NewCFunctionParams
func (recv_ *JSContext) JSNewCFunctionParams(func_idx c.Int, params JSValue) JSValue {
	return 0
}

// llgo:link (*JSContext).JSParse C.JS_Parse
func (recv_ *JSContext) JSParse(input *c.Char, input_len c.Int, filename *c.Char, eval_flags c.Int) JSValue {
	return 0
}

// llgo:link (*JSContext).JSRun C.JS_Run
func (recv_ *JSContext) JSRun(val JSValue) JSValue {
	return 0
}

// llgo:link (*JSContext).JSEval C.JS_Eval
func (recv_ *JSContext) JSEval(input *c.Char, input_len c.Int, filename *c.Char, eval_flags c.Int) JSValue {
	return 0
}

// llgo:link (*JSContext).JSGC C.JS_GC
func (recv_ *JSContext) JSGC() {
}

// llgo:link (*JSContext).JSNewStringLen C.JS_NewStringLen
func (recv_ *JSContext) JSNewStringLen(buf *c.Char, buf_len c.Int) JSValue {
	return 0
}

// llgo:link (*JSContext).JSNewString C.JS_NewString
func (recv_ *JSContext) JSNewString(buf *c.Char) JSValue {
	return 0
}

// llgo:link (*JSContext).JSToCStringLen C.JS_ToCStringLen
func (recv_ *JSContext) JSToCStringLen(plen *c.Int, val JSValue, buf *JSCStringBuf) *c.Char {
	return nil
}

// llgo:link (*JSContext).JSToCString C.JS_ToCString
func (recv_ *JSContext) JSToCString(val JSValue, buf *JSCStringBuf) *c.Char {
	return nil
}

// llgo:link (*JSContext).JSToString C.JS_ToString
func (recv_ *JSContext) JSToString(val JSValue) JSValue {
	return 0
}

// llgo:link (*JSContext).JSToInt32 C.JS_ToInt32
func (recv_ *JSContext) JSToInt32(pres *c.Int, val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSToUint32 C.JS_ToUint32
func (recv_ *JSContext) JSToUint32(pres *c.Uint32T, val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSToInt32Sat C.JS_ToInt32Sat
func (recv_ *JSContext) JSToInt32Sat(pres *c.Int, val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSToNumber C.JS_ToNumber
func (recv_ *JSContext) JSToNumber(pres *c.Double, val JSValue) c.Int {
	return 0
}

// llgo:link (*JSContext).JSGetException C.JS_GetException
func (recv_ *JSContext) JSGetException() JSValue {
	return 0
}

// llgo:link (*JSContext).JSStackCheck C.JS_StackCheck
func (recv_ *JSContext) JSStackCheck(len c.Uint32T) c.Int {
	return 0
}

// llgo:link (*JSContext).JSPushArg C.JS_PushArg
func (recv_ *JSContext) JSPushArg(val JSValue) {
}

// llgo:link (*JSContext).JSCall C.JS_Call
func (recv_ *JSContext) JSCall(call_flags c.Int) JSValue {
	return 0
}

type JSBytecodeHeader struct {
	Magic         c.Uint16T
	Version       c.Uint16T
	BaseAddr      c.UintptrT
	UniqueStrings JSValue
	MainFunc      JSValue
}

/* only used on the host when compiling to file */
// llgo:link (*JSContext).JSPrepareBytecode C.JS_PrepareBytecode
func (recv_ *JSContext) JSPrepareBytecode(hdr *JSBytecodeHeader, pdata_buf **c.Uint8T, pdata_len *c.Uint32T, eval_code JSValue) {
}

/* only used on the host when compiling to file */
// llgo:link (*JSContext).JSRelocateBytecode2 C.JS_RelocateBytecode2
func (recv_ *JSContext) JSRelocateBytecode2(hdr *JSBytecodeHeader, buf *c.Uint8T, buf_len c.Uint32T, new_base_addr c.UintptrT, update_atoms c.Int) c.Int {
	return 0
}

//go:linkname JSIsBytecode C.JS_IsBytecode
func JSIsBytecode(buf *c.Uint8T, buf_len c.Int) c.Int

/* Relocate the bytecode in 'buf' so that it can be executed
   later. Return 0 if OK, != 0 if error */
// llgo:link (*JSContext).JSRelocateBytecode C.JS_RelocateBytecode
func (recv_ *JSContext) JSRelocateBytecode(buf *c.Uint8T, buf_len c.Uint32T) c.Int {
	return 0
}

/* Load the precompiled bytecode from 'buf'. 'buf' must be allocated
   as long as the JSContext exists. Use JS_Run() to execute
   it. warning: the bytecode is not checked so it should come from a
   trusted source. */
// llgo:link (*JSContext).JSLoadBytecode C.JS_LoadBytecode
func (recv_ *JSContext) JSLoadBytecode(buf *c.Uint8T) JSValue {
	return 0
}

/* debug functions */
// llgo:link (*JSContext).JSSetLogFunc C.JS_SetLogFunc
func (recv_ *JSContext) JSSetLogFunc(write_func JSWriteFunc) {
}

// llgo:link (*JSContext).JSPrintValue C.JS_PrintValue
func (recv_ *JSContext) JSPrintValue(val JSValue) {
}

// llgo:link (*JSContext).JSPrintValueF C.JS_PrintValueF
func (recv_ *JSContext) JSPrintValueF(val JSValue, flags c.Int) {
}

// llgo:link (*JSContext).JSDumpValueF C.JS_DumpValueF
func (recv_ *JSContext) JSDumpValueF(str *c.Char, val JSValue, flags c.Int) {
}

// llgo:link (*JSContext).JSDumpValue C.JS_DumpValue
func (recv_ *JSContext) JSDumpValue(str *c.Char, val JSValue) {
}

// llgo:link (*JSContext).JSDumpMemory C.JS_DumpMemory
func (recv_ *JSContext) JSDumpMemory(is_long c.Int) {
}
