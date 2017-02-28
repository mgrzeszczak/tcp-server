package main

import (
	"net"
	"log"
	"github.com/mgrzeszczak/tcp-server/example/data"
)

const (
	conn_type = "tcp"
	host = "localhost:1234"
)

func main(){
	conn, err := net.Dial(conn_type, host)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		log.Printf("Disconnecting %v", host)
		conn.Close()
	}()
	payload := data.New("Ping").Bytes()
	wrote := 0
	for wrote < len(payload){
		n,err := conn.Write(payload[wrote:])
		if err!=nil{
			panic(err)
		}
		wrote+=n
	}
	msg,err := data.Read(conn)
	if err!=nil {
		panic(err)
	}
	log.Printf("Received msg: %s\n",msg.String())
}
