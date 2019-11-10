package hub

import (
  "golang.org/x/net/websocket"
)

type Hub struct {
  Clients map[*websocket.Conn]bool
  Register chan *websocket.Conn
  Unregister chan *websocket.Conn
  Broadcast chan interface{}
  Shutdown chan bool
}

func New() *Hub {
  return &Hub{
    Register: make(chan *websocket.Conn),
    Unregister: make(chan *websocket.Conn),
    Broadcast: make(chan interface{}),
    Shutdown: make(chan bool),
    Clients: make(map[*websocket.Conn]bool),
  }
}

func (h *Hub) Add(c *websocket.Conn) {
  h.Clients[c] = true
}

func (h *Hub) Remove(c *websocket.Conn) {
  c.Close() // this is probably not necessary
  delete(h.Clients, c)
}

func (h *Hub) Publish(msg interface{}) {
  for c, active := range h.Clients {
    if active {
      err := websocket.JSON.Send(c, msg)

      if err != nil {
        // mark as inactive
        h.Clients[c] = !active
      }
    }
  }
}

func (h *Hub) Prune() int {
  // num of clients pruned
  var count int

  for c, active := range h.Clients {
    if !active {
      h.Remove(c)
      count++
    }
  }

  return count
}

func (h *Hub) Close() {
  for c, active := range h.Clients {
    if active {
      c.Close()
    }
  }
}
