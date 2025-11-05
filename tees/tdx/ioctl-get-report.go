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
#cgo LDFLAGS: -ltdx_attest
#include "./tdx_attest.h"
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

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

func getTDXQuote(reportData *TDXReportData) ([]byte, error) {
	// Convert Go struct to C struct
	cReportData := C.tdx_report_data_t{}
	for i := 0; i < TDX_REPORT_DATA_SIZE; i++ {
		cReportData.d[i] = C.uint8_t(reportData.Data[i])
	}

	var cSelectedAttKeyId C.tdx_uuid_t
	var cQuoteBuf *C.uint8_t
	var cQuoteSize C.uint32_t

	result := C.tdx_att_get_quote(&cReportData, nil, 0, &cSelectedAttKeyId, &cQuoteBuf, &cQuoteSize, 0)
	if result != TDX_ATTEST_SUCCESS {
		return nil, fmt.Errorf("failed to get the quote [error code: %d] ", int(result))
	}
	defer C.tdx_att_free_quote(cQuoteBuf)

	// Convert quote buffer to Go slice
	quoteBuf := C.GoBytes(unsafe.Pointer(cQuoteBuf), C.int(cQuoteSize))

	return quoteBuf, nil
}
