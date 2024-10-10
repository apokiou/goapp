package strgen

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"goapp/pkg/util"
	"log"
	"sync"
	"time"
)

type StringGenerator struct {
	strChan     chan<- string  // String output channel.
	quitChannel chan struct{}  // Quit.
	running     sync.WaitGroup // Running.
}

func New(strChan chan<- string) *StringGenerator {
	s := StringGenerator{}
	s.strChan = strChan
	s.quitChannel = make(chan struct{})
	s.running = sync.WaitGroup{}
	return &s
}

func generateRandomHexStrings(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// Start string generator. Stop() must be called at the end.
func (s *StringGenerator) Start() error {
	s.running.Add(1)
	go s.mainLoop()

	return nil
}

func (s *StringGenerator) Stop() {
	close(s.quitChannel)
	s.running.Wait()
}

func (s *StringGenerator) mainLoop() {
	defer s.running.Done()
	iteration := 1
	for {
		select {
		case s.strChan <- util.RandString(10):
		case <-s.quitChannel:
			return
		default:
			hexValue, err := generateRandomHexStrings(10) // Generate 10-character hex
			if err != nil {
				log.Println("Error generating hex:", err)
				return
			}

			s.strChan <- hexValue       // Send hex value to WebSocket handler
			time.Sleep(2 * time.Second) // Adjust delay between sending values
			iteration++
		}
	}
}
