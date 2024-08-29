/*
 * @Author: jffan
 * @Date: 2024-08-01 15:15:30
 * @LastEditTime: 2024-08-15 10:08:17
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\secret\set.go
 * @Description: set secret
 */
package secret

import (
	"encoding/base64"
	"fmt"
	consts "tcas-cli/constants"
	"tcas-cli/manager"
	"tcas-cli/utils/file"
	"tcas-cli/utils/tools"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var secretSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set secret",
	Long:  `set secret`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		logrus.Debugf("secret set url: " + consts.ColorYellow + url + consts.OutReset)
		secretName, _ := cmd.Flags().GetString("name")
		if secretName == "" {
			secretName = tools.GenerateName("secret")
			logrus.Debugf("There is no name set So we generate a secret name (" + consts.ColorYellow + secretName + consts.OutReset + ") for you! ")
		}
		jsonFilePath, _ := cmd.Flags().GetString("file")
		if !file.IsExists(jsonFilePath) {
			logrus.Errorf(consts.ColorRed + "Error: The file you set is not exist! Please check it!" + consts.OutReset)
			return
		}
		logrus.Debugf("read file jsonpath: " + consts.ColorYellow + jsonFilePath + consts.OutReset)
		jsonData, err := file.ReadJSONFile(jsonFilePath)
		if err != nil {
			logrus.Errorf("read json error:", err)
			return
		}
		encodeJsonData := base64.StdEncoding.EncodeToString(jsonData)
		m, err := manager.New(url, "")
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}
		res, err := m.SetSecret(secretName, encodeJsonData)
		if err != nil {
			logrus.Errorf("Request failed: %v", err)
			return
		}
		if res.Code == 200 {
			fmt.Println(consts.ColorGreen + "set secret successful, secret id: " + res.Id + consts.OutReset)
		} else {
			logrus.Errorf(consts.ColorRed + "set secret failed:" + res.Message + consts.OutReset)
		}
	},
}

func init() {
	Cmd.AddCommand(secretSetCmd)
	//set parameter for secret set
	secretSetCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	secretSetCmd.Flags().StringP("name", "n", "", "must, the unique name of the secret")
	secretSetCmd.Flags().StringP("file", "f", "", "must, the path of secret file, only support json format")
	secretSetCmd.MarkFlagRequired("name")
	secretSetCmd.MarkFlagRequired("file")
}
