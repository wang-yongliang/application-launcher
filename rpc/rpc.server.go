package rpc

import (
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/lishimeng/go-log"
)

func (r *BaseRpc) newServer() (err error) {
	log.Debug("init server, type=%s", string(r.serverConn.connType))
	if r.serverConn.connType == Base {
		r.listener, err = net.Listen(r.serverConn.network, r.serverConn.address)
		if err != nil {
			log.Debug("Listen error:", err)
			return err
		}
		log.Info("[BASE]RPC Server is listening on port %s...", r.serverConn.address)
		go r.serveConn() //监听连接
	} else if r.serverConn.connType == HTTP {
		rpc.HandleHTTP() // registers an HTTP handler for RPC
		r.listener, err = net.Listen(r.serverConn.network, r.serverConn.address)
		if err != nil {
			log.Debug("Listen error:", err)
			return err
		}
		log.Info("[JSON]RPC Server is listening on port %s...", r.serverConn.address)
		go http.Serve(r.listener, nil) //监听连接
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var conn io.ReadWriteCloser = struct {
				io.ReadCloser
				io.Writer
			}{
				ReadCloser: r.Body,
				Writer:     w,
			}
			err = rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
			if err != nil {
				log.Debug("ServeRequest error:%s", err)
			}
		})
		log.Info("[JSON-HTTP]RPC Server is listening on port %s...", r.serverConn.address)
		go http.ListenAndServe(r.serverConn.address, nil) //监听连接
		return
	}
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
		log.Info("RPC Server is connected by %s...", conn.RemoteAddr().String())
		go func(conn net.Conn) {
			log.Debug("Accept new client")
			rpc.ServeConn(conn)
		}(conn)
	}
}
