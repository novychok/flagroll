package main

import (
	"context"
	"fmt"

	"github.com/novychok/flagroll/platform/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, cleanup, err := internal.Init()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Start(ctx); err != nil {
		fmt.Printf("app start error: %v\n", err)
	}

}
