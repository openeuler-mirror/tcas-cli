/*
 * @Author: jffan
 * @Date: 2024-08-15 09:16:45
 * @LastEditTime: 2024-08-19 17:03:39
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\manager\request.go
 * @Description: The type of request params
 */
package manager

import "gitee.com/nanhu-lab/tcas-cli/teetool"

type SetPolicyReq struct {
	Name            string `json:"policy_name"`
	Policy          string `json:"policy_rego"`
	AttestationType string `json:"attestation_type"`
}

type SetSecretReq struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type UpdateSecretReq struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

type NodeAttestInfoReq struct {
	Report    teetool.NodeEvidence `json:"report"`
	Nonce     string               `json:"nonce"`
	PolicyIds []string             `json:"policy_ids"`
}

type CertCsrInfoReq struct {
	CommonName  string   `json:"common_name"`
	Expiration  int8     `json:"expiration"`
	IPAddresses []string `json:"ip_addresses"`
}

type AttestCertInfoReq struct {
	Csr        *CertCsrInfoReq    `json:"csr"`
	AttestInfo *NodeAttestInfoReq `json:"attestInfo"`
}
