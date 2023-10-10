package main

import (
  "log"
  "os"
  "os/signal"
  "syscall"

  "github.com/DEEMMOONS/BaseHTTP/internal/server"
)

func handleInterrupt(s *server.Server) {
  c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Server interruption\n")
		s.Down()
		os.Exit(0)
	}()
}

func main() {
  server, err := server.NewServer("config/config.json")
  if err != nil {
    log.Fatal(err)
  }
	handleInterrupt(server)
	if err := server.Up(); err != nil {
		panic(err)
	}
}
