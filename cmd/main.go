package main

import (
	"log"

	"github.com/Minkaill/planner-service.git/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
