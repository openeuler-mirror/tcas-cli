/*
* @Author: jffan
* @Date: 2024-07-31 15:01:17
 * @LastEditTime: 2024-08-21 15:21:00
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\manager\attestmanager.go
* @Description: Request encapsulation
*/
package manager

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"tcas-cli/collectors"
	consts "tcas-cli/constants"
	"tcas-cli/tees"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	APIEndpoint string
	TlsConfig   *tls.Config
	Collectors  map[string]collectors.EvidenceCollector
}

func New(apiEndpoint, caPath string) (*Manager, error) {
	tc := new(tls.Config)
	if caPath != "" {
		certBytes, err := os.ReadFile(caPath)
		if err != nil {
			return nil, err
		}

		CaCertPool := x509.NewCertPool()
		ok := CaCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			return nil, fmt.Errorf("add ca to pool failed")
		}
		tc.RootCAs = CaCertPool
	} else {
		tc.InsecureSkipVerify = true
	}

	c := tees.GetCollectors()

	return &Manager{
		APIEndpoint: apiEndpoint,
		TlsConfig:   tc,
		Collectors:  c,
	}, nil
}

func (m *Manager) newClient(method string, url string) *httplib.BeegoHTTPRequest {
	var client *httplib.BeegoHTTPRequest
	me := strings.ToUpper(method)
	client = httplib.NewBeegoRequest(m.APIEndpoint+url, me)

	if m.TlsConfig != nil {
		client.SetTLSClientConfig(m.TlsConfig)
	}
	return client
}

func (m *Manager) SetPolicy(name, policy, attestationType string) (*PolicySetResponse, error) {

	if name == "" || policy == "" {
		return nil, fmt.Errorf("name or policy is null")
	}
	client := m.newClient("post", PolicyUrl)
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	req := SetPolicyReq{
		Name:            name,
		Policy:          policy,
		AttestationType: attestationType,
	}

	client, err := client.JSONBody(req)
	if err != nil {
		return nil, err
	}

	res := new(PolicySetResponse)
	err = client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request set policy api failed, error: %s ", err)
	}

	return res, nil
}
func (m *Manager) ListPolicy(attestationType string) (*PolicyListResponse, error) {
	if attestationType == "" {
		attestationType = "trust_node"
		fmt.Println("attestationType is null, use default value: `trust_node`")
	}
	client := m.newClient("get", PolicyUrl)
	client.Param("attestation", attestationType)
	res := new(PolicyListResponse)
	err := client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request policy list failed, error: %s ", err)
	}
	return res, nil
}
func (m *Manager) DeletePolicy(policyID string) (*PolicyDeleteResponse, error) {
	if policyID == "" {
		return nil, fmt.Errorf("policyID is null")
	}
	deleteURL := PolicyUrl + "/" + policyID
	client := m.newClient("delete", deleteURL)
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	res := new(PolicyDeleteResponse)

	err := client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request delete policy api failed,error: %s ", err)
	}

	return res, nil
}
func (m *Manager) SetSecret(name, encodeJsonData string) (*SecretSetResponse, error) {
	if name == "" || encodeJsonData == "" {
		return nil, fmt.Errorf("secret name or jsonData is null")
	}
	client := m.newClient("post", SecretUrl)
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	req := SetSecretReq{
		Name:   name,
		Secret: encodeJsonData,
	}

	client, err := client.JSONBody(req)
	if err != nil {
		return nil, err
	}

	res := new(SecretSetResponse)
	err = client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request set secret api failed, error: %s ", err)
	}

	return res, nil
}
func (m *Manager) UpdateSecret(id, encodeJsonData string) (*SecretSetResponse, error) {
	if id == "" || encodeJsonData == "" {
		return nil, fmt.Errorf("secret id or jsonData is null")
	}
	client := m.newClient("put", SecretUrl)
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	req := UpdateSecretReq{
		Id:     id,
		Secret: encodeJsonData,
	}
	client, err := client.JSONBody(req)
	if err != nil {
		return nil, err
	}

	res := new(SecretSetResponse)
	err = client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request update secret api failed, error: %s ", err)
	}

	return res, nil
}
func (m *Manager) ListSecret() (*SecretListResponse, error) {
	client := m.newClient("get", SecretListUrl)
	res := new(SecretListResponse)
	err := client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request secret list failed, error: %s ", err)
	}
	return res, nil
}

