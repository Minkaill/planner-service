package main

import (
	"log"

	"github.com/Minkaill/planner-service.git/pkg/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
