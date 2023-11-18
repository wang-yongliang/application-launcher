package app

import (
	"context"
	"testing"
	"time"

	shutdown "github.com/lishimeng/go-app-shutdown"
	"github.com/wang-yongliang/application-launcher/token"
)

func TestTokenValidatorBuild001(t *testing.T) {

	time.AfterFunc(time.Second*60, func() {
		shutdown.Exit("bye bye")
	})

	_ = New().Start(func(ctx context.Context, builder *ApplicationBuilder) error {

		builder.
			//EnableCache().
			EnableTokenValidator(func(inject TokenValidatorInjectFunc) {
				prov := token.NewRedisStorage(GetCache())
				inject(prov)
			})
		return nil
	}, func(s string) {

	})

}
