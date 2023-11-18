package rpc

import (
	"context"
	"net"
	"net/rpc"

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

func (r *BaseRpc) listenExit() {
	for {
		select {
		case <-r.ctx.Done():
			if r.listener != nil {
				log.Debug("close rpc listener")
				r.listener.Close()
			}
			if r.client != nil {
				log.Debug("close rpc client")
				r.client.Close()
			}
			return
		}
	}
}

func (r *BaseRpc) InitClient() (err error) {
	r.client, err = rpc.Dial(r.clientConn.network, r.clientConn.address)
	if err != nil {
		log.Debug("Dial error:", err)
		return err
	}
	log.Info("RPC Client is connected to %s...", r.clientConn.address)
	return nil
}

func (r *BaseRpc) InitServer() (err error) {
	r.listener, err = net.Listen(r.serverConn.network, r.serverConn.address)
	if err != nil {
		log.Debug("Listen error:", err)
		return err
	}
	log.Info("RPC Server is listening on port %s...", r.serverConn.address)
	go func() {
		defer r.listener.Close()
		conn, err := r.listener.Accept()
		if err != nil {
			log.Debug("Accept error:", err)
			return
		}
		log.Info("RPC Server is connected by %s...", conn.RemoteAddr())
		rpc.ServeConn(conn)
	}()
	return
}

func (r *BaseRpc) Call(serviceMethod string, args any, reply any) (err error) {
	return r.client.Call(serviceMethod, args, reply)
}

func (r *BaseRpc) Go(serviceMethod string, args any, reply any, done interface{}) interface{} {
	call := r.client.Go(serviceMethod, args, reply, done.(chan *rpc.Call))
	return call
}

// func a() {
// 	arith := new(Math)
// 	rpc.Register(arith)

// 	listener, err := net.Listen("tcp", ":1234")
// 	if err != nil {
// 		fmt.Println("Listen error:", err)
// 		return
// 	}

// 	fmt.Println("RPC Server is listening on port 1234...")
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Accept error:", err)
// 			continue
// 		}

// 		go rpc.ServeConn(conn)
// 	}
// }
