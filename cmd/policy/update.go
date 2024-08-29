/*
 * @Author: jffan
 * @Date: 2024-07-31 16:41:04
 * @LastEditTime: 2024-08-06 16:10:50
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\policy\update.go
 * @Description: update policy
 */
package policy

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var policyUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update policy",
	Long:  `update policy`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("policy update called")
	},
}

func init() {
	Cmd.AddCommand(policyUpdateCmd)
	//set parameter for policy delete
	policyUpdateCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
}
