package virtcca

import (
	"encoding/base64"
	"fmt"
	"gitee.com/nanhu-lab/tcas-cli/collectors"
	"github.com/sirupsen/logrus"
)

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) CollectEvidence(userdata []byte) (*collectors.Evidence, error) {
	token, err := GetAttestationToken(userdata)
	if err != nil {
		return nil, fmt.Errorf("get vcca token failed, error: %s", err)
	}

	logrus.Debugf("vcca token: %x", token)

	cert, err := GetDevCert()
	if err != nil {
		return nil, fmt.Errorf("get vcca aik cert failed, error: %s", err)
	}
	logrus.Debugf("vcca aik cert: %x", cert)

	return &collectors.Evidence{
		TeeType:   "virtcca",
		TeeReport: base64.StdEncoding.EncodeToString(token),
		Parameter: struct {
			X5c string `json:"x5c"`
		}{
			base64.StdEncoding.EncodeToString(cert),
		},
	}, nil

}

func (c *Collector) Name() string {
	return "virtcca"
}
