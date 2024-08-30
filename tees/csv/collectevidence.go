package csv

import (
	"encoding/base64"
	"fmt"
	"gitee.com/nanhu-lab/tcas-cli/collectors"
)

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}
func (c *Collector) CollectEvidence(userdata []byte) (*collectors.Evidence, error) {
	report, err := GetReportInByte(userdata)
	if err != nil {
		return nil, fmt.Errorf("get attestation report failed, error: %s", err)
	}

	return &collectors.Evidence{
		TeeType:   "csv",
		TeeReport: base64.StdEncoding.EncodeToString(report),
	}, nil
}

func (c *Collector) Name() string {
	return "csv"
}
