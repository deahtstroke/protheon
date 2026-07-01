package load

import (
	"log/slog"
	"sync/atomic"
)

type DryRunLoader struct {
	inner   Loader
	counter atomic.Int32
}

func NewDryRunLoader(innner Loader) Loader {
	return &DryRunLoader{
		inner: innner,
	}
}

func (l *DryRunLoader) Load(input any) error {
	l.counter.Add(1)
	slog.Info("[dry-run] would load record into target (skipped)", "Record", input, "Counter", l.counter.Load())
	return nil
}

func (l *DryRunLoader) Close() error {
	return l.inner.Close()
}
