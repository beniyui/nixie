package general

import (
	"fmt"
	"sync"
	"time"
)



type spiSession struct {
	msg      string
	frames   []string
	done     chan struct{}
	stopped  chan struct{}
	stopOnce sync.Once
}

func newSpiSession(msg string, frames []string) *spiSession {
	return &spiSession{
		msg:     msg,
		frames:  frames,
		done:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (s *spiSession) start() {
	if len(s.frames) == 0 {
		close(s.stopped) // 回すものが無いなら、stop が待ちぼうけにならないように
		return
	}
	go s.spin()
}

func (s *spiSession) spin() {

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	i := 0
	for {
		select {
		case <-s.done:
			fmt.Print("\r\x1b[2K") // 行を消してから帰る
			// s.stopOnce.Do(func(){close(s.stopped)})
			return
		case <-ticker.C:
			fmt.Printf("\r%s %s", s.msg, s.frames[i%len(s.frames)])
			i++
		}
	}
}

func (s *spiSession) stop() {
	s.stopOnce.Do(func() { close(s.done) })
	// <-s.stopped // 本当に終わるまで待つ
}
