/*
 * @Author: jffan
 * @Date: 2024-07-31 16:34:14
 * @LastEditTime: 2024-08-13 10:54:37
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\policy\list.go
 * @Description: get policy list
 */
package policy

import (
	"encoding/json"
	"fmt"
	consts "tcas-cli/constants"
	"tcas-cli/manager"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var policyListCmd = &cobra.Command{
	Use:   "list",
	Short: "get policy list",
	Long:  `get policy list`,
	Run: func(cmd *cobra.Command, args []string) {
		attestationType, _ := cmd.Flags().GetString("type")
		logrus.Debugf("policy type: " + consts.ColorYellow + attestationType + consts.OutReset)

		url, _ := cmd.Flags().GetString("url")
		logrus.Debugf("policy url: " + consts.ColorYellow + url + consts.OutReset)

		m, err := manager.New(url, "")
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}
		res, err := m.ListPolicy(attestationType)
		if err != nil {
			logrus.Errorf("Request policy list failed: %v", err)
			return
		}
		if res.Code == 200 {
			jsonData, err := json.MarshalIndent(manager.PolicyListJsonFormat{Policies: res.Data}, "", "  ")
			if err != nil {
				logrus.Errorf("Error marshaling JSON:", err)
				return
			} else {
				fmt.Println("------------------policy list start------------------")
				fmt.Println(consts.ColorYellow + string(jsonData) + consts.OutReset)
				fmt.Println("------------------policy list end--------------------")
			}
		} else {
			logrus.Errorf(consts.ColorRed + "request policy list failed:" + res.Message + consts.OutReset)
		}
	},
}

func init() {
	Cmd.AddCommand(policyListCmd)
	//set parameter for policy delete
	policyListCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	policyListCmd.Flags().StringP("type", "t", "trust_node", "optional, the attestation-type of policy, support `trust_node` or `trust_cluster`")
}
