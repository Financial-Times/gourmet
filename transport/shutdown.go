package transport

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdownSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
