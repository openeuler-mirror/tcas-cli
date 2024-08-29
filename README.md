# TCAS Client

## 1. policy manager
### 1.1 set policy 
```shell
./tcasctl policy set -u <url> -n <name> -f <rego file path> -t <attestation-type>
```
+ `-u`: optional, tcas's api url, default is https://api.trustcluster.cc
+ `-n`: must, policy name 
+ `-f`: must, the path of policy file in rego format
+ `-t`: optional, the attestation-type of policy, support trustnode or trustcluster, default is trustnode

successful response
```shell
set policy successful, policy id: cfaaab6d-7a25-436e-a8d9-6357a1e4cb33
```

### 1.2 get policy list
```shell
./tcasctl policy list -u <url> -t <attestation-type>
```
+ `-u`: optional, tcas's api url, default is https://api.trustcluster.cc
+ `-t`: optional, the attestation-type of policy, support trustnode or trustcluster, default is trustnode
successful response: 

```json
{
    "policies": [
        {
            "no": 2,
            "policy_id": "9e434346-682d-4c13-917d-24883ce096d1",
            "policy_rego": "cGFja2FnZSB0cnVzdG5vZGUKCmltcG9ydCBmdXR1cmUua2V5d29yZHMuZXZlcnkKCmRlZmF1bHQgdHJ1c3Rfbm9kZSA6PSBmYWxzZQoKdHJ1c3Rfbm9kZSB7CglhbGxvd2VkX25vZGUoaW5wdXQudHJ1c3Rfbm9kZSkKfQoKYWxsb3dlZF9ub2RlKG5vZGUpIHsKCW5vZGUudGVlLnZpcnRjY2FfcmltID09ICJlYTIxY2NlMWJiODM2Y2E3OTc4NGFjMTFhMmZlMjg4YzY3MWU1ZjYyYzI4ODM3MThlZDhkYTU5OTQ3YWUyOGIxIgoJbm9kZS50ZWUudmlydGNjYV9yZW0wID0gIjM3YTRhMDRmZTRjYjIxYTgwYTgxNDE5ZWI0Zjc1OTJmNzI3MTA2OTcyN2ZiNWViODU5ZmQxYjUwMzE5YmZhZTQiCn0K",
            "policy_name": "test-vcca1",
            "attestation_type": "trust_node",
            "policy_hash": "822e3ca44a68bd28797d61d6428dff04691c1811093d3cd589febd2b00309842",
            "version": 1,
            "createTime": "2024-07-24T09:17:22.4220796Z",
            "updateTime": "2024-07-24T09:17:22.42208149Z"
        },
        {
            "no": 1,
            "policy_id": "4ed9690b-962f-4279-abdb-fdccecba6775",
            "policy_rego": "cGFja2FnZSB0cnVzdG5vZGUKCmltcG9ydCBmdXR1cmUua2V5d29yZHMuZXZlcnkKCmRlZmF1bHQgdHJ1c3Rfbm9kZSA6PSBmYWxzZQoKdHJ1c3Rfbm9kZSB7CglhbGxvd2VkX25vZGUoaW5wdXQudHJ1c3Rfbm9kZSkKfQoKYWxsb3dlZF9ub2RlKG5vZGUpIHsKCW5vZGUudGVlLnZpcnRjY2FfcmltID09ICJlYTIxY2NlMWJiODM2Y2E3OTc4NGFjMTFhMmZlMjg4YzY3MWU1ZjYyYzI4ODM3MThlZDhkYTU5OTQ3YWUyOGIxIgoJbm9kZS50ZWUudmlydGNjYV9yZW0wID0gIjM3YTRhMDRmZTRjYjIxYTgwYTgxNDE5ZWI0Zjc1OTJmNzI3MTA2OTcyN2ZiNWViODU5ZmQxYjUwMzE5YmZhZTQiCn0K",
            "policy_name": "test-vcca",
            "attestation_type": "trust_node",
            "policy_hash": "822e3ca44a68bd28797d61d6428dff04691c1811093d3cd589febd2b00309842",
            "version": 1,
            "createTime": "2024-07-24T09:17:11.51031935Z",
            "updateTime": "2024-07-24T09:17:11.5103208Z"
        }
    ]
}
```

### 1.3 get the detail of the policy  (unsupported now)
```shell
./tcasctl policy detail -u <url> -i <policy_id>
```

+ `-i`: must the id of policy

successful response : 
```json
        {
            "no": 1,
            "policy_id": "4ed9690b-962f-4279-abdb-fdccecba6775",
            "policy_rego": "cGFja2FnZSB0cnVzdG5vZGUKCmltcG9ydCBmdXR1cmUua2V5d29yZHMuZXZlcnkKCmRlZmF1bHQgdHJ1c3Rfbm9kZSA6PSBmYWxzZQoKdHJ1c3Rfbm9kZSB7CglhbGxvd2VkX25vZGUoaW5wdXQudHJ1c3Rfbm9kZSkKfQoKYWxsb3dlZF9ub2RlKG5vZGUpIHsKCW5vZGUudGVlLnZpcnRjY2FfcmltID09ICJlYTIxY2NlMWJiODM2Y2E3OTc4NGFjMTFhMmZlMjg4YzY3MWU1ZjYyYzI4ODM3MThlZDhkYTU5OTQ3YWUyOGIxIgoJbm9kZS50ZWUudmlydGNjYV9yZW0wID0gIjM3YTRhMDRmZTRjYjIxYTgwYTgxNDE5ZWI0Zjc1OTJmNzI3MTA2OTcyN2ZiNWViODU5ZmQxYjUwMzE5YmZhZTQiCn0K",
            "policy_name": "test-vcca",
            "attestation_type": "trust_node",
            "policy_hash": "822e3ca44a68bd28797d61d6428dff04691c1811093d3cd589febd2b00309842",
            "version": 1,
            "createTime": "2024-07-24T09:17:11.51031935Z",
            "updateTime": "2024-07-24T09:17:11.5103208Z"
        }
```

