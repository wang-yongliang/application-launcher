package factory

import (
	"context"

	"github.com/lishimeng/go-log"
	"github.com/wang-yongliang/application-launcher/amqp/rabbit"
	"github.com/wang-yongliang/application-launcher/cache"
	"github.com/wang-yongliang/application-launcher/mqtt"
	"github.com/wang-yongliang/application-launcher/persistence"
	"github.com/wang-yongliang/application-launcher/rpc"
)

const (
	amqpKey  = "amqp_session"
	cacheKey = "cache_redis"
	mqttKey  = "mqtt_redis"
	rpcKey   = "rpc_redis"
)

var globalContext context.Context

//var appCache cache.C

//var amqpSession rabbit.Session

func RegisterCtx(ctx context.Context) {
	globalContext = ctx
}

func GetCtx() (ctx context.Context) {
	ctx = globalContext
	return
}

func RegisterCache(c cache.C) {
	Add(&c, cacheKey)
}

func GetCache() (c cache.C) {
	err := Get(&c, cacheKey)
	if err != nil {
		log.Debug(err)
		c = nil
	}
	return
}

func RegisterAmqp(session rabbit.Session) {
	Add(&session, amqpKey)
}

func GetAmqp() (session rabbit.Session) {
	err := Get(&session, amqpKey)
	if err != nil {
		log.Debug(err)
		session = nil
	}
	return
}

func RegisterMqtt(session mqtt.Session) {
	Add(&session, mqttKey)
}
func GetMqtt() (session mqtt.Session) {
	err := Get(&session, mqttKey)
	if err != nil {
		log.Debug(err)
		session = nil
	}
	return
}

func GetOrm() *persistence.OrmContext {
	return persistence.New()
}
func GetNamedOrm(aliaName string) *persistence.OrmContext {
	return persistence.NewOrm(aliaName)
}

func RegisterRpc(session rpc.Session) {
	Add(&session, rpcKey)
}

func GetRpc() (session rpc.Session) {
	err := Get(&session, rpcKey)
	if err != nil {
		log.Debug(err)
		session = nil
	}
	return
}
