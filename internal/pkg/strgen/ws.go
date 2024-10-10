package strgen

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // Remove the explicit call to WriteHeader
		return
	}
	defer conn.Close()

	w.WriteHeader(http.StatusOK)
}

func handleConnection(conn *websocket.Conn) {
	strChan := make(chan string)
	s := New(strChan)
	s.Start()

	for {
		select {
		case msg := <-strChan:
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
