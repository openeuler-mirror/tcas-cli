package teetool

import (
	"crypto/sha256"
	"fmt"
	"gitee.com/nanhu-lab/tcas-cli/tees"
	"github.com/sirupsen/logrus"
	"strings"
)

func GetNodeEvidenceInfo(tee, devices string, reportData []byte) (*NodeEvidence, error) {
	if tee == "" {
		return nil, fmt.Errorf("tee is null")
	}

	if len(reportData) > 64 {
		return nil, fmt.Errorf("the length of report data exceeds the limit and must be less than 64 bytes")
	}
	m := tees.GetCollectors()
	c, ok := m[tee]
	if !ok {
		return nil, fmt.Errorf("tee: %s not support yet", tee)
	}
	teeReport, err := c.CollectEvidence(reportData)
	if err != nil {
		return nil, fmt.Errorf("collect evidence from %s failed, error: %s", c.Name(), err)
	}

	var deviceList []string
	if devices != "" {
		deviceList = strings.Split(devices, ",")
	}
	deviceReports := make([]*TrustDeviceReport, 0)
	for _, d := range deviceList {
		c, ok := m[d]
		if !ok {
			return nil, fmt.Errorf("device tee: %s not support yet", tee)
		}
		deviceEvidence, err := c.CollectEvidence(reportData)
		if err != nil {
			return nil, fmt.Errorf("collect evidence from %s failed, error: %s", c.Name(), err)
		}

		deviceReports = append(deviceReports, &TrustDeviceReport{
			DeviceType:   d,
			DeviceReport: deviceEvidence.TeeReport,
			Parameter:    deviceEvidence.Parameter,
		})
	}

	return &NodeEvidence{
		Tee:         teeReport.TeeType,
		TeeReport:   teeReport.TeeReport,
		Parameter:   teeReport.Parameter,
		ReportData:  reportData,
		TrustDevice: deviceReports,
	}, nil
}

func GetNodeEvidenceWithUserData(tee, devices, userData, nonce string) (*NodeEvidence, error) {
	var reportData string
	if userData == "" {
		userData = "test report"
		logrus.Warnf("userData is nil, use [%s] as default", userData)
	}

	if nonce == "" {
		logrus.Warnf("nonce is nil")
	}

	shaData := nonce + userData
	reportData = fmt.Sprintf("%x", sha256.Sum256([]byte(shaData)))
	ne, err := GetNodeEvidenceInfo(tee, devices, []byte(reportData))
	if err != nil {
		return nil, err
	}
	ne.UserData = userData

	return ne, nil
}
