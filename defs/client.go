// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defs

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	 "github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)


// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn
	Id   string
	Name string

	// Buffered channel of outbound messages.
	Send chan SendFormat
	Error  error
}

// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnRegister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		sendFormat := SendFormat{}
		if e := c.Conn.ReadJSON(&sendFormat);e!=nil {
			logs.Error(" 读取的json 格 式不正确 error %v  data: %v",e,sendFormat)
			c.Send<-NewSendError(ERROR_BAD_MSG_FORMAT)
			break
		}
		//accept
		c.Hub.Broadcast <- sendFormat
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logs.Error("NextWriter ",err)
				return
			}
			w.Write(message.Marshal())
			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				msg :=<-c.Send
				_, err := w.Write(msg.Marshal())
				if err != nil {
					c.Send<-NewSendSuccess(OKEY_MSG)
				}
			}
			if err := w.Close(); err != nil {
				logs.Error("w.Close",err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logs.Error("ticker.C error",err)
				return
			}
		}
	}
}

func (c *Client) AddClient(){
	_, p, _:= c.Conn.ReadMessage()
	data := SendFormat{}
	e := json.Unmarshal(p, &data)
	if e != nil {
		logs.Debug(e)
	}

	log.Println(data)
	//如果是 type 1 则是新增连接
	if data.Type==1{
		id := Checkuser(data.Name, data.Pwd)
		if id == "" {
			c.Error=errors.New(ERROR_LOGIN)
			return
		}
		c.Id=id
	}
}

