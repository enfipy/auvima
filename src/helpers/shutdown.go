package helpers

import (
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(close func()) {
	quitChan := make(chan os.Signal, 1)

	signal.Notify(quitChan, syscall.SIGTERM)
	signal.Notify(quitChan, syscall.SIGINT)
	signal.Notify(quitChan, syscall.SIGKILL)

	<-quitChan
	close()
}
