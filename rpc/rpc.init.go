package rpc

import (
	"context"

	"github.com/lishimeng/go-log"
)

// type RpcOptions func(r *BaseRpc)

// func WithNetwork(network string) RpcOptions {
// 	return func(r *BaseRpc) {
// 		r.network = network
// 	}
// }

// func WithAddress(address string) RpcOptions {
// 	return func(r *BaseRpc) {
// 		r.address = address
// 	}
// }

// type RpcOptions struct {
// 	Network string
// 	Address string
// }

type RpcServerOptions func(r *BaseRpc)

func WithServerNetwork(network string) RpcServerOptions {
	return func(r *BaseRpc) {
		r.serverConn.network = network
	}
}
func WithServerAddress(address string) RpcServerOptions {
	return func(r *BaseRpc) {
		r.serverConn.address = address
	}
}
func WithServerType(t connType) RpcServerOptions {
	return func(r *BaseRpc) {
		r.serverConn.connType = t
	}
}

type RpcClientOptions func(r *BaseRpc)

func WithClientNetwork(network string) RpcClientOptions {
	return func(r *BaseRpc) {
		r.clientConn.network = network
	}
}
func WithClientAddress(address string) RpcClientOptions {
	return func(r *BaseRpc) {
		r.clientConn.address = address
	}
}
func WithClientType(t connType) RpcClientOptions {
	return func(r *BaseRpc) {
		r.clientConn.connType = t
	}
}

func NewClient(ctx context.Context, rpcOpts ...RpcClientOptions) (session Session, err error) {
	baseRpc := BaseRpc{
		ctx: ctx,
	}
	for _, opt := range rpcOpts {
		opt(&baseRpc)
	}
	if baseRpc.clientConn.network == "" {
		baseRpc.clientConn.network = "tcp"
	}
	if baseRpc.clientConn.address == "" {
		baseRpc.clientConn.address = ":80"
	}
	if baseRpc.clientConn.connType == "" {
		baseRpc.clientConn.connType = JSON_HTTP
	}
	err = baseRpc.newClient()
	if err != nil {
		log.Debug("rpc init client error: %s", err.Error())
		return
	}
	session = &baseRpc
	return
}
func NewServer(ctx context.Context, rpcOpts ...RpcServerOptions) (err error) {
	baseRpc := BaseRpc{
		ctx: ctx,
	}
	for _, opt := range rpcOpts {
		opt(&baseRpc)
	}
	if baseRpc.serverConn.network == "" {
		baseRpc.serverConn.network = "tcp"
	}
	if baseRpc.serverConn.address == "" {
		baseRpc.serverConn.address = ":80"
	}
	if baseRpc.serverConn.connType == "" {
		baseRpc.serverConn.connType = JSON_HTTP
	}
	// baseRpc := NewBaseRpc(ctx, rpcOpts...)
	err = baseRpc.newServer()
	if err != nil {
		log.Debug("rpc init server error: %s", err.Error())
	}
	return
}
