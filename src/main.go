package main

import (
	"log"

	"github.com/google/uuid"
)

func main() {
	log.Printf("Greetings, auvima! Unique id: %s", uuid.New().String())
}
