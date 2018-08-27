package signal

import (
	"context"
	"os/signal"
	"os"
	"syscall"
	"log"
)

func RegisterDoneSignal() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ch := make(chan struct{})
	go func() {
		close(ch)
		for {
			select {
			case sig := <-sigCh:
				switch sig {
				case syscall.SIGINT, syscall.SIGTERM:
					cancel()
					log.Println("receive signal ", sig)
					goto EXIT
				default:
					log.Println("skip signal ", sig)
				}
			}
		}
	EXIT:
	}()
	<-ch
	return ctx
}
