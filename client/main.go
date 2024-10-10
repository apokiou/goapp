package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	numConns := flag.Int("n", 1, "number of parallel connections")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
	log.Printf("connecting to %s\n", u.String())

	conns := make([]*websocket.Conn, *numConns)
	for i := range conns {
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal(err)
		}
		conns[i] = conn
		defer conn.Close()
	}

	for i, conn := range conns {
		go func(i int, conn *websocket.Conn) {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println(err)
					return
				}
				var response struct {
					Iteration int    `json:"iteration"`
					Value     string `json:"value"`
					Conn      string `json:"conn"`
				}
				err = json.Unmarshal(message, &response)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Printf("[conn #%d] iteration: %d, value: %s\n", i, response.Iteration, response.Value)
			}
		}(i, conn)
	}

	time.Sleep(time.Hour) // keep the program running
}
