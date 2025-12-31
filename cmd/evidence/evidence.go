package evidence

import (
	"encoding/json"
	"fmt"
	consts "gitee.com/nanhu-lab/tcas-cli/constants"
	"gitee.com/nanhu-lab/tcas-cli/teetool"
	"os"
	"path/filepath"

	"gitee.com/nanhu-lab/tcas-cli/utils/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
		userdata, _ := cmd.Flags().GetString("userdata")
		reportdata, _ := cmd.Flags().GetString("reportdata")

		outputPath, _ := cmd.Flags().GetString("output")

		var nodeEvidence *teetool.NodeEvidence
		var err error

		if userdata != "" {
			nodeEvidence, err = teetool.GetNodeEvidenceWithUserData(tee, devices, userdata, "")
			if err != nil {
				logrus.Errorf("get evidence info with UserData failed, error: %s", err)
				return
			}
		} else {
			nodeEvidence, err = teetool.GetNodeEvidenceInfo(tee, devices, []byte(reportdata))
			if err != nil {
				logrus.Errorf("get evidence info failed, error: %s", err)
				return
			}
		}

		jsonBytes, err := json.MarshalIndent(nodeEvidence, "", "  ")
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
	Cmd.Flags().StringP("userdata", "u", "", "optional, the base64 encoded userdata")
}
