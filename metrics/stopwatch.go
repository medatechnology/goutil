package metrics

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Like stop watch to monitor the length of a process, usually used between function calls
// Usage:
// print the message and ticker
// swatch := metrics.StartTimeIt("Loading", 50)
// CallFunctionA()
// metrics.StopTimeItPrint(swatch, "Done")
//
// if no need to print (which automatically no ticker), call it like:
// swatch := metrics.StartTimeIt("", 50)
// CallFunctionA()
// metrics.StopTimeIt(swatch) // or metrics.StopTimeItPrint(swatch, "")
//

type TimeIt struct {
	ID         int64 // ID to match the StopTimer
	Message    string
	Start      time.Time
	WithTicker bool
	Ticker     *time.Ticker
	StopChan   chan bool
}

const (
	DOT_TICKER_DURATION = 250 // milliseconds
	DOT_SYMBOL          = "."
)

var watches sync.Map // Replaces the original map

// Start timer (stopwatch) returns ID for the timer
// If message="" (empty string) then it doesn't display
// anything and ticker automatically disabled!
// If ticker speed is 0, then use default value
func StartTimeIt(message string, tickerSpeed int) int64 {
	id := rand.Int63()
	withTicker := false
	if tickerSpeed >= 0 && message != "" {
		if tickerSpeed == 0 {
			tickerSpeed = DOT_TICKER_DURATION
		}
		withTicker = true
	}
	watch := TimeIt{
		ID:         id,
		Message:    message,
		Start:      time.Now(),
		WithTicker: withTicker,
	}

	if withTicker {
		watch.Ticker = time.NewTicker(time.Duration(tickerSpeed) * time.Millisecond)
		watch.StopChan = make(chan bool)

		go func(w TimeIt) {
			for {
				select {
				case <-w.Ticker.C:
					fmt.Print(DOT_SYMBOL)
				case <-w.StopChan:
					return
				}
			}
		}(watch)
	}
	if message != "" {
		fmt.Print(message)
	}

	watches.Store(id, watch) // Thread-safe store
	return id
}

// Stop the timer without printing anything
func StopTimeIt(id int64) time.Duration {
	val, ok := watches.Load(id) // Thread-safe load
	if !ok {
		return 0
	}
	watch := val.(TimeIt) // Type assertion

	elapsed := time.Since(watch.Start)

	if watch.WithTicker {
		watch.Ticker.Stop()
		watch.StopChan <- true
		<-time.After(10 * time.Millisecond)
	}

	watches.Delete(id) // Thread-safe delete
	return elapsed
}

// Stop timer with printing message, usually "Done"
// Output: ..... Done (duration) -- duration is with appropriate unit
// If message == "" then not printing anything
func StopTimeItPrint(id int64, message string) time.Duration {
	elapsed := StopTimeIt(id)
	if message != "" {
		// 	fmt.Printf("Done (%s)\n", elapsed)
		// } else {
		fmt.Printf(" %s (%s)\n", message, elapsed)
	}
	return elapsed
}
