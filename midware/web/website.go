package web

import (
	"github.com/kataras/iris/v12"
	"github.com/wang-yongliang/application-launcher/cms"
	"github.com/wang-yongliang/application-launcher/tool"
)

// Website 注入网站的基本配置
func Website(ctx iris.Context) {
	ws, err := cms.GetWebsiteInfo()
	if err != nil {
		ctx.Next()
		return
	}
	ctx.ViewData(tool.WebsiteCtx, ws)
	ctx.Next()
}
