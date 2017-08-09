/*
Copyright 2017 wechat-go Authors. All Rights Reserved.
MIT License

Copyright (c) 2017

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package youdao // 以插件名命令包名

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"webot-go/service"

	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs" // 导入日志包
)

// Register plugin
// 必须有的插件注册函数
func Register(session *service.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(service.MSG_TEXT, service.Handler(youdao), "youdao")
}

// 消息处理函数
func youdao(session *service.Session, msg *service.ReceivedMessage) {
	uri := "http://fanyi.youdao.com/openapi.do?keyfrom=go-aida&key=145986666&type=data&doctype=json&version=1.1&q=" + url.QueryEscape(msg.Content)
	response, err := http.Get(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		logs.Error(err)
		return
	}
	errorCode, err := jc.GetInt("errorCode")
	if err != nil {
		logs.Error(err)
		return
	}
	if errorCode != 0 {
		logs.Error("youdao API", errorCode)
		return
	}
	trans, err := jc.GetSliceString("translation")
	if err != nil {
		logs.Error(err)
		return
	}
	if len(trans) < 1 {
		return
	}

	session.SendText(msg.At+trans[0], session.Bot.UserName, service.RealTargetUserName(session, msg))
}
