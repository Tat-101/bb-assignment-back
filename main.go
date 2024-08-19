package main

import (
	"github.com/tat-101/bb-assignment-back/config"
	"github.com/tat-101/bb-assignment-back/internal"
)

// TODO: improve log

func main() {
	r := internal.SetupServer()
	r.Run(":" + config.LoadConfig().ServerAddress)
}
