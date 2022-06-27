package main

import (
	"log"
)

func main() {
	router := NewRouter()
	err := router.Run()
	if err != nil {
		log.Fatalf("API start failure: %s", err)
	}
}
