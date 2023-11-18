package rpc

import (
	"context"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/lishimeng/go-log"
)

type Session interface {
	Call(serviceMethod string, args any, reply any) error
	Go(serviceMethod string, args any, reply any, done interface{}) interface{}
}

type conn struct {
	network string
	address string
}
type BaseRpc struct {
	clientConn conn
	serverConn conn
	client     *rpc.Client
	listener   net.Listener
	ctx        context.Context
}

func (r *BaseRpc) Call(serviceMethod string, args any, reply any) (err error) {
	return r.client.Call(serviceMethod, args, reply)
}

func (r *BaseRpc) Go(serviceMethod string, args any, reply any, done interface{}) interface{} {
	call := r.client.Go(serviceMethod, args, reply, done.(chan *rpc.Call))
	return call
}

func (r *BaseRpc) InitClient() (err error) {
	r.client, err = jsonrpc.Dial(r.clientConn.network, r.clientConn.address)
	if err != nil {
		log.Debug("Dial error:", err)
		return
	}
	log.Info("RPC Client is connected to %s...", r.clientConn.address)
	go func() {
		defer r.client.Close()
		<-r.ctx.Done()
		log.Debug("close rpc client")
	}()
	return
}

func (r *BaseRpc) InitServer() (err error) {
	r.listener, err = net.Listen(r.serverConn.network, r.serverConn.address)
	if err != nil {
		log.Debug("Listen error:", err)
		return err
	}
	log.Info("RPC Server is listening on port %s...", r.serverConn.address)
	go r.serveConn()
	go func() {
		defer r.listener.Close()
		<-r.ctx.Done()
		log.Debug("close rpc listener")
	}()
	return
}

func (r *BaseRpc) serveConn() {
	for {
		// waits for and returns the next connection to the listener
		log.Debug("rpc listener waiting for connection...")
		conn, err := r.listener.Accept()
		if err != nil {
			log.Debug("Accept error:%s", err)
			return
		}
		log.Info("RPC Server is connected by %s...", conn.RemoteAddr())
		go func(conn net.Conn) {
			log.Debug("new client")
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
