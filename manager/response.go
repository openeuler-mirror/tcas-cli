/*
 * @Author: jffan
 * @Date: 2024-07-31 15:01:17
 * @LastEditTime: 2024-08-21 10:46:18
 * @LastEditors: jffan
 * @FilePath: \gitee-tcas\manager\response.go
 * @Description: The type of response params
 */
package manager

import "time"

type HttpBaseResponse struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

type PolicySetResponse struct {
	HttpBaseResponse
	PolicyID string `json:"policy_id"`
}

type PolicyDeleteResponse struct {
	HttpBaseResponse
	PolicyID string `json:"policy_id"`
}

type PolicyListResData struct {
	No              int    `json:"no"`
	PolicyId        string `json:"policy_id"`
	PolicyRego      string `json:"policy_rego"`
	PolicyName      string `json:"policy_name"`
	AttestationType string `json:"attestation_type"`
	PolicyHash      string `json:"policy_hash"`
	Version         int    `json:"version"`
	CreateTime      string `json:"createTime"`
	UpdateTime      string `json:"updateTime"`
}

type PolicyListResponse struct {
	HttpBaseResponse
	Data []PolicyListResData `json:"data"`
}

type PolicyListJsonFormat struct {
	Policies []PolicyListResData `json:"policies"`
}

type SecretSetResponse struct {
	HttpBaseResponse
	Id string `json:"id"`
}

type SecretListResData struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type SecretListResponse struct {
	HttpBaseResponse
	Data []SecretListResData `json:"data"`
}

type SecretListJsonFormat struct {
	Secrets []SecretListResData `json:"secrets"`
}

type SecretDeleteResponse struct {
	HttpBaseResponse
	SecretID string `json:"secret_id"`
}
type Nonce struct {
	Nonce       string    `json:"nonce,omitempty"`
	ExpiredTime time.Time `json:"expired,omitempty"`
}
type NonceResponse struct {
	HttpBaseResponse
	Data *Nonce `json:"data"`
}

type TokenResponse struct {
	HttpBaseResponse
	Token string `json:"token"`
}

type CaSingleResponse struct {
	Alg string   `json:"alg"`
	Crv string   `json:"crv"`
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	X   string   `json:"x"`
	X5c []string `json:"x5c"`
	Y   string   `json:"y"`
}

type CaResponse struct {
	Keys []CaSingleResponse `json:"keys"`
}

type AttestCertData struct {
	X5c          []string `json:"x5c"`
	SerialNumber string   `json:"serial_number"`
}

type AttestCertResponse struct {
	HttpBaseResponse
	Data *AttestCertData `json:"data"`
}

type AttestSecretData struct {
	HttpBaseResponse
	Secret interface{} `json:"secret"`
}

type LoginRes struct {
	HttpBaseResponse
	Token  string
	UserId string
}

type User struct {
	No         int        `orm:"auto"`
	Id         string     `orm:"unique;size(255)" json:"userId,omitempty"`
	Display    string     `orm:"size(255)" json:"Display,omitempty"`
	Name       string     `orm:"unique;size(200)" json:"userName,omitempty"`
	Department string     `orm:"size(255)" json:"department,omitempty"`
	Email      string     `orm:"size(200)" json:"email,omitempty"`
	Status     string     `orm:"size(200)" json:"status,omitempty"`
	Comment    string     `orm:"type(text)" json:"comment,omitempty"`
	Roles      string     `orm:"size(32)" json:"roles,omitempty"`
	ApiKey     string     `orm:"size(255);unique" json:"apiKey"`
	LastLogin  *time.Time `orm:"auto_now_add;type(datetime)" json:"lastLogin,omitempty"`
	LastIp     string     `orm:"size(200)" json:"lastIp,omitempty"`
	CreateTime *time.Time `orm:"auto_now_add;type(datetime)" json:"createTime,omitempty"`
	UpdateTime *time.Time `orm:"auto_now;type(datetime)" json:"updateTime,omitempty"`
}

type UserRes struct {
	HttpBaseResponse
	User *User `json:"User"`
}
