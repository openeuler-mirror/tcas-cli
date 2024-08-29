/*
 * @Author: jffan
 * @Date: 2024-08-14 09:33:52
 * @LastEditTime: 2024-08-20 14:35:08
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\verify\cert.go
 * @Description:
 */
package verify

import (
	"crypto/x509"
	"fmt"
	consts "tcas-cli/constants"
	"tcas-cli/manager"
	"tcas-cli/utils/file"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verfifyCertCmd = &cobra.Command{
	Use:   "cert",
	Short: "verify cert",
	Long:  `verify cert`,
	Run: func(cmd *cobra.Command, args []string) {
		var caCert, cert *x509.Certificate
		var err error
		certFilePath, _ := cmd.Flags().GetString("file")
		if !file.IsExists(certFilePath) {
			logrus.Errorf(consts.ColorRed + "The cert file you set does not exist, please check it." + consts.OutReset)
			return
		}
		cert, err = manager.ParseCert(certFilePath)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		caFilePath, _ := cmd.Flags().GetString("ca")
		if caFilePath == "" {
			logrus.Debugf("verify cert ca online")
			url, _ := cmd.Flags().GetString("url")
			logrus.Debugf("verify cert url: " + consts.ColorYellow + url + consts.OutReset)
			m, err := manager.New(url, "")
			if err != nil {
				logrus.Errorf("create attest manager failed, error: %s", err)
				return
			}
			res, err := m.GetRootCert()
			if err != nil {
				logrus.Errorf("Request get ca cert failed: %v", err)
				return
			}
			if res != nil && len(res.Keys) > 0 && len(res.Keys[0].X5c) > 0 {
				pemData, err := manager.X5cToCertPem(res.Keys[0].X5c)
				if err != nil {
					logrus.Errorf(err.Error())
					return
				}
				caCert, err = manager.ParseCert(pemData.Bytes())
				if err != nil {
					logrus.Errorf(err.Error())
					return
				}

			} else {
				logrus.Errorf("No certificate found")
				return
			}
		} else {
			logrus.Debugf("verify cert ca local")
			if !file.IsExists(caFilePath) {
				logrus.Errorf(consts.ColorRed + "Error: The ca file path you set is not exist! Please check it!" + consts.OutReset)
				return
			}
			caCert, err = manager.ParseCert(caFilePath)
			if err != nil {
				logrus.Errorf(err.Error())
				return
			}
		}
		roots := x509.NewCertPool()
		roots.AddCert(caCert)
		opts := x509.VerifyOptions{
			Roots:       roots,
			CurrentTime: time.Now(),
		}
		_, err = cert.Verify(opts)
		if err != nil {
			logrus.Errorf("Failed to verify certificate chain: %v", err)
			return
		}
		fmt.Println(consts.ColorGreen + certFilePath + " verify successful " + consts.OutReset)
	},
}

func init() {
	Cmd.AddCommand(verfifyCertCmd)
	//set parameter for verify cert
	verfifyCertCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	verfifyCertCmd.Flags().StringP("file", "f", "", "must, the path of the cert to be verified")
	verfifyCertCmd.Flags().StringP("ca", "c", "", "optional, the path of the CA certificate file.If not, the CA certificate will be automatically obtained")
	verfifyCertCmd.MarkFlagRequired("file")
}
