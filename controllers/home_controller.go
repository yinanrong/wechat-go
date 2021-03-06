package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"webot-go/service"

	"github.com/songtianyi/rrframework/logs"
)

type HomeController struct {
	controller
}

func NewHomeController() *HomeController {
	h := new(HomeController)
	if h.mux == nil {
		h.mux = make(map[string]func(w http.ResponseWriter, r *http.Request))
	}
	h.mux["/"] = h.index
	h.mux["/qr"] = h.qr
	h.mux["/state"] = h.state
	return h
}
func (c *HomeController) index(w http.ResponseWriter, r *http.Request) {
	index := `<!DOCTYPE html><html><head><meta charset="utf-8" /><title>西风白马啸微信机器人</title><script src="https://unpkg.com/axios/dist/axios.min.js"></script></head><body><div id="img-box"><img id="img" style="height: 301px; width: 301px; margin-left: calc(50% - 100px); margin-top: 100px;" />	</div>	<div style="width: 100%; margin-top: 50px; text-align: center; font-size: 17px;" id="text">请使用微信扫描二维码</div>	<script type="text/javascript">axios.get("/qr").then(function (data) {	var imgBox = document.getElementById("img-box");var img = document.getElementById("img");img.src = data.data.qr;var interval = setInterval(function () {axios.get("/state?uuid=" + data.data.uuid).then(function (res) {if (res.data == 408 || res.data == 500 || res.data == 600 || res.data == 200) {clearInterval(interval);switch (res.data) {case 200:document.getElementById("text").innerHTML = "<font color='green'>登录成功</font>";break;case 201:document.getElementById("text").innerHTML = "请在微信上点击确认";break;case 408:document.getElementById("text").innerHTML = "<font color='red'>扫码超时</font>";break;case 500:document.getElementById("text").innerHTML = "<font color='red'>会话失败</font>";break;case 600:document.getElementById("text").innerHTML = "<font color='red'>已登出</font>";break;}}});}, 600);}).catch(function (err) {alert(err);})</script></body></html>`
	c.View(index, w, r)
}
func (c *HomeController) qr(w http.ResponseWriter, r *http.Request) {
	session, err := service.CreateSession(nil, nil)
	if err != nil {
		logs.Error(err)
		c.BadRequest(w, err)
		return
	}
	session.EnQueue()
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"uuid":"%s","qr":"%s"}`, session.ID, session.Qr())
}
func (c *HomeController) state(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uuid := r.FormValue("uuid")
	if session, ok := service.SessionPop(uuid); ok {
		w.Write([]byte(strconv.Itoa(session.State)))
	} else {
		w.Write([]byte(strconv.Itoa(service.Closed)))
	}
}
