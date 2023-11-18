package app

import (
	"context"
	"fmt"

	rpc2 "net/rpc"

	shutdown "github.com/lishimeng/go-app-shutdown"
	"github.com/wang-yongliang/application-launcher/amqp"
	"github.com/wang-yongliang/application-launcher/application/api"
	"github.com/wang-yongliang/application-launcher/application/repo"
	"github.com/wang-yongliang/application-launcher/cache"
	"github.com/wang-yongliang/application-launcher/factory"
	"github.com/wang-yongliang/application-launcher/midware/auth"
	"github.com/wang-yongliang/application-launcher/mqtt"
	"github.com/wang-yongliang/application-launcher/rpc"
	"github.com/wang-yongliang/application-launcher/server"
	"github.com/wang-yongliang/application-launcher/token"
)

type application struct {
	builder *ApplicationBuilder
}

func New() (instance Application) {
	factory.RegisterCtx(shutdown.Context())
	builder := &ApplicationBuilder{}
	ins := &application{builder: builder}
	instance = ins
	return
}

func (h *application) Start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error, onTerminate ...func(string)) (err error) {

	err = h._start(buildHandler)

	if err == nil {
		shutdown.WaitExit(&shutdown.Configuration{
			BeforeExit: func(s string) {
				if len(onTerminate) > 0 {
					for _, terminateFunc := range onTerminate {
						if terminateFunc != nil {
							terminateFunc(s)
						}
					}
				}
			},
		})
	}
	return
}

func (h *application) _start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error) (err error) {

	if buildHandler == nil {
		err = fmt.Errorf("application builder function nil")
		return
	}
	err = buildHandler(factory.GetCtx(), h.builder)
	if err != nil {
		return
	}

	// 初始化amqp连接
	if h.builder.amqpEnable {
		factory.RegisterAmqp(amqp.New(factory.GetCtx(), h.builder.amqpOptions, h.builder.sessionOptions...))
	}

	if h.builder.dbEnable {
		err = repo.Database(h.builder.dbConfig, h.builder.dbModels...)
		if err != nil {
			return
		}
	}

	if h.builder.cacheEnable {
		factory.RegisterCache(cache.New(factory.GetCtx(), h.builder.redisOpts, h.builder.cacheOpts))
	}

	if h.builder.tokenValidatorEnable {
		h.builder.tokenValidatorBuilder(func(storage token.Storage) {
			if storage != nil {
				auth.TokenStorage = storage
			}
		})
	}

	// 启动amqp业务
	if h.builder.amqpEnable {
		// 在线程中启动每一个handler
		for _, h := range h.builder.amqpHandler {
			go amqp.RegisterHandler(factory.GetAmqp(), h)
		}
	}

	if h.builder.mqttEnable {
		session := mqtt.New(factory.GetCtx(), h.builder.mqttOptions...)
		factory.RegisterMqtt(session)
		err = session.Connect()
		if err != nil {
			return err
		}
	}

	err = h.applyComponents(h.builder.componentsBeforeWebServer)
	if err != nil {
		return err
	}

	if h.builder.webEnable {
		var srv *server.Server
		conf := server.Config{
			Listen: h.builder.webListen,
		}
		if len(h.builder.webLogLevel) > 0 {
			conf.LogLvl = h.builder.webLogLevel
		}
		srv, err = api.Server(conf)
		if h.builder.webStaticEnable {
			err = api.EnableStatic(srv,
				h.builder.assetFile)
			if err != nil {
				return
			}
		}

		err = api.EnableComponents(srv, h.builder.webComponents...)
		if err != nil {
			return
		}

		err = api.EnableMonitors(srv)
		if err != nil {
			return
		}
		err = api.Start(factory.GetCtx(), srv)
		if err != nil {
			return
		}
	}

	if h.builder.rpcServerEnable {
		for _, m := range h.builder.rpcMethods {
			err = rpc2.Register(m)
			if err != nil {
				return
			}
		}
		err = rpc.NewServer(factory.GetCtx(), h.builder.rpcServerOpts...)
		if err != nil {
			return
		}
	}

	if h.builder.rpcCilentEnable {
		var session rpc.Session
		session, err = rpc.NewClient(factory.GetCtx(), h.builder.rpcClientOpts...)
		if err != nil {
			return
		}
		factory.RegisterRpc(session)
	}

	err = h.applyComponents(h.builder.componentsAfterWebServer)
	if err != nil {
		return err
	}

	return
}
