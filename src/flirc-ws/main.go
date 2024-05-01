package main

import (
	"log"
	"os"
)

func main() {

	parseFlags()
	if !IsSocket(*socket) {
		log.Printf("ERR: Invalid socket file  %v", *socket)
		os.Exit(1)
	}
	newHandler().Start()
}
