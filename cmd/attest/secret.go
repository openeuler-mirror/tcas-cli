package attest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	consts "tcas-cli/constants"
	"tcas-cli/manager"
	"tcas-cli/utils/file"
	"tcas-cli/utils/tools"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "attest for getting secret ",
	Long:  `attest for getting secret `,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		tee, _ := cmd.Flags().GetString("tee")
		devices, _ := cmd.Flags().GetString("devices")
		userdata, _ := cmd.Flags().GetString("userdata")
		policies, _ := cmd.Flags().GetString("policies")
		secretID, _ := cmd.Flags().GetString("secretId")
		m, err := manager.New(url, "")
		if err != nil {
			return
		}
		res, err := m.AttestForSecret(tee, userdata, devices, policies, secretID)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		jsonData, err := json.MarshalIndent(res.Secret, "", "  ")
		if err != nil {
			logrus.Errorf("Error marshaling JSON:", err)
			return
		}
		jsonData = append(jsonData, '\n')
		outputPath, _ := cmd.Flags().GetString("output")
		err = file.EnsureDirExists(outputPath)
		if err != nil {
			logrus.Errorf("Error ensuring directory exists: %v", err)
			return
		}
		fileName := tools.GenerateName("secret") + ".json"
		resultPath := filepath.Join(outputPath, fileName)
		err = os.WriteFile(resultPath, jsonData, 0644)
		if err != nil {
			logrus.Errorf("Error writing secret json file: %v", err)
			return
		}
		fmt.Println(consts.ColorGreen + "The secert successfully saved in: " + resultPath + consts.OutReset)
	},
}

func init() {
	Cmd.AddCommand(secretCmd)
	//set parameter for getting secret
	secretCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	secretCmd.Flags().StringP("tee", "t", "", "must, tee type")
	secretCmd.Flags().StringP("devices", "v", "", "optional, the trust devices")
	secretCmd.Flags().StringP("policies", "p", "", "optional, the ids of the policy needed matching")
	secretCmd.Flags().StringP("userdata", "d", "", "optional, the base64 encoded userdata")

	secretCmd.Flags().StringP("secretId", "s", "", "must,the secret ID that needs to be obtained")
	secretCmd.Flags().StringP("output", "o", "./tcas-secret", "optional, the output dir of the secret")

	secretCmd.MarkFlagRequired("tee")
	secretCmd.MarkFlagRequired("secretId")

}