func (m *Manager) DeleteSecret(secretID string) (*SecretDeleteResponse, error) {
	deleteSecretURL := SecretUrl + "/" + secretID
	client := m.newClient("delete", deleteSecretURL)
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	res := new(SecretDeleteResponse)

	err := client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request delete secret api failed,error: %s ", err)
	}

	return res, nil
}

func (m *Manager) GetRootCert() (*CaResponse, error) {
	client := m.newClient("get", CaUrl)
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	res := new(CaResponse)
	err := client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request ca cert failed, error: %s", err)
	}

	return res, nil
}

func (m *Manager) GetNonce() (*NonceResponse, error) {
	client := m.newClient("get", NonceUrl)
	res := new(NonceResponse)
	err := client.ToJSON(res)
	if err != nil {
		return nil, fmt.Errorf("request get nonce failed, error: %s", err)
	}

	return res, nil

}

func (m *Manager) getNodeAttestInfo(tee, runtimedata, devices, policies string) (*NodeAttestInfoReq, error) {
	if tee == "" {
		return nil, fmt.Errorf("tee is null")
	}

	nonceRes, err := m.GetNonce()
	if err != nil {
		return nil, fmt.Errorf("get nonce failed, error: %s", err)
	}
	if nonceRes.Code != 200 {
		return nil, fmt.Errorf("get nonce failed, error: %s", nonceRes.Message)
	}

	logrus.Debugf("nonce is %s", nonceRes.Data.Nonce)
	userdata := fmt.Sprintf("%x", sha256.Sum256([]byte(nonceRes.Data.Nonce+runtimedata)))

	c, ok := m.Collectors[tee]
	if !ok {
		return nil, fmt.Errorf("tee: %s not support yet", tee)
	}
	teeReport, err := c.CollectEvidence([]byte(userdata))
	if err != nil {
		return nil, fmt.Errorf("collect evidence from %s failed, error: %s", c.Name(), err)
	}

	var devicelist []string
	if devices != "" {
		devicelist = strings.Split(devices, ",")
	}
	deviceReports := make([]*TrustDeviceReport, 0)
	for _, d := range devicelist {
		c, ok := m.Collectors[d]
		if !ok {
			return nil, fmt.Errorf("device tee: %s not support yet", tee)
		}
		deviceEvidence, err := c.CollectEvidence([]byte(userdata))
		if err != nil {
			return nil, fmt.Errorf("collect evidence from %s failed, error: %s", c.Name(), err)
		}

		deviceReports = append(deviceReports, &TrustDeviceReport{
			DeviceType:   d,
			DeviceReport: deviceEvidence.TeeReport,
			Parameter:    deviceEvidence.Parameter,
		})
	}

	var policyIds []string
	if policies != "" {
		policyIds = strings.Split(policies, ",")
	}
	return &NodeAttestInfoReq{
		Report: &NodeEvidence{
			Tee:         teeReport.TeeType,
			TeeReport:   teeReport.TeeReport,
			Parameter:   teeReport.Parameter,
			RuntimeData: runtimedata,
			TrustDevice: deviceReports,
		},
		Nonce:     nonceRes.Data.Nonce,
		PolicyIds: policyIds,
	}, nil
}

func (m *Manager) AttestForToken(tee, runtimedata, devices, policies string) (*TokenResponse, error) {
	attestReq, err := m.getNodeAttestInfo(tee, runtimedata, devices, policies)
	if err != nil {
		return nil, fmt.Errorf("get node attestInfo failed, error: %s", err)
	}

	client := m.newClient("post", AttestUrl)
	client, err = client.JSONBody(attestReq)
	if err != nil {
		return nil, fmt.Errorf("set request body failed, error: %s", err)
	}

	tokenRes := new(TokenResponse)
	err = client.ToJSON(tokenRes)
	if err != nil {
		return nil, fmt.Errorf("do request to attest api failed, error: %s", err)
	}

	if tokenRes.Code != 200 {
		return nil, fmt.Errorf("response error: %s", tokenRes.Message)
	}

	return tokenRes, nil
}

