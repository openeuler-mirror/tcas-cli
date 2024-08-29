/*
- @Author: jffan
- @Date: 2024-08-05 09:34:20
- @LastEditTime: 2024-08-14 16:39:35
- @LastEditors: jffan
- @FilePath: \tcas-cli\cmd\ca\ca.go
- @Description: Command to obtain a CA certificate
*/
package ca

import (
	"fmt"
	"os"
	"path/filepath"
	consts "tcas-cli/constants"
	"tcas-cli/manager"
	"tcas-cli/utils/file"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Cmd represents the ca command
var Cmd = &cobra.Command{
	Use:                        "ca",
	Short:                      "get root CA",
	Long:                       "get root CA",
	SuggestionsMinimumDistance: 1,
	DisableSuggestions:         false,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		logrus.Debugf("ca url " + consts.ColorGreen + url + consts.OutReset)
		ouputPath, _ := cmd.Flags().GetString("output")
		logrus.Debugf(consts.ColorYellow + "Save the file in the path （`" + ouputPath + "`）" + consts.OutReset)
		err := file.EnsureDirExists(ouputPath)
		if err != nil {
			logrus.Errorf("Error ensuring directory exists: %v", err)
			return
		}
		//getRootCert
		m, err := manager.New(url, "")
		if err != nil {
			logrus.Errorf("create attest manager failed, error: %s", err)
			return
		}
		res, err := m.GetRootCert()
		if err != nil {
			logrus.Errorf("Request get root ca failed: %v", err)
			return
		}
		if res != nil && len(res.Keys) > 0 && len(res.Keys[0].X5c) > 0 {
			pemData, err := manager.X5cToCertPem(res.Keys[0].X5c)
			if err != nil {
				logrus.Errorf(err.Error())
				return
			}
			fileName := fmt.Sprintf("ca%s.pem", res.Keys[0].Kid)
			resultPath := filepath.Join(ouputPath, fileName)
			err = os.WriteFile(resultPath, []byte(pemData.String()), 0644)
			if err != nil {
				logrus.Errorf("Error writing PEM file: %v", err)
				return
			}
			fmt.Println(consts.ColorGreen + "The root ca save in: " + resultPath + consts.OutReset)
		} else {
			logrus.Errorf("No certificate found")
		}

	},
}

func init() {
	//set parameter for policy delete
	Cmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	Cmd.Flags().StringP("output", "o", "./tcas-certs", "optional, the save path of the ca cert")
}
