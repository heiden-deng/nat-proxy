package main

import (
	"net-proxy/common"
	"os"
	"net"
	"log"
)

var preConn *net.TCPConn

func handleAAccept(listen net.Listener, appChan chan *net.TCPConn)  {
	for  {
		conn,err := listen.Accept()
		if err != nil {
			log.Println("accept app conn err",err)
			continue
		}
		log.Println("accept conn from app side")
		if preConn != nil {
			preConn.Close()
		}
		preConn = conn.(*net.TCPConn)
		appChan <- preConn
	}
}


func main(){
	if len(os.Args) < 2 {
		log.Fatalln("Usage:" + os.Args[0] + " conf_path, missing conf_path param ")
		os.Exit(1)
	}
	cfPath := os.Args[1]
	cf:= common.InitConfig(cfPath)
	appL,err := net.Listen("tcp",cf["app_port"])
	if err != nil {
		log.Println("error listen", err)
		return
	}
	defer appL.Close()
	log.Println(" app side listen ok")
	appAcceptChan := make(chan *net.TCPConn)
	go handleAAccept(appL, appAcceptChan)

	clientL,err := net.Listen("tcp",cf["client_port"])
	if err != nil {
		log.Println("error client listen", err)
		return
	}
	defer clientL.Close()
	log.Println(" client side listen ok")

	//

	for{
		log.Println("Wait app server connect")
		aConn := <- appAcceptChan
		log.Println("Accept connect from app service ,Waiting client connect ..")
		cConn,err := clientL.Accept()
		if err != nil {
			log.Println("accept client connect err",err)
			break
		}
		log.Println("accept connect from client side")
		exitChan := make(chan bool)
		//从客户端读取转发给APP端
		go common.Transfer(cConn.(*net.TCPConn), aConn, exitChan,0)

		//从app端读取转发给客户端
		go common.Transfer(aConn, cConn.(*net.TCPConn), exitChan, 1)
		log.Println("wait client disconnect")
		_ = <- exitChan
		if aConn  != nil{
			aConn.Close()
		}
		if cConn != nil {
			cConn.Close()
		}
		log.Println("Client has disconnect...")
	}

}