### 1.4 delete policy
```shell
./tcasctl policy delete -u <url> -i <policy_id>
```
+ `-i`: must the id of policy

successful response:
```shell
delete policy successful, the policy id is <policy_id>
```

## 2. secret manager 
### 2.1 set secret 
```shell
./tcasctl secret set -u <url> -n <name> -f <secret file> 
```
+ `-f`: must, the path of secret file, only support json format.
+ `-n`: must, the unique name of the secret 

successful response: 
```shell
set secret successful, secret id: <secret_id>
```

### 2.2 update secret 
```shell
./tcasctl secret update -u <url> -f <new secret file> -i <secret id> 
```

+ `-f`: must, the path of new secret file, only support json format.
+ `-i`: must, the id of the old secret

+ successful response:
```shell
update secret successful, secret id: <secret_id>
```


### 2.3 get the secret base info list
```shell
./tcasctl secret list -u <url> 
```

successful response:
```json
{
      "secrets": [
        {
            "id": "0d6f1080-dcf4-4961-8def-9fb1f98d6174",
            "name": "test-vcca",
            "createTime": "2024-07-24T08:41:06.62906502Z",
            "updateTime": "2024-07-24T08:41:06.62906712Z"
        }
    ]
}
```

### 2.4 delete secret 
```shell
./tcasctl secret detele -u <url> -i <secret_id>
```
+ `-i`: must, the id of the old secret

successful response:
```shell
delete secret successful, secret id: <secret_id>
```

## 3 Attest
### 3.1 get token 
```shell
./tcasctl attest token -u <url>  -t <tee> -d <base64 encoded userdata> -p <policy ids> -v <trust-devices>
```
+ `-t`: must, the type of tee, now support csv or virtcca
+ `-d`: optional, the base64 encoded userdata
+ `-p`: optional, the ids of the policy needed matching
+ `-v`: optional, the trust devices

eg: ./tcasctl attest token -u http://127.0.0.1:8081 -t virtcca -d MTIzYWJj 

successful response:
```shell
<token>
```

### 3.2 get secret 
```shell
./tcasctl attest secret -u <url>  -t <tee> -d <base64 encoded userdata> -p <policy ids> -v <trust-devices> -s <secretid> -o <output dir> 
```
+ `-t`: must, the type of tee, now support csv or virtcca
+ `-s`: must,the secret ID that needs to be obtained
+ `-d`: optional, the base64 encoded userdata
+ `-p`: optional, the ids of the policy needed matching 
+ `-v`: optional, the trust devices
+ `-o`:optional, the output dir of the secret,default is ./tcas-secret

successful response:
```shell
{"key":"123"}
```

### 3.3 get cert  
```shell 
./tcasctl attest cert -u <url> -t <tee> -p <policy ids> -v <trust-devices> -k <publickey file> -c <common_name> -e <expiration> -i <ipaddresses> -o <output dir> 
```
+ `-t`: must, the type of tee, now support csv or virtcca
+ `-p`: optional, the ids of the policy needed matching
+ `-v`: optional, the trust devices
+ `-k`: optional, the ecc256 of publickey in pem format, if not present, will generate key pair randomly 
+ `-c`: must, the cert's common_name 
+ `-e`: optional, the cert's expiration time, default: 10 years
+ `-i`: optional, the cert's IP addresses extensions
+ `-o`：optional, the output dir of the cert and keys, default: ./tcas-certs

successful response:
```shell
get cert successful, the <publickey.pem> <privatekey.pem> <serial_number_common_name.crt> save in <output dir>
```

## 4.verify
### 4.1 get root CA
```shell
./tcasctl ca -u <url> -o <path> 
```
`-o`: optional, the save path of the ca cert 

### 4.2 verify token
```shell
./tcasclt verify token -t <token> -f <the path of the cert>
```
- `-t`:must,tcas's token
- `-f`:optionally, the path of the CA certificate. If not, it will be verified online

successful response

```shell
verify token successful, the detail info of the token is as follow:
<token claims>
```

### 4.3 verify cert
```shell
./tcasectl verfiy cert -u <url> -f <the path of the cert> -c <the path of the ca cert>
```
`-f`: must, the path of the cert to be verified

`-c`: optional, the path of the CA certificate file.If not, the CA certificate will be automatically obtained

successful response:
```shell
<the path of the cert> verify successful 
```

## reference
+ https://github.com/edgelesssys/go-tdx-qpl/blob/main/tdx/tdx_linux.go