package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"testing"
)

// Math 远程对象
type Math struct{}

// Args 参数结构体
type Args struct {
	A, B int
}

// Multiply 乘法方法
func (m *Math) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func TestRpcServer(t *testing.T) {
	math := new(Math)
	rpc.Register(math)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}

	fmt.Println("RPC Server is listening on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}

func TestRpcClient(t *testing.T) {
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	defer client.Close()

	args := &Args{7, 8}
	var reply int
	err = client.Call("Math.Multiply", args, &reply)
	if err != nil {
		fmt.Println("Call error:", err)
		return
	}

	fmt.Printf("Math: %d*%d=%d\n", args.A, args.B, reply)
}