func (m *Manager) AttestForCert(tee, eccpemBase64key, devices, policies string, csr *CertCsrInfoReq) (*AttestCertResponse, error) {
	attestReq, err := m.getNodeAttestInfo(tee, eccpemBase64key, devices, policies)
	if err != nil {
		return nil, fmt.Errorf("get node cert attestInfo failed, error: %s", err)
	}

	client := m.newClient("post", AttestCertUrl)
	req := AttestCertInfoReq{
		Csr:        csr,
		AttestInfo: attestReq,
	}
	client, err = client.JSONBody(req)
	if err != nil {
		return nil, fmt.Errorf("set request body failed, error: %s", err)
	}

	certRes := new(AttestCertResponse)
	err = client.ToJSON(certRes)
	if err != nil {
		return nil, fmt.Errorf("do request to attest cert api failed, error: %s", err)
	}

	if certRes.Code != 200 {
		return nil, fmt.Errorf("attest cert response error: %s", certRes.Message)
	}

	return certRes, nil
}

func (m *Manager) AttestForSecret(tee, runtimedata, devices, policies, secretID string) (*AttestSecretData, error) {
	if secretID == "" {
		return nil, fmt.Errorf("secret id is null")
	}
	attestReq, err := m.getNodeAttestInfo(tee, runtimedata, devices, policies)
	if err != nil {
		return nil, fmt.Errorf("get node cert attestInfo failed, error: %s", err)
	}
	client := m.newClient("post", AttestSecretUrl)
	client.Header("SecretId", secretID)
	client, err = client.JSONBody(attestReq)
	if err != nil {
		return nil, fmt.Errorf("set request body failed, error: %s", err)
	}
	secretRes := new(AttestSecretData)
	err = client.ToJSON(secretRes)
	if err != nil {
		return nil, fmt.Errorf("do request to attest secert api failed, error: %s", err)
	}

	if secretRes.Code != 200 {
		return nil, fmt.Errorf("attest secert response error: %s", secretRes.Message)
	}

	return secretRes, nil
}

func X5cToCertPem(x5c []string) (*bytes.Buffer, error) {
	pemData := new(bytes.Buffer)
	if x5c != nil && len(x5c) > 0 {
		for _, x5c := range x5c {
			certBytes, err := base64.StdEncoding.DecodeString(x5c)
			if err != nil {
				return pemData, fmt.Errorf("failed to decode base64 certificate: %s\n", err)
			}
			block := &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(certBytes),
			}
			tempData := pem.EncodeToMemory(block)
			pemData.Write(tempData)
		}
		return pemData, nil
	}
	return pemData, fmt.Errorf("x5c is null")
}

func ParseTokenByPk(publicKey any, tokenString string) (*jwt.Token, error) {
	logrus.Debugf("ca publicKey: %v", publicKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Token validation failed: %v", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			now := time.Now().Unix()
			if int64(exp) < now {
				return nil, fmt.Errorf("Token expired")
			}
		} else {
			return nil, fmt.Errorf("Expiration time claim 'exp' missing")
		}
	} else {
		return nil, fmt.Errorf("Invalid token")
	}
	return token, nil
}

func PrintFormatToken(token *jwt.Token) error {
	jsonHeaderData, err := json.MarshalIndent(token.Header, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshal token header failed: %v", err)
	}
	fmt.Println("------------------Token Info Start------------------")
	fmt.Println(consts.ColorYellow + string(jsonHeaderData) + consts.OutReset)
	jsonClaimsData, err := json.MarshalIndent(token.Claims, "", "  ")
	if err != nil {
		fmt.Println("------------------Token Info End--------------------")
		return fmt.Errorf("Marshal token Claims failed: %v", err)
	}
	fmt.Println(consts.ColorYellow + string(jsonClaimsData) + consts.OutReset)
	fmt.Println("------------------Token Info End--------------------")
	return nil
}

func ParseCert(certData interface{}) (*x509.Certificate, error) {
	var certcontent []byte
	var err error
	resultCert := new(x509.Certificate)
	switch v := certData.(type) {
	case string:
		certcontent, err = os.ReadFile(v)
		if err != nil {
			return resultCert, fmt.Errorf("read cert file failed, error: %s", err)
		}
	case []byte:
		certcontent = v
	default:
		return resultCert, fmt.Errorf("unsupported type: %T", v)
	}
	certBlock, _ := pem.Decode(certcontent)
	logrus.Debugf("certType: %s", certBlock.Type)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return resultCert, fmt.Errorf("Failed to decode PEM block containing certificate")
	}
	resultCert, err = x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return resultCert, fmt.Errorf("Failed to parse cert certificate: %v", err)
	}
	return resultCert, nil
}
