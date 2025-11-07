/*
 * @Author: jffan
 * @Date: 2025-11-03 16:31:36
 * @LastEditTime: 2025-11-05 10:21:27
 * @LastEditors: jffan
 * @FilePath: \tcas-cli\tees\tdx\ioctl-get-report.go
 * @Description:
 */
package tdx

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
#include <stdint.h>

// 定义TDX相关结构体
typedef struct tdx_report_data_t {
    uint8_t d[64];
} tdx_report_data_t;

typedef struct tdx_uuid_t {
    uint8_t d[16];
} tdx_uuid_t;

// 函数指针类型定义
typedef int (*tdx_att_get_quote_func_t)(
    const tdx_report_data_t *report_data,
    const uint8_t *ud_data,
    uint32_t ud_data_size,
    tdx_uuid_t *selected_att_key_id,
    uint8_t **quote_buf,
    uint32_t *quote_size,
    uint32_t flags
);

typedef void (*tdx_att_free_quote_func_t)(uint8_t *quote_buf);

// 全局函数指针
static tdx_att_get_quote_func_t tdx_att_get_quote_ptr = NULL;
static tdx_att_free_quote_func_t tdx_att_free_quote_ptr = NULL;

// 库句柄
static void* tdx_lib_handle = NULL;

// 加载TDX库
int load_tdx_library() {
    if (tdx_lib_handle != NULL && tdx_att_get_quote_ptr != NULL && tdx_att_free_quote_ptr != NULL) {
        return 0; // 已经加载
    }

    // 加载库
    tdx_lib_handle = dlopen("libtdx_attest.so", RTLD_LAZY);
    if (tdx_lib_handle == NULL) {
        // 尝试其他可能的库名称
        tdx_lib_handle = dlopen("libtdx_attest.so.1", RTLD_LAZY);
        if (tdx_lib_handle == NULL) {
            return -1; // 加载失败
        }
    }

    // 获取函数指针
    tdx_att_get_quote_ptr = (tdx_att_get_quote_func_t) dlsym(tdx_lib_handle, "tdx_att_get_quote");
    if (tdx_att_get_quote_ptr == NULL) {
        dlclose(tdx_lib_handle);
        tdx_lib_handle = NULL;
        return -1;
    }

    tdx_att_free_quote_ptr = (tdx_att_free_quote_func_t) dlsym(tdx_lib_handle, "tdx_att_free_quote");
    if (tdx_att_free_quote_ptr == NULL) {
        dlclose(tdx_lib_handle);
        tdx_lib_handle = NULL;
        return -1;
    }

    return 0; // 成功
}

// 卸载TDX库
void unload_tdx_library() {
    if (tdx_lib_handle != NULL) {
        dlclose(tdx_lib_handle);
        tdx_lib_handle = NULL;
    }
    tdx_att_get_quote_ptr = NULL;
    tdx_att_free_quote_ptr = NULL;
}

// 包装函数，通过函数指针调用实际函数
int call_tdx_att_get_quote(
    const tdx_report_data_t *report_data,
    const uint8_t *ud_data,
    uint32_t ud_data_size,
    tdx_uuid_t *selected_att_key_id,
    uint8_t **quote_buf,
    uint32_t *quote_size,
    uint32_t flags) {

    if (tdx_att_get_quote_ptr == NULL) {
        return -1;
    }

    return tdx_att_get_quote_ptr(report_data, ud_data, ud_data_size, selected_att_key_id, quote_buf, quote_size, flags);
}

void call_tdx_att_free_quote(uint8_t *quote_buf) {
    if (tdx_att_free_quote_ptr != NULL) {
        tdx_att_free_quote_ptr(quote_buf);
    }
}

// 检查库是否可用
int is_tdx_library_available() {
    return (load_tdx_library() == 0) ? 1 : 0;
}
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

var tdxMutex sync.Mutex

const (
	TDX_REPORT_DATA_SIZE = 64
	TDX_ATTEST_SUCCESS   = 0
)

type TDXReportData struct {
	Data [TDX_REPORT_DATA_SIZE]uint8
}

type TDXReport struct {
	Data [1024]uint8 // Size may vary based on actual TDX report structure
}

type TDXUUID struct {
	Data [16]uint8
}

// getTDXQuote 通过动态加载的库获取TDX Quote
func getTDXQuote(reportData *TDXReportData) ([]byte, error) {
	tdxMutex.Lock()
	defer tdxMutex.Unlock()
	// 检查库是否可用
	available := C.is_tdx_library_available()
	if available == 0 {
		return nil, fmt.Errorf("TDX library not available on this system")
	}

	// Convert Go struct to C struct
	cReportData := C.tdx_report_data_t{}
	for i := 0; i < TDX_REPORT_DATA_SIZE; i++ {
		cReportData.d[i] = C.uint8_t(reportData.Data[i])
	}

	var cSelectedAttKeyId C.tdx_uuid_t
	var cQuoteBuf *C.uint8_t
	var cQuoteSize C.uint32_t

	// 通过包装函数调用实际的TDX函数
	result := C.call_tdx_att_get_quote(
		&cReportData,
		nil,
		0,
		&cSelectedAttKeyId,
		&cQuoteBuf,
		&cQuoteSize,
		0,
	)

	if result != TDX_ATTEST_SUCCESS {
		return nil, fmt.Errorf("failed to get the quote [error code: %d]", int(result))
	}

	defer C.call_tdx_att_free_quote(cQuoteBuf)

	// Convert quote buffer to Go slice
	quoteBuf := C.GoBytes(unsafe.Pointer(cQuoteBuf), C.int(cQuoteSize))

	return quoteBuf, nil
}

// CloseTDXLibrary 提供关闭TDX库的公共接口
func CloseTDXLibrary() {
	C.unload_tdx_library()
}
