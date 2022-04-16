package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

const defaultAddr string = "127.0.0.1:8080"

func formatPort(portStr string) string {
	port, err := strconv.Atoi(portStr)
	if err != nil || (port > 65535 || port < 1000) {
		log.Printf("%s must be a valid port between 1000-65535", portStr)
		log.Printf("Defaulting to: %s", defaultAddr)
		return defaultAddr
	}
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func main() {
	listenAddr := defaultAddr

	if len(os.Args) > 1 {
		listenAddr = formatPort(os.Args[2])
	}

	log.Printf("Starting QtBot document server on port %s\n", listenAddr)

	documentServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", documentServer)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	go func() {
		signal := <-signals
		log.Printf("SIGNAL RECIEVED: %s", signal)
		os.Exit(1)
	}()

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
