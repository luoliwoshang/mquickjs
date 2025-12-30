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
