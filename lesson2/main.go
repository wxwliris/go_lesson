package main

import (
	"log"
)

func main() {
	server := NewServer()
	app := NewApp(server)
	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
