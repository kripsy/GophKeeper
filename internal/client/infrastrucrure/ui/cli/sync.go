//nolint:gochecknoglobals
package cli

import (
	"fmt"
	"time"
)

const syncDisplayDelay = 700

var syncProgress = []string{"", ".", "..", "...", "..."}

// Sync displays sync status.
func (c *CLI) Sync(stop <-chan struct{}) {
	var counter int
	syncDisplayTicker := time.NewTicker(time.Millisecond * syncDisplayDelay)
	defer syncDisplayTicker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-syncDisplayTicker.C:
			if counter == len(syncProgress) {
				counter = 0
			}
			c.Clear()
			fmt.Print("Synchronization" + syncProgress[counter])
			counter++
		}
	}
}
