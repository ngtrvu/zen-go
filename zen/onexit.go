package zen

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// OnExit interface for anything running when the program exits.
type OnExit interface {
	CleanUp()
}

// OnExitHook contains a list of functions to be run on program exit.
type OnExitHook struct {
	listeners []*OnExit
}

// CleanUp everything
func (onExitHook *OnExitHook) CleanUp() {
	log.Printf("Calling CleanUp on listeners!")
	for _, listener := range onExitHook.listeners {
		(*listener).CleanUp()
	}
}

// New intializes an exit hook
func New() *OnExitHook {
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	listeners := &OnExitHook{
		listeners: []*OnExit{},
	}

	go func() {
		select {
		case <-c:
			listeners.CleanUp()
			cancel()
		case <-ctx.Done():
		}
	}()
	return listeners
}

// AddListener add
func (onExitHook *OnExitHook) AddListener(onExit OnExit) {
	onExitHook.listeners = append(onExitHook.listeners, &onExit)
}
