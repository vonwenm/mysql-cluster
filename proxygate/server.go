package proxygate

import (
	"fmt"
	"time"

	proto1 "github.com/golang/protobuf/proto"
	"github.com/nicholaskh/golib/server"
	log "github.com/nicholaskh/log4go"
	. "github.com/nicholaskh/mysql-cluster/config"
	proto "github.com/nicholaskh/mysql-cluster/proto/go"
)

func LaunchServer() {
	s := server.NewTcpServer("myc")

	s.LaunchTcpServer(Config.ListenAddr, newClientHandler(), time.Minute*2, 200)
}

type ClientHandler struct {
	client *server.Client
}

func newClientHandler() *ClientHandler {
	return &ClientHandler{}
}

func (this *ClientHandler) OnAccept(cli *server.Client) {
	this.client = cli
}

func (this *ClientHandler) OnRead(input string) {
	q := &proto.QueryStruct{}
	err := proto1.Unmarshal([]byte(input), q)
	if err != nil {
		log.Error("parse query error")
	}
	log.Info("sql: %s\npool: %s", q.GetSql(), q.GetPool())
	rows, err := proxyGate.Execute(q)
	if err != nil {
		log.Error(err.Error())
	} else {
		this.client.WriteMsg(fmt.Sprintf("%s\n", rows))
	}
}

func (this *ClientHandler) OnClose() {
	this.client.Close()
}

func handleClient(client *server.Client) {
	client.Conn.Write([]byte("connected"))
}
