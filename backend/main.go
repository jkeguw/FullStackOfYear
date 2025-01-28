package main

import (
	"FullStackOfYear/backend/config"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal("init config failed:", err)
	}
	defer config.Logger.Sync()
}
