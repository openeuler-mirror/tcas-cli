/*
 * @Author: jffan
 * @Date: 2024-07-31 16:39:36
 * @LastEditTime: 2024-08-02 16:56:48
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\policy\detail.go
 * @Description: get the detail of the policy
 */
package policy

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var policyDetailCmd = &cobra.Command{
	Use:   "detail",
	Short: "get the detail of the policy",
	Long:  `get the detail of the policy`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("policy detail called")
	},
}

func init() {
	Cmd.AddCommand(policyDetailCmd)
	//set parameter for policy delete
	policyDetailCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
}
