package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"strings"
	"time"

	"github.com/lishimeng/go-log"
)

type Session interface {
	Call(serviceMethod string, args any, reply any) error
	Go(serviceMethod string, args any, reply any, done interface{}) interface{}
}

type conn struct {
	network  string
	address  string
	connType connType
}

type connType string

const (
	Base      connType = "Base"
	HTTP      connType = "Http"
	JSON_HTTP connType = "Json-Http"
)

type Message struct {
	Id     int    `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}
type Reply struct {
	Id     int    `json:"id"`
	Result any    `json:"result"`
	Error  string `json:"error"`
}

type BaseRpc struct {
	clientConn conn
	serverConn conn
	client     *rpc.Client
	listener   net.Listener
	ctx        context.Context
}

//TODO 重连

func (r *BaseRpc) Call(serviceMethod string, args any, reply any) (err error) {
	if r.clientConn.connType != JSON_HTTP {
		return r.client.Call(serviceMethod, args, reply)
	}
	return r.callJson(serviceMethod, args, reply)
}

func (r *BaseRpc) callJson(serviceMethod string, args any, replyResult any) (err error) {
	params := make([]any, 0)
	params = append(params, args)
	messageId := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000000) // 随机消息ID
	req := Message{Id: messageId, Method: serviceMethod, Params: params}
	bs, err := json.Marshal(req)
	if err != nil {
		return
	}

	body := strings.NewReader(string(bs))
	resp, err := http.Post(r.clientConn.address, "application/json", body)
	if err != nil {
		return
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var reply Reply
	err = json.Unmarshal(buf, &reply)
	if err != nil {
		return
	}
	if len(reply.Error) != 0 {
		return errors.New(reply.Error)
	}
	if reply.Id != messageId {
		return errors.New("message id not match")
	}

	dataBs, err := json.Marshal(reply.Result)
	if err != nil {
		return
	}
	err = json.Unmarshal(dataBs, replyResult)
	if err != nil {
		return
	}
	return
}

func (r *BaseRpc) Go(serviceMethod string, args any, reply any, done interface{}) interface{} {
	call := r.client.Go(serviceMethod, args, reply, done.(chan *rpc.Call))
	return call
}

func (r *BaseRpc) newClient() (err error) {
	log.Debug("init client, type=%s", string(r.clientConn.connType))
	if r.clientConn.connType == Base {
		r.client, err = rpc.Dial(r.clientConn.network, r.clientConn.address)
	} else if r.clientConn.connType == HTTP {
		r.client, err = rpc.DialHTTP(r.clientConn.network, r.clientConn.address)
	} else {
		r.client = &rpc.Client{}
	}
	log.Info("RPC Client is connected to %s...", r.clientConn.address)
	go func() {
		defer r.client.Close()
		<-r.ctx.Done()
		log.Debug("close rpc client")
	}()
	return
}
