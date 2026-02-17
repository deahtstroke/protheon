package cli

import (
	"context"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/containerd/log"
)

type Cli struct {
	events chan tea.Msg

	eventsWg *sync.WaitGroup
	tuiWg    *sync.WaitGroup

	globalCtx    context.Context
	cleanupFuncs []func(context.Context) error
}

func NewCli(ctx context.Context) Cli {
	cli := Cli{
		events:    make(chan tea.Msg, 100),
		tuiWg:     &sync.WaitGroup{},
		eventsWg:  &sync.WaitGroup{},
		globalCtx: ctx,
	}

	// setupEvents
	cli.setupEvents()
	return cli
}

func (cli *Cli) setupEvents() {
	cleanupFunc := func(context.Context) error {
		cli.eventsWg.Wait()
		return nil
	}

	cli.cleanupFuncs = append(cli.cleanupFuncs, cleanupFunc)
}

func (cli *Cli) Subscribe(program *tea.Program) {
	cli.tuiWg.Add(1)

	tuiCtx, tuiCancel := context.WithCancel(cli.globalCtx)
	cli.cleanupFuncs = append(cli.cleanupFuncs, func(ctx context.Context) error {
		tuiCancel()
		cli.tuiWg.Wait()
		return nil
	})

	for {
		select {
		case <-tuiCtx.Done():
			return
		case msg, ok := <-cli.events:
			if !ok {
				return
			}
			program.Send(msg)
		}
	}
}

func (cli *Cli) Shutdown() {
	start := time.Now()
	defer func() {
		log.L.Printf("Shutdown took " + time.Since(start).String())
	}()

	shutdownCtx, cancel := context.WithCancel(cli.globalCtx)
	defer cancel()

	var wg sync.WaitGroup

	for _, cleanup := range cli.cleanupFuncs {
		if cli != nil {
			wg.Go(func() {
				if err := cleanup(shutdownCtx); err != nil {
					log.L.Printf("Failed to cleanup app properly on shutdown: %v", err)
				}
			})
		}
	}

	wg.Wait()
}
