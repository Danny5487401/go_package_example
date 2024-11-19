package _2_clockwork

import (
	"github.com/jonboulle/clockwork"
	"sync"
	"testing"
	"time"
)

func TestMyFunc(t *testing.T) {
	c := clockwork.NewFakeClock()

	// Start our sleepy function
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		myFunc(c)
		wg.Done()
	}()

	// Ensure we wait until myFunc is sleeping
	c.BlockUntil(1)

	// Advance the FakeClock forward in time
	c.Advance(3 * time.Second)

	// Wait until the function completes
	wg.Wait()

}
