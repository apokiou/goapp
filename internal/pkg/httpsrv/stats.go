package httpsrv

import "log"

type sessionStats struct {
	id   string
	sent int
}

func (w *sessionStats) print() {
	message := "message"
	if w.sent != 1 {
		message = "messages"
	}
	log.Printf("session %s has received %d %s\n", w.id, w.sent, message)
}

func (w *sessionStats) inc() {
	w.sent++
}

func (s *Server) incStats(id string) {
	for _, ws := range s.sessionStats {
		if ws.id == id {
			ws.inc()
			return
		}
	}
	s.sessionStats = append(s.sessionStats, sessionStats{id: id, sent: 1})

}
