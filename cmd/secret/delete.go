/*
 * @Author: jffan
 * @Date: 2024-08-01 15:47:27
 * @LastEditTime: 2024-08-14 14:36:14
 * @LastEditors: jffan
 * @FilePath: \tcas-cli\cmd\secret\delete.go
 * @Description: delete secret
 */
package secret

import (
	"fmt"
	consts "tcas-cli/constants"
	"tcas-cli/manager"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var secretDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete secret",
	Long:  `delete secret`,
	Run: func(cmd *cobra.Command, args []string) {
		secretID, _ := cmd.Flags().GetString("id")
		if secretID == "" {
			logrus.Errorf(consts.ColorRed + "secret id is required ! use `--id <secret_id>`" + consts.OutReset)
			return
		}
		url, _ := cmd.Flags().GetString("url")
		logrus.Debugf("secret delete url: " + consts.ColorYellow + url + consts.OutReset)

		m, err := manager.New(url, "")
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}
		res, err := m.DeleteSecret(secretID)
		if err != nil {
			logrus.Errorf("Request failed: %v", err)
			return
		}
		if res.Code == 200 {
			fmt.Println(consts.ColorGreen + "delete secret successful, secret id: " + res.SecretID + consts.OutReset)
		} else {
			logrus.Errorf(consts.ColorRed + "delete secret failed:" + res.Message + consts.OutReset)
		}
	},
}

func init() {
	Cmd.AddCommand(secretDeleteCmd)
	//set parameter for secret delete
	secretDeleteCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	secretDeleteCmd.Flags().StringP("id", "i", "", "the id of the secret")
	secretDeleteCmd.MarkFlagRequired("id")
}
