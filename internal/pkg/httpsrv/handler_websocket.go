package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"goapp/internal/pkg/watcher"

	"github.com/gorilla/websocket"
)

type Response struct {
	Iteration int    `json:"iteration"`
	Value     string `json:"value"`
	Conn      string `json:"conn"`
}

func (s *Server) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	// Create and start a watcher.
	var watch = watcher.New()
	if err := watch.Start(); err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to start watcher: %w", err))
		return
	}
	defer watch.Stop()

	s.addWatcher(watch)
	defer s.removeWatcher(watch)

	// Start WS.
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to upgrade connection: %w", err))
		return
	}
	defer func() { _ = c.Close() }()

	strChan := make(chan string)

	gen := New(strChan)
	gen.Start()

	defer gen.Stop() // Ensure cleanup

	log.Printf("websocket started for watcher %s\n", watch.GetWatcherId())
	defer func() {
		log.Printf("websocket stopped for watcher %s\n", watch.GetWatcherId())
	}()

	// Read done.
	readDoneCh := make(chan struct{})

	// All done.
	doneCh := make(chan struct{})
	defer close(doneCh)

	var messageCount int
	var wg sync.WaitGroup
	var messageCountMutex sync.Mutex
	wg.Add(1)

	go func() {
		defer close(readDoneCh)
		defer wg.Done()
		for {
			select {
			default:
				_, p, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
						log.Printf("failed to read message: %v\n", err)
					}
					return
				}
				messageCountMutex.Lock()
				messageCount++
				messageCountMutex.Unlock()

				var m watcher.CounterReset
				if err := json.Unmarshal(p, &m); err != nil {
					log.Printf("failed to unmarshal message: %v\n", err)
					continue
				}
				watch.ResetCounter()
			case <-doneCh:
				return
			case <-s.quitChannel:
				return
			}
		}
	}()
	defer func() {
		log.Printf("websocket stopped for watcher %s, received %d messages\n", watch.GetWatcherId(), messageCount)
	}()

	for {
		select {
		case cv := <-watch.Recv():
			hexValue := fmt.Sprintf("%X", cv)
			data, _ := json.Marshal(struct {
				Iteration int    `json:"iteration"`
				Value     string `json:"value"`
				Conn      string `json:"conn"`
			}{
				Iteration: cv.Iteration,
				Value:     hexValue,
				Conn:      c.RemoteAddr().String(),
			})
			err = c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("failed to write message: %v\n", err)
				}
				return
			}
		case <-readDoneCh:
			wg.Wait()
			return
		case <-s.quitChannel:
			wg.Wait()
			return
		}
	}

}
