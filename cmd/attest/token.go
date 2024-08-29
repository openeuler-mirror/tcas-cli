package attest

import (
	"fmt"
	"tcas-cli/manager"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "attest for getting token",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		tee, _ := cmd.Flags().GetString("tee")
		devices, _ := cmd.Flags().GetString("devices")
		userdata, _ := cmd.Flags().GetString("userdata")
		policies, _ := cmd.Flags().GetString("policies")

		m, err := manager.New(url, "")
		if err != nil {
			return
		}
		res, err := m.AttestForToken(tee, userdata, devices, policies)
		if err != nil {
			logrus.Errorf("do attest for token failed, error: %s", err)
			return
		}

		fmt.Println(res.Token)
	},
}

func init() {
	Cmd.AddCommand(tokenCmd)
	//set parameter for getting token
	tokenCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	tokenCmd.Flags().StringP("tee", "t", "", "must, tee type")
	tokenCmd.MarkFlagRequired("tee")
	tokenCmd.Flags().StringP("devices", "v", "", "optional, the trust devices")

	tokenCmd.Flags().StringP("userdata", "d", "", "optional, the base64 encoded userdata")
	tokenCmd.Flags().StringP("policies", "p", "", "optional, the ids of the policy needed matching")
}
