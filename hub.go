// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"

)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}


func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			//注册
			h.actionRegister(client)
		case client := <-h.unregister:
			//取消注册
			h.actionUnregister(client)
		case message := <-h.broadcast:
			SendInfo :=sendinfo{}
			e := json.Unmarshal(message, &SendInfo)
			if e != nil {
				log.Println(" 发送失败")
				return
			}
			if SendInfo.Type == 2 && SendInfo.SendId!=""{
				log.Println("SendInfo.SendId",SendInfo.SendId)
				if _,ok:=h.clients[SendInfo.SendId];!ok {
					log.Println("没有找到对应的人")
					return
				}
				h.clients[SendInfo.SendId].send<-[]byte(SendInfo.Content)
			}

			////发送人数
			//for _,client := range h.clients {
			//	select {
			//	case client.send <- message:
			//	default:
			//		close(client.send)
			//		delete(h.clients, client.Id)
			//	}
			//}
		}
	}
}
func (h *Hub) actionRegister(client *Client) {
	//验证失败,不注册
	if client.error!=nil {
		client.send<-[]byte(client.error.Error())
		return
	}
	h.clients[client.Id] = client
	client.send <- []byte("注册成功")
	log.Println("当前在线人数",len(h.clients))
}

func (h *Hub) actionUnregister(client *Client){
	if _, ok := h.clients[client.Id]; ok {
		delete(h.clients, client.Id)
		close(client.send)
	}
}