/*
 * @Author: jffan
 * @Date: 2024-08-01 15:36:12
 * @LastEditTime: 2024-08-13 10:53:43
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\secret\list.go
 * @Description: get the secret base info list
 */
package secret

import (
	"encoding/json"
	"fmt"
	"gitee.com/nanhu-lab/tcas-cli/cmd/initconfig"
	consts "gitee.com/nanhu-lab/tcas-cli/constants"
	"gitee.com/nanhu-lab/tcas-cli/manager"
	"gitee.com/nanhu-lab/tcas-cli/tees"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var secretListCmd = &cobra.Command{
	Use:   "list",
	Short: "get the secret base info list",
	Long:  `get the secret base info list`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		logrus.Debugf("secret url: " + consts.ColorYellow + url + consts.OutReset)
		tcasConfig, err := initconfig.ReadTCASInfoFromConfig(url)
		if err != nil {
			logrus.Errorf("get tcasConfig failed, please run <tcascli init> first")
			return
		}
		m, err := manager.New(tcasConfig, tees.GetCollectors())
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}
		res, err := m.ListSecret()
		if err != nil {
			logrus.Errorf("Request secret list failed: %v", err)
			return
		}
		if res.Code == 200 {
			jsonData, err := json.MarshalIndent(manager.SecretListJsonFormat{Secrets: res.Data}, "", "  ")
			if err != nil {
				logrus.Errorf("Error marshaling JSON:", err)
				return
			} else {
				fmt.Println("------------------secret list start------------------")
				fmt.Println(consts.ColorYellow + string(jsonData) + consts.OutReset)
				fmt.Println("------------------secret list end--------------------")
			}
		} else {
			logrus.Errorf(consts.ColorRed + "request secret list failed:" + res.Message + consts.OutReset)
		}
	},
}

func init() {
	Cmd.AddCommand(secretListCmd)
	//set parameter for secret list
	secretListCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
}
