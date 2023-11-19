package app

import (
	"context"
	"testing"
	"time"

	shutdown "github.com/lishimeng/go-app-shutdown"
	"github.com/lishimeng/go-log"
	"github.com/wang-yongliang/application-launcher/rpc"
)

// Math 远程对象
type Math struct{}

// Args 参数结构体
type Args struct {
	A, B int
}

// Args 参数结构体
type Args2 struct {
	A, B int
}

// Multiply 乘法方法
func (m *Math) Multiply(args *[]Args, reply *int) error {
	t := 0
	for _, arg := range *args {
		t += arg.A * arg.B
	}
	*reply = t
	return nil
}

func TestBase(t *testing.T) {
	time.AfterFunc(time.Second*20, func() {
		shutdown.Exit("bye bye")
	})
	t.Log("start app")
	_ = New().Start(func(ctx context.Context, builder *ApplicationBuilder) error {
		builder.
			EnableRpcServer(rpc.WithServerAddress(":1234"), rpc.WithServerType(rpc.HTTP)).
			EnableRpcClient(rpc.WithClientAddress("localhost:1234"), rpc.WithClientType(rpc.HTTP)).
			RegisterRpcMethods(new(Math)).
			ComponentAfter(setup)
		return nil
	}, func(s string) {
		t.Log(s)
	})
	t.Log("end app")

}

func TestJsonHttpRpc(t *testing.T) {
	time.AfterFunc(time.Second*20, func() {
		shutdown.Exit("bye bye")
	})
	t.Log("start app")
	_ = New().Start(func(ctx context.Context, builder *ApplicationBuilder) error {
		builder.
			EnableRpcServer(rpc.WithServerAddress(":1234"), rpc.WithServerType(rpc.JSON_HTTP)).
			EnableRpcClient(rpc.WithClientAddress("http://localhost:1234"), rpc.WithClientType(rpc.JSON_HTTP)).
			RegisterRpcMethods(new(Math)).
			ComponentAfter(setup)
		return nil
	}, func(s string) {
		t.Log(s)
	})
	t.Log("end app")

}

func setup(context.Context) (err error) {
	args := &Args2{7, 8}
	var params []*Args2
	params = append(params, args)
	var reply int
	log.Info("call Math.Multiply")
	err = GetRpc().Call("Math.Multiply", params, &reply)
	if err != nil {
		log.Info("Call error:", err)
		return
	}
	log.Info("Math: %d*%d=%d\n", args.A, args.B, reply)
	// _, err = rpc2.Dial("tcp", "127.0.0.1:1234")
	// if err != nil {
	// 	log.Info("Call error:", err)
	// 	return
	// }
	// _, err = rpc2.Dial("tcp", "127.0.0.1:1234")
	// if err != nil {
	// 	log.Info("Call error:", err)
	// 	return
	// }
	// _, err = rpc2.Dial("tcp", "127.0.0.1:1234")
	// if err != nil {
	// 	log.Info("Call error:", err)
	// 	return
	// }
	return
}
