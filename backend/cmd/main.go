package main

import (
	"github.com/niko-cb/uct/internal/infrastructure/web/server"
)

func main() {
	server.NewServer().Run()
}
