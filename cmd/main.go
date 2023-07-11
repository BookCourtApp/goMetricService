package main

import (
	"fmt"
	"mserv/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
