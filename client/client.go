package main

import (
	"os"
	"nat-proxy/common"
	"log"
	"net"
	"fmt"
)

func main()  {
	if len(os.Args) < 2 {
		log.Fatalln("Usage:" + os.Args[0] + " conf_path, missing conf_path param ")
		os.Exit(1)
	}
	cfPath := os.Args[1]
	cf := common.InitConfig(cfPath)
	serverAddr,err := net.ResolveTCPAddr("tcp4", cf["server_addr"])
	checkError(err)
	sConn,err := net.DialTCP("tcp",nil, serverAddr)
	checkError(err)
	log.Println("connect to server success,",cf["server_addr"])

	localAddr,err := net.ResolveTCPAddr("tcp4", cf["local_addr"])
	checkError(err)
	aConn,err := net.DialTCP("tcp",nil, localAddr)
	checkError(err)
	log.Println("connect to service success,",cf["local_addr"])
	exitChan := make(chan bool)
	for{

		go common.Transfer(sConn, aConn, exitChan, 3)
		go common.Transfer(aConn, sConn, exitChan, 3)

		log.Println("start transfer data between remote server and local service")
		_ = <-exitChan
		log.Println("client start exit..,close")
		sConn.Close()
		aConn.Close()
		log.Println("reconnect to server and local service ...")
		sConn,err = net.DialTCP("tcp",nil, serverAddr)
		checkError(err)
		log.Println("connect to server success,",cf["server_addr"])

		localAddr,err := net.ResolveTCPAddr("tcp4", cf["local_addr"])
		checkError(err)
		aConn,err = net.DialTCP("tcp",nil, localAddr)
		checkError(err)
		log.Println("connect to service success,",cf["local_addr"])
	}

}



func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}