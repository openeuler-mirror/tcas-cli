package nvidia

import (
	"encoding/json"
	"fmt"
	"os"
	"tcas-cli/collectors"
)

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) CollectEvidence(userdata []byte) (*collectors.Evidence, error) {

	report, err := os.ReadFile("nvidia-report.json")
	if err != nil {
		return nil, fmt.Errorf("open report file failed, error: %s", err)
	}

	ev := new(collectors.Evidence)
	err = json.Unmarshal(report, ev)
	if err != nil {
		return nil, fmt.Errorf("the format of report is not correct")
	}
	return ev, nil

}

func (c *Collector) Name() string {
	return "nvidia"
}
