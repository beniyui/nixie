package general

import (
	"fmt"
	"time"
)


type spiSession struct {
	msg         string
	frames      []string
	doneChannel chan struct{}
}

func (s spiSession) start() {
	go spin(s)
}

func (s spiSession) stop() {
	close(s.doneChannel)
}

func spin(s spiSession) {

	//init
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	done := s.doneChannel
	frames := s.frames
	msg := s.msg

	i := 0
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Printf("\r%s %s", msg, frames[i%len(frames)])
			i++
		}

	}
}
