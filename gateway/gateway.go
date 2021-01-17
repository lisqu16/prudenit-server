package gateway

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		t, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println(string(p))

		/*if err := conn.WriteMessage(t, p); err != nil {
			log.Fatal(err)
			return
		}*/
	}
}

func Gateway(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected")
	reader(ws)
}