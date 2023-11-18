package app

import (
	"context"

	"github.com/wang-yongliang/application-launcher/amqp/rabbit"
	"github.com/wang-yongliang/application-launcher/cache"
	"github.com/wang-yongliang/application-launcher/factory"
	"github.com/wang-yongliang/application-launcher/mqtt"
	"github.com/wang-yongliang/application-launcher/persistence"
	"github.com/wang-yongliang/application-launcher/rpc"
)

type Application interface {
	Start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error, onTerminate ...func(string)) error
}

func GetAmqp() (session rabbit.Session) {
	session = factory.GetAmqp()
	return
}

func GetMqtt() (session mqtt.Session) {
	session = factory.GetMqtt()
	return
}

func GetOrm() *persistence.OrmContext {
	return persistence.New()
}
func GetNamedOrm(aliaName string) *persistence.OrmContext {
	return persistence.NewOrm(aliaName)
}

func GetCache() (c cache.C) {
	c = factory.GetCache()
	return
}

func GetRpc() (session rpc.Session) {
	session = factory.GetRpc()
	return
}
