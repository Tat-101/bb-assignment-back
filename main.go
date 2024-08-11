package main

import (
	"github.com/tat-101/bb-assignment-back/config"
	"github.com/tat-101/bb-assignment-back/database"
)

func main() {
	config.LoadConfig()

	database.Initialize()
}
