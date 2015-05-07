/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
)

var routes *mvc.Route = mvc.NewRoute(nil)

//处理请求
func Handle(ctx *web.Context) {
	routes.Handle(ctx)
}

//注册路由
func RegisterRoutes() {
	bc := &baseC{}
	mc := &mainC{Base: bc} //入口控制器
	lc := &loginC{}
	routes.RegisterController("shop", &shopC{Base: bc})   //商家门店控制器
	routes.RegisterController("goods", &goodsC{Base: bc}) //商品控制器
	routes.RegisterController("comm", &commC{Base: bc})
	routes.RegisterController("order", &orderC{})
	routes.RegisterController("category", &categoryC{Base: bc})
	routes.RegisterController("conf", &configC{Base: bc})
	routes.RegisterController("prom", &promC{Base: bc})
    routes.RegisterController("delivery",&converageAreaC{Base: bc}) //配送区域控制器

	routes.Add("^/export/getExportData$", func(ctx *web.Context) {
		if b, id := chkLogin(ctx); b {
			GetExportData(ctx, id)
		} else {
			redirect(ctx)
		}
	})

	routes.Add("^/login$", func(ctx *web.Context) {
		mvc.Handle(lc, ctx, true)
	})

	routes.Add("^/[^/]*$", func(ctx *web.Context) {
		if b, id := chkLogin(ctx); b {
			mvc.Handle(mc, ctx, true, id)
		} else {
			redirect(ctx)
		}
	})

}
