/*
 * @Author: jffan
 * @Date: 2025-11-03 16:28:00
 * @LastEditTime: 2025-11-05 09:56:46
 * @LastEditors: jffan
 * @FilePath: \tcas-cli\tees\tdx\collectevidence.go
 * @Description:
 */
package tdx

import (
	"encoding/base64"
	"fmt"

	"gitee.com/nanhu-lab/tcas-cli/collectors"
	"github.com/sirupsen/logrus"
)

const TEE_TYPE_NAME = "tdx"

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}
func (c *Collector) CollectEvidence(userdata []byte) (*collectors.Evidence, error) {
	if len(userdata) > TDX_REPORT_DATA_SIZE {
		return nil, fmt.Errorf("userdata size exceeds TDX report data size")
	}
	var reportData TDXReportData
	copy(reportData.Data[:], userdata)
	logrus.Debugf("reportData.Data: %x\n", reportData.Data)
	report, err := getTDXQuote(&reportData)
	if err != nil {
		return nil, fmt.Errorf("get attestation report failed, error: %s", err)
	}

	return &collectors.Evidence{
		TeeType:   TEE_TYPE_NAME,
		TeeReport: base64.StdEncoding.EncodeToString(report),
	}, nil
}
func (c *Collector) Name() string {
	return TEE_TYPE_NAME
}
