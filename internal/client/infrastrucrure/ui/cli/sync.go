package cli

import (
	"fmt"
	"time"
)

var syncProgress = []string{"", ".", "..", "...", "..."}

func (c *CLI) Sync(stop <-chan struct{}) {
	var counter int
	for {
		select {
		case <-stop:
			return
		case <-time.Tick(time.Millisecond * 700):
			if counter == len(syncProgress) {
				counter = 0
			}
			c.Clear()
			fmt.Print("Synchronization" + syncProgress[counter])
			counter++
		}
	}
}
