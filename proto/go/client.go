package proto

import (
	"net"
	"time"

	proto1 "github.com/golang/protobuf/proto"
	log "github.com/nicholaskh/log4go"
)

const (
	RETRY_CNT = 3
)

type Client struct {
	serverAddr string
	net.Conn
	readTimeout time.Duration
}

func NewClient(readTimeout time.Duration) *Client {
	c := new(Client)
	c.readTimeout = readTimeout
	return c
}

func (this *Client) Dial(addr string) error {
	this.serverAddr = addr
	return this.connect()
}

func (this *Client) Query(pool, sql string, args []string) (result string, err error) {
	q := &QueryStruct{Pool: &pool, Sql: &sql, Args: args}

	buf, err := proto1.Marshal(q)
	if err != nil {
		log.Error(err)
		return "", err
	}

	writeSucc := false

	_, err = this.Write(buf)
	if err == nil {
		writeSucc = true
	} else {
		//retry
		for i := 1; i < RETRY_CNT; i++ {
			err = this.connect()
			if err != nil {

			} else {
				_, err := this.Write(buf)
				if err == nil {
					writeSucc = true
					break
				}
			}
		}
	}

	if writeSucc {
		rBuff := make([]byte, 1000)
		this.SetReadDeadline(time.Now().Add(this.readTimeout))
		n, err := this.Read(rBuff)
		if err != nil {
			log.Error("read from server error: %s", err.Error())
		} else {
			result = string(rBuff[:n])
		}
	}
	return
}

func (this *Client) connect() (err error) {
	this.Conn, err = net.Dial("tcp", this.serverAddr)
	if err != nil {
		log.Error("dial server error: %s", err.Error())
	}

	return
}
