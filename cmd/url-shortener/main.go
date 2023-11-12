package main

import (
	"api-shorter/internal/config"
	"fmt"
)

func main() {

	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init config: cleanenv

	// TODO: init logger: slog

	// TODO init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: run server

}
