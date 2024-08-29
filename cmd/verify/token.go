/*
 * @Author: jffan
 * @Date: 2024-08-09 09:29:01
 * @LastEditTime: 2024-08-20 11:02:55
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\verify\token.go
 * @Description: The CA certificate is used to verify the token, which is divided into online verification and local verification
 */
package verify

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	consts "tcas-cli/constants"
	"tcas-cli/manager"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TokenHeader struct {
	Alg string `json:"alg"`
	Jku string `json:"jku"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

var verifyTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "verify token",
	Long:  `verify token`,
	Run: func(cmd *cobra.Command, args []string) {
		var caCert *x509.Certificate
		var err error
		tokenString, _ := cmd.Flags().GetString("token")
		caFilePath, _ := cmd.Flags().GetString("file")
		if caFilePath == "" {
			//online verify
			logrus.Debugf(consts.ColorYellow + "Verify Token Online" + consts.OutReset)
			url, _ := cmd.Flags().GetString("url")
			m, err := manager.New(url, "")
			if err != nil {
				logrus.Errorf("create attest manager failed, error: %s", err)
				return
			}
			res, err := m.GetRootCert()
			if err != nil {
				logrus.Errorf("Request policy list failed: %v", err)
				return
			}
			// parse token header
			parts := strings.Split(tokenString, ".")
			header, err := base64.RawURLEncoding.DecodeString(parts[0])
			if err != nil {
				logrus.Errorf(consts.ColorRed+"Decoding token header failed:%v"+consts.OutReset, err)
				return
			}
			var headerToken TokenHeader
			if err = json.Unmarshal(header, &headerToken); err != nil {
				logrus.Errorf(consts.ColorRed+"Unmarshalling header failed:%v"+consts.OutReset, err)
				return
			}
			for _, key := range res.Keys {
				if key.Kid == headerToken.Kid {
					pemData, err := manager.X5cToCertPem(key.X5c)
					if err != nil {
						logrus.Errorf(err.Error())
						return
					}
					caCert, err = manager.ParseCert(pemData.Bytes())
					if err != nil {
						logrus.Errorf(err.Error())
						return
					}
				}
			}
		} else {
			//local verify
			logrus.Debugf(consts.ColorYellow + "Verify Token Local" + consts.OutReset)
			caCert, err = manager.ParseCert(caFilePath)
			if err != nil {
				logrus.Errorf(err.Error())
				return
			}
		}
		publicKey := caCert.PublicKey
		token, err := manager.ParseTokenByPk(publicKey, tokenString)
		if err != nil {
			fmt.Println(consts.ColorRed + "Verify Token Failed" + consts.OutReset)
			logrus.Errorf(err.Error())
			return
		}
		fmt.Println(consts.ColorGreen + "Verify token successful" + consts.OutReset)
		err = manager.PrintFormatToken(token)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
	},
}

func init() {
	Cmd.AddCommand(verifyTokenCmd)
	//set parameter for secret set
	verifyTokenCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	verifyTokenCmd.Flags().StringP("token", "t", "", "must, tcas's token")
	verifyTokenCmd.Flags().StringP("file", "f", "", "optionally, the path of the CA certificate. If not, it will be verified online")
	verifyTokenCmd.MarkFlagRequired("token")
}
