/*
Copyright © 2025 Daniel Villavicencio<dvm3099@pm.me>
*/
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/deahtstroke/protheon/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cmd.ProtheonMain(ctx)
}
