package main

import "encoding/json"

var h = hub{ //用于管理全部的websocket连接
	c: make(map[*connection]bool), //存在
	b: make(chan []byte),          //要广播的消息
	r: make(chan *connection),     //接收
	u: make(chan *connection),     //移除
}

type hub struct {
	c map[*connection]bool
	b chan []byte
	r chan *connection
	u chan *connection
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.r:
			h.c[c] = true
			c.data.Ip = c.ws.RemoteAddr().String()
			c.data.Type = "handshake"
			c.data.UserList = user_list
			data_b, _ := json.Marshal(c.data)
			c.sc <- data_b
		case c := <-h.u:
			if _, ok := h.c[c]; ok {
				delete(h.c, c)
				close(c.sc)
			}
		case data := <-h.b:
			for c := range h.c {
				select {
				case c.sc <- data: //消息广播
				default:
					delete(h.c, c)
					close(c.sc)
				}
			}
		}
	}
}
