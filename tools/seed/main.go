package main

import (
	"github.com/tat-101/bb-assignment-back/tools/seed/seed"
)

func main() {
	// Seed the database with an admin user
	seed.SeedAdminUser()
}
