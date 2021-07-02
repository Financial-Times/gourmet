package app

import (
	"os"
	"os/signal"
	"syscall"
)

// Terminator - function that determines when an app should be stopped
type Terminator func(chan struct{})

// OSSignalTerminator - Terminator that listens for OS signals SIGINT and
// SIGTERM as notifies the app that it should stop
func OSSignalTerminator(done chan struct{}) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		done <- struct{}{}
	}()
}
