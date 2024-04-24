package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegom0ta/gofiber-study/internal/config"
	db "github.com/diegom0ta/gofiber-study/internal/database"
	"github.com/diegom0ta/gofiber-study/internal/http/server"
)

func main() {
	config.LoadDB()

	db.Connect()

	go func() {
		server.Run(3333)
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Println("Gracefully shutting down started...")

	server.Shutdown()

	if err := db.Disconnect(); err != nil {
		log.Printf("Database disconnection failed: %v", err)
	}

	fmt.Println("All dependencies are shutdown.")

}
