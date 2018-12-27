// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defs

import (
	"encoding/json"


	"github.com/astaxie/beego/logs"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Rnregister chan *Client

}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Rnregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}


func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			//注册
			h.actionRegister(client)
		case client := <-h.Rnregister:
			//取消注册
			h.actionUnregister(client)
		case message := <-h.Broadcast:
			SendInfo := SendFormat{}
			e := json.Unmarshal(message, &SendInfo)
			if e != nil {
				logs.Debug(" 发送失败")
				return
			}
			if SendInfo.Type == 2 && SendInfo.SendId!=""{
				logs.Debug("SendInfo.SendId",SendInfo.SendId)
				if _,ok:=h.clients[SendInfo.SendId];!ok {
					logs.Debug("没有找到对应的人")
					return
				}
				h.clients[SendInfo.SendId].Send<-[]byte(SendInfo.Content)
			}
		}
	}
}
func (h *Hub) actionRegister(client *Client) {
	//验证失败,不注册
	if client.Error!=nil {
		client.Send<-[]byte(client.Error.Error())
		return
	}
	h.clients[client.Id] = client
	client.Send <- []byte("注册成功")
	logs.Debug("当前在线人数",len(h.clients))
}

func (h *Hub) actionUnregister(client *Client){
	if _, ok := h.clients[client.Id]; ok {
		delete(h.clients, client.Id)
		close(client.Send)
	}
}