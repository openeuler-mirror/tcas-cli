/*
 * @Author: jffan
 * @Date: 2024-08-15 09:16:45
 * @LastEditTime: 2024-08-20 15:52:56
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\cmd\attest\cert.go
 * @Description:
 */
package attest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"
	consts "tcas-cli/constants"
	"tcas-cli/manager"
	"tcas-cli/utils/file"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// certCmd represents the cert command
var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "attest for getting cert",
	Long:  `attest for getting cert`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		logrus.Debugf("cert called url %s", url)
		ouputPath, _ := cmd.Flags().GetString("output")
		publicKeyPath, _ := cmd.Flags().GetString("key")
		var publicKeyBase64 string
		if publicKeyPath == "" {
			logrus.Debugf(consts.ColorYellow + "don't provide a public key, we will generate an ECC key for you" + consts.OutReset)
			privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				logrus.Errorf("Error generating private key:", err)
				return
			}
			publicKey := &privateKey.PublicKey
			logrus.Debugf("generate an ECC key in binary format：", publicKey)
			publicKeyBase64, err = base64EccPublicKey(publicKey)
			if err != nil {
				logrus.Errorf(err.Error())
				return
			}
			logrus.Debugf("generate an ECC key in base64 format：", publicKeyBase64)
			//save private key
			err = saveEcckeys(ouputPath, privateKey)
			if err != nil {
				logrus.Errorf(err.Error())
				return
			}
			//save public key
			err = saveEcckeys(ouputPath, publicKey)
			if err != nil {
				logrus.Errorf(err.Error())
				return
			}
		} else {
			logrus.Debugf(consts.ColorYellow + "You give a publickey path, we're going to read it for use!" + consts.OutReset)
			pemData, err := os.ReadFile(publicKeyPath)
			if err != nil {
				logrus.Errorf("Error reading public key file:", err)
				return
			}
			publicKeyBase64 = base64.StdEncoding.EncodeToString(pemData)
		}
		logrus.Debugf("the base64 value of ECC public key in pem format：", publicKeyBase64)
		tee, _ := cmd.Flags().GetString("tee")
		policies, _ := cmd.Flags().GetString("policies")
		devices, _ := cmd.Flags().GetString("devices")
		commonName, _ := cmd.Flags().GetString("commonName")
		years, _ := cmd.Flags().GetInt8("expiration")
		ips, _ := cmd.Flags().GetIPSlice("ips")
		ipStrings := make([]string, len(ips))
		for i, ip := range ips {
			ipStrings[i] = ip.String()
		}
		certCsrInfoReq := &manager.CertCsrInfoReq{
			CommonName:  commonName,
			Expiration:  years,
			IPAddresses: ipStrings,
		}
		m, err := manager.New(url, "")
		if err != nil {
			logrus.Errorf("create get cert manager failed, error: %s", err)
			return
		}
		res, err := m.AttestForCert(tee, publicKeyBase64, devices, policies, certCsrInfoReq)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		if res.Data == nil {
			logrus.Errorf("the response data is nil")
			return
		}
		certPemData, err := manager.X5cToCertPem(res.Data.X5c)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		fileName := fmt.Sprintf("cert-%s.pem", res.Data.SerialNumber)
		err = file.EnsureDirExists(ouputPath)
		if err != nil {
			logrus.Errorf("Error ensuring directory exists: %v", err)
			return
		}
		resultPath := filepath.Join(ouputPath, fileName)
		err = os.WriteFile(resultPath, []byte(certPemData.String()), 0644)
		if err != nil {
			logrus.Errorf("Error writing cert PEM file: %v", err)
			return
		}
		fmt.Println(consts.ColorGreen + "The cert successfully saved in: " + resultPath + consts.OutReset)
	},
}

func base64EccPublicKey(publicKey *ecdsa.PublicKey) (string, error) {
	derPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		logrus.Errorf("x509.MarshalPKIXPublicKey, error: %s", err)
		return "", err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPublicKey,
	}
	publicKeyBase64 := base64.StdEncoding.EncodeToString(pem.EncodeToMemory(block))
	return publicKeyBase64, nil
}

func saveEcckeys(outputPath string, key interface{}) error {
	var typeName string
	var pemType string
	var keyBytes []byte
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		typeName = "private"
		pemType = "EC PRIVATE KEY"
		privateKeyBytes, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return fmt.Errorf("error marshaling EC private key: %v", err)
		}
		keyBytes = privateKeyBytes
	case *ecdsa.PublicKey:
		typeName = "public"
		pemType = "PUBLIC KEY"
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(k)
		if err != nil {
			return fmt.Errorf("error marshaling EC public key: %v", err)
		}
		keyBytes = publicKeyBytes
	default:
		return fmt.Errorf("unsupported key type")
	}
	pemBlock := &pem.Block{
		Type:  pemType,
		Bytes: keyBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	fileName := fmt.Sprintf("%s.pem", typeName)
	err := file.EnsureDirExists(outputPath)
	if err != nil {
		return fmt.Errorf("Error ensuring directory exists: %v", err)
	}
	resultPath := filepath.Join(outputPath, fileName)
	err = os.WriteFile(resultPath, pemBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing PEM file: %v", err)
	}
	fmt.Println(consts.ColorGreen + "The " + typeName + " key successfully saved in: " + resultPath + consts.OutReset)
	return nil
}

func init() {
	Cmd.AddCommand(certCmd)
	//set parameter for getting cert
	var ips []net.IP
	certCmd.Flags().StringP("url", "u", "https://api.trustcluster.cc", "optional, tcas's api url")
	certCmd.Flags().StringP("tee", "t", "", "must, tee type")
	certCmd.Flags().StringP("devices", "v", "", "optional, the trust devices")
	certCmd.Flags().StringP("policies", "p", "", "optional, the ids of the policy needed matching")
	certCmd.Flags().StringP("key", "k", "", "optional, the ecc256 of publickey in pem format, if not present, will generate key pair randomly ")

	certCmd.Flags().StringP("commonName", "c", "", "must, the cert's common_name")
	certCmd.Flags().Int8P("expiration", "e", 10, "optional, the cert's expiration time (years)")
	certCmd.Flags().IPSliceP("ips", "i", ips, "optional, the cert's IP addresses extensions")

	certCmd.Flags().StringP("output", "o", "./tcas-certs", "optional, the output dir of the cert and keys")
	certCmd.MarkFlagRequired("tee")
	certCmd.MarkFlagRequired("commonName")
}
