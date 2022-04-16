package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func formatPort(port string) string {
	return fmt.Sprintf(":%s", port)
}

func main() {
	port := ":8080"

	if len(os.Args) > 1 {
		port = formatPort(os.Args[2])
	}

	fmt.Printf("Starting QtBot document server on port %s\n", port)

	documentServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", documentServer)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
