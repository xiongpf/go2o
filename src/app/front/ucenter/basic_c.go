/**
 * Copyright 2015 @ S1N1 Team.
 * name : basic_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ucenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"html/template"
	"strconv"
)

var _ mvc.Filter = new(baseC)

type basicC struct {
	*baseC
}

func (this *basicC) Profile(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf, _ := this.GetSiteConf(p.Id, p.Secret)
	mm := this.GetMember(ctx)
	js, _ := json.Marshal(mm)
	ctx.App.Template().Execute(ctx.ResponseWriter, func(m *map[string]interface{}) {
		v := *m
		v["partner"] = p
		v["conf"] = conf
		v["partner_host"] = conf.Host
		v["member"] = mm
		v["entity"] = template.JS(js)

	}, "views/ucenter/profile.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *basicC) Pwd(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf, _ := this.GetSiteConf(p.Id, p.Secret)
	mm := this.GetMember(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter, func(m *map[string]interface{}) {
		v := *m
		v["partner"] = p
		v["conf"] = conf
		v["partner_host"] = conf.Host
		v["member"] = mm

	}, "views/ucenter/pwd.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *basicC) Pwd_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	m := this.GetMember(ctx)
	var result gof.JsonResult
	r.ParseForm()
	var oldPwd, newPwd, rePwd string
	oldPwd = r.FormValue("OldPwd")
	newPwd = r.FormValue("NewPwd")
	rePwd = r.FormValue("RePwd")
	var err error
	if newPwd != rePwd {
		err = errors.New("两次密码输入不一致")
	} else {
		err = dps.MemberService.ModifyPassword(m.Id, oldPwd, newPwd)
	}
	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true}
	}
	w.Write(result.Marshal())
}

func (this *basicC) Profile_post(ctx *web.Context) {
	mm := this.GetMember(ctx)
	r, w := ctx.Request, ctx.ResponseWriter
	var result gof.JsonResult
	r.ParseForm()
	clientM := new(member.ValueMember)
	web.ParseFormToEntity(r.Form, clientM)
	clientM.Id = mm.Id
	_, err := goclient.Member.SaveMember(clientM, mm.LoginToken)

	if err != nil {
		result = gof.JsonResult{Result: false, Message: err.Error()}
	} else {
		result = gof.JsonResult{Result: true}
	}
	w.Write(result.Marshal())
}

func (this *basicC) Deliver(ctx *web.Context) {
	p := this.GetPartner(ctx)
	conf, _ := this.GetSiteConf(p.Id, p.Secret)
	m := this.GetMember(ctx)
	ctx.App.Template().Execute(ctx.ResponseWriter, func(mp *map[string]interface{}) {
		v := *mp
		v["partner"] = p
		v["conf"] = conf
		v["partner_host"] = conf.Host
		v["member"] = m

	}, "views/ucenter/deliver.html",
		"views/ucenter/inc/header.html",
		"views/ucenter/inc/menu.html",
		"views/ucenter/inc/footer.html")
}

func (this *basicC) Deliver_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	addrs, err := goclient.Member.GetDeliverAddrs(m.Id, m.LoginToken)
	if err != nil {
		ctx.ResponseWriter.Write([]byte("{error:'错误:" + err.Error() + "'}"))
		return
	}
	js, _ := json.Marshal(addrs)
	ctx.ResponseWriter.Write([]byte(`{"rows":` + string(js) + `}`))
}

func (this *basicC) SaveDeliver_post(ctx *web.Context) {
	m := this.GetMember(ctx)
	w, r := ctx.ResponseWriter, ctx.Request
	r.ParseForm()
	var e member.DeliverAddress
	web.ParseFormToEntity(r.Form, &e)
	e.MemberId = m.Id
	b, err := goclient.Member.SaveDeliverAddr(m.Id, m.LoginToken, &e)
	if err == nil {
		if b {
			w.Write([]byte(`{"result":true}`))
		} else {
			w.Write([]byte(`{"result":false}`))
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
	}
}

func (this *basicC) DeleteDeliver_post(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	m := this.GetMember(ctx)
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))

	b, err := goclient.Member.DeleteDeliverAddr(m.Id, m.LoginToken, id)
	if err == nil {
		if b {
			w.Write([]byte(`{"result":true}`))
		} else {
			w.Write([]byte(`{"result":false}`))
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
	}
}
