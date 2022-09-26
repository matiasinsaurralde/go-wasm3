//
//  m3_env.h
//
//  Created by Steven Massey on 4/19/19.
//  Copyright © 2019 Steven Massey. All rights reserved.
//

#ifndef m3_env_h
#define m3_env_h

#include "wasm3.h"
#include "m3_code.h"
#include "m3_exec.h"
#include "m3_compile.h"

d_m3BeginExternC

typedef struct M3FuncType
{
    struct M3FuncType *     next;

    u32                     numArgs;
    u8                      returnType;
    u8                      argTypes        [3];    // M3FuncType is a dynamically sized object; these are padding
}
M3FuncType;

typedef M3FuncType *        IM3FuncType;

M3Result    AllocFuncType                   (IM3FuncType * o_functionType, u32 i_numArgs);
bool        AreFuncTypesEqual               (const IM3FuncType i_typeA, const IM3FuncType i_typeB);


//---------------------------------------------------------------------------------------------------------------------------------
typedef struct M3Function
{
    struct M3Module *       module;

    M3ImportInfo            import;

    bytes_t                 wasm;
    bytes_t                 wasmEnd;

    cstr_t                  name;

    IM3FuncType             funcType;

    pc_t                    compiled;

#   if (d_m3EnableCodePageRefCounting)
    IM3CodePage *           codePageRefs;           // array of all pages used
    u32                     numCodePageRefs;
#   endif

#   if defined(DEBUG)
    u32                     hits;
#   endif

    u16                     maxStackSlots;

    u16                     numArgSlots;

    u16                     numLocals;          // not including args
    u16                     numLocalBytes;

    void *                  constants;
    u16                     numConstantBytes;

    bool                    ownsWasmCode;
}
M3Function;

void        Function_Release            (IM3Function i_function);
void        Function_FreeCompiledCode   (IM3Function i_function);

cstr_t      GetFunctionImportModuleName (IM3Function i_function);
cstr_t      GetFunctionName             (IM3Function i_function);
u32         GetFunctionNumArgs          (IM3Function i_function);
u32         GetFunctionNumReturns       (IM3Function i_function);
u8          GetFunctionReturnType       (IM3Function i_function);

u32         GetFunctionNumArgsAndLocals (IM3Function i_function);

cstr_t      SPrintFunctionArgList       (IM3Function i_function, m3stack_t i_sp);


//---------------------------------------------------------------------------------------------------------------------------------

typedef struct M3MemoryInfo
{
    u32     initPages;
    u32     maxPages;
}
M3MemoryInfo;


typedef struct M3Memory
{
    M3MemoryHeader *        mallocated;

    u32                     numPages;
    u32                     maxPages;
}
M3Memory;

typedef M3Memory *          IM3Memory;


//---------------------------------------------------------------------------------------------------------------------------------

typedef struct M3DataSegment
{
    const u8 *              initExpr;       // wasm code
    const u8 *              data;

    u32                     initExprSize;
    u32                     memoryRegion;
    u32                     size;
}
M3DataSegment;


void FreeImportInfo (M3ImportInfo * i_info);

//---------------------------------------------------------------------------------------------------------------------------------

typedef struct M3Global
{
    M3ImportInfo            import;

    union
    {
        i64 intValue;
#if d_m3HasFloat
        f64 f64Value;
        f32 f32Value;
#endif
    };

    bytes_t                 initExpr;       // wasm code
    u32                     initExprSize;
    u8                      type;
    bool                    imported;
    bool                    isMutable;
}
M3Global;
typedef M3Global *          IM3Global;


//---------------------------------------------------------------------------------------------------------------------------------
typedef struct M3Module
{
    struct M3Runtime *      runtime;
    struct M3Environment *  environment;

    cstr_t                  name;

    u32                     numFuncTypes;
    IM3FuncType *           funcTypes;          // array of pointers to list of FuncTypes

    u32                     numImports;
    IM3Function *           imports;            // notice: "I" prefix. imports are pointers to functions in another module.

    u32                     numFunctions;
    M3Function *            functions;

    i32                     startFunction;

    u32                     numDataSegments;
    M3DataSegment *         dataSegments;

    u32                     importedGlobals;
    u32                     numGlobals;
    M3Global *              globals;

    u32                     numElementSegments;
    bytes_t                 elementSection;
    bytes_t                 elementSectionEnd;

    IM3Function *           table0;
    u32                     table0Size;

    M3MemoryInfo            memoryInfo;
    bool                    memoryImported;

    bool                    hasWasmCodeCopy;

    struct M3Module *       next;
}
M3Module;

M3Result                    Module_AddGlobal            (IM3Module io_module, IM3Global * o_global, u8 i_type, bool i_mutable, bool i_isImported);

M3Result                    Module_AddFunction          (IM3Module io_module, u32 i_typeIndex, IM3ImportInfo i_importInfo /* can be null */);
IM3Function                 Module_GetFunction          (IM3Module i_module, u32 i_functionIndex);

//---------------------------------------------------------------------------------------------------------------------------------

static const u32 c_m3NumTypesPerPage = 8;

//---------------------------------------------------------------------------------------------------------------------------------

typedef struct M3Environment
{
//    struct M3Runtime *      runtimes;

    IM3FuncType             funcTypes;          // linked list

    M3CodePage *            pagesReleased;
}
M3Environment;

void                        Environment_Release         (IM3Environment i_environment);

// takes ownership of io_funcType and returns a pointer to the persistent version (could be same or different)
void                        Environment_AddFuncType     (IM3Environment i_environment, IM3FuncType * io_funcType);

//---------------------------------------------------------------------------------------------------------------------------------

// OPTZ: function types need to move to the runtime structure so that all modules can share types
// then type equality can be a simple pointer compare for indirect call checks
typedef struct M3Runtime
{
    M3Compilation           compilation;

    IM3Environment          environment;

    M3CodePage *            pagesOpen;      // linked list of code pages with writable space on them
    M3CodePage *            pagesFull;      // linked list of at-capacity pages

    u32                     numCodePages;
    u32                     numActiveCodePages;

    IM3Module               modules;        // linked list of imported modules

    void *                  stack;
    u32                     stackSize;
    u32                     numStackSlots;

    u32                     argc;
    ccstr_t *               argv;

    M3Result                runtimeError;

    M3Memory                memory;
    u32                     memoryLimit;

    M3ErrorInfo             error;
#if d_m3VerboseLogs
    char                    error_message[256];
#endif
    i32                     exit_code;
}
M3Runtime;

void                        InitRuntime                 (IM3Runtime io_runtime, u32 i_stackSizeInBytes);
void                        Runtime_Release             (IM3Runtime io_runtime);

M3Result                    ResizeMemory                (IM3Runtime io_runtime, u32 i_numPages);

typedef void *              (* ModuleVisitor)           (IM3Module i_module, void * i_info);
void *                      ForEachModule               (IM3Runtime i_runtime, ModuleVisitor i_visitor, void * i_info);

void *                      v_FindFunction              (IM3Module i_module, const char * const i_name);

IM3CodePage                 AcquireCodePage             (IM3Runtime io_runtime);
IM3CodePage                 AcquireCodePageWithCapacity (IM3Runtime io_runtime, u32 i_lineCount);
void                        ReleaseCodePage             (IM3Runtime io_runtime, IM3CodePage i_codePage);

M3Result                    m3Error                     (M3Result i_result, IM3Runtime i_runtime, IM3Module i_module, IM3Function i_function, const char * const i_file, u32 i_lineNum, const char * const i_errorMessage, ...);

d_m3EndExternC

#endif // m3_env_h
