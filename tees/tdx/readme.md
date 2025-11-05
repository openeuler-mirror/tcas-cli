# TDX依赖

首先确定当前的主机是否安装有`libtdx_attest`的库，执行如下命令：

```shell
ldconfig -p | grep libtdx_attest
```

如果运行的结果如下所示，代表当前的主机已经安装有此库：

```shell
root@k8s-cluster2-master:~/jffan/tcas-cli# ldconfig -p | grep libtdx_attest
        libtdx_attest.so.1 (libc6,x86-64) => /lib/x86_64-linux-gnu/libtdx_attest.so.1
        libtdx_attest.so (libc6,x86-64) => /lib/x86_64-linux-gnu/libtdx_attest.so
```

如果没有安装，需要安装`libtdx-attest`库，它分为`libtdx-attest1`以及`libtdx-attest-dev`两个库（以ubuntu系统为例）。

- `libtdx-attest1`：是运行时的库
- `libtdx-attest-dev`：是开发库，依赖于`libtdx-attest1`这个库

```shell
root@k8s-cluster2-master:~/jffan/tcas-cli# apt-cache show libtdx-attest-dev
Package: libtdx-attest-dev
Source: sgx-dcap
Priority: optional
Section: devel
Installed-Size: 35
Maintainer: Hector Cao <hector.cao@canonical.com>
Architecture: amd64
Version: 1.21-0ubuntu3.1
Depends: libtdx-attest1 (= 1.21-0ubuntu3.1)
Filename: pool/main/s/sgx-dcap/libtdx-attest-dev_1.21-0ubuntu3.1_amd64.deb
Size: 7362
MD5sum: 6befd67c03375b2a95dd18c1e085460e
SHA1: bf8fc7082710874654c74445b446e698f17e16e3
SHA256: a632dcbc54a584f79003854a77294ee6b3551fad5537660bb05eae840653511b
Description-en: Intel(R) Trust Domain Extensions Attestation For Developers
 This package contains the developement files for libtdx-attest
Description-md5: 51e46a1b3067b55b36238c7158227b66

root@k8s-cluster2-master:~/jffan/tcas-cli# apt-cache show libtdx-attest1
Package: libtdx-attest1
Source: sgx-dcap
Priority: optional
Section: utils
Installed-Size: 37
Maintainer: Hector Cao <hector.cao@canonical.com>
Architecture: amd64
Version: 1.21-0ubuntu3.1
Depends: libc6 (>= 2.38)
Filename: pool/main/s/sgx-dcap/libtdx-attest1_1.21-0ubuntu3.1_amd64.deb
Size: 10262
MD5sum: b85d520eb5ff6f5a98f22ab3e4a0df4f
SHA1: 8f83bb22144fee9bc686764576046997d9c7e548
SHA256: e3ec47684e2456bb1254ca11075f974c838db1e8ad4e8d5abaf5aeb8a9705826
Description-en: Intel(R) Trust Domain Extensions Attestation library
 This library is used in the guest.
Description-md5: c3596e7e0a8d5ce07c54e624ee742ebe
```

安装方式：

```shell
#安装前需要主机有相关的c库依赖，比如gcc、libc6等
apt-get install -y libtdx-attest1 libtdx-attest-dev
```

