package evidence

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	consts "gitee.com/nanhu-lab/tcas-cli/constants"
	"gitee.com/nanhu-lab/tcas-cli/manager"
	"gitee.com/nanhu-lab/tcas-cli/tees"
	"gitee.com/nanhu-lab/tcas-cli/utils/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getNodeEvidenceInfo(tee, reportdata, devices string) (*manager.NodeEvidence, error) {
	if tee == "" {
		return nil, fmt.Errorf("tee is null")
	}
	userdata, err := base64.StdEncoding.DecodeString(reportdata)
	if err != nil {
		return nil, fmt.Errorf("Base64 decoding failed:%s", err)
	}

	if len(userdata) > 64 {
		return nil, fmt.Errorf("The length of the decoded data exceeds the limit and must be less than 64 bytes.")
	}
	m := tees.GetCollectors()
	c, ok := m[tee]
	if !ok {
		return nil, fmt.Errorf("tee: %s not support yet", tee)
	}
	teeReport, err := c.CollectEvidence([]byte(userdata))
	if err != nil {
		return nil, fmt.Errorf("collect evidence from %s failed, error: %s", c.Name(), err)
	}

	var devicelist []string
	if devices != "" {
		devicelist = strings.Split(devices, ",")
	}
	deviceReports := make([]*manager.TrustDeviceReport, 0)
	for _, d := range devicelist {
		c, ok := m[d]
		if !ok {
			return nil, fmt.Errorf("device tee: %s not support yet", tee)
		}
		deviceEvidence, err := c.CollectEvidence([]byte(userdata))
		if err != nil {
			return nil, fmt.Errorf("collect evidence from %s failed, error: %s", c.Name(), err)
		}

		deviceReports = append(deviceReports, &manager.TrustDeviceReport{
			DeviceType:   d,
			DeviceReport: deviceEvidence.TeeReport,
			Parameter:    deviceEvidence.Parameter,
		})
	}

	return &manager.NodeEvidence{
		Tee:         teeReport.TeeType,
		TeeReport:   teeReport.TeeReport,
		Parameter:   teeReport.Parameter,
		RuntimeData: reportdata,
		TrustDevice: deviceReports,
	}, nil
}

var Cmd = &cobra.Command{
	Use:   "evidence",
	Short: "get evidence",
	Long:  "obtain the tee report in JSON format",
	//SilenceUsage:               true,
	SuggestionsMinimumDistance: 1,
	DisableSuggestions:         false,
	Run: func(cmd *cobra.Command, args []string) {
		tee, _ := cmd.Flags().GetString("tee")
		devices, _ := cmd.Flags().GetString("devices")
		reportdata, _ := cmd.Flags().GetString("reportdata")

		outputPath, _ := cmd.Flags().GetString("output")

		res, err := getNodeEvidenceInfo(tee, reportdata, devices)
		if err != nil {
			logrus.Errorf("get evidence info failed, error: %s", err)
			return
		}
		jsonBytes, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			logrus.Errorf("failed to marshal response to JSON: %s", err)
			return
		}
		fmt.Println("------------------evidence info start------------------")
		fmt.Println(consts.ColorYellow + string(jsonBytes) + consts.OutReset)
		fmt.Println("------------------evidence info end--------------------")
		fileName := fmt.Sprintf("%s-evidence.json", tee)
		err = file.EnsureDirExists(outputPath)
		if err != nil {
			fmt.Errorf("Error ensuring directory exists: %v", err)
			return
		}
		resultPath := filepath.Join(outputPath, fileName)
		err = os.WriteFile(resultPath, jsonBytes, 0644)
		if err != nil {
			logrus.Errorf("failed to write evidence to file: %s", err)
			return
		}
		fmt.Println(consts.ColorGreen + "The evidence successfully saved in: " + resultPath + consts.OutReset)
	},
}

func init() {
	Cmd.Flags().StringP("tee", "t", "", "must, tee type")
	Cmd.MarkFlagRequired("tee")
	Cmd.Flags().StringP("output", "o", "./tee-reports", "optional, the save path of the tees report")
	Cmd.Flags().StringP("devices", "v", "", "optional, the trust devices")
	Cmd.Flags().StringP("reportdata", "r", "", "optional, the base64 encoded repotdata")
}
