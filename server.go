package tcp_server

import (
	"fmt"
	"os"
	"os/signal"
	"net"
	"io"
)

const (
	conn_type = "tcp"
	unknown_error = "Unknown error occured."
)

/****************************************
		PUBLIC
 ****************************************/
type TcpEventHandler interface {
	OnOpen(Client)
	OnClose(ClientData)
	OnError(ClientData,error)
	OnMessage(interface{},Client)
}

type Message interface {
	Bytes() []byte
}

type MessageReader interface {
	Read(io.Reader) (interface{},error)
}

type Client interface {
	ClientData
	Send(Message)
	Close()
}

type ClientData interface {
	Address() net.Addr
	Id() int
}

func Run(port int, tcpServer TcpEventHandler, msgReader MessageReader) error {
	return run(port,tcpServer,msgReader)
}

/****************************************
		PRIVATE
 ****************************************/
type connection struct {
	addr net.Addr
	conn net.Conn
	done chan<- *connection
	id int
	closed bool
}

func (c *connection) Address() net.Addr {
	return c.addr
}
func (c *connection) Id() int {
	return c.id
}
func (c *connection) Send(msg Message) {
	if c.closed {
		return
	}
	bytes := msg.Bytes()
	out := c.conn
	wrote := 0
	for wrote<len(bytes){
		n,err := out.Write(bytes[wrote:])
		if err!=nil{
			panic(err)
		}
		wrote+=n
	}
}
func (c *connection) Close(){
	if c.closed {
		return
	}
	c.closed = true
	c.conn.Close()
}

func run(port int, tcpServer TcpEventHandler, msgReader MessageReader) error {
	signals := make(chan os.Signal)
	defer close(signals)
	signal.Notify(signals,os.Interrupt)

	addr := fmt.Sprintf(":%d",port)

	listener, err := net.Listen(conn_type,addr)
	if err != nil {
		return err
	}

	lChan := make(chan bool)
	defer func(){
		listener.Close()
		<-lChan
	}()

	newConnChan := listen(listener,lChan)
	doneConnChan := make(chan *connection)

	connCount := 0
	connMap := make(map[int]*connection)


	defer func(){
		close(doneConnChan)
	}()

	for {
		select {
		case c:=<-doneConnChan:
			delete(connMap,c.id)
		case c:=<-newConnChan:
			connCount++
			conn := connection{
				conn : c,
				id : connCount,
				done : doneConnChan,
				addr : c.RemoteAddr(),
			}
			connMap[connCount] = &conn
			go handle(&conn,msgReader,tcpServer)
		case <-signals:
			return nil
		}
	}
}

func handle(conn *connection, msgReader MessageReader, tcpServer TcpEventHandler){
	defer func() {
		tcpServer.OnClose(conn)
		conn.done <- conn
		conn.conn.Close()
	}()
	defer func() {
		e := recover()
		if e==nil {
			return
		}
		if err,ok := e.(error);ok {
			tcpServer.OnError(conn,err)
		} else {
			tcpServer.OnError(conn,fmt.Errorf(unknown_error))
		}
	}()
	tcpServer.OnOpen(conn)
	for {
		m,err := msgReader.Read(conn.conn)
		if err != nil {
			if conn.closed || err == io.EOF {
				return
			}
			panic(err)
		}
		tcpServer.OnMessage(m,conn)
	}
}

func listen(listener net.Listener, done chan<- bool) <-chan net.Conn {
	c := make(chan net.Conn)
	go func(){
		defer func(){
			close(c)
			done <- true
			close(done)
		}()
		for {
			conn,err := listener.Accept()
			if err!=nil {
				return
			}
			c <- conn
		}
	}()
	return c
}