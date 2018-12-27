// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defs

import (
	"github.com/astaxie/beego/logs"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	Broadcast chan SendFormat

	// Register requests from the clients.
	Register chan *Client

	// UnRegister requests from clients.
	UnRegister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan SendFormat),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			//注册
			h.switchType(SendFormat{Type: TYPE_LOGIN}, client)

		case client := <-h.UnRegister:
			//取消注册
			h.doUnregister(client)
		case message := <-h.Broadcast:
			h.switchType(message, nil)

		}
	}
}

func (h *Hub) doUnregister(client *Client) {
	if _, ok := h.clients[client.Id]; ok {
		delete(h.clients, client.Id)
		close(client.Send)
	}
}
func (h *Hub) switchType(message SendFormat, client *Client) {
	switch message.Type {
	case TYPE_LOGIN:
		h.doRegister(client)
	case TYPE_SEND_ID:
		h.doSendID(message)
	}
}

func (h *Hub) doRegister(client *Client) {

	if client == nil { //注册失败时,再次注册,client 为空,则走此流程
		return
	} else if client.Error != nil {
		logs.Error("注册失败 ", client.Error)
		client.Send <- NewSendError(client.Error.Error())
	} else {
		//成功
		h.clients[client.Id] = client
		client.Send <- NewSendSuccess(OKEY_LOGIN)
		logs.Debug("当前在线人数", len(h.clients))
	}
}

func (h *Hub) doSendID(message SendFormat) {
	if message.SendId == "" {
		h.clients[message.ClientId].Send <- message
		logs.Debug("没有 send_id")
		return
	}

	logs.Debug("message.SendId", message.SendId)

	if _, ok := h.clients[message.SendId]; !ok {
		message.Response = ERROR_NO_CLIENT_ID
		h.clients[message.ClientId].Send <- message
		logs.Debug("没有找到对应的人")
		return
	}
	h.clients[message.SendId].Send <- message
}
