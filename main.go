package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	httpServer := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	fmt.Println("Starting server on port 8080")
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

}
