package app

import (
	"context"
	"testing"
	"time"

	shutdown "github.com/lishimeng/go-app-shutdown"
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
func (m *Math) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
func TestServer(t *testing.T) {
	time.AfterFunc(time.Second*20, func() {
		shutdown.Exit("bye bye")
	})
	t.Log("start app")
	_ = New().Start(func(ctx context.Context, builder *ApplicationBuilder) error {

		builder.
			EnableRpcServer(rpc.WithServerAddress(":1234")).
			EnableRpcClient(rpc.WithClientAddress("127.0.0.1:1234")).
			RegisterRpcMethods(new(Math)).
			ComponentAfter(func(ctx context.Context) (err error) {
				args := &Args2{7, 8}
				var reply int

				GetRpc().Call("Math.Multiply", args, &reply)
				if err != nil {
					t.Log("Call error:", err)
					return
				}
				t.Logf("Math: %d*%d=%d\n", args.A, args.B, reply)
				return
			})
		return nil
	}, func(s string) {
		t.Log(s)
	})
	t.Log("end app")
}
