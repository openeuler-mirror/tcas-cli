/*
 * @Author: jffan
 * @Date: 2024-08-01 15:18:37
 * @LastEditTime: 2024-08-15 16:36:47
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\secret\update.go
 * @Description: update secret
 */
package secret

import (
	"encoding/base64"
	"fmt"
	consts "tcas-cli/constants"
	"tcas-cli/manager"
	"tcas-cli/utils/file"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var secretUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update secret",
	Long:  `update secret`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		logrus.Debugf("secret update url: " + consts.ColorYellow + url + consts.OutReset)
		secretId, _ := cmd.Flags().GetString("id")
		logrus.Debugf("secret update id: " + consts.ColorYellow + secretId + consts.OutReset)
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
		newEncodeJsonData := base64.StdEncoding.EncodeToString(jsonData)
		m, err := manager.New(url, "")
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}
		res, err := m.UpdateSecret(secretId, newEncodeJsonData)
		if err != nil {
			logrus.Errorf("Request failed: %v", err)
			return
		}
		if res.Code == 200 {
			fmt.Println(consts.ColorGreen + "update secret successful, secret id: " + res.Id + consts.OutReset)
		} else {
			logrus.Errorf(consts.ColorRed + "update secret failed:" + res.Message + consts.OutReset)
		}
	},
}

func init() {
	Cmd.AddCommand(secretUpdateCmd)
	//set parameter for secret update
	secretUpdateCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	secretUpdateCmd.Flags().StringP("file", "f", "", "must, the path of new secret file, only support json format")
	secretUpdateCmd.Flags().StringP("id", "i", "", "must, the id of the old secret")
	secretUpdateCmd.MarkFlagRequired("id")
	secretUpdateCmd.MarkFlagRequired("file")
}
