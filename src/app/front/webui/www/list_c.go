/**
 * Copyright 2015 @ S1N1 Team.
 * name : list_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package www

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof/web"
	"go2o/src/app/cache/apicache"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/goclient"
	"html/template"
	"strconv"
)

type listC struct {
	*baseC
}

func (this *listC) Index(ctx *web.Context) {
	_, w := ctx.Request, ctx.ResponseWriter
	p, _ := this.GetPartner(ctx)
	mm := this.GetMember(ctx)
	if b, siteConf := GetSiteConf(w, p); b {
		categories := apicache.GetCategories(ctx.App, p.Id, p.Secret)
		ctx.App.Template().Execute(w, func(m *map[string]interface{}) {
			(*m)["partner"] = p
			(*m)["title"] = "在线订餐-" + p.Name
			(*m)["categories"] = template.HTML(categories)
			(*m)["member"] = mm
			(*m)["conf"] = siteConf
		},
			"views/web/www/list.html",
			"views/web/www/inc/header.html",
			"views/web/www/inc/footer.html")
	}
}

func (this *listC) GetList(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	p, _ := this.GetPartner(ctx)
	const getNum int = -1 //-1表示全部
	categoryId, err := strconv.Atoi(r.URL.Query().Get("cid"))
	if err != nil {
		w.Write([]byte(`{"error":"yes"}`))
		return
	}
	items, err := goclient.Partner.GetItems(p.Id, p.Secret, categoryId, getNum)
	if err != nil {
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	buf := bytes.NewBufferString("<ul>")

	for _, v := range items {

		buf.WriteString(fmt.Sprintf(`
			<li>
				<div class="gs_goodss">
                        <img src="%s" alt="%s"/>
                        <h3 class="name">%s%s</h3>
                        <span class="srice">原价:￥%s</span>
                        <span class="sprice">优惠价:￥%s</span>
                        <a href="javascript:cart.add(%d,1);" class="add">&nbsp;</a>
                </div>
             </li>
		`, format.GetGoodsImageUrl(v.Image), v.Name, v.Name, v.SmallTitle, format.FormatFloat(v.Price),
			format.FormatFloat(v.SalePrice),
			v.Id))
	}
	buf.WriteString("</ul>")
	w.Write(buf.Bytes())
}
