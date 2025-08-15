<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/joho/godotenv](#githubcomjohogodotenv)
  - [.env](#env)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第三方使用-->flannel](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8--flannel)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/joho/godotenv

godotenv库从.env文件中读取配置， 然后存储到程序的环境变量中。



## .env

在软件开发中，环境变量是一种常见的配置选项，它们通常用于存储与环境相关的信息，如数据库连接字符串、API密钥等。
.env文件是一种专门用于存储这类环境变量的简单配置文件。


.env文件是一个纯文本文件，每一行都是一个键值对，用等号=分隔。例如：


```env
DB_HOST=localhost
DB_PORT=3306
API_KEY=123456
```


## 源码分析


加载配置
```go
// github.com/joho/godotenv@v1.5.1/godotenv.go

func Load(filenames ...string) (err error) {
	// 如果传入为空,默认加载 .env 文件
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}


func loadFile(filename string, overload bool) error {
	// 读取文件
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			// 设置环境变量
			_ = os.Setenv(key, value)
		}
	}

	return nil
}
```

## 第三方使用-->flannel

```go
func recycleIPTables(nw ip.IP4Net, lease *subnet.Lease) error {
	prevNetwork := ReadCIDRFromSubnetFile(opts.subnetFile, "FLANNEL_NETWORK")
	prevSubnet := ReadCIDRFromSubnetFile(opts.subnetFile, "FLANNEL_SUBNET")
	// recycle iptables rules only when network configured or subnet leased is not equal to current one.
	if prevNetwork != nw && prevSubnet != lease.Subnet {
		log.Infof("Current network or subnet (%v, %v) is not equal to previous one (%v, %v), trying to recycle old iptables rules", nw, lease.Subnet, prevNetwork, prevSubnet)
		lease := &subnet.Lease{
			Subnet: prevSubnet,
		}
		if err := network.DeleteIP4Tables(network.MasqRules(prevNetwork, lease)); err != nil {
			return err
		}
	}
	return nil
}


// 从文件中读取 cidr 信息
func ReadCIDRFromSubnetFile(path string, CIDRKey string) ip.IP4Net {
	var prevCIDR ip.IP4Net
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		prevSubnetVals, err := godotenv.Read(path)
		if err != nil {
			log.Errorf("Couldn't fetch previous %s from subnet file at %s: %s", CIDRKey, path, err)
		} else if prevCIDRString, ok := prevSubnetVals[CIDRKey]; ok {
			err = prevCIDR.UnmarshalJSON([]byte(prevCIDRString))
			if err != nil {
				log.Errorf("Couldn't parse previous %s from subnet file at %s: %s", CIDRKey, path, err)
			}
		}
	}
	return prevCIDR
}

```


## 参考
- [Go每日一库之11：godotenv](https://darjun.github.io/2020/02/12/godailylib/godotenv/)