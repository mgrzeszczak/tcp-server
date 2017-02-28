package main

import (
	"github.com/mgrzeszczak/tcp-server"
	"io"
	"log"
	"github.com/mgrzeszczak/tcp-server/example/data"
)

const port = 1234

type eventHandler struct {

}

func (eh *eventHandler) OnOpen(c tcp_server.Client){
	log.Printf("New client @%s\n",c.Address().String())
}
func (eh *eventHandler) OnMessage(msg interface{}, c tcp_server.Client){
	log.Printf("New message from @%s : %s\n",c.Address().String(),msg.(*data.StringMsg).String())
	c.Send(data.New("Pong"))
}
func (eh *eventHandler) OnError(c tcp_server.ClientData, e error) {
	log.Printf("Error from @%s : %s\n",c.Address().String(),e.Error())
}
func (eh *eventHandler) OnClose(c tcp_server.ClientData) {
	log.Printf("Client @%s disconnected\n",c.Address().String())
}
func (eh *eventHandler) Read(r io.Reader) (interface{},error) {
	return data.Read(r)
}

func main() {
	ev := &eventHandler{}
	log.Printf("Starting server @ port: %d\n",port)
	tcp_server.Run(port,ev,ev)
	log.Println("Server closed")
}
