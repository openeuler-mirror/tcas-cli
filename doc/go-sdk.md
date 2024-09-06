# Go-SDK
## 0 准备
安装Go 1.19以上版本，并获取go module
```shell
go get gitee.com/nanhu-lab/tcas-cli/manager
```

## 1.基础调用
获取nonce
```golang
package main

import (
	"fmt"
	"gitee.com/nanhu-lab/tcas-cli/manager"
)

func main() {
	m, _ := manager.New("https://api.trustcluster.cc", "", nil)
	nonce, _ := m.GetNonce()
	fmt.Println(nonce.Data)
}

```