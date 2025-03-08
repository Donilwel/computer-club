package main

import (
	"computer-club/di"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container := di.NewContainer()
	container.RunServer(ctx)
}
