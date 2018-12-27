// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"gowebsocket/defs"
	"gowebsocket/handers"
	"net/http"

	"github.com/astaxie/beego/logs"
)

var addr = flag.String("addr", ":8080", "http service address")

func init(){


	flag.Parse()
	logs.Info("服务器开启 端口为 port:%v",*addr)
}

func main() {

	hub := defs.NewHub()
	go hub.Run()
	http.HandleFunc("/", handers.ServeHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handers.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		logs.Error("ListenAndServe: ", err)
	}
}
