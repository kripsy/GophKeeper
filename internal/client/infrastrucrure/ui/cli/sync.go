package cli

import (
	"fmt"
	"time"
)

var syncProgress = []string{"", ".", "..", "...", "..."}

func (c CLI) Sync(done <-chan struct{}) {
	var counter int
	for {
		select {
		case <-done:
			return
		default:
			if counter == len(syncProgress) {
				counter = 0
			}
			c.Clear()
			fmt.Print("Synchronization" + syncProgress[counter])
			counter++
			time.Sleep(time.Millisecond * 700)
		}
	}
}
