package signals

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Handler struct {
	cancel context.CancelFunc
}

func NewHandler(cancel context.CancelFunc) *Handler {
	return &Handler{
		cancel: cancel,
	}
}

func (h *Handler) SetupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, stopping...")
		h.cancel()
	}()
}
