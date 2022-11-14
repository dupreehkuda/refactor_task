package main

import (
	"github.com/dupreehkuda/refactor_task/internal/service"
)

func main() {
	api := service.NewByConfig()
	api.Run()
}
